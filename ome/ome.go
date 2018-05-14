package ome

import (
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/cal"

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

	renLedgerMu      *sync.Mutex
	renLedger        cal.RenLedger
	renLedgerPointer int
	renLedgerLimit   int
}

func NewOme(delegate Delegate, ranker Ranker, storer Storer, renLedger cal.RenLedger) Ome {
	return Ome{
		delegate: delegate,
		ranker:   ranker,
		storer:   storer,

		renLedgerMu:      new(sync.Mutex),
		renLedger:        renLedger,
		renLedgerPointer: 0,
		renLedgerLimit:   100, // TODO: Configuration option for this limit
	}
}

func (engine *Ome) SyncRenLedger() error {
	engine.renLedgerMu.Lock()
	defer engine.renLedgerMu.Unlock()

	orderIDs, err := engine.renLedger.Orders(engine.renLedgerPointer, engine.renLedgerLimit)
	if err != nil {
		return err
	}

	err = nil
	for i, id := range orderIDs {
		if errLocal := engine.ranker.Insert(id, engine.renLedgerPointer+i); errLocal != nil && err == nil {
			err = errLocal
		}
	}
	engine.renLedgerPointer += len(orderIDs)

	return err
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
