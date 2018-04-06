package orderbook

import (
	"sync"
)

// An OrderBookCache is responsible for store the orders and their
// status in the cache.
type OrderBookCache struct {
	mu *sync.RWMutex

	orders map[string]Message
}

// NewOrderBookCache creates a new OrderBookCache
func NewOrderBookCache() OrderBookCache {
	return OrderBookCache{
		mu:     new(sync.RWMutex),
		orders: map[string]Message{},
	}
}

// Open is called when we first receive the order fragment.
// It will create the order record and make its status 'open'.
func (orderBookCache *OrderBookCache) Open(message Message) {
	orderBookCache.storeOrderMessage(message)
}

// Match will change the order status to 'unconfirmed' if the order
// is valid and it's status is 'open'.
func (orderBookCache *OrderBookCache) Match(message Message) {
	orderBookCache.storeOrderMessage(message)
}

// Confirm will change the order status to 'confirmed' if the order
// is valid and it's status is 'unconfirmed'.
func (orderBookCache *OrderBookCache) Confirm(message Message) {
	orderBookCache.storeOrderMessage(message)
}

// Release will change the order status to 'open' if the order
// is valid and it's status is 'unconfirmed'.
func (orderBookCache *OrderBookCache) Release(message Message) {
	orderBookCache.storeOrderMessage(message)
}

// Settle will change the order status to 'settled' if the order
// is valid and it's status is 'confirmed'.
func (orderBookCache *OrderBookCache) Settle(message Message) {
	orderBookCache.storeOrderMessage(message)
}

// Blocks will gather all the orders records and returns them in
// the format of orderbook.Message
func (orderBookCache *OrderBookCache) Blocks() []Message {
	orderBookCache.mu.RLock()
	defer orderBookCache.mu.RUnlock()

	blocks := make([]Message, len(orderBookCache.orders))
	i := 0
	for _, ord := range orderBookCache.orders {
		blocks[i] = ord
		i++
	}

	return blocks
}

func (orderBookCache *OrderBookCache) storeOrderMessage(message Message) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	// Store the order message if we haven't seen the order before.
	if _, ok := orderBookCache.orders[string(message.Ord.ID)]; !ok {
		orderBookCache.orders[string(message.Ord.ID)] = message
		return
	}

	// Merge order by the priority of the order status
	if message.Status < orderBookCache.orders[string(message.Ord.ID)].Status {
		orderBookCache.orders[string(message.Ord.ID)] = message
		return
	}

}
