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
func (cache *Cache) Open(entry Entry) error {
	cache.ordersMu.Lock()
	defer cache.ordersMu.Unlock()

	// Check if the order has already been opened
	if _, ok := cache.orders[string(entry.Order.ID)]; ok {
		return fmt.Errorf("can't open already existing order")
	} else if _, ok := cache.cancels[string(entry.Order.ID)]; ok {
		return fmt.Errorf("can't open already existing order")
	}

	entry.Status = order.Open
	cache.storeOrderMessage(entry)

	return nil
}

// Match will change the order status to 'unconfirmed' if the order
// is valid and it's status is 'open'.
func (cache *Cache) Match(entry Entry) error {
	cache.ordersMu.Lock()
	defer cache.ordersMu.Unlock()

	// Check if the order has been cancelled by the trader.
	if _, ok := cache.cancels[string(entry.Order.ID)]; ok {
		return fmt.Errorf("trying to match canceled order")
	}

	if cachedOrder, ok := cache.orders[string(entry.Order.ID)]; !ok || cachedOrder.Status != order.Open {
		return fmt.Errorf("can only match orders with status Open")
	}

	entry.Status = order.Unconfirmed
	cache.storeOrderMessage(entry)

	return nil

}

// Confirm will change the order status to 'confirmed' if the order
// is valid and it's status is 'unconfirmed'.
func (cache *Cache) Confirm(entry Entry) error {
	cache.ordersMu.Lock()
	defer cache.ordersMu.Unlock()

	// Check if the order has been cancelled by the trader.
	if _, ok := cache.cancels[string(entry.Order.ID)]; ok {
		return fmt.Errorf("trying to confirm canceled order")
	}

	if cachedOrder, ok := cache.orders[string(entry.Order.ID)]; !ok || cachedOrder.Status != order.Unconfirmed {
		return fmt.Errorf("can only confirm orders with status Unconfirmed")
	}

	// Check if the order has been cancelled by the trader.
	entry.Status = order.Confirmed
	cache.storeOrderMessage(entry)

	return nil

}

// Release will change the order status to 'open' if the order
// is valid and it's status is 'unconfirmed'.
func (cache *Cache) Release(entry Entry) error {
	cache.ordersMu.Lock()
	cache.cancelMu.RLock()
	defer cache.ordersMu.Unlock()
	defer cache.cancelMu.RUnlock()

	// Check if the order has been cancelled by the trader.
	if _, ok := cache.cancels[string(entry.Order.ID)]; ok {
		return fmt.Errorf("trying to release canceled order")
	}

	if cachedOrder, ok := cache.orders[string(entry.Order.ID)]; !ok || cachedOrder.Status != order.Unconfirmed {
		return fmt.Errorf("can only release orders with status Unconfirmed, order had status %v", cachedOrder.Status)
	}

	// Check if the order has been cancelled by the trader.
	entry.Status = order.Open
	cache.storeOrderMessage(entry)

	return nil
}

// Settle will change the order status to 'settled' if the order
// is valid and it's status is 'confirmed'.
func (cache *Cache) Settle(entry Entry) error {
	cache.ordersMu.Lock()
	defer cache.ordersMu.Unlock()

	if cachedOrder, ok := cache.orders[string(entry.Order.ID)]; !ok || cachedOrder.Status != order.Confirmed {
		return fmt.Errorf("can only settle orders with status Confirmed")
	}

	entry.Status = order.Settled
	cache.storeOrderMessage(entry)

	return nil
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
	}

	// msg.Status is Open or Unconfirmed
	cache.cancels[string(id)] = struct{}{}
	delete(cache.orders, string(id))

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
	cache.orders[string(entry.Order.ID)] = entry
}
