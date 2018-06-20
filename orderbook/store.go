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

// Storer for the orders and order fragments that are received by the
// Orderbook.
type Storer interface {
	PutOrder(order order.Order) error
	DeleteOrder(id order.ID) error
	Order(id order.ID) (order.Order, error)
	Orders() (OrderIterator, error)

	PutOrderFragment(orderFragment order.Fragment) error
	DeleteOrderFragment(id order.ID) error
	OrderFragment(id order.ID) (order.Fragment, error)
	OrderFragments() (OrderFragmentIterator, error)

	PutBuyPointer(pointer Pointer) error
	BuyPointer() (Pointer, error)

	PutOrderPointer(pointer Pointer) error
	OrderPointer() (Pointer, error)
}

// OrderIterator is used to iterate over an order.order collection.
type OrderIterator interface {

	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the order.Order collection, otherwise false.
	Next() bool

	// Cursor returns the order.Order at the current cursor location. Returns
	// an error if the cursor is out of range.
	Cursor() (order.Order, error)
}

// OrderFragmentIterator is used to iterate over an order.Fragment collection.
type OrderFragmentIterator interface {

	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the order.Fragment collection, otherwise false.
	Next() bool

	// Cursor returns the order.Fragment at the current cursor location.
	// Returns an error if the cursor is out of range.
	Cursor() (order.Fragment, error)
}

// Pointer points to the last order.Order that was successfully synchronised.
type Pointer int
