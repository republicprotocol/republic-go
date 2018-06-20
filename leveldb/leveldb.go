package leveldb

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"path"

	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/order"
	"github.com/syndtr/goleveldb/leveldb"
)

// Store is an implementation of storage interfaces using LevelDB to load and
// store data to persistent storage. It uses a single database instance and
// high-order bytes for separating data into tables. LevelDB provides a basic
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

// InsertOrderFragment implements the orderbook.Storer interface.
func (store *Store) InsertOrderFragment(orderFragment order.Fragment) error {
	data, err := json.Marshal(orderFragment)
	if err != nil {
		return err
	}
	return store.orderFragments.Put(orderFragment.OrderID[:], data, nil)
}

// InsertOrder implements the orderbook.Storer interface.
func (store *Store) InsertOrder(order order.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return store.orders.Put(order.ID[:], data, nil)
}

// InsertComputation implements the ome.Storer interface.
func (store *Store) InsertComputation(computation ome.Computation) error {
	data, err := json.Marshal(computation)
	if err != nil {
		return err
	}
	return store.computations.Put(computation.ID[:], data, nil)
}

// OrderFragment implements the orderbook.Storer interface.
func (store *Store) OrderFragment(id order.ID) (order.Fragment, error) {
	orderFragment := order.Fragment{}
	data, err := store.orderFragments.Get(id[:], nil)
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

// Order implements the orderbook.Storer interface.
func (store *Store) Order(id order.ID) (order.Order, error) {
	order := order.Order{}
	data, err := store.orders.Get(id[:], nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = orderbook.ErrOrderNotFound
		}
		return order, err
	}
	if err := json.Unmarshal(data, &order); err != nil {
		return order, err
	}
	return order, nil
}

// Orders implements the orderbook.Storer interface.
func (store *Store) Orders() ([]order.Order, error) {
	ords := []order.Order{}
	iter := store.orders.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		data := iter.Value()

		ord := order.Order{}
		if err := json.Unmarshal(data, &ord); err != nil {
			return ords, err
		}
		ords = append(ords, ord)
	}
	return ords, iter.Error()
}

// Computation implements the ome.Storer interface.
func (store *Store) Computation(id ome.ComputationID) (ome.Computation, error) {
	computation := ome.Computation{}
	data, err := store.computations.Get(id[:], nil)
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
func (store *Store) Computations() (ome.Computations, error) {
	coms := ome.Computations{}
	iter := store.computations.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		data := iter.Value()

		com := ome.Computation{}
		if err := json.Unmarshal(data, &com); err != nil {
			return coms, err
		}
		coms = append(coms, com)
	}
	return coms, iter.Error()
}

// RemoveOrderFragment implements the orderbook.Storer interface.
func (store *Store) RemoveOrderFragment(id order.ID) error {
	return store.orderFragments.Delete(id[:], nil)
}

// RemoveOrder implements the orderbook.Storer interface.
func (store *Store) RemoveOrder(id order.ID) error {
	return store.orders.Delete(id[:], nil)
}

// RemoveComputation implements the ome.Storer interface.
func (store *Store) RemoveComputation(id ome.ComputationID) error {
	return store.computations.Delete(id[:], nil)
}

// InsertBuyPointer implements the orderbook.SyncStorer interface.
func (store *Store) InsertBuyPointer(pointer orderbook.SyncPointer) error {
	data := [4]byte{}
	binary.PutVarint(data[:], int64(pointer))
	return store.sync.Put([]byte("buy"), data[:], nil)
}

// InsertSellPointer implements the orderbook.SyncStorer interface.
func (store *Store) InsertSellPointer(pointer orderbook.SyncPointer) error {
	data := [4]byte{}
	binary.PutVarint(data[:], int64(pointer))
	return store.sync.Put([]byte("sell"), data[:], nil)
}

// BuyPointer implements the orderbook.SyncStorer interface.
func (store *Store) BuyPointer() (orderbook.SyncPointer, error) {
	data, err := store.sync.Get([]byte("buy"), nil)
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
	return orderbook.SyncPointer(pointer), nil
}

// SellPointer implements the orderbook.SyncStorer interface.
func (store *Store) SellPointer() (orderbook.SyncPointer, error) {
	data, err := store.sync.Get([]byte("sell"), nil)
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
	return orderbook.SyncPointer(pointer), nil
}
