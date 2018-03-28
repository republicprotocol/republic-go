package orderbook

import (
	"github.com/republicprotocol/republic-go/dispatch"
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

// The broadcaster is the subject in the observer design pattern
type Broadcaster interface {
	Subscribe(id string, listener chan *rpc.SyncBlock) error
	Unsubscribe(id string)
}

// An OrderBook is responsible for store the historical orders both in
// cache and in disk. It also streams the newly received orders to its
// subscriber.
type OrderBook struct {
	orderBookCache    OrderBookCache
	orderBookDB       OrderBookDB
	orderBookStreamer OrderBookNotifier
}

// NewOrderBook creates a new OrderBook with the given connection limits.
func NewOrderBook(maxConnections int) *OrderBook {
	return &OrderBook{
		orderBookCache:    NewOrderBookCache(),
		orderBookDB:       NewOrderBookDB(),
		orderBookStreamer: NewOrderBookStreamer(maxConnections),
	}
}

// Sync will stream the order history to the message queue provided.
func (orderBook OrderBook) Sync(queue dispatch.MessageQueue) error {
	blocks := orderBook.orderBookCache.Blocks()
	for _, block := range blocks {
		orderBook.orderBookStreamer.Send(block)
	}
	return nil
}

// Open is called when we first receive the order fragment.
func (orderBook OrderBook) Open(ord *order.Order) {
	orderBook.orderBookCache.Open(ord)
	orderBook.orderBookDB.Open(ord)
	orderBook.orderBookStreamer.Open(ord)
}

// Match is called when we discover a match for the order.
func (orderBook OrderBook) Match(ord *order.Order) {
	orderBook.orderBookCache.Match(ord)
	orderBook.orderBookDB.Match(ord)
	orderBook.orderBookStreamer.Match(ord)
}

// Confirm is called when the order has been confirmed by the hyperdrive.
func (orderBook OrderBook) Confirm(ord *order.Order) {
	orderBook.orderBookCache.Confirm(ord)
	orderBook.orderBookDB.Confirm(ord)
	orderBook.orderBookStreamer.Confirm(ord)
}

// Release is called when the order has been denied by the hyperdrive.
func (orderBook OrderBook) Release(ord *order.Order) {
	orderBook.orderBookCache.Release(ord)
	orderBook.orderBookDB.Release(ord)
	orderBook.orderBookStreamer.Release(ord)
}

// Release is called when the order is settled.
func (orderBook OrderBook) Settle(ord *order.Order) {
	orderBook.orderBookCache.Settle(ord)
	orderBook.orderBookDB.Settle(ord)
	orderBook.orderBookStreamer.Settle(ord)
}

