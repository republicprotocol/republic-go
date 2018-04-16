package orderbook

import (
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

// Cache is responsible for store the orders and their
// status in the cache.
type Cache struct {
	ordersMu *sync.RWMutex
	orders   map[string]Entry

	cancelMu *sync.RWMutex
	cancels  map[string]struct{}
}

// NewCache creates a new Cache
func NewCache() Cache {
	return Cache{
		ordersMu: new(sync.RWMutex),
		orders:   map[string]Entry{},

		cancelMu: new(sync.RWMutex),
		cancels:  map[string]struct{}{},
	}
}

// Open is called when we first receive the order fragment.
// It will create the order record and make its status 'open'.
func (cache *Cache) Open(entry Entry) {
	cache.storeOrderMessage(entry)
}

// Match will change the order status to 'unconfirmed' if the order
// is valid and it's status is 'open'.
func (cache *Cache) Match(entry Entry) {
	cache.storeOrderMessage(entry)
}

// Confirm will change the order status to 'confirmed' if the order
// is valid and it's status is 'unconfirmed'.
func (cache *Cache) Confirm(entry Entry) {
	cache.storeOrderMessage(entry)
}

// Release will change the order status to 'open' if the order
// is valid and it's status is 'unconfirmed'.
func (cache *Cache) Release(entry Entry) {
	cache.ordersMu.Lock()
	cache.cancelMu.RLock()
	defer cache.ordersMu.Unlock()
	defer cache.cancelMu.RUnlock()

	// Check if the order has been cancelled by the trader.
	if _, ok := cache.cancels[string(entry.Order.ID)]; ok {
		delete(cache.orders, string(entry.Order.ID))
	} else {
		cache.storeOrderMessage(entry)
	}
}

// Settle will change the order status to 'settled' if the order
// is valid and it's status is 'confirmed'.
func (cache *Cache) Settle(entry Entry) {
	cache.storeOrderMessage(entry)
}

// Cancel is called when trader wants to cancel the order.
// Order can only be cancelled when its status is unconfirmed or open.
func (cache *Cache) Cancel(id order.ID) error {
	cache.ordersMu.RLock()
	cache.cancelMu.Lock()
	defer cache.ordersMu.RUnlock()
	defer cache.cancelMu.Unlock()

	msg, ok := cache.orders[string(id)]
	if !ok {
		return fmt.Errorf("order does not exist")
	}
	if msg.Status > order.Unconfirmed {
		return fmt.Errorf("too late too cancel the order")
	} else if msg.Status == order.Unconfirmed {
		cache.cancels[string(id)] = struct{}{}
	} else if msg.Status == order.Open {
		delete(cache.orders, string(id))
	}

	return nil
}

// Blocks will gather all the orders records and returns them in
// the format of orderbook.Message
func (cache *Cache) Blocks() []Entry {
	cache.ordersMu.RLock()
	defer cache.ordersMu.RUnlock()

	blocks := make([]Entry, len(cache.orders))
	i := 0
	for _, ord := range cache.orders {
		blocks[i] = ord
		i++
	}

	return blocks
}

func (cache *Cache) storeOrderMessage(entry Entry) {
	cache.ordersMu.Lock()
	defer cache.ordersMu.Unlock()

	// Store the order entry if we haven't seen the order before.
	if _, ok := cache.orders[string(entry.Order.ID)]; !ok {
		cache.orders[string(entry.Order.ID)] = entry
		return
	}

	// Merge order by the priority of the order status
	if entry.Status < cache.orders[string(entry.Order.ID)].Status {
		cache.orders[string(entry.Order.ID)] = entry
		return
	}
}
