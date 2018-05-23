package ome

import (
	"fmt"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/smpc"
)

type Ome interface {
	cal.EpochListener

	// Sync the Omer with the orderbook.Orderbooker so that it can discover new
	// orders, purge confirmed orders, and reprioritize order matching
	// computations.
	Sync() error
}

type ome struct {
	ranker   Ranker
	computer Computer

	ξ         cal.Epoch
	orderbook orderbook.Orderbook
	smpcer    smpc.Smpcer
}

func NewOme(ranker Ranker, computer Computer, orderbook orderbook.Orderbook, smpcer smpc.Smpcer) Ome {
	return &ome{
		ranker:   ranker,
		computer: computer,

		ξ:         cal.Epoch{},
		orderbook: orderbook,
		smpcer:    smpcer,
	}
}

// OnChangeEpoch implements the cal.EpochListener interface.
func (ome *ome) OnChangeEpoch(ξ cal.Epoch) {
	ome.ξ = ξ
}

// Sync implements the Omer interface.
func (ome *ome) Sync() error {
	changeset, err := ome.orderbook.Sync()
	if err != nil {
		return fmt.Errorf("cannot sync orderbook: %v", err)
	}
	if err := ome.syncRanker(changeset); err != nil {
		return fmt.Errorf("cannot sync ranker: %v", err)
	}
	if err := ome.syncComputer(changeset); err != nil {
		return fmt.Errorf("cannot sync computer: %v", err)
	}
	return nil
}

func (ome *ome) syncRanker(changeset orderbook.ChangeSet) error {
	for _, change := range changeset {
		switch change.OrderStatus {
		case order.Open:
			if change.OrderParity == order.ParityBuy {
				ome.ranker.InsertBuy(PriorityOrder{
					ID:       change.OrderID,
					Priority: Priority(change.OrderPriority),
				})
			} else {
				ome.ranker.InsertSell(PriorityOrder{
					ID:       change.OrderID,
					Priority: Priority(change.OrderPriority),
				})
			}
		case order.Canceled, order.Settled, order.Confirmed:
			ome.ranker.Remove(change.OrderID)
		}
	}
	return nil
}

func (ome *ome) syncComputer(changeset orderbook.ChangeSet) error {
	buffer := [128]Computation{}
	n := ome.ranker.Computations(buffer[:])
	ome.computer.Compute(ome.ξ.Hash, buffer[:n])
	return nil
}
