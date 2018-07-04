package ome

import (
	"sync"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
)

type ComputationGenerator interface {
	OnChangeEpoch(epoch registry.Epoch)
}

type computationGenerator struct {
	matMu                *sync.Mutex
	matCurrDone          chan struct{}
	matCurrNotifications chan orderbook.Notification
	matPrevDone          chan struct{}
	matPrevNotifications chan orderbook.Notification
}

func (gen *computationGenerator) Generate(done <-chan struct{}, notifications <-chan orderbook.Notification) (<-chan Computation, <-chan error) {
	computations := make(chan Computation)
	errs := make(chan error)

	// Handle all incoming notifications by mapping them to the appropriate
	// computation matrix
	go func() {
		defer close(errs)
		for {
			select {
			case <-done:
				return
			case notification, ok := <-notifications:
				if !ok {
					return
				}
				gen.handleNotification(notification, done, errs)
			}
		}
	}()

	// Merge all outputs from the computation matrices
	go func() {
		defer close(computations)
		gen.mergeComputations(done, computations)
	}()
	go func() {
		defer close(errs)
		gen.mergeErrors(done, errs)
	}()

	return computations, errs
}

// OnChangeEpoch implements the ComputationGenerator interface.
func (gen *computationGenerator) OnChangeEpoch(epoch registry.Epoch) {
	gen.matMu.Lock()
	defer gen.matMu.Unlock()

	// Close the previous epoch
	if gen.matPrevDone == nil {
		close(gen.matPrevDone)
		close(gen.matPrevNotifications)
	}
	// Transition the current epoch to be the previous epoch and create a new
	// epoch setup
	gen.matPrevDone = gen.matCurrDone
	gen.matPrevNotifications = gen.matCurrNotifications
	gen.matCurrDone = make(chan struct{})
	gen.matCurrNotifications = make(chan orderbook.Notification)

	// Start the mat for this epoch
	mat := newComputationMatrix()
	computations, errs := mat.generate(gen.matCurrDone, gen.matCurrNotifications)

	// Signal that the outputs of this mat should be accepted by the merger
	select {
	case <-gen.done:
	case gen.computationsMerger <- computations:
	}
	select {
	case <-gen.done:
	case gen.errMerger <- errs:
	}
}

func (gen *computationGenerator) handleNotification(notification orderbook.Notification, done <-chan struct{}) {
	switch notification := notification.(type) {
	case orderbook.NotificationOpenOrder:
		gen.handleNotificationOpenOrder(notification, done, errs)
	default:
		select {
		case <-done:
			return
		case gen.matCurrNotifications <- notification:
		}
		select {
		case <-done:
			return
		case gen.matPrevNotifications <- notification:
		}
	}
}

func (gen *computationGenerator) handleNotificationOpenOrder(notification orderbook.NotificationOpenOrder, done <-chan struct{}) {
	switch notification.OrderFragment.Depth {
	case 0:
		select {
		case <-done:
			return
		case gen.matCurrNotifications <- notification:
		}
	case 1:
		select {
		case <-done:
			return
		case gen.matPrevNotifications <- notification:
		}
	}
}

func (gen *computationGenerator) mergeComputations(done <-chan struct{}, computations chan<- Computation) {
	for {
		select {
		case <-done:
			return
		case ch, ok := <-gen.computationMerger:
			if !ok {
				return
			}
			go func() {
				for {
					select {
					case <-done:
						return
					case computation, ok := <-ch:
						if !ok {
							return
						}
						select {
						case <-done:
						case computations <- computation:
						}
					}
				}
			}()
		}
	}
}

func (gen *computationGenerator) mergeErrors(done <-chan struct{}, errs chan<- error) {
	for {
		select {
		case <-done:
			return
		case ch, ok := <-gen.errMerger:
			if !ok {
				return
			}
			go func() {
				for {
					select {
					case <-done:
						return
					case err, ok := <-ch:
						if !ok {
							return
						}
						select {
						case <-done:
						case errs <- err:
						}
					}
				}
			}()
		}
	}
}

type computationMatrix struct {
	buyOrderFragments  []order.Fragment
	sellOrderFragments []order.Fragment
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
	case orderbook.NotificationOpenOrder:
		mat.handleNotificationOpenOrder(notification, done, computations, errs)
	case orderbook.NotificationConfirmOrder:
		mat.handleNotificationConfirmOrder(notification)
	case orderbook.NotificationCancelOrder:
		mat.handleNotificationCancelOrder(notification)
	default:
		select {
		case <-done:
		case errs <- orderbook.ErrUnexpectedNotificationType:
		}
	}
}

func (mat *computationMatrix) handleNotificationOpenOrder(notification orderbook.NotificationOpenOrder, done <-chan struct{}, computations chan<- Computation, errs chan<- error) {

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
		// TODO: Only compute a subset of the computations generated to
		// increase system-wide parallelism
		// TODO: Check that a trader is not matching against themselves

		var computation Computation
		if notification.OrderFragment.OrderParity == order.ParityBuy {
			computation = NewComputation(notification.OrderFragment.OrderID, cmpOrderFragments[i].OrderID, ComputationStateNil, false)
		} else {
			computation = NewComputation(cmpOrderFragments[i].OrderID, notification.OrderFragment.OrderID, ComputationStateNil, false)
		}

		select {
		case <-done:
			return
		case computations <- computation:
		}
	}
}

func (mat *computationMatrix) handleNotificationConfirmOrder(notification orderbook.NotificationConfirmOrder) {
	mat.removeOrderFragment(notification.OrderID)
}

func (mat *computationMatrix) handleNotificationCancelOrder(notification orderbook.NotificationCancelOrder) {
	mat.removeOrderFragment(notification.OrderID)
}

func (mat *computationMatrix) removeOrderFragment(orderID order.ID) {
	for i := range mat.buyOrderFragments {
		if mat.buyOrderFragments[i].OrderID == orderID {
			mat.buyOrderFragments = append(mat.buyOrderFragments[:i], mat.buyOrderFragments[i+1:])
			return
		}
	}
	for i := range mat.sellOrderFragments {
		if mat.sellOrderFragments[i].OrderID == orderID {
			mat.sellOrderFragments = append(mat.sellOrderFragments[:i], mat.sellOrderFragments[i+1:])
			return
		}
	}
}
