package leveldb

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// OrderbookOrderValue is the storage format for orders being stored in
// LevelDB. It contains additional timestamping information so that LevelDB can
// provide pruning.
type OrderbookOrderValue struct {
	Timestamp   time.Time    `json:"timestamp"`
	Status      order.Status `json:"status"`
	Trader      string       `json:"trader"`
	BlockNumber uint64       `json:"blockNumber"`
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
	offset := len(OrderbookOrderTableBegin)
	length := len(orderID)
	copy(orderID[:], iter.inner.Key()[offset:offset+length])
	log.Println("key", iter.inner.Key())
	log.Println("orderID", orderID)

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
	return table.db.Put(table.key(id[:]), data, nil)
}

// DeleteOrder implements the orderbook.OrderStorer interface.
func (table *OrderbookOrderTable) DeleteOrder(id order.ID) error {
	return table.db.Delete(table.key(id[:]), nil)
}

// Order implements the orderbook.OrderStorer interface.
func (table *OrderbookOrderTable) Order(id order.ID) (order.Status, string, uint64, error) {
	data, err := table.db.Get(table.key(id[:]), nil)
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
	iter := table.db.NewIterator(&util.Range{Start: table.key(OrderbookOrderIterBegin), Limit: table.key(OrderbookOrderIterEnd)}, nil)
	return newOrderbookOrderIterator(iter), nil
}

// Prune iterates over all orders and deletes those that have expired.
func (table *OrderbookOrderTable) Prune() (err error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(OrderbookOrderIterBegin), Limit: table.key(OrderbookOrderIterEnd)}, nil)
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

func (table *OrderbookOrderTable) key(k []byte) []byte {
	return append(append(OrderbookOrderTableBegin, k...), OrderbookOrderTablePadding...)
}

// OrderbookOrderFragmentValue is the storage format for order fragments being
// stored in LevelDB. It contains additional timestamping information so that
// LevelDB can provide pruning.
type OrderbookOrderFragmentValue struct {
	Timestamp     time.Time      `json:"timestamp"`
	OrderFragment order.Fragment `json:"orderFragment"`
}

// OrderbookOrderFragmentIterator implements the
// orderbook.OrderFragmentIterator using a LevelDB iterator.
type OrderbookOrderFragmentIterator struct {
	inner iterator.Iterator
}

func newOrderbookOrderFragmentIterator(iter iterator.Iterator) *OrderbookOrderFragmentIterator {
	return &OrderbookOrderFragmentIterator{
		inner: iter,
	}
}

// Next implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderbookOrderFragmentIterator) Next() bool {
	return iter.inner.Next()
}

// Cursor implements the orderbook.OrderFragmentIterator interface.
func (iter *OrderbookOrderFragmentIterator) Cursor() (order.Fragment, error) {
	if !iter.inner.Valid() {
		return order.Fragment{}, orderbook.ErrCursorOutOfRange
	}
	value := OrderbookOrderFragmentValue{}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &value); err != nil {
		return order.Fragment{}, err
	}
	return value.OrderFragment, iter.inner.Error()
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

// OrderbookOrderFragmentTable implements the orderbook.OrderFragmentStorer interface.
type OrderbookOrderFragmentTable struct {
	db *leveldb.DB
}

// NewOrderbookOrderFragmentTable returns a new OrderbookOrderFragmentTable that uses a LevelDB
// instance to store and load values from the disk.
func NewOrderbookOrderFragmentTable(db *leveldb.DB) *OrderbookOrderFragmentTable {
	return &OrderbookOrderFragmentTable{db: db}
}

