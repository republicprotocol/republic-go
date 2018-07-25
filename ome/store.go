package ome

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
)

// ErrComputationNotFound is returned when the Storer cannot find a Computation
// associated with a ComputationID.
var ErrComputationNotFound = errors.New("computation not found")

// ErrOrderFragmentNotFound is returned when attempting to read an order that
// cannot be found.
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

// ErrCursorOutOfRange is returned when an iterator cursor is used to read a
// value outside the range of the iterator.
var ErrCursorOutOfRange = errors.New("cursor out of range")

// ComputationStorer for the Computations that are synchronised.
type ComputationStorer interface {
	PutComputation(computation Computation) error
	DeleteComputation(id ComputationID) error
	Computation(id ComputationID) (Computation, error)
	Computations() (ComputationIterator, error)
}

// ComputationIterator is used to iterate over a Computation collection.
type ComputationIterator interface {

	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the Computation collection, otherwise false.
	Next() bool

	// Cursor returns the Computation at the current cursor location. Returns
	// an error if the cursor is out of range.
	Cursor() (Computation, error)

	// Collect all Computations in the iterator into a slice.
	Collect() ([]Computation, error)

	// Release the resources allocated by the iterator.
	Release()
}

// OrderFragmentStorer for the order.Fragments that are received.
type OrderFragmentStorer interface {
	PutBuyOrderFragment(epoch registry.Epoch, orderFragment order.Fragment, trader string) error
	DeleteBuyOrderFragment(epoch registry.Epoch, id order.ID) error
	BuyOrderFragment(epoch registry.Epoch, id order.ID) (order.Fragment, string, error)
	BuyOrderFragments(epoch registry.Epoch) (OrderFragmentIterator, error)

	PutSellOrderFragment(epoch registry.Epoch, orderFragment order.Fragment, trader string) error
	DeleteSellOrderFragment(epoch registry.Epoch, id order.ID) error
	SellOrderFragment(epoch registry.Epoch, id order.ID) (order.Fragment, string, error)
	SellOrderFragments(epoch registry.Epoch) (OrderFragmentIterator, error)
}

// OrderFragmentIterator is used to iterate over an order.Fragment collection.
type OrderFragmentIterator interface {

	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the order.Fragment collection, otherwise false.
	Next() bool

	// Cursor returns the order.Fragment at the current cursor location.
	// Returns an error if the cursor is out of range.
	Cursor() (order.Fragment, string, error)

	// Collect all order.Fragments in the iterator into a slice.
	Collect() ([]order.Fragment, []string, error)

	// Release the resources allocated by the iterator.
	Release()
}
