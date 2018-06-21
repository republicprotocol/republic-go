package leveldb

import (
	"encoding/json"

	"github.com/republicprotocol/republic-go/order"

	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type ChangeIterator struct {
	inner iterator.Iterator
}

// Next implements the orderbook.ChangeIterator interface.
func (iter *ChangeIterator) Next() bool {
	return iter.inner.Next()
}

// Cursor implements the orderbook.ChangeIterator interface.
func (iter *ChangeIterator) Cursor() (orderbook.Change, error) {
	data := iter.Value()
	change := orderbook.Change{}
	err := json.Unmarshal(data, &change)
	return change, err
}

// Collect implements the orderbook.ChangeIterator interface.
func (iter *ChangeIterator) Collect() ([]orderbook.Change, error) {
	changes := []orderbook.Change{}
	iter := store.db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		data := iter.Value()
		change := orderbook.Change{}
		if err := json.Unmarshal(data, &change); err != nil {
			return changes, err
		}
		changes = append(changes, change)
	}
	return changes, iter.Error()
}

type OrderFragmentIterator struct {
	inner iterator.Iterator
}

// Next implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderFragmentIterator) Next() bool {
	return iter.inner.Next()
}

// Cursor implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderFragmentIterator) Cursor() (order.Fragment, error) {
	data := iter.Value()
	fragment := order.Fragment{}
	err := json.Unmarshal(data, &fragment)
	return change, err
}

// Collect implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderFragmentIterator) Collect() ([]order.Fragment, error) {
	fragments := []order.Fragment{}
	iter := store.db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		data := iter.Value()
		change := order.Fragment{}
		if err := json.Unmarshal(data, &change); err != nil {
			return fragments, err
		}
		fragments = append(fragments, change)
	}
	return fragments, iter.Error()
}
