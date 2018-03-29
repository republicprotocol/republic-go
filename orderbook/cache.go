package orderbook

import (
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

// An OrderBookCache is responsible for store the orders and their
// status in the cache.
type OrderBookCache struct {
	mu *sync.RWMutex

	orders map[string]order.Order
	status map[string]order.Status
}

// NewOrderBookCache creates a new OrderBookCache
func NewOrderBookCache() OrderBookCache {
	return OrderBookCache{
		mu:     new(sync.RWMutex),
		orders: map[string]order.Order{},
		status: map[string]order.Status{},
	}
}

// Open is called when we first receive the order fragment.
// It will create the order record and make its status 'open'.
func (orderBookCache *OrderBookCache) Open(ord order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	if _, ok := orderBookCache.status[string(ord.ID)]; !ok {
		orderBookCache.orders[string(ord.ID)] = ord
		orderBookCache.status[string(ord.ID)] = order.Open
	}
}

// Match will change the order status to 'unconfirmed' if the order
// is valid and it's status is 'open'.
func (orderBookCache *OrderBookCache) Match(ord order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	if status, ok := orderBookCache.status[string(ord.ID)]; ok && status == order.Open {
		orderBookCache.orders[string(ord.ID)] = ord
		orderBookCache.status[string(ord.ID)] = order.Unconfirmed
	}
}

// Confirm will change the order status to 'confirmed' if the order
// is valid and it's status is 'unconfirmed'.
func (orderBookCache *OrderBookCache) Confirm(ord order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	if status, ok := orderBookCache.status[string(ord.ID)]; ok && status == order.Unconfirmed {
		orderBookCache.orders[string(ord.ID)] = ord
		orderBookCache.status[string(ord.ID)] = order.Confirmed
	}
}

// Release will change the order status to 'open' if the order
// is valid and it's status is 'unconfirmed'.
func (orderBookCache *OrderBookCache) Release(ord order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	if status, ok := orderBookCache.status[string(ord.ID)]; ok && status == order.Unconfirmed {
		orderBookCache.orders[string(ord.ID)] = ord
		orderBookCache.status[string(ord.ID)] = order.Open
	}
}

// Settle will change the order status to 'settled' if the order
// is valid and it's status is 'confirmed'.
func (orderBookCache *OrderBookCache) Settle(ord order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	if status, ok := orderBookCache.status[string(ord.ID)]; ok && status == order.Confirmed {
		orderBookCache.orders[string(ord.ID)] = ord
		orderBookCache.status[string(ord.ID)] = order.Settled
	}
}

// Blocks will gather all the orders records and returns them in
// the format of orderbook.Message
func (orderBookCache *OrderBookCache) Blocks() []Message {
	orderBookCache.mu.RLock()
	defer orderBookCache.mu.RUnlock()

	blocks := make([]Message, 0)
	for _, ord := range orderBookCache.orders {
		status, ok := orderBookCache.status[string(ord.ID)]
		if ok {
			block := NewMessage(ord, status, nil)
			blocks = append(blocks, block)
		}
	}

	return blocks
}


