package ome

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/republicprotocol/republic-go/smpc"
)

// ComputationBacklogExpiry defines how long the Ome will wait for an
// order.Fragment before rejecting a Computation.
const ComputationBacklogExpiry = 5 * time.Minute

// OmeBufferLimit defines the buffer size used by the Ome when reading data.
const OmeBufferLimit = 1024

// An Ome runs the logic for a single node in the secure order matching engine.
type Ome interface {

	// Run the secure order matching engine until the done channel is closed.
	Run(done <-chan struct{}) <-chan error

	// OnChangeEpoch should be called whenever a new Epoch is observed.
	OnChangeEpoch(registry.Epoch)
}

type ome struct {
	addr             identity.Address
	gen              ComputationGenerator
	matcher          Matcher
	confirmer        Confirmer
	settler          Settler
	computationStore ComputationStorer
	orderbook        orderbook.Orderbook
	smpcer           smpc.Smpcer

	computationBacklogMu *sync.RWMutex
	computationBacklog   map[ComputationID]Computation

	epochMu   *sync.RWMutex
	epochCurr *registry.Epoch
	epochPrev *registry.Epoch
}

// NewOme returns an Ome that uses an order.Orderbook to synchronize changes
// from the Ethereum blockchain, and an smpc.Smpcer to run the secure
// multi-party computations necessary for the secure order matching engine.
func NewOme(addr identity.Address, gen ComputationGenerator, matcher Matcher, confirmer Confirmer, settler Settler, computationStore ComputationStorer, orderbook orderbook.Orderbook, smpcer smpc.Smpcer, epoch registry.Epoch) Ome {
	ome := &ome{
		addr:             addr,
		gen:              gen,
		matcher:          matcher,
		confirmer:        confirmer,
		settler:          settler,
		computationStore: computationStore,
		orderbook:        orderbook,
		smpcer:           smpcer,

		computationBacklogMu: new(sync.RWMutex),
		computationBacklog:   map[ComputationID]Computation{},

		epochMu:   new(sync.RWMutex),
		epochCurr: nil,
		epochPrev: nil,
	}
	ome.OnChangeEpoch(epoch)
	return ome
}

// Run implements the Ome interface.
func (ome *ome) Run(done <-chan struct{}) <-chan error {
	matches := make(chan Computation, OmeBufferLimit)
	errs := make(chan error, OmeBufferLimit)

	var wg sync.WaitGroup

	notifications, orderbookErrs := ome.orderbook.Sync(done)
	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Forward(done, orderbookErrs, errs)
	}()

	computations, genErrs := ome.gen.Generate(done, notifications)
	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Forward(done, genErrs, errs)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case computation, ok := <-computations:
				if !ok {
					return
				}
				if err := ome.sendComputationToMatcher(computation, done, matches); err != nil {
					select {
					case <-done:
						return
					case errs <- err:
					}
				}
			}
		}
	}()

	// Sync the Confirmer to the Settler
	wg.Add(1)
	go func() {
		defer wg.Done()
		ome.syncConfirmerToSettler(done, matches, errs)
	}()

	// Retry Computations that failed due to a missing order.Fragment
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(14 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
			case <-ticker.C:
			}

			ome.syncOrderFragmentBacklog(done, matches)
		}
	}()

	// Cleanup
	go func() {
		defer close(errs)
		wg.Wait()
	}()

	return errs
}

// OnChangeEpoch updates the Ome to the next Epoch. This will cause
// cascading changes throughout the Ome, most notably it will connect to a new
// Smpc network that will handle future Computations.
func (ome *ome) OnChangeEpoch(epoch registry.Epoch) {
	ome.epochMu.Lock()
	defer ome.epochMu.Unlock()

	// Do not update if the epoch has not actually changed
	if ome.epochCurr != nil && bytes.Equal(epoch.Hash[:], ome.epochCurr.Hash[:]) {
		return
	}

	go func() {
		// Connect to the new network
		pod, err := epoch.Pod(ome.addr)
		if err != nil {
			logger.Error(fmt.Sprintf("cannot find pod: %v", err))
			return
		}
		ome.smpcer.Connect(epoch.Hash, pod.Darknodes)

		ome.orderbook.OnChangeEpoch(epoch)
		ome.gen.OnChangeEpoch(epoch)

		// Wait for some time to allow for the connections to begin
		time.Sleep(14 * time.Second)

		// Replace the previous epoch and disconnect from it
		if ome.epochPrev != nil {
			ome.smpcer.Disconnect(ome.epochPrev.Hash)
		}
		ome.epochPrev = ome.epochCurr
		ome.epochCurr = &epoch
	}()
}

func (ome *ome) syncConfirmerToSettler(done <-chan struct{}, matches <-chan Computation, errs chan<- error) {
	confirmations, confirmationErrs := ome.confirmer.Confirm(done, matches)
	for {
		select {
		case <-done:
			return

		case confirmation, ok := <-confirmations:
			if !ok {
				return
			}
			ome.sendComputationToSettler(confirmation)

		case err, ok := <-confirmationErrs:
			if !ok {
				return
			}
			select {
			case <-done:
			case errs <- err:
			}
		}
	}
}

func (ome *ome) syncOrderFragmentBacklog(done <-chan struct{}, matches chan<- Computation) {
	ome.computationBacklogMu.Lock()
	defer ome.computationBacklogMu.Unlock()

	buffer := [OmeBufferLimit]Computation{}
	bufferN := 0

	// Build a buffer of Computations that will be retried
	for _, com := range ome.computationBacklog {
		delete(ome.computationBacklog, com.ID)
		// Check for expiry of the Computation
		if com.Timestamp.Add(ComputationBacklogExpiry).Before(time.Now()) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("⧖ expired backlog computation buy = %v, sell = %v", com.Buy, com.Sell))
			com.State = ComputationStateRejected
			if err := ome.computationStore.PutComputation(com); err != nil {
				logger.Error(fmt.Sprintf("cannot store expired computation buy = %v, sell = %v: %v", com.Buy, com.Sell, err))
			}
			continue
		}
		// Add this Computation to the buffer
		buffer[bufferN] = com
		if bufferN++; bufferN >= OmeBufferLimit {
			break
		}
	}

	// Retry each of the Computations in the buffer
	if bufferN > 0 {
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("retrying %v computations", bufferN))
		for i := 0; i < bufferN; i++ {
			if err := ome.sendComputationToMatcher(buffer[i], done, matches); err != nil {
				ome.computationBacklog[buffer[i].ID] = buffer[i]
			}
		}
	}
}

func (ome *ome) sendComputationToMatcher(com Computation, done <-chan struct{}, matches chan<- Computation) error {
	logger.Compute(logger.LevelDebug, fmt.Sprintf("resolving buy = %v, sell = %v at epoch = %v", com.Buy, com.Sell, base64.StdEncoding.EncodeToString(com.Epoch[:8])))
	ome.matcher.Resolve(com, func(com Computation) {
		if !com.Match {
			return
		}
		ome.sendComputationToConfirmer(com, done, matches)
	})
	return nil
}

func (ome *ome) sendComputationToConfirmer(com Computation, done <-chan struct{}, matches chan<- Computation) {
	select {
	case <-done:
	case matches <- com:
	}
}

func (ome *ome) sendComputationToSettler(com Computation) {
	logger.Compute(logger.LevelDebug, fmt.Sprintf("settling buy = %v, sell = %v", com.Buy, com.Sell))
	if err := ome.settler.Settle(com); err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot settle: %v", err))
	}
}
