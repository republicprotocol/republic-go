package leveldb

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"path"
	"sync"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Table prefixes for the OrderbookStore.
var (
	TableOrderbookOrdersBegin         = []byte{0x01, 0x00}
	TableOrderbookOrdersEnd           = []byte{0x01, 0xFF}
	TableOrderbookOrderFragmentsBegin = []byte{0x02, 0x00}
	TableOrderbookOrderFragmentsEnd   = []byte{0x02, 0xFF}
)

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

// PutOrder implements the orderbook.OrderStorer interface.
func (store *OrderbookStore) PutOrder(id order.ID, status order.Status) error {
	data := [4]byte{}
	binary.PutUvarint(data[:], uint64(status))
	return store.db.Put(append(TableOrderbookOrdersBegin, id[:]...), data[:], nil)
}

// DeleteOrder implements the orderbook.OrderStorer interface.
func (store *OrderbookStore) DeleteOrder(id order.ID) error {
	return store.db.Delete(append(TableOrderbookOrdersBegin, id[:]...), nil)
}

// Order implements the orderbook.OrderStorer interface.
func (store *OrderbookStore) Order(id order.ID) (order.Status, error) {
	data, err := store.db.Get(append(TableOrderbookOrdersBegin, id[:]...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = orderbook.ErrOrderNotFound
		}
		return order.Nil, err
	}
	orderStatus, err := binary.ReadUvarint(bytes.NewBuffer(data[:]))
	if err != nil {
		return order.Nil, err
	}

	// TODO: Casting the orderStatus to a uint8 from a uint64 can cause an
	// overflow error. The value of orderStatus should be checked before the
	// cast, and an appropriate error should be returned in the case of
	// erroneous values.

	return order.Status(orderStatus), nil
}

// Orders implements the orderbook.OrderStorer interface.
func (store *OrderbookStore) Orders() (orderbook.OrderIterator, error) {
	iter := store.db.NewIterator(&util.Range{Start: TableOrderbookOrdersBegin, Limit: TableOrderbookOrdersEnd}, nil)
	return newOrderbookOrderIterator(iter), nil
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
	data, err := store.db.Get(append(TableOrderbookOrderFragmentsBegin, id[:]...), nil)
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

// OrderbookOrderIterator is a LevelDB implementation of the order iterator
// interface defined by the orderbook package.
type OrderbookOrderIterator struct {
	inner iterator.Iterator
	next  bool
}

func newOrderbookOrderIterator(iter iterator.Iterator) *OrderbookOrderIterator {
	return &OrderbookOrderIterator{
		inner: iter,
		next:  false,
	}
}

// Next implements the orderbook.OrderIterator interface.
func (iter *OrderbookOrderIterator) Next() bool {
	iter.next = iter.inner.Next()
	return iter.next
}

// Cursor implements the orderbook.OrderIterator interface.
func (iter *OrderbookOrderIterator) Cursor() (order.ID, order.Status, error) {
	orderID := order.ID{}
	if !iter.next {
		return orderID, order.Nil, orderbook.ErrCursorOutOfRange
	}

	// Copy the key into the order ID making sure to ignore the table prefix
	copy(orderID[:], iter.inner.Key()[len(TableOrderbookOrdersBegin):])

	// Read the data as an order status
	data := iter.inner.Value()
	orderStatus, err := binary.ReadUvarint(bytes.NewBuffer(data[:]))
	if err != nil {
		return orderID, order.Nil, err
	}

	// TODO: Casting the orderStatus to a uint8 from a uint64 can cause an
	// overflow error. The value of orderStatus should be checked before the
	// cast, and an appropriate error should be returned in the case of
	// erroneous values.

	return orderID, order.Status(orderStatus), iter.inner.Error()
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
