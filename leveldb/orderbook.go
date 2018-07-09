package leveldb

import (
	"encoding/json"
	"path"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Constants for use in the OrderbookOrderTable.
var (
	OrderbookOrderTableBegin = []byte{0x01, 0x00}
	OrderbookOrderTableEnd   = []byte{0x01, 0xFF}
	OrderbookOrderIterBegin  = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	OrderbookOrderIterEnd    = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	OrderbookOrderExpiry     = 72 * time.Hour
)

// OrderbookOrderValue is the storage format for values being store in LevelDB.
// It contains metadata, such as the timestamp, so that LevelDB can provide
// features such as expiration.
type OrderbookOrderValue struct {
	Timestamp   time.Time    `json:"timestamp"`
	Status      order.Status `json:"status"`
	Trader      string       `json:"trader"`
	BlockNumber uint64       `json:"blockNumber"`
}

// OrderbookOrderTable implements the orderbook.OrderStorer interface.
type OrderbookOrderTable struct {
	db *leveldb.DB
}

// NewOrderbookOrderTable returns a new OrderbookOrderTable that uses a LevelDB
// instance to store and load values from the disk.
func NewOrderbookOrderTable(db *leveldb.DB) *OrderbookOrderTable {
	return &OrderbookOrderTable{db: db}
}

// PutOrder implements the orderbook.OrderStorer interface.
func (table *OrderbookOrderTable) PutOrder(id order.ID, status order.Status, trader string, blockNumber uint64) error {
	value := OrderbookOrderValue{
		Timestamp:   time.Now(),
		Status:      status,
		Trader:      trader,
		BlockNumber: blockNumber,
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return table.db.Put(append(OrderbookOrderTableBegin, id[:]...), data, nil)
}

// DeleteOrder implements the orderbook.OrderStorer interface.
func (table *OrderbookOrderTable) DeleteOrder(id order.ID) error {
	return table.db.Delete(append(OrderbookOrderTableBegin, id[:]...), nil)
}

// Order implements the orderbook.OrderStorer interface.
func (table *OrderbookOrderTable) Order(id order.ID) (order.Status, string, uint64, error) {
	data, err := table.db.Get(append(OrderbookOrderTableBegin, id[:]...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = orderbook.ErrOrderNotFound
		}
		return order.Nil, "", 0, err
	}

	value := OrderbookOrderValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return order.Nil, "", 0, err
	}
	return value.Status, value.Trader, value.BlockNumber, nil
}

// Orders implements the orderbook.OrderStorer interface.
func (table *OrderbookOrderTable) Orders() (orderbook.OrderIterator, error) {
	iter := table.db.NewIterator(&util.Range{Start: append(OrderbookOrderTableBegin, OrderbookOrderIterBegin...), Limit: append(OrderbookOrderTableEnd, OrderbookOrderIterEnd...)}, nil)
	return newOrderbookOrderIterator(iter), nil
}

// Prune iterates over all orders and deletes those that have expired.
func (table *OrderbookOrderTable) Prune() (err error) {
	iter := table.db.NewIterator(&util.Range{Start: append(OrderbookOrderTableBegin, OrderbookOrderIterBegin...), Limit: append(OrderbookOrderTableEnd, OrderbookOrderIterEnd...)}, nil)
	defer iter.Release()

	now := time.Now()
	for iter.Next() {
		key := iter.Key()
		value := OrderbookOrderValue{}
		if localErr := json.Unmarshal(iter.Value(), &value); localErr != nil {
			err = localErr
			continue
		}
		if value.Timestamp.Add(OrderbookOrderExpiry).Before(now) {
			if localErr := table.db.Delete(key, nil); localErr != nil {
				err = localErr
			}
		}
	}
	return err
}

// OrderbookOrderIterator implements the orderbook.OrderIterator using a
// LevelDB iterator.
type OrderbookOrderIterator struct {
	inner iterator.Iterator
}

func newOrderbookOrderIterator(iter iterator.Iterator) *OrderbookOrderIterator {
	return &OrderbookOrderIterator{
		inner: iter,
	}
}

// Next implements the orderbook.OrderIterator interface.
func (iter *OrderbookOrderIterator) Next() bool {
	return iter.inner.Next()
}

// Cursor implements the orderbook.OrderIterator interface.
func (iter *OrderbookOrderIterator) Cursor() (order.ID, order.Status, error) {
	orderID := order.ID{}
	if !iter.inner.Valid() {
		return orderID, order.Nil, orderbook.ErrCursorOutOfRange
	}

	// Copy the key into the order ID making sure to ignore the table prefix
	copy(orderID[:], iter.inner.Key()[len(OrderbookOrderTableBegin):])

	value := OrderbookOrderValue{}
	if err := json.Unmarshal(iter.inner.Value(), &value); err != nil {
		return orderID, order.Nil, err
	}
	return orderID, value.Status, iter.inner.Error()
}

// Collect implements the orderbook.OrderIterator interface.
func (iter *OrderbookOrderIterator) Collect() ([]order.ID, []order.Status, error) {
	orderIDs := []order.ID{}
	orderStatuses := []order.Status{}
	for iter.Next() {
		orderID, orderStatus, err := iter.Cursor()
		if err != nil {
			return orderIDs, orderStatuses, err
		}
		orderIDs = append(orderIDs, orderID)
		orderStatuses = append(orderStatuses, orderStatus)
	}
	return orderIDs, orderStatuses, iter.inner.Error()
}

// Release implements the orderbook.OrderIterator interface.
func (iter *OrderbookOrderIterator) Release() {
	iter.inner.Release()
}

// OrderbookStore is a LevelDB implementation of the storage interfaces defined
// by the orderbook package. Orders and order fragments are stored in
// persistent storage.
type OrderbookStore struct {
	db *leveldb.DB
}

func NewOrderbookStore(dir string) (*OrderbookStore, error) {
	db, err := leveldb.OpenFile(path.Join(dir, "orderbook"), nil)
	if err != nil {
		return nil, err
	}
	return &OrderbookStore{
		db: db,
	}, nil
}

// Close the internal LevelDB database. A call to OrderbookStore.Close must
// happen to guarantee that all resources are freed correctly.
func (store *OrderbookStore) Close() error {
	return store.db.Close()
}

// PutOrderFragment implements the orderbook.OrderFragmentStorer interface.
func (store *OrderbookStore) PutOrderFragment(epoch registry.Epoch, orderFragment order.Fragment) error {
	data, err := json.Marshal(orderFragment)
	if err != nil {
		return err
	}
	return store.db.Put(append(append(TableOrderbookOrderFragmentsBegin, epoch.Hash[:]...), orderFragment.OrderID[:]...), data, nil)
}

// DeleteOrderFragment implements the orderbook.OrderFragmentStorer interface.
func (store *OrderbookStore) DeleteOrderFragment(epoch registry.Epoch, id order.ID) error {
	return store.db.Delete(append(append(TableOrderbookOrderFragmentsBegin, epoch.Hash[:]...), id[:]...), nil)
}

// OrderFragment implements the orderbook.OrderFragmentStorer interface.
func (store *OrderbookStore) OrderFragment(epoch registry.Epoch, id order.ID) (order.Fragment, error) {
	orderFragment := order.Fragment{}
	data, err := store.db.Get(append(append(TableOrderbookOrderFragmentsBegin, epoch.Hash[:]...), id[:]...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = orderbook.ErrOrderFragmentNotFound
		}
		return orderFragment, err
	}
	if err := json.Unmarshal(data, &orderFragment); err != nil {
		return orderFragment, err
	}
	return orderFragment, nil
}

// OrderFragments implements the orderbook.OrderFragmentStorer interface.
func (store *OrderbookStore) OrderFragments(epoch registry.Epoch) (orderbook.OrderFragmentIterator, error) {
	iter := store.db.NewIterator(&util.Range{Start: TableOrderbookOrderFragmentsBegin, Limit: TableOrderbookOrderFragmentsEnd}, nil)
	return newOrderbookOrderFragmentIterator(iter), nil
}

type OrderbookPointerStore struct {
	pointerMu *sync.RWMutex
	pointer   orderbook.Pointer
}

func NewOrderbookPointerStore() *OrderbookPointerStore {
	return &OrderbookPointerStore{
		pointerMu: new(sync.RWMutex),
		pointer:   orderbook.Pointer(0),
	}
}

func (store *OrderbookPointerStore) PutPointer(pointer orderbook.Pointer) error {
	store.pointerMu.Lock()
	defer store.pointerMu.Unlock()

	store.pointer = pointer
	return nil
}

func (store *OrderbookPointerStore) Pointer() (orderbook.Pointer, error) {
	store.pointerMu.RLock()
	defer store.pointerMu.RUnlock()

	return store.pointer, nil
}

func (store *OrderbookPointerStore) Clone() (orderbook.PointerStorer, error) {
	store.pointerMu.Lock()
	defer store.pointerMu.Unlock()

	return &OrderbookPointerStore{
		pointerMu: new(sync.RWMutex),
		pointer:   store.pointer,
	}, nil
}

// OrderbookOrderFragmentIterator is a LevelDB implementation of the order
// fragment iterator interface defined by the orderbook package.
type OrderbookOrderFragmentIterator struct {
	inner iterator.Iterator
	next  bool
}

func newOrderbookOrderFragmentIterator(iter iterator.Iterator) *OrderbookOrderFragmentIterator {
	return &OrderbookOrderFragmentIterator{
		inner: iter,
		next:  false,
	}
}

// Next implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderbookOrderFragmentIterator) Next() bool {
	iter.next = iter.inner.Next()
	return iter.next
}

// Cursor implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderbookOrderFragmentIterator) Cursor() (order.Fragment, error) {
	orderFragment := order.Fragment{}
	if !iter.next {
		return orderFragment, orderbook.ErrCursorOutOfRange
	}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &orderFragment); err != nil {
		return orderFragment, err
	}
	return orderFragment, iter.inner.Error()
}

// Collect implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderbookOrderFragmentIterator) Collect() ([]order.Fragment, error) {
	orderFragments := []order.Fragment{}
	for iter.Next() {
		orderFragment, err := iter.Cursor()
		if err != nil {
			return orderFragments, err
		}
		orderFragments = append(orderFragments, orderFragment)
	}
	return orderFragments, iter.inner.Error()
}

// Release implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderbookOrderFragmentIterator) Release() {
	iter.inner.Release()
}
