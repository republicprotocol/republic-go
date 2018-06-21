package orderbook

import (
	"fmt"

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
	storer SyncStorer

	renLedger      cal.RenLedger
	renLedgerLimit int
}

// NewSyncer returns a new Syncer that will sync a bounded number of orders
// from a cal.RenLedger. It uses a SyncStorer to prevent re-syncing the entire
// cal.RenLedger when it reboots.
func NewSyncer(storer SyncStorer, renLedger cal.RenLedger, renLedgerLimit int) Syncer {
	return &syncer{
		storer: storer,

		renLedger:      renLedger,
		renLedgerLimit: renLedgerLimit,
	}
}

// Sync implements the Syncer interface.
func (syncer *syncer) Sync() (ChangeSet, error) {
	changeset := syncer.purge()

	buyPointer, err := syncer.storer.BuyPointer()
	if err != nil {
		return changeset, err
	}
	sellPointer, err := syncer.storer.SellPointer()
	if err != nil {
		return changeset, err
	}

	buyOrderIDs, buyErr := syncer.renLedger.BuyOrders(int(buyPointer), syncer.renLedgerLimit)
	if buyErr == nil {
		for _, ord := range buyOrderIDs {
			status, err := syncer.renLedger.Status(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order status: %v", err))
				buyErr = err
				continue
			}
			blockNumber, err := syncer.renLedger.BlockNumber(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order block: %v", err))
				buyErr = err
				continue
			}
			trader, err := syncer.renLedger.Trader(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order trader: %v", err))
				buyErr = err
				continue
			}

			buyPointer++
			change := NewChange(ord, order.ParityBuy, status, Priority(buyPointer), trader, blockNumber)
			changeset = append(changeset, change)
			if err := syncer.storer.PutChange(change); err != nil {
				logger.Error(fmt.Sprintf("cannot store synchronised order: %v", err))
			}
		}
		if buyErr == nil {
			syncer.storer.PutBuyPointer(buyPointer)
		}
	}

	// Get new sell orders from the ledger
	sellOrderIDs, sellErr := syncer.renLedger.SellOrders(int(sellPointer), syncer.renLedgerLimit)
	if sellErr == nil {
		for _, ord := range sellOrderIDs {

			status, err := syncer.renLedger.Status(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order status: %v", err))
				sellErr = err
				continue
			}
			blockNumber, err := syncer.renLedger.BlockNumber(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order block: %v", err))
				sellErr = err
				continue
			}
			trader, err := syncer.renLedger.Trader(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order trader: %v", err))
				sellErr = err
				continue
			}

			sellPointer++
			change := NewChange(ord, order.ParitySell, status, Priority(sellPointer), trader, blockNumber)
			changeset = append(changeset, change)
			if err := syncer.storer.PutChange(change); err != nil {
				logger.Error(fmt.Sprintf("cannot store synchronised order: %v", err))
			}
		}
		if sellErr == nil {
			syncer.storer.PutSellPointer(sellPointer)
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

		changesIter, err := syncer.storer.Changes()
		if err != nil {
			logger.Error(fmt.Sprintf("cannot build changes iterator for purging: %v", err))
			return
		}
		changesCollection, err := changesIter.Collect()
		if err != nil {
			logger.Error(fmt.Sprintf("cannot build changes collection for purging: %v", err))
			return
		}

		dispatch.ForAll(changesCollection, func(i int) {
			change := changesCollection[i]

			status, err := syncer.renLedger.Status(change.OrderID)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync change status: %v", err))
				return
			}
			if status == order.Open {
				return
			}

			blockNumber, err := syncer.renLedger.BlockNumber(change.OrderID)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync change block: %v", err))
				return
			}

			change.OrderStatus = status
			change.BlockNumber = blockNumber
			changes <- change

			if err := syncer.storer.DeleteChange(change.OrderID); err != nil {
				logger.Error(fmt.Sprintf("cannot delete synchronised change: %v", err))
			}
		})
	}()

	changeset := make([]Change, 0, 128)
	for change := range changes {
		changeset = append(changeset, change)
	}
	return changeset
}
