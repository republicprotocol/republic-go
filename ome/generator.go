package ome

import (
	"log"
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
)

type ComputationGenerator interface {
	Generate(done <-chan struct{}, notifications <-chan orderbook.Notification) (<-chan Computation, <-chan error)
	OnChangeEpoch(epoch registry.Epoch)
}

type computationGenerator struct {
	doneMu *sync.Mutex
	done   chan struct{}

	matMu                *sync.Mutex
	matCurrDone          chan struct{}
	matCurrNotifications chan orderbook.Notification
	matPrevDone          chan struct{}
	matPrevNotifications chan orderbook.Notification

	broadcastComputations chan (<-chan Computation)
	broadcastErrs         chan (<-chan error)

	orderFragmentStore OrderFragmentStorer
}

func NewComputationGenerator(orderFragmentStore OrderFragmentStorer) ComputationGenerator {
	return &computationGenerator{
		doneMu: new(sync.Mutex),
		done:   nil,

		matMu:                new(sync.Mutex),
		matCurrDone:          nil,
		matCurrNotifications: nil,
		matPrevDone:          nil,
		matPrevNotifications: nil,

		broadcastComputations: make(chan (<-chan Computation)),
		broadcastErrs:         make(chan (<-chan error)),

		orderFragmentStore: orderFragmentStore,
	}
}

func (gen *computationGenerator) Generate(done <-chan struct{}, notifications <-chan orderbook.Notification) (<-chan Computation, <-chan error) {
	computations := make(chan Computation)
	errs := make(chan error)

	// Handle all incoming notifications by mapping them to the appropriate
	// computation matrix
	go func() {
		for {
			select {
			case <-done:
				return
			case notification, ok := <-notifications:
				if !ok {
					return
				}
				gen.routeNotification(notification, done)
			}
		}
	}()

	go func() {
		// Wait for all goroutines to finish and then cleanup
		defer close(computations)
		defer close(errs)

		// Merge all of the channels on the broadcast channel into the output
		// channel
		dispatch.CoBegin(
			func() {
				dispatch.Merge(done, gen.broadcastComputations, computations)
			},
			func() {
				dispatch.Merge(done, gen.broadcastErrs, errs)
			})
	}()

	return computations, errs
}

// OnChangeEpoch implements the ComputationGenerator interface.
func (gen *computationGenerator) OnChangeEpoch(epoch registry.Epoch) {
	gen.matMu.Lock()
	defer gen.matMu.Unlock()

	// Transition the current epoch into the previous epoch and setup a new
	// current epoch
	if gen.matPrevDone != nil {
		close(gen.matPrevDone)
		close(gen.matPrevNotifications)
	}
	gen.matPrevDone = gen.matCurrDone
	gen.matPrevNotifications = gen.matCurrNotifications
	gen.matCurrDone = make(chan struct{})
	gen.matCurrNotifications = make(chan orderbook.Notification)

	mat := newComputationMatrix(epoch, gen.orderFragmentStore)
	computations, errs := mat.generate(gen.matCurrDone, gen.matCurrNotifications)

	go func() {
		select {
		case <-gen.done:
		case gen.broadcastComputations <- computations:
		}

		select {
		case <-gen.done:
		case gen.broadcastErrs <- errs:
		}
	}()
}

func (gen *computationGenerator) routeNotification(notification orderbook.Notification, done <-chan struct{}) {
	switch notification := notification.(type) {

	// Notifications that open orders need to be routed to the appropriate
	// computation matrix
	case orderbook.NotificationOpenOrder:
		gen.routeNotificationOpenOrder(notification, done)

	// All computation matrices receive all other notifications
	default:
		gen.matMu.Lock()
		defer gen.matMu.Unlock()

		if gen.matCurrNotifications != nil {
			select {
			case <-done:
				return
			case gen.matCurrNotifications <- notification:
			}
		}
		if gen.matPrevNotifications != nil {
			select {
			case <-done:
				return
			case gen.matPrevNotifications <- notification:
			}
		}
	}
}

func (gen *computationGenerator) routeNotificationOpenOrder(notification orderbook.NotificationOpenOrder, done <-chan struct{}) {
	gen.matMu.Lock()
	defer gen.matMu.Unlock()

	switch notification.OrderFragment.EpochDepth {
	case 0:
		if gen.matCurrNotifications != nil {
			select {
			case <-done:
				return
			case gen.matCurrNotifications <- notification:
			}
		}
	case 1:
		if gen.matPrevNotifications != nil {
			select {
			case <-done:
				return
			case gen.matPrevNotifications <- notification:
			}
		}
	}
}

type computationMatrix struct {
	epoch              registry.Epoch
	orderFragmentStore OrderFragmentStorer
}

