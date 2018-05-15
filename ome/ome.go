package ome

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
)

var ErrNotFoundInStore = errors.New("not found in store")

type Omer interface {
	SyncRenLedger() error
}

type Ome struct {
	delegate Delegate
	ranker   Ranker
	storer   Storer
	syncer   Syncer
}

func NewOme(delegate Delegate, ranker Ranker, storer Storer, syncer Syncer) Ome {
	return Ome{
		delegate: delegate,
		ranker:   ranker,
		storer:   storer,
		syncer:   syncer,
	}
}

func (engine *Ome) OpenOrder(orderFragment order.Fragment) error {
	return engine.storer.Insert(orderFragment)
}

type Delegate interface {
	ConfirmOrderMatch(orders []order.Order)
}

type Storer interface {
	Insert(orderFragment order.Fragment) error
	Get(id order.ID) (order.Fragment, error)
}

type Ranker interface {
	Insert(id order.ID, priority int) error
}
