package orderbook

import (
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
)

type OrderBookSyncer interface {
	Open(message Message) error
	Match(message Message) error
	Confirm(message Message) error
	Release(message Message) error
	Settle(message Message) error
}

// Broadcaster is the subject in the observer design pattern
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
		orderBookDB:    OrderBookDB{},
		splitter:       dispatch.NewSplitter(maxConnections),
	}
}

// Subscribe will start listening to the orderbook for updates.
func (orderBook OrderBook) Subscribe(id string, queue dispatch.MessageQueue) error {
	var err error
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()

		err = orderBook.splitter.RunMessageQueue(id, queue)
	}()

	blocks := orderBook.orderBookCache.Blocks()
	for _, block := range blocks {
		err := queue.Send(block)
		if err != nil {
			return err
		}
	}

	wg.Wait()
	return err
}

// Unsubscribe will stop listening to the orderbook for updates
func (orderBook OrderBook) Unsubscribe(id string) {
	orderBook.splitter.ShutdownMessageQueue(id)
}

// Open is called when we first receive the order fragment.
func (orderBook OrderBook) Open(message *Message) error {
	orderBook.orderBookCache.Open(message)
	orderBook.orderBookDB.Open(message)
	return orderBook.splitter.Send(message)
}

// Match is called when we discover a match for the order.
func (orderBook OrderBook) Match(message *Message) error {
	orderBook.orderBookCache.Match(message)
	orderBook.orderBookDB.Match(message)
	return orderBook.splitter.Send(message)
}

// Confirm is called when the order has been confirmed by the hyperdrive.
func (orderBook OrderBook) Confirm(message *Message) error {
	orderBook.orderBookCache.Confirm(message)
	orderBook.orderBookDB.Confirm(message)
	return orderBook.splitter.Send(message)
}

// Release is called when the order has been denied by the hyperdrive.
func (orderBook OrderBook) Release(message *Message) error {
	orderBook.orderBookCache.Release(message)
	orderBook.orderBookDB.Release(message)
	return orderBook.splitter.Send(message)
}

// Settle is called when the order is settled.
func (orderBook OrderBook) Settle(message *Message) error {
	orderBook.orderBookCache.Settle(message)
	orderBook.orderBookDB.Settle(message)
	return orderBook.splitter.Send(message)
}
