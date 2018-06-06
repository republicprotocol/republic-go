package orderbook

import (
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
)

type ChangeSet []Change

type Change struct {
	OrderID       order.ID
	OrderParity   order.Parity
	OrderStatus   order.Status
	OrderPriority uint64
}

func NewChange(id order.ID, parity order.Parity, status order.Status, priority uint64) Change {
	return Change{
		OrderID:       id,
		OrderParity:   parity,
		OrderStatus:   status,
		OrderPriority: priority,
	}
}

type Syncer interface {

	// Sync orders and order states from the Ren Ledger to this local
	// Orderbooker. Returns a list of changes that were made to this local
	// Orderbooker during the synchronization.
	Sync() (ChangeSet, error)
}

type syncer struct {
	renLedger        cal.RenLedger
	renLedgerLimit   int
	buyOrderPointer  int
	sellOrderPointer int

	ordersMu   *sync.RWMutex
	buyOrders  map[int]order.ID
	sellOrders map[int]order.ID
}

func NewSyncer(renLedger cal.RenLedger, limit int) Syncer {
	return &syncer{
		renLedger:        renLedger,
		renLedgerLimit:   limit,
		buyOrderPointer:  0,
		sellOrderPointer: 0,
		ordersMu:         new(sync.RWMutex),
		buyOrders:        map[int]order.ID{},
		sellOrders:       map[int]order.ID{},
	}
}

func (syncer *syncer) Sync() (ChangeSet, error) {
	changeset := syncer.purge()

	buyOrderIDs, buyErr := syncer.renLedger.BuyOrders(syncer.buyOrderPointer, syncer.renLedgerLimit)
	if buyErr == nil {
		syncer.buyOrderPointer += len(buyOrderIDs)
		for i, ord := range buyOrderIDs {

			status, err := syncer.renLedger.Status(ord)
			if err != nil {
				log.Println("cannot sync order status", err)
				continue
			}

			change := NewChange(ord, order.ParityBuy, status, uint64(syncer.buyOrderPointer+i))
			changeset = append(changeset, change)
			syncer.buyOrders[syncer.buyOrderPointer+i] = ord
		}
	}

	// Get new sell orders from the ledger
	sellOrderIDs, sellErr := syncer.renLedger.SellOrders(syncer.sellOrderPointer, syncer.renLedgerLimit)
	if sellErr == nil {
		syncer.sellOrderPointer += len(sellOrderIDs)
		for i, ord := range sellOrderIDs {

			status, err := syncer.renLedger.Status(ord)
			if err != nil {
				log.Println("cannot sync order status", err)
				continue
			}

			change := NewChange(ord, order.ParitySell, status, uint64(syncer.sellOrderPointer+i))
			changeset = append(changeset, change)
			syncer.sellOrders[syncer.sellOrderPointer+i] = ord
		}
	}
	if buyErr != nil && sellErr != nil {
		return changeset, fmt.Errorf("buy err = %v, sell err = %v", buyErr, sellErr)
	}
	return changeset, nil
}

func (syncer *syncer) purge() ChangeSet {
	changes := make(chan Change, 100)

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

	changeset := make([]Change, 0, 100)
	for change := range changes {
		changeset = append(changeset, change)
	}
	return changeset
}
