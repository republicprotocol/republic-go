package ome

import (
	"context"
	"log"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc"
)

type Syncer interface {
	Address() identity.Address
	SyncRenLedger(done <-chan struct{}, pod <-smpc.InstConnect,  ranker Ranker) error
	ConfirmOrder(id order.ID, matches []order.ID) error
	Threshold() int
}

type OrderWithPriority struct {
	ID       order.ID
	Priority uint64
}

type syncer struct {
	address                 identity.Address
	pool                    cal.Darkpool
	renLedger               cal.RenLedger
	renLedgerLimit          int
	renLedgerSyncedInterval int
	buyOrderPointer         int
	sellOrderPointer        int
	buyOrders               map[int]order.ID
	sellOrders              map[int]order.ID
}

func NewSyncer(address identity.Address, renLedger cal.RenLedger, limit, interval int) syncer {
	return syncer{
		address:                 address,
		renLedger:               renLedger,
		renLedgerLimit:          limit,
		renLedgerSyncedInterval: interval,
		buyOrderPointer:         0,
		sellOrderPointer:        0,
		buyOrders:               map[int]order.ID{},
		sellOrders:              map[int]order.ID{},
	}
}

func (syncer *syncer) Address() identity.Address {
	return syncer.address
}

func (syncer *syncer) SyncRenLedger(done <-chan struct{}, ranker Ranker) error {

	pods, err := syncer.pool.Pods()
	if err != nil {
		return
	}
	poolIndex, _, err := syncer.pool.Pool(syncer.address.ID())
	if err != nil {
		return
	}
	for {
		select {
		case <- ctx.Done():
			return
		default:
			syncer.Prune(ranker)
		}

		// For each new buy order, create order pairs with the sell order
		orderIDs, err := syncer.renLedger.BuyOrders(syncer.buyOrderPointer, syncer.renLedgerLimit)
		if err != nil {
			errs <- err
			return
		}
		syncer.buyOrderPointer += len(orderIDs)

		for index, id := range orderIDs {
			for m, n := range syncer.sellOrders {
				if (index+syncer.buyOrderPointer+m)%len(pods) == poolIndex {
					orderPair := NewOrderPair(id, []order.ID{n}, index+syncer.buyOrderPointer+m)
					ranker.Insert(orderPair)
				}
			}
		}

		// For each new sell order, create order pairs with the buy order
		orderIDs, err = syncer.renLedger.SellOrders(syncer.sellOrderPointer, syncer.renLedgerLimit)
		if err != nil {
			errs <- err
			return
		}
		syncer.sellOrderPointer += len(orderIDs)

		for index, id := range orderIDs {
			for m, n := range syncer.buyOrders {
				if (index+syncer.sellOrderPointer+m)%len(pods) == poolIndex {
					orderPair := NewOrderPair(id, []order.ID{n}, index+syncer.sellOrderPointer+m)
					ranker.Insert(orderPair)
				}
			}
		}

		time.Sleep(time.Duration(syncer.renLedgerSyncedInterval) * time.Second)
	}


	return errs
}

func (syncer *syncer) ConfirmOrder(id order.ID, matches []order.ID) error {
	return syncer.renLedger.ConfirmOrder(id, matches)
}

func (syncer *syncer) Prune(ranker Ranker) {
	dispatch.Dispatch(
		func() {
			dispatch.CoForAll(syncer.sellOrders, func(key int) {
				status, err := syncer.renLedger.Status(syncer.buyOrders[key])
				if err != nil {
					log.Println("fail to check order status", err)
					return
				}

				if status == order.Canceled || status == order.Confirmed || status == order.Settled {
					delete(syncer.buyOrders, key)
					ranker.Remove(syncer.buyOrders[key])
				}
			})
		},
		func() {
			dispatch.CoForAll(syncer.sellOrders, func(key int) {
				status, err := syncer.renLedger.Status(syncer.sellOrders[key])
				if err != nil {
					log.Println("fail to check order status", err)
					return
				}

				if status == order.Canceled || status == order.Confirmed || status == order.Settled {
					delete(syncer.sellOrders, key)
					ranker.Remove(syncer.sellOrders[key])
				}
			})
		},
	)
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
