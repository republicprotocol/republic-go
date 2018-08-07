package leveldb

import (
	"sync"

	"github.com/republicprotocol/republic-go/oracle"
)

// midpointPriceStorer implements MidpointPriceStorer interface with an
// in-memory map implementation. Data will be lost every time it restarts.
type midpointPriceStorer struct {
	mu     *sync.Mutex
	prices oracle.MidpointPrice
}

// NewMidpointPriceStorer returns a new MidpointPriceStorer.
func NewMidpointPriceStorer() oracle.MidpointPriceStorer {
	return &midpointPriceStorer{
		mu: new(sync.Mutex),
	}
}

// PutMidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) PutMidpointPrice(midPointPrice oracle.MidpointPrice) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.prices = midPointPrice

	return nil
}

// MidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrice() (oracle.MidpointPrice, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	return storer.prices, nil
}
