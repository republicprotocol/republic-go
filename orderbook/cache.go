package orderbook

import (
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

// An OrderBookCache is responsible for store the orders and their
// status in the cache.
type OrderBookCache struct {
	ordersMu *sync.RWMutex
	orders   map[string]*Message

	cancelMu *sync.RWMutex
	cancels  map[string]struct{}
}

// NewOrderBookCache creates a new OrderBookCache
func NewOrderBookCache() OrderBookCache {
	return OrderBookCache{
		ordersMu: new(sync.RWMutex),
		orders:   map[string]*Message{},

		cancelMu: new(sync.RWMutex),
		cancels:  map[string]struct{}{},
	}
}

// Open is called when we first receive the order fragment.
// It will create the order record and make its status 'open'.
func (orderBookCache *OrderBookCache) Open(message *Message) {
	orderBookCache.storeOrderMessage(message)
}

// Match will change the order status to 'unconfirmed' if the order
// is valid and it's status is 'open'.
func (orderBookCache *OrderBookCache) Match(message *Message) {
	orderBookCache.storeOrderMessage(message)
}

// Confirm will change the order status to 'confirmed' if the order
// is valid and it's status is 'unconfirmed'.
func (orderBookCache *OrderBookCache) Confirm(message *Message) {
	orderBookCache.storeOrderMessage(message)
}

// Release will change the order status to 'open' if the order
// is valid and it's status is 'unconfirmed'.
func (orderBookCache *OrderBookCache) Release(message *Message) {
	orderBookCache.cancelMu.RLock()
	defer orderBookCache.cancelMu.RUnlock()

	// Check if the order has been cancelled by the trader.
	if _, ok := orderBookCache.cancels[string(message.Ord.ID)] ; ok {
		orderBookCache.ordersMu.RLock()
		defer orderBookCache.ordersMu.RUnlock()
		delete(orderBookCache.orders, string(message.Ord.ID))
	} else {
		orderBookCache.storeOrderMessage(message)
	}
}

// Settle will change the order status to 'settled' if the order
// is valid and it's status is 'confirmed'.
func (orderBookCache *OrderBookCache) Settle(message *Message) {
	orderBookCache.storeOrderMessage(message)
}

// Cancel is called when trader wants to cancel the order.
// Order can only be cancelled when its status is unconfirmed or open.
func (orderBookCache *OrderBookCache) Cancel(id order.ID) error{
	orderBookCache.ordersMu.RLock()
	defer orderBookCache.ordersMu.RUnlock()

	msg ,ok := orderBookCache.orders[string(id)]
	if !ok{
		return fmt.Errorf("order does not exist")
	}

	if msg.Status > order.Unconfirmed {
		return fmt.Errorf("too late too cancel the order")
	} else if msg.Status == order.Unconfirmed {
		orderBookCache.cancelMu.Lock()
		defer orderBookCache.cancelMu.Unlock()

		orderBookCache.cancels[string(id)] = struct{}{}
	} else if msg.Status == order.Open {
		delete(orderBookCache.orders, string(id))
	}

	return nil
}

// Blocks will gather all the orders records and returns them in
// the format of orderbook.Message
func (orderBookCache *OrderBookCache) Blocks() []*Message {
	orderBookCache.ordersMu.RLock()
	defer orderBookCache.ordersMu.RUnlock()

	blocks := make([]*Message, len(orderBookCache.orders))
	i := 0
	for _, ord := range orderBookCache.orders {
		blocks[i] = ord
		i++
	}

	return blocks
}

func (orderBookCache *OrderBookCache) storeOrderMessage(message *Message) {
	orderBookCache.ordersMu.Lock()
	defer orderBookCache.ordersMu.Unlock()

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