func newComputationMatrix(epoch registry.Epoch, orderFragmentStore OrderFragmentStorer) *computationMatrix {
	return &computationMatrix{
		epoch:              epoch,
		orderFragmentStore: orderFragmentStore,
	}
}

func (mat *computationMatrix) generate(done <-chan struct{}, notifications <-chan orderbook.Notification) (<-chan Computation, <-chan error) {
	computations := make(chan Computation)
	errs := make(chan error)

	go func() {
		defer close(computations)
		defer close(errs)

		for {
			select {
			case <-done:
				return
			case notification, ok := <-notifications:
				if !ok {
					return
				}
				mat.handleNotification(notification, done, computations, errs)
			}
		}
	}()

	return computations, errs
}

func (mat *computationMatrix) handleNotification(notification orderbook.Notification, done <-chan struct{}, computations chan<- Computation, errs chan<- error) {
	switch notification := notification.(type) {
	// Notifications that open orders result in the insertion of that order
	// into the matrix
	case orderbook.NotificationOpenOrder:
		mat.insertOrderFragment(notification, done, computations, errs)

	// Notifications that close an order result in the removal of that order
	// from storage
	case orderbook.NotificationConfirmOrder:
		mat.removeOrderFragment(notification.OrderID)
	case orderbook.NotificationCancelOrder:
		mat.removeOrderFragment(notification.OrderID)
	default:
		select {
		case <-done:
		case errs <- orderbook.ErrUnexpectedNotificationType:
		}
	}
}

func (mat *computationMatrix) insertOrderFragment(notification orderbook.NotificationOpenOrder, done <-chan struct{}, computations chan<- Computation, errs chan<- error) {

	// Store the order.Fragment and get the opposing list so that computations
	// can be generated
	var oppositeOrderFragmentIter OrderFragmentIterator
	var err error

	if notification.OrderFragment.OrderParity == order.ParityBuy {
		if err := mat.orderFragmentStore.PutBuyOrderFragment(mat.epoch, notification.OrderFragment, notification.Trader); err != nil {
			log.Printf("[error] (generator) cannot store buy order fragment = %v: %v", notification.OrderID, err)
			return
		}
		oppositeOrderFragmentIter, err = mat.orderFragmentStore.SellOrderFragments(mat.epoch)
		if err != nil {
			log.Printf("[error] (generator) cannot load buy order fragment iterator: %v", err)
			return
		}
		defer oppositeOrderFragmentIter.Release()
	} else {
		if err := mat.orderFragmentStore.PutSellOrderFragment(mat.epoch, notification.OrderFragment, notification.Trader); err != nil {
			log.Printf("[error] (generator) cannot store sell order fragment = %v: %v", notification.OrderID, err)
			return
		}
		oppositeOrderFragmentIter, err = mat.orderFragmentStore.BuyOrderFragments(mat.epoch)
		if err != nil {
			log.Printf("[error] (generator) cannot load sell order fragment iterator: %v", err)
			return
		}
		defer oppositeOrderFragmentIter.Release()
	}

	// Iterate through the opposing list and generate computations
	for oppositeOrderFragmentIter.Next() {
		orderFragment, trader, err := oppositeOrderFragmentIter.Cursor()
		if err != nil {
			log.Printf("[error] (generator) cannot load cursor: %v", err)
			continue
		}

		// TODO: Order fragments are opened on a hierarchical path through the
		// Darknode pods. Pods should prioritise computations for which they
		// are the first pod along that path (this is how pods break ties for
		// computations that otherwise have the same priority).

		// Traders should not match against themselves
		if trader == notification.Trader {
			continue
		}

		// TODO: Check that at least one of the orders in the pairing was
		// opened during this matrix epoch. Otherwise, orders that are opened
		// in the same epoch will be matched twice. Once in the current epoch,
		// and once in the previous epoch.

		var computation Computation
		if notification.OrderFragment.OrderParity == order.ParityBuy {
			computation = NewComputation(mat.epoch.Hash, notification.OrderFragment, orderFragment, ComputationStateNil, false)
		} else {
			computation = NewComputation(mat.epoch.Hash, orderFragment, notification.OrderFragment, ComputationStateNil, false)
		}

		select {
		case <-done:
			return
		case computations <- computation:
		}
	}
}

func (mat *computationMatrix) removeOrderFragment(orderID order.ID) {
	if err := mat.orderFragmentStore.DeleteBuyOrderFragment(mat.epoch, orderID); err != nil {
		log.Printf("[error] (generator) cannot delete order fragment = %v; %v", orderID, err)
	}
	if err := mat.orderFragmentStore.DeleteSellOrderFragment(mat.epoch, orderID); err != nil {
		log.Printf("[error] (generator) cannot delete order fragment = %v; %v", orderID, err)
	}
}
