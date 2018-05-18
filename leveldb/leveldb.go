package leveldb

import (
	"encoding/json"
	"path"

	"github.com/republicprotocol/republic-go/order"
	"github.com/syndtr/goleveldb/leveldb"
)

// Store is an implementation of the orderbook.Storer interface. That uses
// LevelDB to load and store data to persitent storage.
type Store struct {
	orderFragments *leveldb.DB
	orders         *leveldb.DB
}

func NewStore(dir string) (Store, error) {
	orderFragments, err := leveldb.OpenFile(path.Join(dir, "orderFragments"), nil)
	if err != nil {
		return Store{}, err
	}
	orders, err := leveldb.OpenFile(path.Join(dir, "orders"), nil)
	if err != nil {
		return Store{}, err
	}

	return Store{
		orderFragments: orderFragments,
		orders:         orders,
	}, nil
}

func (store *Store) Close() error {
	if err := store.orderFragments.Close(); err != nil {
		return err
	}
	return store.orders.Close()
}

func (store *Store) InsertOrderFragment(orderFragment order.Fragment) error {
	data, err := json.Marshal(orderFragment)
	if err != nil {
		return err
	}
	return store.orderFragments.Put(orderFragment.OrderID[:], data, nil)
}

func (store *Store) InsertOrder(order order.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return store.orders.Put(order.ID[:], data, nil)
}

func (store *Store) OrderFragment(id order.ID) (order.Fragment, error) {
	orderFragment := order.Fragment{}
	data, err := store.orderFragments.Get(id[:], nil)
	if err != nil {
		return orderFragment, err
	}
	if err := json.Unmarshal(data, &orderFragment); err != nil {
		return orderFragment, err
	}
	return orderFragment, nil
}

func (store *Store) Order(id order.ID) (order.Order, error) {
	order := order.Order{}
	data, err := store.orders.Get(id[:], nil)
	if err != nil {
		return order, err
	}
	if err := json.Unmarshal(data, &order); err != nil {
		return order, err
	}
	return order, nil
}

func (store *Store) RemoveOrderFragment(id order.ID) error {
	return store.orderFragments.Delete(id[:], nil)
}

func (store *Store) RemoveOrder(id order.ID) error {
	return store.orders.Delete(id[:], nil)
}
