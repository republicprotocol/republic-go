package ome

import (
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
}

func NewComputationGenerator() ComputationGenerator {
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

	// Merge all of the channels on the broadcast channel into the output
	// channel
	go func() {
		defer close(computations)
		dispatch.Merge(done, gen.broadcastComputations, computations)
	}()
	go func() {
		defer close(errs)
		dispatch.Merge(done, gen.broadcastErrs, errs)
	}()

	return computations, errs
}

// OnChangeEpoch implements the ComputationGenerator interface.
func (gen *computationGenerator) OnChangeEpoch(epoch registry.Epoch) {
	gen.matMu.Lock()
	defer gen.matMu.Unlock()

	// Transition the current epoch into the previous epoch and setup a new
	// current epoch
	if gen.matPrevDone == nil {
		close(gen.matPrevDone)
		close(gen.matPrevNotifications)
	}
	gen.matPrevDone = gen.matCurrDone
	gen.matPrevNotifications = gen.matCurrNotifications
	gen.matCurrDone = make(chan struct{})
	gen.matCurrNotifications = make(chan orderbook.Notification)

	mat := newComputationMatrix(epoch)
	computations, errs := mat.generate(gen.matCurrDone, gen.matCurrNotifications)

	select {
	case <-gen.done:
	case gen.broadcastComputations <- computations:
	}

	select {
	case <-gen.done:
	case gen.broadcastErrs <- errs:
	}
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

	switch notification.OrderFragment.Depth {
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
	buyOrderFragments  []order.Fragment
	sellOrderFragments []order.Fragment
}

func newComputationMatrix(epoch registry.Epoch) *computationMatrix {
	return &computationMatrix{
		epoch:              epoch,
		buyOrderFragments:  []order.Fragment{},
		sellOrderFragments: []order.Fragment{},
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
	var cmpOrderFragments []order.Fragment
	if notification.OrderFragment.OrderParity == order.ParityBuy {
		mat.buyOrderFragments = append(mat.buyOrderFragments, notification.OrderFragment)
		cmpOrderFragments = mat.sellOrderFragments
	} else {
		mat.sellOrderFragments = append(mat.sellOrderFragments, notification.OrderFragment)
		cmpOrderFragments = mat.buyOrderFragments
	}

	// Iterate through the opposing list and generate computations
	for i := range cmpOrderFragments {

		// TODO: Order fragments are opened on a hierarchical path through the
		// Darknode pods. Pods should prioritise computations for which they
		// are the first pod along that path (this is how pods break ties for
		// computations that otherwise have the same priority).

		// FIXME: Order fragments from the same trader should not be matched
		// against each other.

		var computation Computation
		if notification.OrderFragment.OrderParity == order.ParityBuy {
			computation = NewComputation(mat.epoch.Hash, notification.OrderFragment.OrderID, cmpOrderFragments[i].OrderID, ComputationStateNil, false)
		} else {
			computation = NewComputation(mat.epoch.Hash, cmpOrderFragments[i].OrderID, notification.OrderFragment.OrderID, ComputationStateNil, false)
		}

		select {
		case <-done:
			return
		case computations <- computation:
		}
	}
}

func (mat *computationMatrix) removeOrderFragment(orderID order.ID) {
	for i := range mat.buyOrderFragments {
		if mat.buyOrderFragments[i].OrderID == orderID {
			mat.buyOrderFragments = append(mat.buyOrderFragments[:i], mat.buyOrderFragments[i+1:]...)
			return
		}
	}
	for i := range mat.sellOrderFragments {
		if mat.sellOrderFragments[i].OrderID == orderID {
			mat.sellOrderFragments = append(mat.sellOrderFragments[:i], mat.sellOrderFragments[i+1:]...)
			return
		}
	}
}
