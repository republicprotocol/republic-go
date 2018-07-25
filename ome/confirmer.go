package ome

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
)

// ErrOrderNotConfirmed is an internal error that is returned when an
// order.Order is not confirmed.
var ErrOrderNotConfirmed = errors.New("order not confirmed")

// A Confirmer consumes Computations that been resolved to a match, usually by
// a call to Matcher.Resolve, and reaches consensus with other Darknodes that
// this match has happened. This prevents the occurrence of conflicting matches
// across the parallel resolving done by different Darknode pods.
type Confirmer interface {

	// Confirm Computations that have resolved in a match by reaching consensus
	// with other Darknodes. The input channel will be consumed and an output
	// channel of confirmed Computations is produced. Stop the Confirmer by
	// closing the done channel.
	Confirm(done <-chan struct{}, coms <-chan Computation) (<-chan Computation, <-chan error)
}

type confirmer struct {
	computationStore ComputationStorer

	contract              ContractBinder
	orderbookPollInterval time.Duration
	orderbookBlockDepth   uint

	confirmingMu         *sync.Mutex
	confirmingBuyOrders  map[order.ID]struct{}
	confirmingSellOrders map[order.ID]struct{}
	confirmed            map[order.ID]time.Time
}

// NewConfirmer returns a Confirmer that submits Computations to the
// Orderbook for confirmation. It polls the Orderbook on an interval
// and checks for consensus on confirmations by waiting until a submitted
// Computation has been confirmed has the confirmation has passed the block
// depth limit.
func NewConfirmer(computationStore ComputationStorer, contract ContractBinder, orderbookPollInterval time.Duration, orderbookBlockDepth uint) Confirmer {
	return &confirmer{
		computationStore: computationStore,

		contract:              contract,
		orderbookPollInterval: orderbookPollInterval,
		orderbookBlockDepth:   orderbookBlockDepth,

		confirmingMu:         new(sync.Mutex),
		confirmingBuyOrders:  map[order.ID]struct{}{},
		confirmingSellOrders: map[order.ID]struct{}{},
		confirmed:            map[order.ID]time.Time{},
	}
}

// Confirm implements the Confirmer interface.
func (confirmer *confirmer) Confirm(done <-chan struct{}, coms <-chan Computation) (<-chan Computation, <-chan error) {
	confirmations := make(chan Computation, 64)
	errs := make(chan error, 1)

	// Two background processes are run that must end before cleanup can
	// happen safely
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-done:
				return
			case com, ok := <-coms:
				if !ok {
					return
				}

				// Check that these orders have not already been confirmed
				if _, ok := confirmer.confirmed[com.Buy.OrderID]; ok {
					continue
				}
				if _, ok := confirmer.confirmed[com.Sell.OrderID]; ok {
					continue
				}

				// Confirm Computations on the blockchain and register them for
				// observation (we need to wait for finality)
				if err := confirmer.beginConfirmation(com); err != nil {
					// An error in confirmation should not stop the
					// Confirmer from monitoring the Computation for
					// confirmation (another node might have succeeded), so
					// we pass through
					logger.Error(err.Error())
				}

				// Wait for the confirmation of these orders to pass the depth
				// limit
				confirmer.confirmingMu.Lock()
				confirmer.confirmingBuyOrders[com.Buy.OrderID] = struct{}{}
				confirmer.confirmingSellOrders[com.Sell.OrderID] = struct{}{}
				confirmer.confirmingMu.Unlock()
			}
		}
	}()

	// Periodically poll the orderbook to observe the state of
	// confirmations that have passed the block depth limit
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(confirmer.orderbookPollInterval)
		defer ticker.Stop()

		for {
			select {

			// Graceful shutdown when the done channel is closed
			case <-done:
				return

			case <-ticker.C:
				confirmer.confirmingMu.Lock()
				confirmer.checkOrdersForConfirmationFinality(order.ParityBuy, done, confirmations, errs)
				confirmer.checkOrdersForConfirmationFinality(order.ParitySell, done, confirmations, errs)
				// Clean up confirmed orders that are old enough to forget about
				for key, t := range confirmer.confirmed {
					if time.Since(t) > time.Hour {
						delete(confirmer.confirmed, key)
					}
				}
				confirmer.confirmingMu.Unlock()
			}
		}
	}()

	// Cleanup
	go func() {
		defer close(confirmations)
		defer close(errs)
		wg.Wait()
	}()

	return confirmations, errs
}

