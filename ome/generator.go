package ome

import (
	"log"
	"sort"
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
)

type computationWeight struct {
	computation Computation
	weight      uint64
}

type ComputationGenerator interface {
	Generate(done <-chan struct{}, notifications <-chan orderbook.Notification) (<-chan Computation, <-chan error)
	OnChangeEpoch(epoch registry.Epoch)
}

type computationGenerator struct {
	doneMu *sync.Mutex
	done   chan struct{}

	addr identity.Address

	matMu                *sync.Mutex
	matCurrDone          chan struct{}
	matCurrNotifications chan orderbook.Notification
	matPrevDone          chan struct{}
	matPrevNotifications chan orderbook.Notification

	broadcastComputations chan (<-chan Computation)
	broadcastErrs         chan (<-chan error)

	orderFragmentStore OrderFragmentStorer
}

func NewComputationGenerator(addr identity.Address, orderFragmentStore OrderFragmentStorer) ComputationGenerator {
	return &computationGenerator{
		doneMu: new(sync.Mutex),
		done:   nil,

		addr: addr,

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

	mat := newComputationMatrix(gen.addr, epoch, gen.orderFragmentStore)
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
	pod                *registry.Pod
	epoch              registry.Epoch
	orderFragmentStore OrderFragmentStorer

	sortedComputationsMu     *sync.Mutex
	sortedComputations       []computationWeight
	sortedComputationsSignal chan struct{}
}

func newComputationMatrix(addr identity.Address, epoch registry.Epoch, orderFragmentStore OrderFragmentStorer) *computationMatrix {
	mat := &computationMatrix{
		epoch:              epoch,
		orderFragmentStore: orderFragmentStore,

		sortedComputationsMu:     new(sync.Mutex),
		sortedComputations:       []computationWeight{},
		sortedComputationsSignal: make(chan struct{}),
	}
	pod, err := epoch.Pod(addr)
	if err != nil {
		mat.pod = nil
	} else {
		mat.pod = &pod
	}
	return mat
}

func (mat *computationMatrix) generate(done <-chan struct{}, notifications <-chan orderbook.Notification) (<-chan Computation, <-chan error) {
	computations := make(chan Computation)
	errs := make(chan error)

	go func() {
		defer close(computations)
		defer close(errs)

		dispatch.CoBegin(
			func() {
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
			},
			func() {
				for {
					select {
					case <-done:
						return
					case <-mat.sortedComputationsSignal:

						mat.sortedComputationsMu.Lock()
						if len(mat.sortedComputations) == 0 {
							mat.sortedComputationsMu.Unlock()
							continue
						}
						computationWeight := mat.sortedComputations[0]
						mat.sortedComputations = mat.sortedComputations[1:]
						mat.sortedComputationsMu.Unlock()

						select {
						case <-done:
							return
						case computations <- computationWeight.computation:
						}
					}
				}
			})
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
	// If we are not part of a pod during this epoch then we cannot process
	// computations
	if mat.pod == nil {
		return
	}

	// Store the order.Fragment and get the opposing list so that computations
	// can be generated
	var oppositeOrderFragmentIter OrderFragmentIterator
	var err error

	if notification.OrderFragment.OrderParity == order.ParityBuy {
		if err := mat.orderFragmentStore.PutBuyOrderFragment(mat.epoch, notification.OrderFragment, notification.Trader, uint64(notification.Priority)); err != nil {
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
		if err := mat.orderFragmentStore.PutSellOrderFragment(mat.epoch, notification.OrderFragment, notification.Trader, uint64(notification.Priority)); err != nil {
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

	mat.sortedComputationsMu.Lock()
	// Iterate through the opposing list and generate computations
	didGenerateNewComputation := false
	for oppositeOrderFragmentIter.Next() {
		orderFragment, trader, priority, err := oppositeOrderFragmentIter.Cursor()
		if err != nil {
			log.Printf("[error] (generator) cannot load cursor: %v", err)
			continue
		}

		if !isCompatible(notification, orderFragment, trader, priority) {
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

		// Get the priority adjustment based on the distance of our pod from
		// the first pod that could possibly do this computation
		buyPath := mat.epoch.Pods.PathOfOrder(computation.Buy.OrderID)
		sellPath := mat.epoch.Pods.PathOfOrder(computation.Sell.OrderID)
		commonPath := buyPath.Ancestor(sellPath)
		index, ok := commonPath.IndexOfPod(mat.pod)
		if !ok {
			log.Printf("[error] (generator) received orders with divergent paths")
			continue
		}
		adjustment := uint64(len(commonPath) - (index + 1))
		comWeight := computationWeight{weight: uint64(notification.Priority) + priority + adjustment, computation: computation}
		func() {
			// Insert sort into the list of sorted computations
			didGenerateNewComputation = true
			if len(mat.sortedComputations) == 0 {
				mat.sortedComputations = append(mat.sortedComputations, comWeight)
				return
			}
			n := sort.Search(len(mat.sortedComputations), func(i int) bool {
				return comWeight.weight >= mat.sortedComputations[i].weight
			})
			mat.sortedComputations = append(mat.sortedComputations[:n], append([]computationWeight{comWeight}, mat.sortedComputations[n:]...)...)
		}()
	}
	mat.sortedComputationsMu.Unlock()
	if didGenerateNewComputation {
		select {
		case <-done:
		case mat.sortedComputationsSignal <- struct{}{}:
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

// isCompatible checks if the notification's order is compatible with another order based
// on the following conditions:
// 1. If the trader is the same as the notification's trader, the 2 orders are incompatible.
// 2. If both orders are Fill-or-Kill (FOK), they are incompatible.
// 3. If one of the orders is a FOK, then both orders are incompatible if the other order
//    is of a higher priority.
func isCompatible(notification orderbook.NotificationOpenOrder, orderFragment order.Fragment, trader string, priority uint64) bool {

	// Traders should not match against themselves
	if trader == notification.Trader {
		return false
	}

	switch orderFragment.OrderType {

	// The orderFragment is a Fill-or-Kill order
	case order.TypeMidpointFOK, order.TypeLimitFOK:
		switch notification.OrderFragment.OrderType {
		case order.TypeMidpointFOK, order.TypeLimitFOK:
			// Both orders are FOK, thus, incompatible.
			return false
		default:
			// Does notification.OrderFragment, which is not an FOK order, have a higher
			// priority than the FOK order ?
			if uint64(notification.Priority) > priority {
				return false
			}
			return true
		}

	// The orderFragment is not a Fill-or-Kill order
	default:
		switch notification.OrderFragment.OrderType {
		case order.TypeMidpointFOK, order.TypeLimitFOK:
			// Does notification.OrderFragment, which is an FOK order, have a lower
			// priority than the other order ?
			if priority > uint64(notification.Priority) {
				return false
			}
			return true
		default:
			return true
		}
	}
}
