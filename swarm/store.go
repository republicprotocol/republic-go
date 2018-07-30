package swarm

import (
	"errors"

	"github.com/republicprotocol/republic-go/identity"
)

// ErrMultiAddressNotFound is returned from a query when no
// identity.MultiAddress can be found for the identity.Address.
var ErrMultiAddressNotFound = errors.New("multiaddress not found")

// ErrCursorOutOfRange is returned when an iterator cursor is used to read a
// value outside the range of the iterator.
var ErrCursorOutOfRange = errors.New("cursor out of range")

// MultiAddressStorer for the identity.MultiAddresses that are registered with
// the dark node registry.
type MultiAddressStorer interface {
	PutMultiAddress(multiaddress identity.MultiAddress) (bool, error)
	MultiAddress(address identity.Address) (identity.MultiAddress, error)
	MultiAddresses() (MultiAddressIterator, error)
}

// MultiAddressIterator is used to iterate over an identity.MultiAddress collection.
type MultiAddressIterator interface {

	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the identity.MultiAddress collection, otherwise false.
	Next() bool

	// Cursor returns the identity.MultiAddress at the current cursor location.
	// Returns an error if the cursor is out of range.
	Cursor() (identity.MultiAddress, error)

	// Collect all identity.MultiAddresses in the iterator into slices.
	Collect() ([]identity.MultiAddress, error)

	// Release the resources allocated by the iterator.
	Release()
}
