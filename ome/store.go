package ome

import (
	"encoding/json"
	"os"
	"path"

	"github.com/republicprotocol/republic-go/order"
	"github.com/syndtr/goleveldb/leveldb"
)

type Storer interface {
	Insert(orderFragment order.Fragment) error
	Get(id order.ID) (order.Fragment, error)
}

// storer is an levelDB implementation of the Storer interface
type storer struct {
	orderFragments *leveldb.DB
}

func NewStore() (storer, error) {
	dbPath := path.Join(os.Getenv("HOME"), ".darknode", "fragments")

	orderFragments, err := leveldb.OpenFile(path.Join(dbPath, "fragments"), nil)
	if err != nil {
		return storer{}, err
	}

	return storer{
		orderFragments: orderFragments,
	}, nil
}

func (storer storer) Insert(fragment order.Fragment) error {
	encoded, err := json.Marshal(fragment)
	if err != nil {
		return err
	}

	return storer.orderFragments.Put(fragment.OrderID[:], encoded, nil)
}

func (storer storer) Get(id order.ID) (order.Fragment, error) {
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
