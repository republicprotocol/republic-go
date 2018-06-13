package testutils

import (
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

// Syncer is a mock implementation of the orderbook.Syncer
type Syncer struct {
	hasSynced       bool
	numberOfMatches int
	orders          []order.Order
}

// NewSyncer returns a mock implementation of an orderbook.Syncer interface.
func NewSyncer(numberOfMatches int) orderbook.Syncer {
	return &Syncer{
		hasSynced:       false,
		numberOfMatches: numberOfMatches,
		orders:          []order.Order{},
	}
}

// Sync returns the first 5 orders in the order list and the syncer as synced.
func (syncer *Syncer) Sync() (orderbook.ChangeSet, error) {
	if !syncer.hasSynced {
		changes := make(orderbook.ChangeSet, 5)
		i := 0
		for _, ord := range syncer.orders {
			changes[i] = orderbook.Change{
				OrderID:       ord.ID,
				OrderParity:   ord.Parity,
				OrderPriority: uint64(i),
				OrderStatus:   order.Open,
			}
			i++
		}
		syncer.hasSynced = true
		return changes, nil
	}

	return orderbook.ChangeSet{}, nil
}

// ConfirmOrderMatch confirms the two orders are a match and increment the
// match counter by 1.
func (syncer *Syncer) ConfirmOrderMatch(order.ID, order.ID) error {
	syncer.numberOfMatches++
	return nil
}
