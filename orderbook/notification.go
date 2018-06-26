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

// NotificationSyncOrderFragment is used to signal the synchronisation
// of an order.Fragment from an external user, or from another Darknode.
type NotificationSyncOrderFragment struct {
	OrderFragment order.Fragment
}

// IsNotification implements the Notification interface.
func (notification NotificationSyncOrderFragment) IsNotification() {}

// NotificationSyncOrder is used to signal the synchronisation of an
// order.Order from the Ethereum blockchain. This happens multiple times during
// the lifetime of an order.Order as the order.Status is updated.
type NotificationSyncOrder struct {
	OrderID     order.ID
	OrderStatus order.Status
}

// IsNotification implements the Notification interface.
func (notification NotificationSyncOrder) IsNotification() {}

// NotificationOpenOrder is used to signal the opening of an order.Order. This
// happens once the order.Fragment has been synchronised and the order.Order
// has been synchronised to the order.Open status.
type NotificationOpenOrder struct {
	OrderFragment order.Fragment
	OrderID       order.ID
	OrderStatus   order.Status
}

// IsNotification implements the Notification interface.
func (notification NotificationOpenOrder) IsNotification() {}
