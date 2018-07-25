package leveldb

import (
	"encoding/json"
	"time"

	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"

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
	db     *leveldb.DB
	expiry time.Duration
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
func (table *SomerComputationTable) Computations() (ome.ComputationIterator, error) {
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
		if value.Timestamp.Add(table.expiry).Before(now) {
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

// SomerOrderFragmentValue is the storage format for computations being stored in
// LevelDB. It contains additional timestamping information so that LevelDB can
// provide pruning.
type SomerOrderFragmentValue struct {
	Timestamp     time.Time      `json:"timestamp"`
	OrderFragment order.Fragment `json:"orderFragment"`
	Trader        string         `json:"trader"`
}

// SomerOrderFragmentIterator implements the ome.OrderFragmentIterator using a
// LevelDB iterator.
type SomerOrderFragmentIterator struct {
	inner iterator.Iterator
}

func newSomerOrderFragmentIterator(iter iterator.Iterator) *SomerOrderFragmentIterator {
	return &SomerOrderFragmentIterator{
		inner: iter,
	}
}

// Next implements the ome.OrderFragmentIterator interface.
func (iter *SomerOrderFragmentIterator) Next() bool {
	return iter.inner.Next()
}

// Cursor implements the ome.OrderFragmentIterator interface.
func (iter *SomerOrderFragmentIterator) Cursor() (order.Fragment, string, error) {
	if !iter.inner.Valid() {
		return order.Fragment{}, "", ome.ErrCursorOutOfRange
	}
	value := SomerOrderFragmentValue{}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &value); err != nil {
		return order.Fragment{}, "", err
	}
	return value.OrderFragment, value.Trader, iter.inner.Error()
}

// Collect implements the ome.OrderFragmentIterator interface.
func (iter *SomerOrderFragmentIterator) Collect() ([]order.Fragment, []string, error) {
	orderFragments := []order.Fragment{}
	traders := []string{}
	for iter.Next() {
		orderFragment, trader, err := iter.Cursor()
		if err != nil {
			return orderFragments, traders, err
		}
		orderFragments = append(orderFragments, orderFragment)
		traders = append(traders, trader)
	}
	return orderFragments, traders, iter.inner.Error()
}

// Release implements the ome.OrderFragmentIterator interface.
func (iter *SomerOrderFragmentIterator) Release() {
	iter.inner.Release()
}

// SomerOrderFragmentTable implements the ome.OrderFragmentStorer interface using
// LevelDB.
type SomerOrderFragmentTable struct {
	db     *leveldb.DB
	expiry time.Duration
}

// NewSomerOrderFragmentTable returns a new SomerOrderFragmentTable that uses the
// given LevelDB instance to store and load values from the disk.
func NewSomerOrderFragmentTable(db *leveldb.DB, expiry time.Duration) *SomerOrderFragmentTable {
	return &SomerOrderFragmentTable{
		db:     db,
		expiry: expiry,
	}
}

// PutBuyOrderFragment implements the ome.OrderFragmentStorer interface.
func (table *SomerOrderFragmentTable) PutBuyOrderFragment(epoch registry.Epoch, orderFragment order.Fragment, trader string) error {
	value := SomerOrderFragmentValue{
		Timestamp:     time.Now(),
		OrderFragment: orderFragment,
		Trader:        trader,
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return table.db.Put(table.buyKey(epoch.Hash[:], orderFragment.OrderID[:]), data, nil)
}

// DeleteBuyOrderFragment implements the ome.OrderFragmentStorer interface.
func (table *SomerOrderFragmentTable) DeleteBuyOrderFragment(epoch registry.Epoch, id order.ID) error {
	return table.db.Delete(table.buyKey(epoch.Hash[:], id[:]), nil)
}

// BuyOrderFragment implements the ome.OrderFragmentStorer interface.
func (table *SomerOrderFragmentTable) BuyOrderFragment(epoch registry.Epoch, id order.ID) (order.Fragment, string, error) {
	data, err := table.db.Get(table.buyKey(epoch.Hash[:], id[:]), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = ome.ErrOrderFragmentNotFound
		}
		return order.Fragment{}, "", err
	}

	value := SomerOrderFragmentValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return order.Fragment{}, "", err
	}
	return value.OrderFragment, value.Trader, nil
}

// BuyOrderFragments implements the ome.OrderFragmentStorer interface.
func (table *SomerOrderFragmentTable) BuyOrderFragments(epoch registry.Epoch) (ome.OrderFragmentIterator, error) {
	iter := table.db.NewIterator(&util.Range{Start: table.buyKey(epoch.Hash[:], SomerBuyOrderFragmentIterBegin), Limit: table.buyKey(epoch.Hash[:], SomerBuyOrderFragmentIterEnd)}, nil)
	return newSomerOrderFragmentIterator(iter), nil
}

// PutSellOrderFragment implements the ome.OrderFragmentStorer interface.
func (table *SomerOrderFragmentTable) PutSellOrderFragment(epoch registry.Epoch, orderFragment order.Fragment, trader string) error {
	value := SomerOrderFragmentValue{
		Timestamp:     time.Now(),
		OrderFragment: orderFragment,
		Trader:        trader,
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return table.db.Put(table.sellKey(epoch.Hash[:], orderFragment.OrderID[:]), data, nil)
}

// DeleteSellOrderFragment implements the ome.OrderFragmentStorer interface.
func (table *SomerOrderFragmentTable) DeleteSellOrderFragment(epoch registry.Epoch, id order.ID) error {
	return table.db.Delete(table.sellKey(epoch.Hash[:], id[:]), nil)
}

// SellOrderFragment implements the ome.OrderFragmentStorer interface.
func (table *SomerOrderFragmentTable) SellOrderFragment(epoch registry.Epoch, id order.ID) (order.Fragment, string, error) {
	data, err := table.db.Get(table.sellKey(epoch.Hash[:], id[:]), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = ome.ErrOrderFragmentNotFound
		}
		return order.Fragment{}, "", err
	}

	value := SomerOrderFragmentValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return order.Fragment{}, "", err
	}
	return value.OrderFragment, value.Trader, nil
}

