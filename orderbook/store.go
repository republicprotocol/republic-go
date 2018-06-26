package orderbook

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
)

// ErrOrderNotFound is returned when attempting to read an order that cannot be
// found.
var ErrOrderNotFound = errors.New("order not found")

// ErrOrderFragmentNotFound is returned when attempting to read an order that
// cannot be found.
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

// ErrChangeNotFound is returned when attempting to read a change that cannot
// be found.
var ErrChangeNotFound = errors.New("order not found")

// ErrCursorOutOfRange is returned when an iterator cursor is used to read a
// value outside the range of the iterator.
var ErrCursorOutOfRange = errors.New("cursor out of range")

// ChangeStorer for the Changes that are synchronised.
type ChangeStorer interface {
	PutChange(change Change) error
	DeleteChange(id order.ID) error
	Change(id order.ID) (Change, error)
	Changes() (ChangeIterator, error)
}

// ChangeIterator is used to iterate over a Change collection.
type ChangeIterator interface {

	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the Change collection, otherwise false.
	Next() bool

	// Cursor returns the Change at the current cursor location. Returns
	// an error if the cursor is out of range.
	Cursor() (Change, error)

	// Collect all Changes in the iterator into a slice.
	Collect() ([]Change, error)

	// Release the resources allocated by the iterator.
	Release()
}

// OrderStorer for the order.Orders that are synchronised from the Ethereum
// blockchain.
type OrderStorer interface {
	PutOrder(id order.ID, status order.Status) error
	DeleteOrder(id order.ID)
	Order(id order.ID) (order.Status, error)
	Orders() (OrderIterator, error)
}

// OrderIterator is used to iterate over an order.Order collection.
type OrderIterator interface {

	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the order.Order collection, otherwise false.
	Next() bool

	// Cursor returns the order.Order at the current cursor location.
	// Returns an error if the cursor is out of range.
	Cursor() (order.ID, order.Status, error)

	// Collect all order.IDs and order.Statuses in the iterator into slices.
	Collect() ([]order.ID, []order.Status, error)

	// Release the resources allocated by the iterator.
	Release()
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

	// Release the resources allocated by the iterator.
	Release()
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

// SyncStorer combines the ChangeStorer interface and the PointerStorer
// interface into a unified interface that is convenient for synchronisation.
type SyncStorer interface {
	ChangeStorer
	PointerStorer
}
