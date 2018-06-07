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

// Store is an implementation of the orderbook.Storer interface that uses
// LevelDB to load and store data to persitent storage.
type Store struct {
	orderFragments *leveldb.DB
	orders         *leveldb.DB
	computations   *leveldb.DB
	sync           *leveldb.DB
}

// NewStore returns a LevelDB implemntation of an orderbooker.Storer. It stores
// and loads order fragments, and orders, using the specified directory.
func NewStore(dir string) (Store, error) {
	orderFragments, err := leveldb.OpenFile(path.Join(dir, "orderFragments"), nil)
	if err != nil {
		return Store{}, err
	}
	orders, err := leveldb.OpenFile(path.Join(dir, "orders"), nil)
	if err != nil {
		return Store{}, err
	}
	computations, err := leveldb.OpenFile(path.Join(dir, "computations"), nil)
	if err != nil {
		return Store{}, err
	}
	sync, err := leveldb.OpenFile(path.Join(dir, "sync"), nil)
	if err != nil {
		return Store{}, err
	}

	return Store{
		orderFragments: orderFragments,
		orders:         orders,
		computations:   computations,
		sync:           sync,
	}, nil
}

// Close the internal LevelDB handles. Resources will be leaked if this is not
// called when the Store is no longer needed.
func (store *Store) Close() error {
	if err := store.orderFragments.Close(); err != nil {
		return err
	}
	if err := store.orders.Close(); err != nil {
		return err
	}

	return store.computations.Close()
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

// RemoveOrderFragment implements the orderbook.Storer interface.
func (store *Store) RemoveOrderFragment(id order.ID) error {
	return store.orderFragments.Delete(id[:], nil)
}

// RemoveOrder implements the orderbook.Storer interface.
func (store *Store) RemoveOrder(id order.ID) error {
	return store.orders.Delete(id[:], nil)
}

func (store *Store) InsertComputation(computation ome.Computation) error {
	data, err := json.Marshal(computation)
	if err != nil {
		return err
	}
	return store.computations.Put(computation.ID[:], data, nil)
}

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
		return 0, err
	}
	pointer, err := binary.ReadVarint(bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	return orderbook.SyncPointer(pointer), nil
}
