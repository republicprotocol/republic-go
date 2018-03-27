package orderbook

import (
	"sync"

	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

type OrderBookSyncer interface {
	Open(ord *order.Order)
	Match(ord *order.Order)
	Confirm(ord *order.Order)
	Release(ord *order.Order)
	Settle(ord *order.Order)
}

// An OrderBookCache is responsible for store the orders and their
// status in the cache.
type OrderBookCache struct {
	mu *sync.RWMutex

	orders map[string]*order.Order
	status map[string]order.Status
}

// NewOrderBookCache creates a new OrderBookCache
func NewOrderBookCache() OrderBookCache {
	return OrderBookCache{
		mu:     new(sync.RWMutex),
		orders: map[string]*order.Order{},
		status: map[string]order.Status{},
	}
}

// Open is called when we first receive the order fragment.
// It will create the order record and make its status 'open'.
func (orderBookCache *OrderBookCache) Open(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if _, ok := orderBookCache.status[string(ord.ID)]; !ok {
		orderBookCache.status[string(ord.ID)] = order.Open
	}
}

// Match will change the order status to 'unconfirmed' if the order
// is valid and it's status is 'open'.
func (orderBookCache *OrderBookCache) Match(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Open {
		orderBookCache.status[string(ord.ID)] = order.Unconfirmed
	}
}

// Confirm will change the order status to 'confirmed' if the order
// is valid and it's status is 'unconfirmed'.
func (orderBookCache *OrderBookCache) Confirm(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status, ok := orderBookCache.status[string(ord.ID)]; ok && status == order.Unconfirmed {
		orderBookCache.status[string(ord.ID)] = order.Confirmed
	}
}

// Release will change the order status to 'open' if the order
// is valid and it's status is 'unconfirmed'.
func (orderBookCache *OrderBookCache) Release(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	if status, ok := orderBookCache.status[string(ord.ID)]; ok && status == order.Unconfirmed {
		orderBookCache.orders[string(ord.ID)] = ord
		orderBookCache.status[string(ord.ID)] = order.Open
	}
}

// Settle will change the order status to 'settled' if the order
// is valid and it's status is 'confirmed'.
func (orderBookCache *OrderBookCache) Settle(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Confirmed {
		orderBookCache.status[string(ord.ID)] = order.Settled
	}
}

// Blocks will gather all the orders records and returns them in
// the format of a rpc.SyncBlock
func (orderBookCache *OrderBookCache) Blocks() []*rpc.SyncBlock {
	orderBookCache.mu.RLock()
	defer orderBookCache.mu.RUnlock()

	blocks := make([]*rpc.SyncBlock, 0)
	for _, ord := range orderBookCache.orders {
		status, ok := orderBookCache.status[string(ord.ID)]
		if ok {
			block := orderToSyncBlock(ord, status)
			blocks = append(blocks, block)
		}
	}

	return blocks
}
