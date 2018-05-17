package orderbook

import (
	"encoding/json"
	"errors"

	"github.com/republicprotocol/republic-go/order"
	"github.com/syndtr/goleveldb/leveldb"
)

var ErrNotFoundInStore = errors.New("not found in store")

type Storer interface {
	InsertOrderFragment(orderFragment order.Fragment) error
	InsertOrder(order order.Order) error
	OrderFragment(id order.ID) (order.Fragment, error)
	Order(id order.ID) (order.Order, error)
	RemoveOrderFragment(id order.ID) error
	RemoveOrder(id order.ID) error
	Stop() error
}

type nilStore struct {
}

// NewNilStore returns an implementation of the Storer interface that does not
// store any data.
func NewNilStore() Storer {
	return &nilStore{}
}

func (store *nilStore) InsertOrderFragment(orderFragment order.Fragment) error {
	return nil
}

func (store *nilStore) InsertOrder(order order.Order) error {
	return nil
}

func (store *nilStore) OrderFragment(id order.ID) (order.Fragment, error) {
	return order.Fragment{}, ErrOrderFragmentNotFound
}

func (store *nilStore) Order(id order.ID) (order.Order, error) {
	return order.Order{}, ErrOrderNotFound
}

func (store *nilStore) RemoveOrderFragment(id order.ID) error {
	return nil
}

func (store *nilStore) RemoveOrder(id order.ID) error {
	return nil
}

func (store *nilStore) Stop() error {
	return nil
}

// LevelDBStorer is an levelDB implementation of the Storer interface
type LevelDBStorer struct {
	orderFragments *leveldb.DB
	orders         *leveldb.DB
}

func NewLevelDBStorer(dbPath string) (LevelDBStorer, error) {
	orderFragments, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return LevelDBStorer{}, err
	}

	return LevelDBStorer{
		orderFragments: orderFragments,
	}, nil
}

func (storer LevelDBStorer) Stop() error {
	return storer.orderFragments.Close()
}

func (storer LevelDBStorer) Insert(fragment order.Fragment) error {
	encoded, err := json.Marshal(fragment)
	if err != nil {
		return err
	}

	return storer.orderFragments.Put(fragment.OrderID[:], encoded, nil)
}

func (storer LevelDBStorer) Get(id order.ID) (order.Fragment, error) {
	data, err := storer.orderFragments.Get(id[:], nil)
	if err != nil {
		return order.Fragment{}, err
	}

	var fragment order.Fragment
	err = json.Unmarshal(data, fragment)
	if err != nil {
		return order.Fragment{}, nil
	}

	return fragment, nil
}

func (storer LevelDBStorer) Delete(id order.ID) error {
	return storer.orderFragments.Delete(id[:], nil)
}
