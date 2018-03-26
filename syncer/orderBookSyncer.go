package syncer

import (
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)



type OrderBookCache struct {
	mu *sync.RWMutex

	orders map[string]*order.Order
	status map[string]order.Status
}

func NewOrderBookCache() *OrderBookCache {
	return &OrderBookCache{
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

func (orderBookCache *OrderBookCache) Orders() []*rpc.SyncBlock {
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

// todo : finish this
type OrderBookDB struct {
}

func NewOrderBookDB() OrderBookDB {
	return OrderBookDB {
	}
}

func (orderBookDB *OrderBookDB) Open(ord *order.Order){
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Match(ord *order.Order){
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Confirm(ord *order.Order){
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Release(ord *order.Order){
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Settle(ord *order.Order){
	// TODO: Implement key/value file store
}

type OrderBookStreamer struct {
	subscribersMu *sync.Mutex
	subscribers    map[string]chan *rpc.SyncBlock

	maxConnections int
}

func NewOrderBookStreamer(maxConnection int) OrderBookStreamer {
	return OrderBookStreamer{
		subscribersMu:             new(sync.Mutex),
		subscribers:    map[string]chan *rpc.SyncBlock{},

		maxConnections: maxConnection,
	}
}

func (orderBookStreamer *OrderBookStreamer) Subscribe(id string, listener chan *rpc.SyncBlock) error {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	if len(orderBookStreamer.subscribers) >= orderBookStreamer.maxConnections {
		return fmt.Errorf("cannot subscribe %s: connection limit reached", id)
	}
	orderBookStreamer.subscribers[id] = listener

	return nil
}

func (orderBookStreamer *OrderBookStreamer) Unsubscribe(id string) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	delete(orderBookStreamer.subscribers, id)
}

func (orderBookStreamer *OrderBookStreamer) Open(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		// Allow back-pressure to cause blocking (this is meant to be mitigated
		// by dropping dead clients, or reducing the maximum connections)
		subscriber <- orderToSyncBlock(ord, order.Open)
	}
}

func (orderBookStreamer *OrderBookStreamer) Match(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- orderToSyncBlock(ord, order.Unconfirmed)
	}
}

func (orderBookStreamer *OrderBookStreamer) Confirm(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- orderToSyncBlock(ord, order.Confirmed)
	}
}

func (orderBookStreamer *OrderBookStreamer) Release(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- orderToSyncBlock(ord, order.Open)
	}
}

func (orderBookStreamer *OrderBookStreamer) Settle(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- orderToSyncBlock(ord, order.Settled)
	}
}

func orderToSyncBlock(ord *order.Order, status order.Status) *rpc.SyncBlock{
	block := new(rpc.SyncBlock)
	block.Timestamp = time.Now().Unix()
	block.Signature = []byte{} // todo : will be finished later
	switch status{
	case order.Open:
		block.OrderBlock = &rpc.SyncBlock_Open{
			Open: rpc.SerializeOrder(ord),
		}
	case order.Unconfirmed:
		block.OrderBlock = &rpc.SyncBlock_Unconfirmed{
			Unconfirmed: rpc.SerializeOrder(ord),
		}
	case order.Confirmed:
		block.OrderBlock = &rpc.SyncBlock_Confirmed{
			Confirmed: rpc.SerializeOrder(ord),
		}
	case order.Settled:
		block.OrderBlock = &rpc.SyncBlock_Settled{
			Settled: rpc.SerializeOrder(ord),
		}
	default:
		return nil
	}

	return block
}