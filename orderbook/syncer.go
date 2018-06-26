package orderbook

import (
	"bytes"
	"fmt"

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
	OrderID       order.ID     `json:"orderId"`
	OrderParity   order.Parity `json:"orderParity"`
	OrderStatus   order.Status `json:"orderStatus"`
	OrderPriority Priority     `json:"orderPriority"`
	Trader        string       `json:"trader"`
	BlockNumber   uint         `json:"blockNumber"`
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

// Equal returns an equality check between two Changes.
func (change *Change) Equal(other *Change) bool {
	return bytes.Equal(change.OrderID[:], other.OrderID[:]) &&
		change.OrderParity == other.OrderParity &&
		change.OrderStatus == other.OrderStatus &&
		change.OrderPriority == other.OrderPriority &&
		change.Trader == other.Trader &&
		change.BlockNumber == other.BlockNumber
}

// A Syncer is used to synchronize orders, and changes to orders, to local
// storage.
type Syncer interface {

	// Sync orders and order states from the Orderbook to this local
	// Orderbooker. Returns a list of changes that were made to this local
	// Orderbooker during the synchronization.
	Sync(done <-chan struct{}) (<-chan Notification, <-chan error)
}

type syncer struct {
	storer   SyncStorer
	contract ContractBinder
	limit    int
}

// NewSyncer returns a new Syncer that will sync a bounded number of orders
// from the ContractBinder. It uses a SyncStorer to prevent re-syncing the entire
// ContractBinder when it reboots.
func NewSyncer(storer SyncStorer, contract ContractBinder, limit int) Syncer {
	return &syncer{
		storer: storer,

		contract: contract,
		limit:    limit,
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

	buyOrderIDs, buyErr := syncer.contract.BuyOrders(int(buyPointer), syncer.limit)
	if buyErr == nil {
		for _, ord := range buyOrderIDs {
			depth, err := syncer.contract.Depth(ord)
			if err == nil && depth > 6000 {
				buyPointer++
				continue
			}
			status, err := syncer.contract.Status(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order status: %v", err))
				buyErr = err
				buyPointer++
				continue
			}
			blockNumber, err := syncer.contract.BlockNumber(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order block: %v", err))
				buyErr = err
				buyPointer++
				continue
			}
			trader, err := syncer.contract.Trader(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order trader: %v", err))
				buyErr = err
				buyPointer++
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
	sellOrderIDs, sellErr := syncer.contract.SellOrders(int(sellPointer), syncer.limit)
	if sellErr == nil {
		for _, ord := range sellOrderIDs {
			depth, err := syncer.contract.Depth(ord)
			if err == nil && depth > 6000 {
				sellPointer++
				continue
			}
			status, err := syncer.contract.Status(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order status: %v", err))
				sellErr = err
				sellPointer++
				continue
			}
			blockNumber, err := syncer.contract.BlockNumber(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order block: %v", err))
				sellErr = err
				sellPointer++
				continue
			}
			trader, err := syncer.contract.Trader(ord)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync order trader: %v", err))
				sellErr = err
				sellPointer++
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
		defer changesIter.Release()
		changesCollection, err := changesIter.Collect()
		if err != nil {
			logger.Error(fmt.Sprintf("cannot build changes collection for purging: %v", err))
			return
		}

		dispatch.ForAll(changesCollection, func(i int) {
			change := changesCollection[i]

			status, err := syncer.contract.Status(change.OrderID)
			if err != nil {
				logger.Error(fmt.Sprintf("cannot sync change status: %v", err))
				return
			}
			if status == order.Open {
				return
			}

			blockNumber, err := syncer.contract.BlockNumber(change.OrderID)
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
