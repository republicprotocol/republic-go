package leveldb

import (
	"encoding/json"
	"time"

	"github.com/republicprotocol/republic-go/ome"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// SomerComputationValue is the storage format for computations being store in
// LevelDB. It contains additional timestamping information so that LevelDB can
// provide pruning.
type SomerComputationValue struct {
	Timestamp   time.Time       `json:"timestamp"`
	Computation ome.Computation `json:"computation"`
}

// SomerComputationIterator implements the ome.ComputationIterator using a
// LevelDB iterator.
type SomerComputationIterator struct {
	inner iterator.Iterator
}

func newSomerComputationIterator(iter iterator.Iterator) *SomerComputationIterator {
	return &SomerComputationIterator{
		inner: iter,
	}
}

// Next implements the ome.ComputationIterator interface.
func (iter *SomerComputationIterator) Next() bool {
	return iter.inner.Next()
}

// Cursor implements the ome.ComputationIterator interface.
func (iter *SomerComputationIterator) Cursor() (ome.Computation, error) {
	if !iter.inner.Valid() {
		return ome.Computation{}, ome.ErrCursorOutOfRange
	}
	value := SomerComputationValue{}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &value); err != nil {
		return ome.Computation{}, err
	}
	return value.Computation, iter.inner.Error()
}

// Collect implements the ome.ComputationIterator interface.
func (iter *SomerComputationIterator) Collect() ([]ome.Computation, error) {
	computations := []ome.Computation{}
	for iter.Next() {
		computation, err := iter.Cursor()
		if err != nil {
			return computations, err
		}
		computations = append(computations, computation)
	}
	return computations, iter.inner.Error()
}

// Release implements the ome.ComputationIterator interface.
func (iter *SomerComputationIterator) Release() {
	iter.inner.Release()
}

// SomerComputationTable implements the ome.ComputationStorer interface using
// LevelDB.
type SomerComputationTable struct {
	db *leveldb.DB
}

// NewSomerComputationTable returns a new SomerComputationTable that uses the
// given LevelDB instance to store and load values from the disk.
func NewSomerComputationTable(db *leveldb.DB) *SomerComputationTable {
	return &SomerComputationTable{db: db}
}

// PutComputation implements the ome.ComputationStorer interface.
func (table *SomerComputationTable) PutComputation(computation ome.Computation) error {
	value := SomerComputationValue{
		Timestamp:   time.Now(),
		Computation: computation,
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return table.db.Put(table.key(computation.ID[:]), data, nil)
}

// DeleteComputation implements the ome.ComputationStorer interface.
func (table *SomerComputationTable) DeleteComputation(id ome.ComputationID) error {
	return table.db.Delete(table.key(id[:]), nil)
}

// Computation implements the ome.ComputationStorer interface.
func (table *SomerComputationTable) Computation(id ome.ComputationID) (ome.Computation, error) {
	data, err := table.db.Get(table.key(id[:]), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = ome.ErrComputationNotFound
		}
		return ome.Computation{}, err
	}

	value := SomerComputationValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return ome.Computation{}, err
	}
	return value.Computation, nil
}

// Computations implements the ome.ComputationStorer interface.
func (table *SomerComputationTable) Computations() (ome.Computations, error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(SomerComputationIterBegin), Limit: table.key(SomerComputationIterEnd)}, nil)
	return newSomerComputationIterator(iter), nil
}

// Prune iterates over all computations and deletes those that have expired.
func (table *SomerComputationTable) Prune() (err error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(SomerComputationIterBegin), Limit: table.key(SomerComputationIterEnd)}, nil)
	defer iter.Release()

	now := time.Now()
	for iter.Next() {
		key := iter.Key()
		value := SomerComputationValue{}
		if localErr := json.Unmarshal(iter.Value(), &value); localErr != nil {
			err = localErr
			continue
		}
		if value.Timestamp.Add(SomerComputationExpiry).Before(now) {
			if localErr := table.db.Delete(key, nil); localErr != nil {
				err = localErr
			}
		}
	}
	return err
}

func (table *SomerComputationTable) key(k []byte) []byte {
	return append(append(SomerComputationTableBegin, k...), SomerComputationTablePadding...)
}
