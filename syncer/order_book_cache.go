package syncer

import (
	"sync"

	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

type OrderBookCache struct {
	mu *sync.RWMutex

	orders map[string]*order.Order
	status map[string]order.Status
}

func NewOrderBookCache() OrderBookCache {
	return OrderBookCache{
		mu:     new(sync.RWMutex),
		orders: map[string]*order.Order{},
		status: map[string]order.Status{},
	}
}

func (orderBookCache *OrderBookCache) Open(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if _, ok := orderBookCache.status[string(ord.ID)]; !ok {
		orderBookCache.status[string(ord.ID)] = order.Open
	}
	// todo : do we need to something here ?
}

func (orderBookCache *OrderBookCache) Match(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Open {
		orderBookCache.status[string(ord.ID)] = order.Unconfirmed
	}
}

func (orderBookCache *OrderBookCache) Confirm(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Unconfirmed {
		orderBookCache.status[string(ord.ID)] = order.Confirmed
	}
}

func (orderBookCache *OrderBookCache) Release(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Unconfirmed {
		orderBookCache.status[string(ord.ID)] = order.Open
	}
}

func (orderBookCache *OrderBookCache) Settle(ord *order.Order) {
	orderBookCache.mu.Lock()
	defer orderBookCache.mu.Unlock()

	orderBookCache.orders[string(ord.ID)] = ord
	if status := orderBookCache.status[string(ord.ID)]; status == order.Confirmed {
		orderBookCache.status[string(ord.ID)] = order.Settled
	}
}

func (orderBookCache *OrderBookCache) Blocks() []*rpc.SyncBlock {
	orderBookCache.mu.RLock()
	defer orderBookCache.mu.RUnlock()

	blocks := make ([]*rpc.SyncBlock, 0)
	for _ , ord := range orderBookCache.orders {
		status, ok := orderBookCache.status[string(ord.ID)]
		if ok {
			block := orderToSyncBlock(ord , status)
			blocks = append(blocks, block)
		}
	}

	return blocks
}

