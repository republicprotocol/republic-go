package leveldb

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"path"

	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Key prefixes to partition data into tables.
var (
	TableOrderbookPointers  = []byte{0x01, 0x0, 0x0, 0x0}
	TableOrderbookChanges   = []byte{0x02, 0x0, 0x0, 0x0}
	TableOrderbookFragments = []byte{0x03, 0x0, 0x0, 0x0}
	TableOmeComputations    = []byte{0x04, 0x0, 0x0, 0x0}
)

// Key postfixes used by iterators to define table ranges.
var (
	TableIterStart = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	TableIterEnd   = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

// Store is an implementation of storage interfaces using LevelDB to load and
// store data to persistent storage. It uses a single database instance and
// high-order bytes for separating data into tables. All keys are 40 bytes
// long, 8 table prefix bytes and a 32 byte key. LevelDB provides a basic type
// type of in-memory caching but has no optimisations that are specific to the
// data.
type Store struct {
	db *leveldb.DB
}

// NewStore returns a LevelDB implementation of storage interfaces. A call to
// Store.Close is required to free resources allocated by the Store.
func NewStore(dir string) (Store, error) {
	db, err := leveldb.OpenFile(path.Join(dir, "db"), nil)
	if err != nil {
		return Store{}, err
	}
	return Store{
		db: db,
	}, nil
}

// Close the internal LevelDB database.
func (store *Store) Close() error {
	return store.db.Close()
}

// PutBuyPointer implements the orderbook.PointerStorer interface.
func (store *Store) PutBuyPointer(pointer orderbook.Pointer) error {
	buy := [32]byte{0}
	data := [4]byte{}
	binary.PutVarint(data[:], int64(pointer))
	return store.db.Put(append(TableOrderbookPointers, buy[:]...), data[:], nil)
}

// BuyPointer implements the orderbook.PointerStorer interface.
func (store *Store) BuyPointer() (orderbook.Pointer, error) {
	buy := [32]byte{0}
	data, err := store.db.Get(append(TableOrderbookPointers, buy[:]...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}
	pointer, err := binary.ReadVarint(bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	return orderbook.Pointer(pointer), nil
}

// PutSellPointer implements the orderbook.PointerStorer interface.
func (store *Store) PutSellPointer(pointer orderbook.Pointer) error {
	sell := [32]byte{1}
	data := [4]byte{}
	binary.PutVarint(data[:], int64(pointer))
	return store.db.Put(append(TableOrderbookPointers, sell[:]...), data[:], nil)
}

// SellPointer implements the orderbook.PointerStorer interface.
func (store *Store) SellPointer() (orderbook.Pointer, error) {
	sell := [32]byte{1}
	data, err := store.db.Get(append(TableOrderbookPointers, sell[:]...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}
	pointer, err := binary.ReadVarint(bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	return orderbook.Pointer(pointer), nil
}

// PutChange implements the orderbook.ChangeStorer interface.
func (store *Store) PutChange(change orderbook.Change) error {
	data, err := json.Marshal(change)
	if err != nil {
		return err
	}
	return store.db.Put(append(TableOrderbookChanges, change.OrderID[:]...), data, nil)
}

// DeleteChange implements the orderbook.ChangeStorer interface.
func (store *Store) DeleteChange(id order.ID) error {
	return store.db.Delete(append(TableOrderbookChanges, id[:]...), nil)
}

// Change implements the orderbook.ChangeStorer interface.
func (store *Store) Change(id order.ID) (orderbook.Change, error) {
	change := orderbook.Change{}
	data, err := store.db.Get(append(TableOrderbookChanges, id[:]...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = orderbook.ErrChangeNotFound
		}
		return change, err
	}
	if err := json.Unmarshal(data, &change); err != nil {
		return change, err
	}
	return change, nil
}

// Changes implements the orderbook.ChangeStorer interface.
func (store *Store) Changes() (orderbook.ChangeIterator, error) {
	iter := store.db.NewIterator(&util.Range{Start: append(TableOrderbookChanges, TableIterStart...), Limit: append(TableOrderbookChanges, TableIterEnd...)}, nil)
	return newChangeIterator(iter), nil
}

// PutOrderFragment implements the orderbook.OrderFragmentStorer interface.
func (store *Store) PutOrderFragment(fragment order.Fragment) error {
	data, err := json.Marshal(fragment)
	if err != nil {
		return err
	}
	return store.db.Put(append(TableOrderbookFragments, fragment.OrderID[:]...), data, nil)
}

// DeleteOrderFragment implements the orderbook.OrderFragmentStorer interface.
func (store *Store) DeleteOrderFragment(id order.ID) error {
	return store.db.Delete(append(TableOrderbookFragments, id[:]...), nil)
}

// OrderFragment implements the orderbook.OrderFragmentStorer interface.
func (store *Store) OrderFragment(id order.ID) (order.Fragment, error) {
	change := order.Fragment{}
	data, err := store.db.Get(append(TableOrderbookFragments, id[:]...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = orderbook.ErrOrderFragmentNotFound
		}
		return change, err
	}
	if err := json.Unmarshal(data, &change); err != nil {
		return change, err
	}
	return change, nil
}

// OrderFragments implements the orderbook.OrderFragmentStorer interface.
func (store *Store) OrderFragments() (orderbook.OrderFragmentIterator, error) {
	iter := store.db.NewIterator(&util.Range{Start: append(TableOrderbookFragments, TableIterStart...), Limit: append(TableOrderbookFragments, TableIterEnd...)}, nil)
	return newOrderFragmentIterator(iter), nil
}

// PutComputation implements the ome.Storer interface.
func (store *Store) PutComputation(computation ome.Computation) error {
	data, err := json.Marshal(computation)
	if err != nil {
		return err
	}
	return store.db.Put(append(TableOmeComputations, computation.ID[:]...), data, nil)
}

// DeleteComputation implements the ome.Storer interface.
func (store *Store) DeleteComputation(id ome.ComputationID) error {
	return store.db.Delete(append(TableOmeComputations, id[:]...), nil)
}

// Computation implements the ome.Storer interface.
func (store *Store) Computation(id ome.ComputationID) (ome.Computation, error) {
	computation := ome.Computation{}
	data, err := store.db.Get(append(TableOmeComputations, id[:]...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = ome.ErrComputationNotFound
		}
		return computation, err
	}
	if err := json.Unmarshal(data, &computation); err != nil {
		return computation, err
	}
	return computation, nil
}

// Computations implements the ome.Storer interface.
func (store *Store) Computations() (ome.ComputationIterator, error) {
	iter := store.db.NewIterator(&util.Range{Start: append(TableOmeComputations, TableIterStart...), Limit: append(TableOmeComputations, TableIterEnd...)}, nil)
	return newComputationIterator(iter), nil
}
