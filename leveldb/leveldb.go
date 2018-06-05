package leveldb

import (
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
	computaions, err := leveldb.OpenFile(path.Join(dir, "computations"), nil)
	if err != nil {
		return Store{}, err
	}

	return Store{
		orderFragments: orderFragments,
		orders:         orders,
		computations:   computaions,
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

func (store *Store) InsertComputation(computations ome.Computation) error {
	data, err := json.Marshal(computations)
	if err != nil {
		return err
	}

	return store.computations.Put(computations.ID()[:], data, nil)
}

func (store *Store) Computation(id [32]byte) (ome.Computation, error) {
	computation := ome.Computation{}
	data, err := store.computations.Get(id[:], nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = orderbook.ErrOrderNotFound
		}
		return computation, err
	}
	if err := json.Unmarshal(data, &computation); err != nil {
		return computation, err
	}
	return computation, nil
}
