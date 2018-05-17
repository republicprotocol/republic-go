package orderbook

import (
	"log"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
)

type Syncer interface {
	Sync() ([]OrderUpdate, error)
}

type syncer struct {
	renLedger        cal.RenLedger
	renLedgerLimit   int
	buyOrderPointer  int
	sellOrderPointer int
	buyOrders        map[int]order.ID
	sellOrders       map[int]order.ID
}

func NewSyncer(renLedger cal.RenLedger, limit int) Syncer {
	return &syncer{
		renLedger:        renLedger,
		renLedgerLimit:   limit,
		buyOrderPointer:  0,
		sellOrderPointer: 0,
		buyOrders:        map[int]order.ID{},
		sellOrders:       map[int]order.ID{},
	}
}

func (syncer *syncer) Sync() ([]OrderUpdate, error) {
	// Prune the orders first
	updates := syncer.Prune()

	// Get new buy orders from the ledger
	buyOrderIDs, err := syncer.renLedger.BuyOrders(syncer.buyOrderPointer, syncer.renLedgerLimit)
	if err != nil {
		return nil, err
	}
	syncer.buyOrderPointer += len(buyOrderIDs)
	for i, ord := range buyOrderIDs {
		update := NewOrderChange(ord, order.ParityBuy, order.Open, uint64(syncer.buyOrderPointer+i))
		updates = append(updates, update)
	}
	// Get new sell orders from the ledger
	sellOrderIDs, err := syncer.renLedger.SellOrders(syncer.sellOrderPointer, syncer.renLedgerLimit)
	if err != nil {
		return nil, err
	}
	syncer.sellOrderPointer += len(sellOrderIDs)
	for i, ord := range sellOrderIDs {
		update := NewOrderChange(ord, order.ParitySell, order.Open, uint64(syncer.sellOrderPointer+i))
		updates = append(updates, update)
	}

	return updates, nil
}

func (syncer *syncer) Prune() []OrderUpdate {
	orderChanges := make(chan OrderUpdate, 100)
	defer close(orderChanges)

	go func() {
		dispatch.Dispatch(
			func() {
				dispatch.CoForAll(syncer.sellOrders, func(key int) {
					status, err := syncer.renLedger.Status(syncer.buyOrders[key])
					if err != nil {
						log.Println("fail to check order status", err)
						return
					}
					priority, err := syncer.renLedger.Priority(syncer.buyOrders[key])
					if err != nil {
						log.Println("fail to check order priority", err)
						return
					}
					if status != order.Open {
						update := NewOrderChange(syncer.buyOrders[key], order.ParityBuy, status, priority)

						orderChanges <- update
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
					priority, err := syncer.renLedger.Priority(syncer.sellOrders[key])
					if err != nil {
						log.Println("fail to check order priority", err)
						return
					}
					if status != order.Open {
						update := NewOrderChange(syncer.sellOrders[key], order.ParitySell, status, priority)
						orderChanges <- update
					}
				})
			},
		)
	}()

	updates := make([]OrderUpdate, 0)
	for update := range orderChanges {
		updates = append(updates, update)
	}

	return updates
}

type OrderUpdate struct {
	ID       order.ID
	Parity   order.Parity
	Status   order.Status
	Priority uint64
}

func NewOrderChange(id order.ID, parity order.Parity, status order.Status, priority uint64) OrderUpdate {
	return OrderUpdate{
		ID:       id,
		Parity:   parity,
		Status:   status,
		Priority: priority,
	}
}
