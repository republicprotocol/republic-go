package ome

import (
	"errors"
	"log"

	"github.com/republicprotocol/republic-go/order"
)

var ErrNotFoundInStore = errors.New("not found in store")

type Omer interface {
	SyncRenLedger() error
	ConfirmOrder(id order.ID, matches []order.ID) error
	//todo : connect with smpc
}

type Ome struct {
	storer Storer
	syncer Syncer
	ranker Ranker
}

func NewOme(ranker Ranker, storer Storer, syncer Syncer) Ome {
	return Ome{
		ranker: ranker,
		storer: storer,
		syncer: syncer,
	}
}

func (engine *Ome) OpenOrder(orderFragment order.Fragment) error {
	return engine.storer.Insert(orderFragment)
}

func (engine *Ome) SyncRenLedger() error {
	done := make(chan struct{})
	defer close(done)

	go func() {
		// todo : handle error
		engine.syncer.SyncRenLedger(done, engine.ranker)
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// todo : pass the order the orderPairs to smpc and get result
				orderPairs := engine.ranker.Get(50)
				log.Print(orderPairs)
			}
		}
	}()

	return nil
}

func (engine *Ome) ConfirmOrder(id order.ID, matches []order.ID) error {
	if err := engine.syncer.ConfirmOrder(id, matches); err != nil {
		return err
	}

	// Remove the order from the ranker
	engine.ranker.Remove(id)
	for _, match := range matches {
		engine.ranker.Remove(match)
	}

	return nil
}
