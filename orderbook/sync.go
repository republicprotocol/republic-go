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

// A Priority is an unsigned integer representing logical time priority. The
// lower the number, the higher the priority.
type Priority uint64

// Change represents a change found by the Syncer. It stores all the relevant
// information for the order.Order that was changed.
type Change struct {
	OrderID       order.ID
	OrderParity   order.Parity
	OrderStatus   order.Status
	OrderPriority Priority
	Trader        string
	BlockNumber   uint
}

// NewChange returns a Change object with the respective data stored inside it.
func NewChange(id order.ID, parity order.Parity, status order.Status, priority Priority, trader string, blockNumber uint) Change {
	return Change{
		OrderID:       id,
		OrderParity:   parity,
		OrderStatus:   status,
		OrderPriority: priority,
		Trader:        trader,
		BlockNumber:   blockNumber,
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

	syncBuyPointer  int
	syncSellPointer int

	ordersMu   *sync.RWMutex
	buyOrders  map[int]order.ID
	sellOrders map[int]order.ID
}

// NewSyncer returns a new Syncer that will sync a bounded number of orders
// from a cal.RenLedger. It uses a SyncStorer to prevent re-syncing the entire
// cal.RenLedger when it reboots.
func NewSyncer(renLedger cal.RenLedger, renLedgerLimit int) Syncer {

	syncer := &syncer{
		renLedger:      renLedger,
		renLedgerLimit: renLedgerLimit,

		syncBuyPointer:  0,
		syncSellPointer: 0,

		ordersMu:   new(sync.RWMutex),
		buyOrders:  map[int]order.ID{},
		sellOrders: map[int]order.ID{},
	}

	logger.Info(fmt.Sprintf("buy pointer: %v", syncer.syncBuyPointer))
	logger.Info(fmt.Sprintf("sell pointer: %v", syncer.syncSellPointer))

	return syncer
}

// Sync implements the Syncer interface.
func (syncer *syncer) Sync() (ChangeSet, error) {
	changeset := syncer.purge()

	buyOrderIDs, buyErr := syncer.renLedger.BuyOrders(syncer.syncBuyPointer, syncer.renLedgerLimit)
	if buyErr == nil {
		for _, ord := range buyOrderIDs {
			depth, err := syncer.renLedger.Depth(ord)
			if err == nil && depth > 6000 {
				syncer.syncBuyPointer++
				continue
			}
			status, err := syncer.renLedger.Status(ord)
			if err != nil {
				log.Println("cannot sync order status", err)
				continue
			}
			blockNumber, err := syncer.renLedger.BlockNumber(ord)
			if err != nil {
				log.Println("cannot sync order blocknumber", err)
				continue
			}
			trader, err := syncer.renLedger.Trader(ord)
			if err != nil {
				log.Println("cannot sync order owner", err)
				continue
			}

			syncer.syncBuyPointer++
			change := NewChange(ord, order.ParityBuy, status, Priority(syncer.syncBuyPointer), trader, blockNumber)
			changeset = append(changeset, change)
			syncer.buyOrders[syncer.syncBuyPointer] = ord
		}
	}

	// Get new sell orders from the ledger
	sellOrderIDs, sellErr := syncer.renLedger.SellOrders(syncer.syncSellPointer, syncer.renLedgerLimit)
	if sellErr == nil {
		for _, ord := range sellOrderIDs {
			depth, err := syncer.renLedger.Depth(ord)
			if err == nil && depth > 6000 {
				syncer.syncSellPointer++
				continue
			}
			status, err := syncer.renLedger.Status(ord)
			if err != nil {
				log.Println("cannot sync order status", err)
				continue
			}
			blockNumber, err := syncer.renLedger.BlockNumber(ord)
			if err != nil {
				log.Println("cannot sync order blocknumber", err)
				continue
			}
			trader, err := syncer.renLedger.Trader(ord)
			if err != nil {
				log.Println("cannot sync order owner", err)
				continue
			}

			syncer.syncSellPointer++
			change := NewChange(ord, order.ParitySell, status, Priority(syncer.syncSellPointer), trader, blockNumber)
			changeset = append(changeset, change)
			syncer.sellOrders[syncer.syncSellPointer] = ord
		}
	}

	logger.Info(fmt.Sprintf("updated buy pointer: %v", syncer.syncBuyPointer))
	logger.Info(fmt.Sprintf("updated sell pointer: %v", syncer.syncSellPointer))

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
				dispatch.ForAll(syncer.buyOrders, func(key int) {
					syncer.ordersMu.RLock()
					buyOrder := syncer.buyOrders[key]
					syncer.ordersMu.RUnlock()

					status, err := syncer.renLedger.Status(buyOrder)
					if err != nil {
						logger.Error(fmt.Sprintf("failed to check order status %v", err))
						return
					}
					if status == order.Open {
						return
					}

					blockNumber, err := syncer.renLedger.BlockNumber(buyOrder)
					if err != nil {
						log.Println("cannot sync order status", err)
						return
					}
					priority, err := syncer.renLedger.Priority(buyOrder)
					if err != nil {
						logger.Error(fmt.Sprintf("failed to check order priority %v", err))
						return
					}
					trader, err := syncer.renLedger.Trader(buyOrder)
					if err != nil {
						log.Println("cannot sync order owner", err)
						return
					}

					changes <- NewChange(buyOrder, order.ParityBuy, status, Priority(priority), trader, blockNumber)

					syncer.ordersMu.Lock()
					delete(syncer.buyOrders, key)
					syncer.ordersMu.Unlock()
				})
			},
			func() {
				// Purge all sell orders
				dispatch.ForAll(syncer.sellOrders, func(key int) {
					syncer.ordersMu.RLock()
					sellOrder := syncer.sellOrders[key]
					syncer.ordersMu.RUnlock()

					status, err := syncer.renLedger.Status(sellOrder)
					if err != nil {
						logger.Error(fmt.Sprintf("failed to check order status: %v", err))
						return
					}
					if status == order.Open {
						return
					}

					blockNumber, err := syncer.renLedger.BlockNumber(sellOrder)
					if err != nil {
						log.Println("cannot sync order status", err)
						return
					}
					priority, err := syncer.renLedger.Priority(sellOrder)
					if err != nil {
						logger.Error(fmt.Sprintf("failed to check order priority: %v", err))
						return
					}
					trader, err := syncer.renLedger.Trader(sellOrder)
					if err != nil {
						log.Println("cannot sync order owner", err)
						return
					}

					changes <- NewChange(sellOrder, order.ParitySell, status, Priority(priority), trader, blockNumber)

					syncer.ordersMu.Lock()
					delete(syncer.sellOrders, key)
					syncer.ordersMu.Unlock()
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
