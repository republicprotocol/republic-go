package orderbook

import (
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
	Blocks() []Entry
	Order(id order.ID) Entry
}

// Broadcaster is the subject in the observer design pattern
type Broadcaster interface {
	Subscribe(ch interface{}) error
	Unsubscribe(ch interface{})
	Close()
}

// An Orderbook is responsible for store the historical orders both in cache
// and in disk. It also streams the newly received orders to its subscriber.
type Orderbook struct {
	cache    Syncer
	database Syncer
	splitter dispatch.Splitter
	splitCh  chan Entry
}

// NewOrderbook creates a new Orderbook with the given logger and splitter
func NewOrderbook(maxConnections int) Orderbook {
	splitter := dispatch.NewSplitter(maxConnections)
	splitCh := make(chan Entry)
	go splitter.Split(splitCh)

	cache := NewCache()
	database := Database{}

	return Orderbook{
		cache:    &cache,
		database: &database,
		splitter: splitter,
		splitCh:  splitCh,
	}
}

// Subscribe will start listening to the orderbook for updates. The channel
// must not be closed until after the Unsubscribe method is called.
func (orderbook Orderbook) Subscribe(ch chan Entry) error {
	if err := orderbook.splitter.Subscribe(ch); err != nil {
		return err
	}

	blocks := orderbook.cache.Blocks()
	for _, block := range blocks {
		ch <- block
	}

	return nil
}

// Unsubscribe will stop listening to the orderbook for updates
func (orderbook Orderbook) Unsubscribe(ch interface{}) {
	orderbook.splitter.Unsubscribe(ch)
}

// Close will close the splitCh
func (orderbook Orderbook) Close() {
	close(orderbook.splitCh)
}

// Open is called when we first receive the order fragment.
func (orderbook Orderbook) Open(entry Entry) error {
	if err := orderbook.cache.Open(entry); err != nil {
		return err
	}
	// orderbook.database.Open(entry)
	orderbook.splitCh <- entry
	return nil
}

// Match is called when we discover a match for the order.
func (orderbook Orderbook) Match(entry Entry) error {
	if err := orderbook.cache.Match(entry); err != nil {
		return err
	}
	// orderbook.database.Match(entry)
	orderbook.splitCh <- entry
	return nil
}

// Confirm is called when the order has been confirmed by the hyperdrive.
func (orderbook Orderbook) Confirm(entry Entry) error {
	if err := orderbook.cache.Confirm(entry); err != nil {
		return err
	}
	// orderbook.database.Confirm(entry)
	orderbook.splitCh <- entry
	return nil
}

// Release is called when the order has been denied by the hyperdrive.
func (orderbook Orderbook) Release(entry Entry) error {
	if err := orderbook.cache.Release(entry); err != nil {
		return err
	}
	// orderbook.database.Release(entry)
	orderbook.splitCh <- entry
	return nil
}

// Settle is called when the order is settled.
func (orderbook Orderbook) Settle(entry Entry) error {
	if err := orderbook.cache.Settle(entry); err != nil {
		return err
	}
	// orderbook.database.Settle(entry)
	orderbook.splitCh <- entry
	return nil
}

// Cancel is called when the order is canceled.
func (orderbook Orderbook) Cancel(id order.ID) error {
	if err := orderbook.cache.Cancel(id); err != nil {
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

// Blocks will gather all the order records and returns them in
// the format of orderbook.Entry
func (orderbook Orderbook) Blocks() []Entry {
	return orderbook.cache.Blocks()
}

// Order retrieves information regarding an order.
func (orderbook Orderbook) Order(id order.ID) Entry {
	return orderbook.cache.Order(id)
}
