package orderbook

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
)

// ErrOrderFragmentNotFound is return when attempting to load an order that
// cannot be found.
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

// ErrOrderNotFound is return when attempting to load an order that cannot be
// found.
var ErrOrderNotFound = errors.New("order not found")

// Storer for the orders and order fragments that are received by the
// Orderbook.
type Storer interface {
	InsertOrderFragment(orderFragment order.Fragment) error
	InsertOrder(order order.Order) error
	OrderFragment(id order.ID) (order.Fragment, error)
	Order(id order.ID) (order.Order, error)
	Orders() ([]order.Order, error)
	RemoveOrderFragment(id order.ID) error
	RemoveOrder(id order.ID) error
}

// SyncPointer points to the last order.Order that was successfully
// synchronized by the Syncer. It is primarily used to prevent re-syncing all
// orders aftter a reboot.
type SyncPointer = int

// SyncStorer exposes functionality for storing and loading synchronization
// data that the Syncer uses to keep track of where it is in the
// synchronization process.
type SyncStorer interface {

	// InsertBuyPointer into the Storer. The prevents the Syncer needing to
	// re-sync all buy orders after a reboot.
	InsertBuyPointer(SyncPointer) error

	// InsertSellPointer into the Storer. The prevents the Syncer needing to
	// re-sync all sell orders after a reboot.
	InsertSellPointer(SyncPointer) error

	// BuyPointer returns the SyncPointer stored in the Storer. It defaults to
	// zero.
	BuyPointer() (SyncPointer, error)

	// SellPointer returns the SyncPointer stored in the Storer. It defaults to
	// zero.
	SellPointer() (SyncPointer, error)
}
