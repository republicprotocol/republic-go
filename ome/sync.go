package ome

import (
	"sort"
	"sync"
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

	ordersMu *sync.RWMutex
	orders   map[int]order.ID

	queueMu *sync.Mutex
	queue   PriorityQueue
}

func NewSyncer(renLedger cal.RenLedger, limit, interval int) Syncer {
	return Syncer{
		renLedger:               renLedger,
		renLedgerSettledPointer: 0,
		renLedgerSyncedPointer:  0,
		renLedgerLimit:          limit,
		renLedgerSyncedInterval: interval,

		ordersMu: new(sync.RWMutex),
		orders:   map[int]order.ID{},

		queueMu: new(sync.Mutex),
		queue:   NewPriorityQueue(),
	}
}

func (syncer *Syncer) SyncRenLedger(done <-chan struct{}) <-chan error {
	errs := make(chan error, 1)

	go syncer.Prune(done)

	go func() {
		defer close(errs)

		for {
			orderIDs, err := syncer.renLedger.Orders(syncer.renLedgerSyncedPointer, syncer.renLedgerLimit)
			if err != nil {
				errs <- err
				return
			}

			syncer.ordersMu.Lock()
			for i := range orderIDs {
				syncer.orders[syncer.renLedgerSyncedPointer+i] = orderIDs[i]
			}

			for i, j := range syncer.orders {
				for m, n := range syncer.orders {
					if i == m {
						continue
					}
					orderPair := NewOrderPair(j, n, i+m)
					syncer.queue.Insert(orderPair)
				}
			}
			syncer.ordersMu.Unlock()

			time.Sleep(time.Duration(syncer.renLedgerSyncedInterval) * time.Second)
		}
	}()

	return errs
}

func (syncer *Syncer) Prune(done <-chan struct{})  {
	for {
		select {
		case <-done:
			return
		default:
			syncer.ordersMu.Lock()
			dispatch.CoForAll(syncer.orders, func(key int) {
				status := syncer.renLedger.CheckStatus(syncer.orders[key])
				if status == order.Canceled || status == order.Confirmed || status == order.Settled {
					delete(syncer.orders, key)
		syncer.queue.Remove(syncer.orders[key])
				}
			})
			syncer.ordersMu.Unlock()
		}

		time.Sleep(10 * time.Second)
	}
}

type orderPair struct {
	one      order.ID
	another  order.ID
	priority int
}

func NewOrderPair(one, another order.ID, priority int) orderPair {
	return orderPair{
		one:      one,
		another:  another,
		priority: priority,
	}
}

type PriorityQueue struct {
	mu *sync.Mutex
	pairs []orderPair
}

func NewPriorityQueue() PriorityQueue {
	return PriorityQueue{
		mu : new(sync.Mutex),
		pairs: []orderPair{},
	}
}

func (queue PriorityQueue) Insert(pair orderPair) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	index := sort.Search(len(queue.pairs), func(i int) bool {
		return queue.pairs[i].priority > pair.priority
	})
	queue.pairs = append(queue.pairs[:index-1], append([]orderPair{pair}, queue.pairs[index:]...)...)
}

func (queue PriorityQueue) Remove(id order.ID) {
	queue.mu.Lock()
	defer queue.mu.Unlock()

	for i := 0; i < len(queue.pairs); i++ {
		if queue.pairs[i].one == id || queue.pairs[i].another == id {
			// a:= []{1,2,3}
			// log.Println(a[3])   # this will panic
			// log.Println(a[3:])  # this will not panic ,amazing!
			queue.pairs = append(queue.pairs[:i], queue.pairs[i+1:]...)
		}
	}
}
