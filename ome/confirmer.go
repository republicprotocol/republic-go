package ome

import (
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/order"
)

// A Confirmer consumers Computations that have resulted in an order match and
// reaches consensus with other Darknodes on this order match. This prevents
// the occurrence of conflicting order matches.
type Confirmer interface {

	// ConfirmOrderMatches by consuming order matches from an input channel and
	// producing confirmed order matches to an output channel. A call to
	// Confirmer.ConfirmOrderMatches should use a background goroutine to avoid
	// blocking the caller.
	ConfirmOrderMatches(done <-chan struct{}, orderMatches <-chan Computation) (<-chan Computation, <-chan error)
}

type confirmer struct {
	renLedgerDepth        uint
	renLedgerPollInterval time.Duration
	renLedger             cal.RenLedger

	confirmingMu         *sync.Mutex
	confirmingBuyOrders  map[order.ID]struct{}
	confirmingSellOrders map[order.ID]struct{}
}

// NewConfirmer returns a Confirmer that polls the cal.RenLedger once per
// interval and checks for confirmed order matches that have passed the block
// depth limit. These confirmations will not be reshuffle (with high
// probability), depending on the block depth limit.
func NewConfirmer(renLedgerDepth uint, renLedgerPollInterval time.Duration, renLedger cal.RenLedger) Confirmer {
	return &confirmer{
		renLedgerDepth:        renLedgerDepth,
		renLedgerPollInterval: renLedgerPollInterval,
		renLedger:             renLedger,

		confirmingMu:         new(sync.Mutex),
		confirmingBuyOrders:  map[order.ID]struct{}{},
		confirmingSellOrders: map[order.ID]struct{}{},
	}
}

// ConfirmOrderMatches implements the Confirmer interface.
func (confirmer *confirmer) ConfirmOrderMatches(done <-chan struct{}, orderMatches <-chan Computation) (<-chan Computation, <-chan error) {
	confirmedOrderMatches := make(chan Computation)
	errs := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	// Confirm order matches on the blockchain and register them for
	// observation (we need to wait for finality)
	go func() {
		defer wg.Done()

		for {
			select {

			// Graceful shutdown when the done channel is closed
			case <-done:
				return

			case orderMatch, ok := <-orderMatches:
				// Graceful shutdown when the input channel is closed
				if !ok {
					return
				}

				if err := confirmer.beginConfirmOrder(orderMatch); err != nil {
					select {
					case <-done:
						return
					case errs <- err:
					}
				}

				// Wait for the confirmation of these orders to pass the depth
				// limit
				func() {
					confirmer.confirmingMu.Lock()
					defer confirmer.confirmingMu.Unlock()

					confirmer.confirmingBuyOrders[orderMatch.Buy] = struct{}{}
					confirmer.confirmingSellOrders[orderMatch.Sell] = struct{}{}
				}()
			}
		}
	}()

	// Periodically poll the cal.RenLedger to observe the state of orders that
	// have been confirmed passed the block depth limit
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(confirmer.renLedgerPollInterval)
		defer ticker.Stop()

		for {
			select {

			// Graceful shutdown when the done channel is closed
			case <-done:
				return

			case <-ticker.C:

				// Create a closure so that we can scope the defer statement to
				// this iteration of the for loop
				func() {
					confirmer.confirmingMu.Lock()
					defer confirmer.confirmingMu.Unlock()

					confirmer.checkOrdersForConfirmationFinality(order.ParityBuy, done, confirmedOrderMatches, errs)
					confirmer.checkOrdersForConfirmationFinality(order.ParitySell, done, confirmedOrderMatches, errs)
				}()
			}
		}
	}()

	// Goroutine responsible for cleanup
	go func() {
		defer close(confirmedOrderMatches)
		defer close(errs)
		wg.Wait()
	}()

	return confirmedOrderMatches, errs
}

func (confirmer *confirmer) beginConfirmOrder(orderMatch Computation) error {
	// TODO: Check the error and if it failed due to an ephemeral error then
	// try again.
	return confirmer.renLedger.ConfirmOrder(orderMatch.Buy, orderMatch.Sell)
}

func (confirmer *confirmer) checkOrdersForConfirmationFinality(orderParity order.Parity, done <-chan struct{}, confirmedOrderMatches chan<- Computation, errs chan<- error) {

	var confirmingOrders map[order.ID]struct{}
	if orderParity == order.ParityBuy {
		confirmingOrders = confirmer.confirmingBuyOrders
	} else {
		confirmingOrders = confirmer.confirmingSellOrders
	}

	for ord := range confirmingOrders {
		matchingOrd, err := confirmer.checkOrderForConfirmationFinality(ord, order.ParityBuy)
		if err != nil {
			select {
			case <-done:
				return
			case errs <- err:
				continue
			}
		}
		if matchingOrd == nil {
			continue
		}

		if orderParity == order.ParityBuy {
			delete(confirmer.confirmingBuyOrders, ord)
			delete(confirmer.confirmingSellOrders, *matchingOrd)
		} else {
			delete(confirmer.confirmingSellOrders, ord)
			delete(confirmer.confirmingBuyOrders, *matchingOrd)
		}

		confirmedOrderMatch := Computation{}
		if orderParity == order.ParityBuy {
			confirmedOrderMatch.Buy = ord
			confirmedOrderMatch.Sell = *matchingOrd
		} else {
			confirmedOrderMatch.Sell = ord
			confirmedOrderMatch.Buy = *matchingOrd
		}
		select {
		case <-done:
			return
		case confirmedOrderMatches <- confirmedOrderMatch:
		}
	}
}

func (confirmer *confirmer) checkOrderForConfirmationFinality(ord order.ID, orderParity order.Parity) (*order.ID, error) {
	// Ignore orders that are not pass the depth limit
	depth, err := confirmer.renLedger.Depth(ord)
	if err != nil {
		return nil, err
	}
	if depth < confirmer.renLedgerDepth {
		return nil, nil
	}

	// Purge orders that are not confirmed
	status, err := confirmer.renLedger.Status(ord)
	if err != nil {
		return nil, err
	}
	if status != cal.StatusConfirmed {
		log.Println("order status is ", status)
		if orderParity == order.ParityBuy {
			delete(confirmer.confirmingBuyOrders, ord)
		} else {
			delete(confirmer.confirmingSellOrders, ord)
		}
		return nil, nil
	}

	// Output confirmed order matches
	matchingOrd, err := confirmer.renLedger.OrderMatch(ord)
	if err != nil {
		return nil, err
	}
	return &matchingOrd, nil
}
