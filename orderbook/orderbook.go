package orderbook

import (
	"reflect"
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
	Subscribe(ch interface{}) error
	Unsubscribe(ch interface{})
}

// An Orderbook is responsible for store the historical orders both in cache
// and in disk. It also streams the newly received orders to its subscriber.
type Orderbook struct {
	cache    Cache
	database Database
	splitter dispatch.Splitter
	splitCh  chan Entry
}

// NewOrderbook creates a new Orderbook with the given logger and splitter
func NewOrderbook(maxConnections int) Orderbook {
	splitter := dispatch.NewSplitter(maxConnections)
	splitCh := make(chan Entry)
	go splitter.Split(splitCh)

	return Orderbook{
		cache:    NewCache(),
		database: Database{},
		splitter: splitter,
		splitCh:  splitCh,
	}
}

// Subscribe will start listening to the orderbook for updates.
func (orderbook Orderbook) Subscribe(ch interface{}) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		orderbook.splitter.Subscribe(ch)
	}()

	blocks := orderbook.cache.Blocks()
	for _, block := range blocks {
		dispatch.SendToInterface(ch, reflect.ValueOf(block))
	}

	wg.Wait()
	return err
}

// Unsubscribe will stop listening to the orderbook for updates
func (orderbook Orderbook) Unsubscribe(ch interface{}) {
	orderbook.splitter.Unsubscribe(ch)
}

// Open is called when we first receive the order fragment.
func (orderbook Orderbook) Open(entry Entry) error {
	orderbook.cache.Open(entry)
	// orderbook.database.Open(entry)
	orderbook.splitCh <- entry
	return nil
}

// Match is called when we discover a match for the order.
func (orderbook Orderbook) Match(entry Entry) error {
	orderbook.cache.Match(entry)
	// orderbook.database.Match(entry)
	orderbook.splitCh <- entry
	return nil
}

// Confirm is called when the order has been confirmed by the hyperdrive.
func (orderbook Orderbook) Confirm(entry Entry) error {
	orderbook.cache.Confirm(entry)
	// orderbook.database.Confirm(entry)
	orderbook.splitCh <- entry
	return nil
}

// Release is called when the order has been denied by the hyperdrive.
func (orderbook Orderbook) Release(entry Entry) error {
	orderbook.cache.Release(entry)
	// orderbook.database.Release(entry)
	orderbook.splitCh <- entry
	return nil
}

// Settle is called when the order is settled.
func (orderbook Orderbook) Settle(entry Entry) error {
	orderbook.cache.Settle(entry)
	// orderbook.database.Settle(entry)
	orderbook.splitCh <- entry
	return nil
}

// Cancel is called when the order is canceled.
func (orderbook Orderbook) Cancel(id order.ID) error {
	err := orderbook.cache.Cancel(id)
	if err != nil {
		return err
	}
	// err = orderbook.database.Cancel(id)
	// if err != nil {
	// 	return err
	// }

	entry := NewEntry(order.Order{ID: id}, order.Canceled, [32]byte{})
	orderbook.splitCh <- entry
	return nil
}

// Order retrieves information regarding an order.
func (orderbook Orderbook) Order(id order.ID) Entry {
	return orderbook.cache.orders[string(id)]
}
