package leveldb

import (
	"encoding/json"

	"github.com/republicprotocol/republic-go/ome"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// ComputationIterator implements the ome.ComputationIterator
// interface using a LevelDB iterator. It is not safe for concurrent use.
type ComputationIterator struct {
	inner iterator.Iterator
	next  bool
}

func newComputationIterator(iter iterator.Iterator) *ComputationIterator {
	return &ComputationIterator{
		inner: iter,
		next:  false,
	}
}

// Next implements the ome.ComputationIterator interface.
func (iter *ComputationIterator) Next() bool {
	iter.next = iter.inner.Next()
	return iter.next
}

// Cursor implements the ome.ComputationIterator interface.
func (iter *ComputationIterator) Cursor() (ome.Computation, error) {
	com := ome.Computation{}
	if !iter.next {
		return com, ome.ErrCursorOutOfRange
	}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &com); err != nil {
		return com, err
	}
	return com, iter.inner.Error()
}

// Collect implements the ome.ComputationIterator interface.
func (iter *ComputationIterator) Collect() ([]ome.Computation, error) {
	coms := []ome.Computation{}
	for iter.Next() {
		com, err := iter.Cursor()
		if err != nil {
			return coms, err
		}
		coms = append(coms, com)
	}
	return coms, iter.inner.Error()
}

// Release implements the ome.ComputationIterator interface.
func (iter *ComputationIterator) Release() {
	iter.inner.Release()
}