func (confirmer *confirmer) beginConfirmation(orderMatch Computation) error {
	if err := confirmer.contract.ConfirmOrder(orderMatch.Buy.OrderID, orderMatch.Sell.OrderID); err != nil {
		return fmt.Errorf("cannot confirm computation buy = %v, sell = %v: %v", orderMatch.Buy.OrderID, orderMatch.Sell.OrderID, err)
	}
	return nil
}

func (confirmer *confirmer) checkOrdersForConfirmationFinality(orderParity order.Parity, done <-chan struct{}, confirmations chan<- Computation, errs chan<- error) {
	var confirmingOrders map[order.ID]struct{}
	if orderParity == order.ParityBuy {
		confirmingOrders = confirmer.confirmingBuyOrders
	} else {
		confirmingOrders = confirmer.confirmingSellOrders
	}
	for ord := range confirmingOrders {
		ordMatch, err := confirmer.checkOrderForConfirmationFinality(ord, orderParity)
		if err != nil {
			if err == ErrOrderNotConfirmed {
				continue
			}
			select {
			case <-done:
				return
			case errs <- err:
				continue
			}
		}

		com, err := confirmer.computationFromOrders(orderParity, ord, ordMatch)
		if err != nil {
			select {
			case <-done:
				return
			case errs <- err:
				if orderParity == order.ParityBuy {
					delete(confirmer.confirmingBuyOrders, ord)
					delete(confirmer.confirmingSellOrders, ordMatch)
					continue
				}
				delete(confirmer.confirmingBuyOrders, ordMatch)
				delete(confirmer.confirmingSellOrders, ord)
				continue
			}
		}
		if err := confirmer.computationStore.PutComputation(com); err != nil {
			select {
			case <-done:
				return
			case errs <- err:
			}
		}

		// Check that these orders have not already been output
		if _, ok := confirmer.confirmed[com.Buy.OrderID]; ok {
			continue
		}
		if _, ok := confirmer.confirmed[com.Sell.OrderID]; ok {
			continue
		}
		select {
		case <-done:
			return
		case confirmations <- com:
			delete(confirmer.confirmingBuyOrders, com.Buy.OrderID)
			delete(confirmer.confirmingSellOrders, com.Sell.OrderID)
			confirmer.confirmed[com.Buy.OrderID] = time.Now()
			confirmer.confirmed[com.Sell.OrderID] = time.Now()
		}
	}
}

func (confirmer *confirmer) checkOrderForConfirmationFinality(ord order.ID, orderParity order.Parity) (order.ID, error) {
	// Ignore orders that are not pass the depth limit
	depth, err := confirmer.contract.Depth(ord)
	if err != nil {
		return order.ID{}, err
	}
	if depth < confirmer.orderbookBlockDepth {
		return order.ID{}, ErrOrderNotConfirmed
	}

	// Purge orders that are not confirmed
	status, err := confirmer.contract.Status(ord)
	if err != nil {
		return order.ID{}, err
	}
	if status != StatusConfirmed {
		if orderParity == order.ParityBuy {
			delete(confirmer.confirmingBuyOrders, ord)
		} else {
			delete(confirmer.confirmingSellOrders, ord)
		}
		return order.ID{}, ErrOrderNotConfirmed
	}

	match, err := confirmer.contract.OrderMatch(ord)
	if err != nil {
		return order.ID{}, err
	}
	return match, nil
}

func (confirmer *confirmer) computationFromOrders(orderParity order.Parity, ord, ordMatch order.ID) (Computation, error) {
	var comIDDepth0 ComputationID
	var comIDDepth1 ComputationID
	if orderParity == order.ParityBuy {
		comIDDepth0 = NewComputationID(ord, ordMatch, 0)
		comIDDepth1 = NewComputationID(ord, ordMatch, 1)
	} else {
		comIDDepth0 = NewComputationID(ordMatch, ord, 0)
		comIDDepth1 = NewComputationID(ordMatch, ord, 1)
	}

	com, err := confirmer.computationStore.Computation(comIDDepth0)
	if err != nil {
		com, err = confirmer.computationStore.Computation(comIDDepth1)
	}
	if err != nil {
		return com, err
	}

	com.State = ComputationStateAccepted
	com.Timestamp = time.Now()
	return com, nil
}