// SellOrderFragments implements the ome.OrderFragmentStorer interface.
func (table *SomerOrderFragmentTable) SellOrderFragments(epoch registry.Epoch) (ome.OrderFragmentIterator, error) {
	iter := table.db.NewIterator(&util.Range{Start: table.sellKey(epoch.Hash[:], SomerSellOrderFragmentIterBegin), Limit: table.sellKey(epoch.Hash[:], SomerSellOrderFragmentIterEnd)}, nil)
	return newSomerOrderFragmentIterator(iter), nil
}

// Prune iterates over all order fragments and deletes those that have expired.
func (table *SomerOrderFragmentTable) Prune() (err error) {
	buyIter := table.db.NewIterator(&util.Range{Start: table.buyKey(SomerBuyOrderFragmentIterBegin, SomerBuyOrderFragmentIterBegin), Limit: table.buyKey(SomerBuyOrderFragmentIterEnd, SomerBuyOrderFragmentIterEnd)}, nil)
	defer buyIter.Release()

	now := time.Now()
	for buyIter.Next() {
		key := buyIter.Key()
		value := SomerOrderFragmentValue{}
		if localErr := json.Unmarshal(buyIter.Value(), &value); localErr != nil {
			err = localErr
			continue
		}
		if value.Timestamp.Add(table.expiry).Before(now) {
			if localErr := table.db.Delete(key, nil); localErr != nil {
				err = localErr
			}
		}
	}

	sellIter := table.db.NewIterator(&util.Range{Start: table.sellKey(SomerSellOrderFragmentIterBegin, SomerSellOrderFragmentIterBegin), Limit: table.sellKey(SomerSellOrderFragmentIterEnd, SomerSellOrderFragmentIterEnd)}, nil)
	defer sellIter.Release()

	for sellIter.Next() {
		key := sellIter.Key()
		value := SomerOrderFragmentValue{}
		if localErr := json.Unmarshal(sellIter.Value(), &value); localErr != nil {
			err = localErr
			continue
		}
		if value.Timestamp.Add(table.expiry).Before(now) {
			if localErr := table.db.Delete(key, nil); localErr != nil {
				err = localErr
			}
		}
	}
	return err
}

func (table *SomerOrderFragmentTable) buyKey(epoch, orderID []byte) []byte {
	return append(append(append(SomerBuyOrderFragmentTableBegin, epoch...), orderID...), SomerBuyOrderFragmentTablePadding...)
}

func (table *SomerOrderFragmentTable) sellKey(epoch, orderID []byte) []byte {
	return append(append(append(SomerSellOrderFragmentTableBegin, epoch...), orderID...), SomerSellOrderFragmentTablePadding...)
}
