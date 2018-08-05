package oracle

import (
	"errors"
	"sync"
)

// ErrCursorOutOfRange is returned when an iterator cursor is used to read a
// value outside the range of the iterator.
var ErrCursorOutOfRange = errors.New("cursor out of range")

// MidpointPriceStorer is used for retrieving/storing MidpointPrice.
type MidpointPriceStorer interface {

	// PutMidpointPrice stores the given midPointPrice into the storer.
	// It doesn't do any nonce check and will always overwrite previous record.
	PutMidpointPrice(MidpointPrice) error

	// MidpointPrice returns the latest mid-point price of the given token.
	MidpointPrice() (MidpointPrice, error)
}

// midpointPriceStorer implements MidpointPriceStorer interface with an
// in-memory map implementation. Data will be lost every time it restarts.
type midpointPriceStorer struct {
	mu     *sync.Mutex
	prices MidpointPrice
}

// NewMidpointPriceStorer returns a new MidpointPriceStorer.
func NewMidpointPriceStorer() MidpointPriceStorer {
	return &midpointPriceStorer{
		mu: new(sync.Mutex),
	}
}

// PutMidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) PutMidpointPrice(midPointPrice MidpointPrice) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.prices = midPointPrice

	return nil
}

// MidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrice() (MidpointPrice, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	return storer.prices, nil
}
