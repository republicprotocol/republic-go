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
	ComputationStateAccepted
	ComputationStateRejected
	ComputationStateSettled
)

// Computations is an alias type.
type Computations []Computation

// A Computation is a combination of a buy order.Order and a sell order.Order.
type Computation struct {
	ID    ComputationID
	State ComputationState

	Buy     order.ID
	Sell    order.ID
	IsMatch bool
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

	// Sync the orderbook
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
			ome.syncOrderbook(done, errs)
			timeNextSync := timeBeginSync.Add(14 * time.Second)
			if time.Now().After(timeNextSync) {
				continue
			}
			time.Sleep(timeNextSync.Sub(time.Now()))
		}
	}()

	// Sync the matcher
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
			if wait := ome.syncMatcher(done, matches, errs); !wait {
				continue
			}
			timeNextSync := timeBeginSync.Add(14 * time.Second)
			if time.Now().After(timeNextSync) {
				continue
			}
			time.Sleep(timeNextSync.Sub(time.Now()))
		}
	}()

	// Sync the confirmer
	wg.Add(1)
	go func() {
		defer wg.Done()
		ome.syncConfirmer(done, matches, errs)
	}()

	// Cleanup
	go func() {
		defer close(errs)
		wg.Wait()
	}()

	return errs
}

func (ome *ome) syncOrderbook(done <-chan struct{}, errs chan<- error) {
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

	if len(changeset) == 0 {
		return
	}
	if err := ome.syncRanker(changeset); err != nil {
		select {
		case <-done:
			return
		case errs <- fmt.Errorf("cannot sync ranker: %v", err):
			return
		}
	}
}

func (ome *ome) syncRanker(changeset orderbook.ChangeSet) error {
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
	return nil
}

func (ome *ome) syncMatcher(done <-chan struct{}, matches chan<- Computation, errs chan<- error) bool {
	ome.ξMu.RLock()
	ξ := ome.ξ.Hash
	ome.ξMu.RUnlock()

	buffer := [128]Computation{}
	n := ome.ranker.Computations(buffer[:])

	for i := 0; i < n; i++ {
		logger.Compute(logger.LevelDebug, fmt.Sprintf("resolving buy = %v, sell = %v", buffer[i].Buy, buffer[i].Sell))
		err := ome.matcher.Resolve(ξ, buffer[i], func(com Computation) {
			select {
			case <-done:
			case matches <- com:
			}
		})
		if err != nil {
			select {
			case <-done:
				return false
			case errs <- err:
			}
		}
	}
	return n != 128
}

func (ome *ome) syncConfirmer(done <-chan struct{}, matches <-chan Computation, errs chan<- error) {
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
			if err := ome.settler.Settle(ξ, confirmation); err != nil {
				logger.Network(logger.LevelError, fmt.Sprintf("cannot settle: %v", err))
			}
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
