package leveldb

import (
	"encoding/json"

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
func (iter *ChangeIterator) Cursor() (Change, error) {
	data := iter.Value()
	change := orderbook.Change{}
	err := json.Unmarshal(data, &change)
	return change, err
}

// Collect implements the orderbook.ChangeIterator interface.
func (iter *ChangeIterator) Collect() ([]Change, error) {
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
