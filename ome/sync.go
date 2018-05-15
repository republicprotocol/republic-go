package ome

import (
	"log"
	"sort"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
)

type Syncer struct {
	renLedger               cal.RenLedger
	renLedgerSettledPointer int
	renLedgerSyncedPointer  int
	renLedgerLimit          int
	renLedgerSyncedInterval int
	orders                  map[int]order.ID
	queue                   *PriorityQueue
}

func NewSyncer(renLedger cal.RenLedger, limit, interval int) Syncer {
	return Syncer{
		renLedger:               renLedger,
		renLedgerSettledPointer: 0,
		renLedgerSyncedPointer:  0,
		renLedgerLimit:          limit,
		renLedgerSyncedInterval: interval,
		orders:                  map[int]order.ID{},
		queue:                   NewPriorityQueue(),
	}
}

func (syncer *Syncer) SyncRenLedger(done <-chan struct{}) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		for {
			syncer.Prune()

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
					orderPair := NewOrderPair(j, n, syncer.renLedgerSyncedPointer+i+m)
					syncer.queue.Insert(orderPair)
				}
				syncer.orders[syncer.renLedgerSyncedPointer+i] = orderIDs[i]
			}
			syncer.renLedgerSyncedPointer += len(orderIDs)

			time.Sleep(time.Duration(syncer.renLedgerSyncedInterval) * time.Second)
		}
	}()

	return errs
}

func (syncer *Syncer) Prune() {
	dispatch.CoForAll(syncer.orders, func(key int) {
		status, err := syncer.renLedger.CheckStatus(syncer.orders[key])
		if err != nil {
			log.Println("fail to check order status", err)
			return
		}

		if status == order.Canceled || status == order.Confirmed || status == order.Settled {
			delete(syncer.orders, key)
			syncer.queue.Remove(syncer.orders[key])
		}
	})
}

type OrderPair struct {
	one      order.ID
	another  order.ID
	priority int
}

func NewOrderPair(one, another order.ID, priority int) OrderPair {
	return OrderPair{
		one:      one,
		another:  another,
		priority: priority,
	}
}

type PriorityQueue struct {
	pairs []OrderPair
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		pairs: []OrderPair{},
	}
}

func (queue *PriorityQueue) Insert(pair OrderPair) {
	index := sort.Search(len(queue.pairs), func(i int) bool {
		return queue.pairs[i].priority > pair.priority
	})
	queue.pairs = append(queue.pairs[:index-1], append([]OrderPair{pair}, queue.pairs[index:]...)...)
}

func (queue *PriorityQueue) Remove(id order.ID) {
	for i := 0; i < len(queue.pairs); i++ {
		if queue.pairs[i].one == id || queue.pairs[i].another == id {
			// a:= []{1,2,3}
			// log.Println(a[3])   # this will panic
			// log.Println(a[3:])  # this will not panic ,amazing!
			queue.pairs = append(queue.pairs[:i], queue.pairs[i+1:]...)
		}
	}
}
