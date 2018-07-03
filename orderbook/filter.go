package orderbook

import (
	"github.com/republicprotocol/republic-go/order"
)

type merger struct {
	orderStore         OrderStorer
	orderFragmentStore OrderFragmentStorer
	buffer             int
}

func newMerger(orderStore OrderStorer, orderFragmentStore OrderFragmentStorer, buffer int) *merger {
	return &merger{
		orderStore:         orderStore,
		orderFragmentStore: orderFragmentStore,
		buffer:             buffer,
	}
}

func (merger *merger) merge(done <-chan struct{}, notificationsIn <-chan Notification, orderFragments <-chan order.Fragment) (<-chan Notification, <-chan error) {
	notifications := make(chan Notification, merger.buffer)
	errs := make(chan error, merger.buffer)

	go func() {
		defer close(notifications)
		defer close(errs)

		for {
			select {
			case <-done:
				return
			case notification, ok := <-notificationsIn:
				if !ok {
					return
				}
				merger.handleNotification(notification, done, notifications, errs)
			case orderFragment, ok := <-orderFragments:
				if !ok {
					return
				}
				merger.handleOrderFragment(orderFragment, done, notifications, errs)
			}
		}
	}()

	return notifications, errs
}

func (merger *merger) handleNotification(notification Notification, done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {
	if notification == nil {
		return
	}

	// Check the notification type
	switch notification := notification.(type) {
	case NotificationOpenOrder:
		merger.handleNotificationOpenOrder(notification, done, notifications, errs)
	default:
		select {
		case <-done:
			return
		case notifications <- notification:
		}
	}
}

func (merger *merger) handleNotificationOpenOrder(notification NotificationOpenOrder, done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {
	// Store the order.Order with the order.Open status
	if err := merger.orderStore.PutOrder(notification.OrderID, order.Open); err != nil {
		select {
		case <-done:
		case errs <- err:
		}
		// Continue handling the notification
	}

	// Check for a respective order.Fragment
	orderFragment, err := merger.orderFragmentStore.OrderFragment(notification.OrderID)
	if err != nil {
		if err == ErrOrderFragmentNotFound {
			// If there is no respective order.Fragment, do nothing and
			// return
			return
		}
		select {
		case <-done:
		case errs <- err:
		}
		return
	}

	notification.OrderFragment = orderFragment
	select {
	case <-done:
	case notifications <- notification:
	}
}

func (merger *merger) handleOrderFragment(orderFragment order.Fragment, done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {
	// Store the order.Fragment
	if err := merger.orderFragmentStore.PutOrderFragment(orderFragment); err != nil {
		select {
		case <-done:
		case errs <- err:
		}
		// Continue handling the notification
	}

	// Check for a respective order.Status
	orderStatus, err := merger.orderStore.Order(orderFragment.OrderID)
	if err != nil {
		if err == ErrOrderNotFound {
			// If there is no respective order.Order, do nothing and return
			return
		}
		select {
		case <-done:
		case errs <- err:
		}
		return
	}

	// Emit a notification for this order.Order and order.Fragment
	notification := NotificationOpenOrder{
		OrderID:       orderFragment.OrderID,
		OrderFragment: orderFragment,
	}
	select {
	case <-done:
	case notifications <- notification:
	}
}
