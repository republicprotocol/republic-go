package ome

import (
	"log"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
)

type Syncer interface {
	SyncRenLedger(done <-chan struct{}, ranker Ranker) <-chan error
	ConfirmOrder(id order.ID, matches []order.ID) error
}

type OrderWithPriority struct {
	ID       order.ID
	Priority uint64
}

type syncer struct {
	renLedger               cal.RenLedger
	renLedgerSyncedPointer  int
	renLedgerLimit          int
	renLedgerSyncedInterval int
	orders                  map[int]order.ID
}

func NewSyncer(renLedger cal.RenLedger, limit, interval int) syncer {
	return syncer{
		renLedger:               renLedger,
		renLedgerSyncedPointer:  0,
		renLedgerLimit:          limit,
		renLedgerSyncedInterval: interval,
		orders:                  map[int]order.ID{},
	}
}

func (syncer *syncer) SyncRenLedger(done <-chan struct{}, ranker Ranker) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		for {
			syncer.Prune(ranker)

			orderIDs, err := syncer.renLedger.Orders(syncer.renLedgerSyncedPointer, syncer.renLedgerLimit)
			if err != nil {
				errs <- err
				return
			}

			for i, j := range orderIDs {
				for m, n := range syncer.orders {
					if i == m {
						continue
					}
					orderPair := NewOrderPair(j, []order.ID{n}, syncer.renLedgerSyncedPointer+i+m)
					ranker.Insert(orderPair)
				}
				syncer.orders[syncer.renLedgerSyncedPointer+i] = orderIDs[i]
			}
			syncer.renLedgerSyncedPointer += len(orderIDs)

			time.Sleep(time.Duration(syncer.renLedgerSyncedInterval) * time.Second)
		}
	}()

	return errs
}

func (syncer *syncer) ConfirmOrder(id order.ID, matches []order.ID) error {
	return syncer.renLedger.ConfirmOrder(id, matches)
}

func (syncer *syncer) Prune(ranker Ranker) {
	dispatch.CoForAll(syncer.orders, func(key int) {
		status, err := syncer.renLedger.Status(syncer.orders[key])
		if err != nil {
			log.Println("fail to check order status", err)
			return
		}

		if status == order.Canceled || status == order.Confirmed || status == order.Settled {
			delete(syncer.orders, key)
			ranker.Remove(syncer.orders[key])
		}
	})
}

type OrderPair struct {
	orderID  order.ID
	matches  []order.ID
	priority int
}

func NewOrderPair(id order.ID, matches []order.ID, priority int) OrderPair {
	return OrderPair{
		orderID:  id,
		matches:  matches,
		priority: priority,
	}
}
