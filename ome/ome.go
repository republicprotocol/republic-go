package ome

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/orderbook"
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

	// OnChangeEpoch should be called whenever a new cal.Epoch is observed.
	OnChangeEpoch(cal.Epoch)
}

type ome struct {
	ranker    Ranker
	matcher   Matcher
	confirmer Confirmer
	settler   Settler
	storer    Storer
	orderbook orderbook.Orderbook
	smpcer    smpc.Smpcer

	computationBacklogMu *sync.RWMutex
	computationBacklog   map[ComputationID]Computation

	epochMu   *sync.RWMutex
	epochCurr *cal.Epoch
	epochPrev *cal.Epoch
}

// NewOme returns an Ome that uses an order.Orderbook to synchronize changes
// from the Ethereum blockchain, and an smpc.Smpcer to run the secure
// multi-party computations necessary for the secure order matching engine.
func NewOme(ranker Ranker, matcher Matcher, confirmer Confirmer, settler Settler, storer Storer, orderbook orderbook.Orderbook, smpcer smpc.Smpcer) Ome {
	return &ome{
		ranker:    ranker,
		matcher:   matcher,
		confirmer: confirmer,
		settler:   settler,
		storer:    storer,
		orderbook: orderbook,
		smpcer:    smpcer,

		computationBacklogMu: new(sync.RWMutex),
		computationBacklog:   map[ComputationID]Computation{},

		epochMu:   new(sync.RWMutex),
		epochCurr: nil,
		epochPrev: nil,
	}
}

// Run implements the Ome interface.
func (ome *ome) Run(done <-chan struct{}) <-chan error {
	matches := make(chan Computation, OmeBufferLimit)
	errs := make(chan error, OmeBufferLimit)

	var wg sync.WaitGroup

	// Sync the orderbook.Orderbook to the Ranker
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
			}
			timeBeginSync := time.Now()

			ome.syncOrderbookToRanker(done, errs)

			timeNextSync := timeBeginSync.Add(14 * time.Second)
			if time.Now().After(timeNextSync) {
				continue
			}
			time.Sleep(timeNextSync.Sub(time.Now()))
		}
	}()

	// Sync the Ranker
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-done:
				return
			default:
			}
			timeBeginSync := time.Now()

			if wait := ome.syncRanker(done, matches, errs); !wait {
				continue
			}

			timeNextSync := timeBeginSync.Add(14 * time.Second)
			if time.Now().After(timeNextSync) {
				continue
			}
			time.Sleep(timeNextSync.Sub(time.Now()))
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

// OnChangeEpoch updates the Ome to the next cal.Epoch. This will cause
// cascading changes throughout the Ome, most notably it will connect to a new
// Smpc network that will handle future Computations.
func (ome *ome) OnChangeEpoch(epoch cal.Epoch) {
	ome.epochMu.Lock()
	defer ome.epochMu.Unlock()

	// Do not update if the epoch has not actually changed
	if bytes.Equal(epoch.Hash[:], ome.epochCurr.Hash[:]) {
		return
	}

	// Replace the previous epoch
	if ome.epochPrev != nil {
		ome.smpcer.Disconnect(ome.epochPrev.Hash)
	}
	ome.epochPrev = ome.epochCurr

	// Replace the current epoch
	ome.epochCurr = &epoch
	ome.smpcer.Connect(ome.epochCurr.Hash, ome.epochCurr.Darknodes)

	// Notify the Ranker
	ome.ranker.OnChangeEpoch(epoch)
}

func (ome *ome) syncOrderbookToRanker(done <-chan struct{}, errs chan<- error) {
	changeset, err := ome.orderbook.Sync()
	if err != nil {
		select {
		case <-done:
			return
		case errs <- fmt.Errorf("cannot sync orderbook: %v", err):
			return
		}
	}
	logger.Network(logger.LevelDebug, fmt.Sprintf("sync orderbook: %v changes in changeset", len(changeset)))

	for _, change := range changeset {
		ome.ranker.InsertChange(change)
	}
}

func (ome *ome) syncRanker(done <-chan struct{}, matches chan<- Computation, errs chan<- error) bool {
	buffer := [OmeBufferLimit]Computation{}
	n := ome.ranker.Computations(buffer[:])

	for i := 0; i < n; i++ {
		switch buffer[i].State {
		case ComputationStateNil:
			if err := ome.sendComputationToMatcher(buffer[i], done, matches); err != nil {
				ome.computationBacklogMu.Lock()
				ome.computationBacklog[buffer[i].ID] = buffer[i]
				ome.computationBacklogMu.Unlock()
			}

		case ComputationStateMatched:
			ome.sendComputationToConfirmer(buffer[i], done, matches)

		case ComputationStateAccepted:
			ome.sendComputationToSettler(buffer[i])

		default:
			logger.Error(fmt.Sprintf("unexpected state for computation buy = %v, sell = %v: %v", buffer[i].Buy, buffer[i].Sell, buffer[i].State))
		}

	}
	return n != OmeBufferLimit
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
			if err := ome.storer.InsertComputation(com); err != nil {
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
	for i := 0; i < bufferN; i++ {
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("retrying computation buy = %v, sell = %v", buffer[i].Buy, buffer[i].Sell))
		if err := ome.sendComputationToMatcher(buffer[i], done, matches); err != nil {
			logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("cannot resolve computation buy = %v, sell = %v: %v", buffer[i].Buy, buffer[i].Sell, err))
			ome.computationBacklog[buffer[i].ID] = buffer[i]
		}
	}
}

func (ome *ome) sendComputationToMatcher(com Computation, done <-chan struct{}, matches chan<- Computation) error {
	buyFragment, err := ome.storer.OrderFragment(com.Buy)
	if err != nil {
		return err
	}
	sellFragment, err := ome.storer.OrderFragment(com.Sell)
	if err != nil {
		return err
	}

	logger.Compute(logger.LevelDebug, fmt.Sprintf("resolving buy = %v, sell = %v", com.Buy, com.Sell))
	ome.matcher.Resolve(com, buyFragment, sellFragment, func(com Computation) {
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
