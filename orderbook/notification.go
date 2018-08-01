package orderbook

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
)

// ErrUnexpectedNotificationType is returned when a component receives a type
// of Notification that it cannot process.
var ErrUnexpectedNotificationType = errors.New("unexpected notification type")

// A Notification is emitted by the Notifier to notify other components that
// there has been a change to an order.Order or order.Fragment. This interface
// should not be implemented outside of this package.
type Notification interface {

	// Types that can be used as a Notification must implement this pass
	// through method. It only exists to restrict Notifications to types that
	// have been explicitly marked as compatible to avoid programmer error.
	IsNotification()
}

// Notifications is a slice.
type Notifications []Notification

// NotificationOpenOrder is used to signal the opening of an order.Order. This
// happens when a rader opens an order.Order and the order.Fragment has been
// received.
type NotificationOpenOrder struct {
	OrderID       order.ID
	OrderFragment order.Fragment
	Trader        string
	Priority      uint
}

// IsNotification implements the Notification interface.
func (notification NotificationOpenOrder) IsNotification() {}

// IsCompatible checks the compatibility of the notification's order with another
// order.
// 1. If the trader is the same as the notification's trader, the 2 orders are
// incompatible. (Traders should not match against themselves)
// 2. If both orders are Fill-or-Kill (FOK), they are incompatible.
// 3. If one of the orders is a FOK, then both orders are incompatible if the other order
// is of a higher priority.
func (notification NotificationOpenOrder) IsCompatible(orderFragment order.Fragment, trader string, priority uint64) bool {
	if trader == notification.Trader {
		return false
	}

	if orderFragment.OrderType >= order.TypeFOK {
		if notification.OrderFragment.OrderType >= order.TypeFOK || uint64(notification.Priority) > priority {
			return false
		}
		return true
	}
	if notification.OrderFragment.OrderType >= order.TypeFOK && priority > uint64(notification.Priority) {
		return false
	}
	return true
}

// NotificationConfirmOrder is used to signal the confirmation of an
// order.ID. This happens when an order.Order has been matched with another
// order.Order, and consensus has been reached for the match.
type NotificationConfirmOrder struct {
	OrderID order.ID
}

// IsNotification implements the Notification interface.
func (notification NotificationConfirmOrder) IsNotification() {}

// NotificationCancelOrder is used to signal the cancelation of an order.ID.
// This happens when a trader cancels an order.Order before it has been
// confirmed.
type NotificationCancelOrder struct {
	OrderID order.ID
}

// IsNotification implements the Notification interface.
func (notification NotificationCancelOrder) IsNotification() {}
