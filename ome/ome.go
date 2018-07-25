package ome

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
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

	epochMu   *sync.RWMutex
	epochCurr *registry.Epoch
	epochPrev *registry.Epoch
}

// NewOme returns an Ome that uses an order.Orderbook to synchronize changes
// from the Ethereum blockchain, and an smpc.Smpcer to run the secure
// multi-party computations necessary for the secure order matching engine.
func NewOme(addr identity.Address, gen ComputationGenerator, matcher Matcher, confirmer Confirmer, settler Settler, computationStore ComputationStorer, orderbook orderbook.Orderbook, smpcer smpc.Smpcer, epochPrev registry.Epoch) Ome {
	ome := &ome{
		addr:             addr,
		gen:              gen,
		matcher:          matcher,
		confirmer:        confirmer,
		settler:          settler,
		computationStore: computationStore,
		orderbook:        orderbook,
		smpcer:           smpcer,

		epochMu:   new(sync.RWMutex),
		epochCurr: nil,
		epochPrev: nil,
	}
	ome.OnChangeEpoch(epochPrev)
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
	ome.epochMu.RLock()
	defer ome.epochMu.RUnlock()

	// Do not update if the epoch has not actually changed
	if ome.epochCurr != nil && bytes.Equal(epoch.Hash[:], ome.epochCurr.Hash[:]) {
		return
	}

	go func() {
		ome.epochMu.Lock()
		defer ome.epochMu.Unlock()

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
		time.Sleep(time.Second)

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

func (ome *ome) sendComputationToMatcher(com Computation, done <-chan struct{}, matches chan<- Computation) error {
	logger.Compute(logger.LevelDebug, fmt.Sprintf("resolving buy = %v, sell = %v at epoch = %v", com.Buy.OrderID, com.Sell.OrderID, base64.StdEncoding.EncodeToString(com.Epoch[:8])))
	ome.matcher.Resolve(com, func(com Computation) {
		if !com.Match {
			return
		}
		log.Printf("[debug] (resolve) ✔ buy = %v, sell = %v", com.Buy.OrderID, com.Sell.OrderID)
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
	logger.Compute(logger.LevelDebug, fmt.Sprintf("settling buy = %v, sell = %v", com.Buy.OrderID, com.Sell.OrderID))
	if err := ome.settler.Settle(com); err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot settle: %v", err))
	}
}
