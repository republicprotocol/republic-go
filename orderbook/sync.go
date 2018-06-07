package orderbook

import (
	"fmt"
	"log"
	"sync"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
)

// ChangeSet is an alias type.
type ChangeSet []Change

// Change represents a change found by the Syncer. It stores all the relevant
// information for the order.Order that was changed.
type Change struct {
	OrderID       order.ID
	OrderParity   order.Parity
	OrderStatus   order.Status
	OrderPriority uint64
}

// NewChange returns a Change object with the respective data stored inside it.
func NewChange(id order.ID, parity order.Parity, status order.Status, priority uint64) Change {
	return Change{
		OrderID:       id,
		OrderParity:   parity,
		OrderStatus:   status,
		OrderPriority: priority,
	}
}

// A Syncer is used to synchronize orders, and changes to orders, to local
// storage.
type Syncer interface {

	// Sync orders and order states from the Ren Ledger to this local
	// Orderbooker. Returns a list of changes that were made to this local
	// Orderbooker during the synchronization.
	Sync() (ChangeSet, error)
}

type syncer struct {
	renLedger      cal.RenLedger
	renLedgerLimit int

	syncStorer      SyncStorer
	syncBuyPointer  int
	syncSellPointer int

	ordersMu   *sync.RWMutex
	buyOrders  map[int]order.ID
	sellOrders map[int]order.ID
}

// NewSyncer returns a new Syncer that will sync a bounded number of orders
// from a cal.RenLedger. It uses a SyncStorer to prevent re-syncing the entire
// cal.RenLedger when it reboots.
func NewSyncer(syncStorer SyncStorer, renLedger cal.RenLedger, renLedgerLimit int) Syncer {
	syncer := &syncer{
		renLedger:      renLedger,
		renLedgerLimit: renLedgerLimit,

		syncStorer:      syncStorer,
		syncBuyPointer:  0,
		syncSellPointer: 0,

		ordersMu:   new(sync.RWMutex),
		buyOrders:  map[int]order.ID{},
		sellOrders: map[int]order.ID{},
	}

	var err error
	if syncer.syncBuyPointer, err = syncer.syncStorer.BuyPointer(); err != nil {
		logger.Error(fmt.Sprintf("cannot load buy pointer: %v", err))
	}
	if syncer.syncSellPointer, err = syncer.syncStorer.SellPointer(); err != nil {
		logger.Error(fmt.Sprintf("cannot load sell pointer: %v", err))
	}

	return syncer
}

// Sync implements the Syncer interface.
func (syncer *syncer) Sync() (ChangeSet, error) {
	changeset := syncer.purge()

	buyOrderIDs, buyErr := syncer.renLedger.BuyOrders(syncer.syncBuyPointer, syncer.renLedgerLimit)
	if buyErr == nil {
		syncer.syncBuyPointer += len(buyOrderIDs)
		for i, ord := range buyOrderIDs {

			status, err := syncer.renLedger.Status(ord)
			if err != nil {
				log.Println("cannot sync order status", err)
				continue
			}

			change := NewChange(ord, order.ParityBuy, status, uint64(syncer.syncBuyPointer+i))
			changeset = append(changeset, change)
			syncer.buyOrders[syncer.syncBuyPointer+i] = ord
		}
		if err := syncer.syncStorer.InsertBuyPointer(syncer.syncBuyPointer); err != nil {
			logger.Error("cannot insert buy pointer")
		}
	}

	// Get new sell orders from the ledger
	sellOrderIDs, sellErr := syncer.renLedger.SellOrders(syncer.syncSellPointer, syncer.renLedgerLimit)
	if sellErr == nil {
		syncer.syncSellPointer += len(sellOrderIDs)
		for i, ord := range sellOrderIDs {

			status, err := syncer.renLedger.Status(ord)
			if err != nil {
				log.Println("cannot sync order status", err)
				continue
			}

			change := NewChange(ord, order.ParitySell, status, uint64(syncer.syncSellPointer+i))
			changeset = append(changeset, change)
			syncer.sellOrders[syncer.syncSellPointer+i] = ord
		}
		if err := syncer.syncStorer.InsertSellPointer(syncer.syncSellPointer); err != nil {
			logger.Error("cannot insert sell pointer")
		}
	}
	if buyErr != nil && sellErr != nil {
		return changeset, fmt.Errorf("buy err = %v, sell err = %v", buyErr, sellErr)
	}
	return changeset, nil
}

func (syncer *syncer) purge() ChangeSet {
	changes := make(chan Change, 128)

	go func() {
		defer close(changes)

		dispatch.CoBegin(
			func() {
				// Purge all buy orders by iterating over them and reading
				// their status and priority from the Ren Ledger
				dispatch.CoForAll(syncer.buyOrders, func(key int) {
					syncer.ordersMu.RLock()
					buyOrder := syncer.buyOrders[key]
					syncer.ordersMu.RUnlock()

					status, err := syncer.renLedger.Status(buyOrder)
					if err != nil {
						logger.Error(fmt.Sprintf("Failed to check order status %v", err))
						return
					}

					priority, err := syncer.renLedger.Priority(buyOrder)
					if err != nil {
						logger.Error(fmt.Sprintf("Failed to check order priority %v", err))
						return
					}

					if status != order.Open {
						changes <- NewChange(buyOrder, order.ParityBuy, status, priority)

						syncer.ordersMu.Lock()
						delete(syncer.buyOrders, key)
						syncer.ordersMu.Unlock()
					}
				})
			},
			func() {
				// Purge all sell orders
				dispatch.CoForAll(syncer.sellOrders, func(key int) {
					syncer.ordersMu.RLock()
					sellOrder := syncer.sellOrders[key]
					syncer.ordersMu.RUnlock()

					status, err := syncer.renLedger.Status(sellOrder)
					if err != nil {
						logger.Error(fmt.Sprintf("Failed to check order status: %v", err))
						return
					}

					priority, err := syncer.renLedger.Priority(sellOrder)
					if err != nil {
						logger.Error(fmt.Sprintf("Failed to check order priority: %v", err))
						return
					}

					if status != order.Open {
						changes <- NewChange(sellOrder, order.ParitySell, status, priority)

						syncer.ordersMu.Lock()
						delete(syncer.sellOrders, key)
						syncer.ordersMu.Unlock()
					}
				})
			},
		)
	}()

	changeset := make([]Change, 0, 128)
	for change := range changes {
		changeset = append(changeset, change)
	}
	return changeset
}
