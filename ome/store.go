package ome

import (
	"encoding/json"
	"errors"

	"github.com/republicprotocol/republic-go/order"
	"github.com/syndtr/goleveldb/leveldb"
)

var ErrNotFoundInStore = errors.New("not found in store")

type Storer interface {
	Insert(orderFragment order.Fragment) error
	Get(id order.ID) (order.Fragment, error)
	Delete(id order.ID) error
	Stop() error
}

// LevelDBStorer is an levelDB implementation of the Storer interface
type LevelDBStorer struct {
	orderFragments *leveldb.DB
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

func (storer LevelDBStorer) Delete(id order.ID) error  {
	return storer.orderFragments.Delete(id[:], nil)
}
