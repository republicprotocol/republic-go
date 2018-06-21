package orderbook

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
)

// ErrOrderFragmentNotFound is returned when attempting to read an order that
// cannot be found.
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

// ErrOrderNotFound is returned whened attempting to read an order that cannot
// be found.
var ErrOrderNotFound = errors.New("order not found")

// ErrCursorOutOfRange is returned when an iterator cursor is used to read a
// value outside the range of the iterator.
var ErrCursorOutOfRange = errors.New("cursor out of range")

// OrderStorer for the order.Orders that are synchronised.
type OrderStorer interface {
	PutOrder(order order.Order) error
	DeleteOrder(id order.ID) error
	Order(id order.ID) (order.Order, error)
	Orders() (OrderIterator, error)
}

// OrderIterator is used to iterate over an order.order collection.
type OrderIterator interface {

	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the order.Order collection, otherwise false.
	Next() bool

	// Cursor returns the order.Order at the current cursor location. Returns
	// an error if the cursor is out of range.
	Cursor() (order.Order, error)

	// Collect all order.Orders in the iterator into a slice.
	Collect() ([]order.Order, error)
}

// OrderFragmentStorer for the order.Fragments that are received.
type OrderFragmentStorer interface {
	PutOrderFragment(orderFragment order.Fragment) error
	DeleteOrderFragment(id order.ID) error
	OrderFragment(id order.ID) (order.Fragment, error)
	OrderFragments() (OrderFragmentIterator, error)
}

// OrderFragmentIterator is used to iterate over an order.Fragment collection.
type OrderFragmentIterator interface {

	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the order.Fragment collection, otherwise false.
	Next() bool

	// Cursor returns the order.Fragment at the current cursor location.
	// Returns an error if the cursor is out of range.
	Cursor() (order.Fragment, error)

	// Collect all order.Fragments in the iterator into a slice.
	Collect() ([]order.Fragment, error)
}

// PointerStorer for the synchronisation pointers used to track the progress
// of synchronising order.Orders.
type PointerStorer interface {
	PutBuyPointer(pointer Pointer) error
	BuyPointer() (Pointer, error)

	PutSellPointer(pointer Pointer) error
	SellPointer() (Pointer, error)
}

// Pointer points to the last order.Order that was successfully synchronised.
type Pointer int

// SyncStorer combines the OrderStorer interface and the PointerStorer
// interface into a unified set of storage functions that are useful for
// synchronisation.
type SyncStorer interface {
	OrderStorer
	PointerStorer
}
