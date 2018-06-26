package orderbook

import (
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

// A Filter accepts a channel of Notifications and filters them into a new
// channel of Notifications.
type Filter interface {
	Filter(done <-chan struct{}, notifications <-chan Notification) (<-chan Notification, <-chan error)
}

type mergeFilter struct {
	orderMu            *sync.Mutex
	orderStore         OrderStorer
	orderFragmentStore OrderFragmentStorer

	buffer int
}

// NewMergeFilter returns a Filter that captures NotificationSyncOrders and
// NotificationSyncOrderFragments. It uses an OrderStorer and an
// OrderFragmentStorer to store open order.Orders and order.Fragments until
// it can merge them into a unified NotificationOpenOrder. It does not need
// exclusive ownership over these stores. All other other Notifications are
// forwarded to the output channel without modification. This is particularly
// convenient for waiting until both an order.Order status is opened, and its
// respective order.Fragment is received, before processing an order.Order
// further.
func NewMergeFilter(orderStore OrderStorer, orderFragmentStore OrderFragmentStorer, buffer int) Filter {
	return &mergeFilter{
		orderMu:            new(sync.Mutex),
		orderStore:         orderStore,
		orderFragmentStore: orderFragmentStore,

		buffer: buffer,
	}
}

// Filter implements the Filter interface.
func (filter *mergeFilter) Filter(done <-chan struct{}, in <-chan Notification) (<-chan Notification, <-chan error) {
	out := make(chan Notification, filter.buffer)
	errs := make(chan error, filter.buffer)

	go func() {
		defer close(out)
		defer close(errs)

		for {
			select {
			case <-done:
				return
			case notification, ok := <-in:
				if !ok {
					return
				}
				filter.handleInternalNotification(notification, done, out, errs)
			}
		}
	}()

	return out, errs
}

func (filter *mergeFilter) handleInternalNotification(notification Notification, done <-chan struct{}, out chan<- Notification, errs chan<- error) {
	if notification == nil {
		return
	}

	// Check the notification type
	switch notification := notification.(type) {
	case NotificationSyncOrder:
		filter.handleNotificationSyncOrder(notification, done, out, errs)
	case NotificationSyncOrderFragment:
		filter.handleNotificationSyncOrderFragment(notification, done, out, errs)
	default:
		select {
		case <-done:
			return
		case errs <- ErrUnexpectedNotificationType:
		}
	}
}

func (filter *mergeFilter) handleNotificationSyncOrder(notification NotificationSyncOrder, done <-chan struct{}, out chan<- Notification, errs chan<- error) {
	// Store the order.Order with the order.Open status
	if err := filter.orderStore.PutOrder(notification.OrderID, notification.OrderStatus); err != nil {
		select {
		case <-done:
		case errs <- err:
		}
		// Continue handling the notification
	}

	if notification.OrderStatus != order.Open {
		select {
		case <-done:
		case out <- notification:
		}
		return
	}

	// Check for a respective order.Fragment
	orderFragment, err := filter.orderFragmentStore.OrderFragment(notification.OrderID)
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
	notificationOpenOrder := NotificationOpenOrder{
		OrderFragment: orderFragment,
		OrderID:       notification.OrderID,
		OrderStatus:   notification.OrderStatus,
	}
	select {
	case <-done:
	case out <- notificationOpenOrder:
	}
}

func (filter *mergeFilter) handleNotificationSyncOrderFragment(notification NotificationSyncOrderFragment, done <-chan struct{}, out chan<- Notification, errs chan<- error) {

	// Store the order.Fragment
	if err := filter.orderFragmentStore.PutOrderFragment(notification.OrderFragment); err != nil {
		select {
		case <-done:
		case errs <- err:
		}
		// Continue handling the notification
	}

	// Check for a respective order.Status
	orderStatus, err := filter.orderStore.Order(notification.OrderFragment.OrderID)
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
	notificationOut := NotificationOpenOrder{
		OrderFragment: notification.OrderFragment,
		OrderID:       notification.OrderFragment.OrderID,
		OrderStatus:   orderStatus,
	}
	select {
	case <-done:
	case out <- notificationOut:
	}

}
