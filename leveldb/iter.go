package leveldb

import (
	"encoding/json"

	"github.com/republicprotocol/republic-go/order"

	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// ChangeIterator implements the orderbook.ChangeIterator interface using a
// LevelDB iterator. It is not safe for concurrent use.
type ChangeIterator struct {
	inner iterator.Iterator
	next  bool
}

func newChangeIterator(iter iterator.Iterator) *ChangeIterator {
	return &ChangeIterator{
		inner: iter,
		next:  false,
	}
}

// Next implements the orderbook.ChangeIterator interface.
func (iter *ChangeIterator) Next() bool {
	iter.next = iter.inner.Next()
	return iter.next
}

// Cursor implements the orderbook.ChangeIterator interface.
func (iter *ChangeIterator) Cursor() (orderbook.Change, error) {
	change := orderbook.Change{}
	if !iter.next {
		return change, orderbook.ErrCursorOutOfRange
	}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &change); err != nil {
		return change, err
	}
	return change, iter.inner.Error()
}

// Collect implements the orderbook.ChangeIterator interface.
func (iter *ChangeIterator) Collect() ([]orderbook.Change, error) {
	changes := []orderbook.Change{}
	for iter.Next() {
		change, err := iter.Cursor()
		if err != nil {
			return changes, err
		}
		changes = append(changes, change)
	}
	return changes, iter.inner.Error()
}

// Release implements the orderbook.ChangeIterator interface.
func (iter *ChangeIterator) Release() {
	iter.inner.Release()
}

// OrderFragmentIterator implements the orderbook.OrderFragmentIterator
// interface using a LevelDB iterator.
type OrderFragmentIterator struct {
	inner iterator.Iterator
	next  bool
}

func newOrderFragmentIterator(iter iterator.Iterator) *OrderFragmentIterator {
	return &OrderFragmentIterator{
		inner: iter,
		next:  false,
	}
}

// Next implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderFragmentIterator) Next() bool {
	iter.next = iter.inner.Next()
	return iter.next
}

// Cursor implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderFragmentIterator) Cursor() (order.Fragment, error) {
	fragment := order.Fragment{}
	if !iter.next {
		return fragment, orderbook.ErrCursorOutOfRange
	}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &fragment); err != nil {
		return fragment, err
	}
	return fragment, iter.inner.Error()
}

// Collect implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderFragmentIterator) Collect() ([]order.Fragment, error) {
	fragments := []order.Fragment{}
	for iter.Next() {
		fragment, err := iter.Cursor()
		if err != nil {
			return fragments, err
		}
		fragments = append(fragments, fragment)
	}
	return fragments, iter.inner.Error()
}

// Release implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderFragmentIterator) Release() {
	iter.inner.Release()
}
