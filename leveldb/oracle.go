package leveldb

import (
	"sync"

	"github.com/republicprotocol/republic-go/oracle"
	"github.com/republicprotocol/republic-go/order"
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

// MidpointPrices implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrices() (oracle.MidpointPrice, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	return storer.prices, nil
}

// MidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrice(token order.Tokens) (uint64, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	if price, ok := storer.prices.Prices[uint64(token)]; ok {
		return price, nil
	}
	return 0, oracle.ErrMidpointPriceNotFound
}