// PutOrderFragment implements the orderbook.OrderFragmentStorer interface.
func (table *OrderbookOrderFragmentTable) PutOrderFragment(epoch registry.Epoch, orderFragment order.Fragment) error {
	value := OrderbookOrderFragmentValue{
		Timestamp:     time.Now(),
		OrderFragment: orderFragment,
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return table.db.Put(table.key(epoch.Hash[:], orderFragment.OrderID[:]), data, nil)
}

// DeleteOrderFragment implements the orderbook.OrderFragmentStorer interface.
func (table *OrderbookOrderFragmentTable) DeleteOrderFragment(epoch registry.Epoch, id order.ID) error {
	return table.db.Delete(table.key(epoch.Hash[:], id[:]), nil)
}

// OrderFragment implements the orderbook.OrderFragmentStorer interface.
func (table *OrderbookOrderFragmentTable) OrderFragment(epoch registry.Epoch, id order.ID) (order.Fragment, error) {
	data, err := table.db.Get(table.key(epoch.Hash[:], id[:]), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = orderbook.ErrOrderFragmentNotFound
		}
		return order.Fragment{}, err
	}

	value := OrderbookOrderFragmentValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return order.Fragment{}, err
	}
	return value.OrderFragment, nil
}

// OrderFragments implements the orderbook.OrderFragmentStorer interface.
func (table *OrderbookOrderFragmentTable) OrderFragments(epoch registry.Epoch) (orderbook.OrderFragmentIterator, error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(epoch.Hash[:], OrderbookOrderFragmentIterBegin), Limit: table.key(epoch.Hash[:], OrderbookOrderFragmentIterEnd)}, nil)
	return newOrderbookOrderFragmentIterator(iter), nil
}

// Prune iterates over all orders and deletes those that have expired.
func (table *OrderbookOrderFragmentTable) Prune() (err error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(OrderbookOrderFragmentIterBegin, OrderbookOrderFragmentIterBegin), Limit: table.key(OrderbookOrderFragmentIterEnd, OrderbookOrderFragmentIterEnd)}, nil)
	defer iter.Release()

	now := time.Now()
	for iter.Next() {
		key := iter.Key()
		value := OrderbookOrderFragmentValue{}
		if localErr := json.Unmarshal(iter.Value(), &value); localErr != nil {
			err = localErr
			continue
		}
		if value.Timestamp.Add(OrderbookOrderFragmentExpiry).Before(now) {
			if localErr := table.db.Delete(key, nil); localErr != nil {
				err = localErr
			}
		}
	}
	return err
}

func (table *OrderbookOrderFragmentTable) key(epoch, orderID []byte) []byte {
	return append(append(append(OrderbookOrderFragmentTableBegin, epoch...), orderID...), OrderbookOrderFragmentTablePadding...)
}

// OrderbookPointerTable implements the orderbook.PointerStorer using in-memory
// storage. Data stored is not persistent across reboots.
type OrderbookPointerTable struct {
	pointerMu *sync.RWMutex
	pointer   orderbook.Pointer
}

// NewOrderbookPointerTable returns a new OrderbookPointerTable with the
// orderbook.Pointer initialised to zero.
func NewOrderbookPointerTable() *OrderbookPointerTable {
	return &OrderbookPointerTable{
		pointerMu: new(sync.RWMutex),
		pointer:   orderbook.Pointer(0),
	}
}

// PutPointer implements the orderbook.PointerStorer interface.
func (store *OrderbookPointerTable) PutPointer(pointer orderbook.Pointer) error {
	store.pointerMu.Lock()
	defer store.pointerMu.Unlock()

	store.pointer = pointer
	return nil
}

// Pointer implements the orderbook.PointerStorer interface.
func (store *OrderbookPointerTable) Pointer() (orderbook.Pointer, error) {
	store.pointerMu.RLock()
	defer store.pointerMu.RUnlock()

	return store.pointer, nil
}

// Clone implements the orderbook.PointerStorer interface.
func (store *OrderbookPointerTable) Clone() (orderbook.PointerStorer, error) {
	store.pointerMu.Lock()
	defer store.pointerMu.Unlock()

	return &OrderbookPointerTable{
		pointerMu: new(sync.RWMutex),
		pointer:   store.pointer,
	}, nil
}
