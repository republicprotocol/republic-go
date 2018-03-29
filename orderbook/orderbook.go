package orderbook

import (
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
)

type OrderBookSyncer interface {
	Open(ord order.Order) error
	Match(ord order.Order) error
	Confirm(ord order.Order) error
	Release(ord order.Order) error
	Settle(ord order.Order) error
}

// The broadcaster is the subject in the observer design pattern
type Broadcaster interface {
	Subscribe(id string, queue dispatch.MessageQueue) error
	Unsubscribe(id string)
}

// An OrderBook is responsible for store the historical orders both in
// cache and in disk. It also streams the newly received orders to its
// subscriber.
type OrderBook struct {
	orderBookCache OrderBookCache
	orderBookDB    OrderBookDB
	splitter       dispatch.Splitter
}

// NewOrderBook creates a new OrderBook with the given logger and splitter
func NewOrderBook(maxConnections int) *OrderBook {
	return &OrderBook{
		orderBookCache: NewOrderBookCache(),
		orderBookDB:    NewOrderBookDB(),
		splitter:       dispatch.NewSplitter(maxConnections),
	}
}

// Sync will stream the order history to the message queue provided.
func (orderBook OrderBook) SyncHistory(queue dispatch.MessageQueue) error {
	blocks := orderBook.orderBookCache.Blocks()
	for _, block := range blocks {
		orderBook.splitter.Send(block)
	}
	return nil
}

// Subscribe will start listening to the orderbook for updates.
func (orderBook OrderBook) Subscribe(id string, queue dispatch.MessageQueue) error {
	return orderBook.splitter.RunMessageQueue(id, queue)
}

// Unsubscribe will stop listening to the orderbook for updates
func (orderBook OrderBook) Unsubscribe(id string){
	orderBook.splitter.ShutdownMessageQueue(id)
}

// Open is called when we first receive the order fragment.
func (orderBook OrderBook) Open(ord order.Order) error{
	orderBook.orderBookCache.Open(ord)
	orderBook.orderBookDB.Open(ord)
	message := NewMessage(nil, ord, order.Open)
	return orderBook.splitter.Send(message)
}

// Match is called when we discover a match for the order.
func (orderBook OrderBook) Match(ord order.Order) error{
	orderBook.orderBookCache.Match(ord)
	orderBook.orderBookDB.Match(ord)
	message := NewMessage(nil, ord, order.Unconfirmed)
	return orderBook.splitter.Send(message)
}

// Confirm is called when the order has been confirmed by the hyperdrive.
func (orderBook OrderBook) Confirm(ord order.Order) error{
	orderBook.orderBookCache.Confirm(ord)
	orderBook.orderBookDB.Confirm(ord)
	message := NewMessage(nil, ord, order.Confirmed)
	return orderBook.splitter.Send(message)
}

// Release is called when the order has been denied by the hyperdrive.
func (orderBook OrderBook) Release(ord order.Order) error{
	orderBook.orderBookCache.Release(ord)
	orderBook.orderBookDB.Release(ord)
	message := NewMessage(nil, ord, order.Open)
	return orderBook.splitter.Send(message)
}

// Release is called when the order is settled.
func (orderBook OrderBook) Settle(ord order.Order)error {
	orderBook.orderBookCache.Settle(ord)
	orderBook.orderBookDB.Settle(ord)
	message := NewMessage(nil, ord, order.Settled)
	return orderBook.splitter.Send(message)
}
