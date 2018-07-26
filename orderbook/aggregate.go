package orderbook

import (
	"log"

	"github.com/republicprotocol/republic-go/identity"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
)

// An Aggregator collects order statuses and order fragments, emitting a
// notification when both have been received for a given order.
type Aggregator interface {

	// InsertOrder status into the Aggregator. Returns a Notification
	// if the respective order fragment has already been inserted, and the
	// order was inserted with the open status.
	InsertOrder(orderID order.ID, orderStatus order.Status, trader string) (Notification, error)

	// InsertOrderFragment into the Aggregator. Returns a Notification
	// if the respective order is currently inserted with the open status.
	InsertOrderFragment(orderFragment order.Fragment) (Notification, error)
}

type aggregator struct {
	addr               identity.Address
	epoch              registry.Epoch
	orderStore         OrderStorer
	orderFragmentStore OrderFragmentStorer
	pod                *registry.Pod
}

func NewAggregator(addr identity.Address, epoch registry.Epoch, orderStore OrderStorer, orderFragmentStore OrderFragmentStorer) Aggregator {
	agg := &aggregator{
		addr:               addr,
		epoch:              epoch,
		orderStore:         orderStore,
		orderFragmentStore: orderFragmentStore,
	}
	pod, err := epoch.Pod(addr)
	if err != nil {
		agg.pod = nil
	} else {
		agg.pod = &pod
	}
	return agg
}

// InsertOrder implements the Aggregator interface.
func (agg *aggregator) InsertOrder(orderID order.ID, orderStatus order.Status, trader string) (Notification, error) {
	if !agg.isInPathOfEpoch(orderID) {
		return nil, nil
	}

	if orderStatus != order.Open {
		// The order is no longer open
		if err := agg.orderStore.DeleteOrder(orderID); err != nil {
			log.Printf("[error] (sync) cannot delete order: %v", err)
		}
		return nil, nil
	}
	// Store the order
	if err := agg.orderStore.PutOrder(orderID, orderStatus, trader); err != nil {
		return nil, err
	}
	// Fetch the order fragment
	orderFragment, err := agg.orderFragmentStore.OrderFragment(agg.epoch, orderID)
	if err != nil {
		if err == ErrOrderFragmentNotFound {
			// No order fragment was found
			return nil, nil
		}
		return nil, err
	}
	// Produce notification
	log.Printf("[info] (sync) order = %v", orderID)
	return NotificationOpenOrder{
		OrderID:       orderID,
		OrderFragment: orderFragment,
		Trader:        trader,
	}, nil
}

// InsertOrderFragment implements the Aggregator interface.
func (agg *aggregator) InsertOrderFragment(orderFragment order.Fragment) (Notification, error) {
	if !agg.isInPathOfEpoch(orderFragment.OrderID) {
		return nil, nil
	}

	// Store the order fragment
	if err := agg.orderFragmentStore.PutOrderFragment(agg.epoch, orderFragment); err != nil {
		return nil, err
	}
	// Fetch the order
	orderStatus, trader, err := agg.orderStore.Order(orderFragment.OrderID)
	if err != nil {
		if err == ErrOrderNotFound {
			// No order was found
			return nil, nil
		}
		return nil, err
	}
	if orderStatus != order.Open {
		// The order was found but is no longer open
		if err := agg.orderStore.DeleteOrder(orderFragment.OrderID); err != nil {
			log.Printf("[error] (sync) cannot delete order: %v", err)
		}
		if err := agg.orderFragmentStore.DeleteOrderFragment(agg.epoch, orderFragment.OrderID); err != nil {
			log.Printf("[error] (sync) cannot delete order fragment: %v", err)
		}
		return nil, nil
	}
	// Produce notification
	log.Printf("[sync] order = %v", orderFragment.OrderID)
	return NotificationOpenOrder{
		OrderID:       orderFragment.OrderID,
		OrderFragment: orderFragment,
		Trader:        trader,
	}, nil
}

func (agg *aggregator) isInPathOfEpoch(orderID order.ID) bool {
	if agg.pod == nil || agg.epoch.Pods == nil || len(agg.epoch.Pods) == 0 {
		return false
	}
	index, ok := agg.epoch.Pods.PathOfOrder(orderID).IndexOfPod(agg.pod)
	return index >= 0 && ok
}
