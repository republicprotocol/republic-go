package orderbook

import (
	"errors"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
)

// ErrWriteToClosedOrderbook is returned when an attempt to updated the
// Orderbook is made after a call to Orderbook.Close.
var ErrWriteToClosedOrderbook = errors.New("write to closed orderbook")

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

// An Orderbook is responsible for store the historical orders both in cache
// and in disk. It also streams the newly received orders to its subscriber.
type Orderbook struct {
	cache    Syncer
	database Syncer

	broadcaster       *dispatch.Broadcaster
	broadcasterChDone chan struct{}
	broadcasterCh     chan interface{}
}

// NewOrderbook creates a new Orderbook with the given logger and splitter
func NewOrderbook() Orderbook {
	cache := NewCache()
	database := Database{}

	broadcaster := dispatch.NewBroadcaster()
	broadcasterChDone := make(chan struct{})
	broadcasterCh := make(chan interface{})
	go broadcaster.Broadcast(broadcasterChDone, broadcasterCh)

	return Orderbook{
		cache:    &cache,
		database: &database,

		broadcaster:       broadcaster,
		broadcasterCh:     broadcasterCh,
		broadcasterChDone: broadcasterChDone,
	}
}

// Close the Orderbook. All listeners will eventually be closed and no more
// listeners will be accepted.
func (orderbook *Orderbook) Close() {
	orderbook.broadcaster.Close()
	close(orderbook.broadcasterChDone)
}

// Listen to the orderbook for updates. Calls to Orderbook.Listen are
// non-blocking, and the background worker is terminated when the done
// channel is closed. A read-only channel of entries is returned, and will be
// closed when no more data will be written to it.
func (orderbook *Orderbook) Listen(done <-chan struct{}) <-chan Entry {
	listener := orderbook.broadcaster.Listen(done)
	subscriber := make(chan Entry)

	go func() {
		defer close(subscriber)
		dispatch.CoBegin(func() {
			for {
				select {
				case <-done:
					return
				case <-orderbook.broadcasterChDone:
					return
				case msg, ok := <-listener:
					if !ok {
						return
					}
					if msg, ok := msg.(Entry); ok {
						select {
						case <-done:
							return
						case <-orderbook.broadcasterChDone:
							return
						case subscriber <- msg:
						}
					}
				}
			}
		}, func() {
			blocks := orderbook.cache.Blocks()
			for _, block := range blocks {
				select {
				case <-done:
					return
				case <-orderbook.broadcasterChDone:
					return
				case subscriber <- block:
				}
			}
		})
	}()

	return subscriber
}

// Open is called when we first receive the order fragment.
func (orderbook *Orderbook) Open(entry Entry) error {
	if err := orderbook.cache.Open(entry); err != nil {
		return err
	}
	// orderbook.database.Open(entry)

	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- entry:
		return nil
	}
}

// Match is called when we discover a match for the order.
func (orderbook *Orderbook) Match(entry Entry) error {
	if err := orderbook.cache.Match(entry); err != nil {
		return err
	}
	// orderbook.database.Match(entry)

	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- entry:
		return nil
	}
}

// Confirm is called when the order has been confirmed by the hyperdrive.
func (orderbook *Orderbook) Confirm(entry Entry) error {
	if err := orderbook.cache.Confirm(entry); err != nil {
		return err
	}
	// orderbook.database.Confirm(entry)

	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- entry:
		return nil
	}
}

// Release is called when the order has been denied by the hyperdrive.
func (orderbook *Orderbook) Release(entry Entry) error {
	if err := orderbook.cache.Release(entry); err != nil {
		return err
	}
	// orderbook.database.Release(entry)

	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- entry:
		return nil
	}
}

// Settle is called when the order is settled.
func (orderbook *Orderbook) Settle(entry Entry) error {
	if err := orderbook.cache.Settle(entry); err != nil {
		return err
	}
	// orderbook.database.Settle(entry)

	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- entry:
		return nil
	}
}

// Cancel is called when the order is canceled.
func (orderbook *Orderbook) Cancel(id order.ID) error {
	if err := orderbook.cache.Cancel(id); err != nil {
		return err
	}
	// err = orderbook.database.Cancel(id)
	// if err != nil {
	// 	return err
	// }

	entry := NewEntry(order.Order{ID: id}, order.Canceled)

	select {
	case <-orderbook.broadcasterChDone:
		return ErrWriteToClosedOrderbook
	case orderbook.broadcasterCh <- entry:
		return nil
	}
}

// Blocks will gather all the order records and returns them in
// the format of orderbook.Entry
func (orderbook *Orderbook) Blocks() []Entry {
	return orderbook.cache.Blocks()
}

// Order retrieves information regarding an order.
func (orderbook *Orderbook) Order(id order.ID) Entry {
	return orderbook.cache.Order(id)
}
