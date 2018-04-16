package orderbook

import (
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
)

type Syncer interface {
	Open(entry Entry) error
	Match(entry Entry) error
	Confirm(entry Entry) error
	Release(entry Entry) error
	Settle(entry Entry) error
	Cancel(id order.ID) error
}

// Broadcaster is the subject in the observer design pattern
type Broadcaster interface {
	Subscribe(id string, queue dispatch.MessageQueue) error
	Unsubscribe(id string)
}

// An Orderbook is responsible for store the historical orders both in cache
// and in disk. It also streams the newly received orders to its subscriber.
type Orderbook struct {
	cache    Cache
	database Database
	splitter dispatch.Splitter
}

// NewOrderBook creates a new OrderBook with the given logger and splitter
func NewOrderBook(maxConnections int) *OrderBook {
	return &OrderBook{
		cache:    NewOrderBookCache(),
		database: Database{},
		splitter: dispatch.NewSplitter(maxConnections),
	}
}

// Subscribe will start listening to the orderbook for updates.
func (orderBook OrderBook) Subscribe(id string, queue dispatch.MessageQueue) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err = orderBook.splitter.RunMessageQueue(id, queue)
	}()

	blocks := orderBook.cache.Blocks()
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
func (orderBook OrderBook) Open(entry Entry) error {
	orderBook.cache.Open(entry)
	orderBook.database.Open(entry)
	return orderBook.splitter.Send(entry)
}

// Match is called when we discover a match for the order.
func (orderBook OrderBook) Match(entry Entry) error {
	orderBook.cache.Match(entry)
	orderBook.database.Match(entry)
	return orderBook.splitter.Send(entry)
}

// Confirm is called when the order has been confirmed by the hyperdrive.
func (orderBook OrderBook) Confirm(entry Entry) error {
	orderBook.cache.Confirm(entry)
	orderBook.database.Confirm(entry)
	return orderBook.splitter.Send(entry)
}

// Release is called when the order has been denied by the hyperdrive.
func (orderBook OrderBook) Release(entry Entry) error {
	orderBook.cache.Release(entry)
	orderBook.database.Release(entry)
	return orderBook.splitter.Send(entry)
}

// Settle is called when the order is settled.
func (orderBook OrderBook) Settle(entry Entry) error {
	orderBook.cache.Settle(entry)
	orderBook.database.Settle(entry)
	return orderBook.splitter.Send(entry)
}

// Cancel is called when the order is canceled.
func (orderBook OrderBook) Cancel(id order.ID) error {
	err := orderBook.cache.Cancel(id)
	if err != nil {
		return err
	}
	err = orderBook.database.Cancel(id)
	if err != nil {
		return err
	}

	return orderBook.splitter.Send(NewMessage(order.Order{ID: id}, order.Canceled, [32]byte{}))
}

// Order retrieves information regarding an order.
func (orderBook OrderBook) Order(id order.ID) Entry {
	return orderBook.cache.orders[string(id)]
}
