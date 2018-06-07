package ome

import (
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/logger"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/smpc"
)

// ComputationID is used to distinguish between different combinations of
// orders that are being matched against each other.
type ComputationID [32]byte

// ComputationState is used to track the state of a Computation as it changes
// over its lifetime. This prevents duplicated work in the system.
type ComputationState int

// Values for a ComputationState
const (
	ComputationStateNil = iota
	ComputationStateMatched
	ComputationStateMismatched
	ComputationStateAccepted
	ComputationStateRejected
	ComputationStateSettled
)

// A Priority is an unsigned integer representing logical time priority. The
// lower the number, the higher the priority.
type Priority uint64

// Computations is an alias type.
type Computations []Computation

// A Computation is a combination of a buy order.Order and a sell order.Order.
type Computation struct {
	ID       ComputationID
	State    ComputationState
	Priority Priority
	Match    bool

	Buy  order.ID
	Sell order.ID
}

// NewComputation returns a pending Computation between a buy order.Order and a
// sell order.Order. It initialized the ComputationID to the Keccak256 hash of
// the buy order.ID and the sell order.ID.
func NewComputation(buy, sell order.ID) Computation {
	com := Computation{
		Buy:  buy,
		Sell: sell,
	}
	copy(com.ID[:], crypto.Keccak256(buy[:], sell[:]))
	return com
}

// An Ome runs the logic for a single node in the secure order matching engine.
type Ome interface {
	cal.EpochListener

	// Run the Ome in the background. Stop the Ome by closing the done channel.
	Run(done <-chan struct{}) <-chan error
}

type ome struct {
	ranker    Ranker
	matcher   Matcher
	confirmer Confirmer
	settler   Settler
	orderbook orderbook.Orderbook
	smpcer    smpc.Smpcer

	ξMu *sync.RWMutex
	ξ   cal.Epoch
}

func NewOme(ranker Ranker, matcher Matcher, confirmer Confirmer, settler Settler, orderbook orderbook.Orderbook, smpcer smpc.Smpcer) Ome {
	return &ome{
		ranker:    ranker,
		matcher:   matcher,
		confirmer: confirmer,
		settler:   settler,
		orderbook: orderbook,
		smpcer:    smpcer,

		ξMu: new(sync.RWMutex),
		ξ:   cal.Epoch{},
	}
}

// OnChangeEpoch implements the cal.EpochListener interface.
func (ome *ome) OnChangeEpoch(ξ cal.Epoch) {
	ome.ξMu.Lock()
	defer ome.ξMu.Unlock()

	ome.smpcer.Disconnect(ome.ξ.Hash)
	ome.ξ = ξ
	ome.smpcer.Connect(ome.ξ.Hash, ome.ξ.Darknodes)
}

// Run implements the Ome interface.
func (ome *ome) Run(done <-chan struct{}) <-chan error {
	matches := make(chan Computation, 64)
	errs := make(chan error, 64)

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

	// Sync the Ranker to the Matcher
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

	// Cleanup
	go func() {
		defer close(errs)
		wg.Wait()
	}()

	return errs
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
		switch change.OrderStatus {
		case order.Open:
			if change.OrderParity == order.ParityBuy {
				ome.ranker.InsertBuy(PriorityOrder{
					Priority: Priority(change.OrderPriority),
					Order:    change.OrderID,
				})
			} else {
				ome.ranker.InsertSell(PriorityOrder{
					Priority: Priority(change.OrderPriority),
					Order:    change.OrderID,
				})
			}
		case order.Canceled, order.Confirmed:
			ome.ranker.Remove(change.OrderID)
		}
	}
}

func (ome *ome) syncRanker(done <-chan struct{}, matches chan<- Computation, errs chan<- error) bool {
	buffer := [128]Computation{}
	n := ome.ranker.Computations(buffer[:])

	ome.ξMu.RLock()
	ξ := ome.ξ.Hash
	ome.ξMu.RUnlock()

	for i := 0; i < n; i++ {
		switch buffer[i].State {
		case ComputationStateNil:
			ome.sendComputationToMatcher(ξ, buffer[i], done, matches, errs)

		case ComputationStateMatched:
			ome.sendComputationToConfirmer(buffer[i], done, matches)

		case ComputationStateAccepted:
			ome.sendComputationToSettler(ξ, buffer[i])

		default:
			logger.Error(fmt.Sprintf("unexpected state for computation buy = %v, sell = %v: %v", buffer[i].Buy, buffer[i].Sell, buffer[i].State))
		}

	}
	return n != 128
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

			ome.ξMu.RLock()
			ξ := ome.ξ.Hash
			ome.ξMu.RUnlock()
			ome.sendComputationToSettler(ξ, confirmation)

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

func (ome *ome) sendComputationToMatcher(ξ [32]byte, com Computation, done <-chan struct{}, matches chan<- Computation, errs chan<- error) {
	logger.Compute(logger.LevelDebug, fmt.Sprintf("resolving buy = %v, sell = %v", com.Buy, com.Sell))
	err := ome.matcher.Resolve(ξ, com, func(com Computation) {
		if !com.Match {
			return
		}
		ome.sendComputationToConfirmer(com, done, matches)
	})
	if err != nil {
		select {
		case <-done:
		case errs <- err:
		}
	}
}

func (ome *ome) sendComputationToConfirmer(com Computation, done <-chan struct{}, matches chan<- Computation) {
	select {
	case <-done:
	case matches <- com:
	}
}

func (ome *ome) sendComputationToSettler(ξ [32]byte, com Computation) {
	logger.Compute(logger.LevelDebug, fmt.Sprintf("settling buy = %v, sell = %v", com.Buy, com.Sell))
	if err := ome.settler.Settle(ξ, com); err != nil {
		logger.Network(logger.LevelError, fmt.Sprintf("cannot settle: %v", err))
	}
}
