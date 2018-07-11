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
func NewSyncer(numberOfOrders int) *Syncer {
	return &Syncer{
		hasSynced:       false,
		numberOfMatches: 0,
		orders:          make([]order.Order, 0, numberOfOrders),
	}
}

// Sync returns the first 5 orders in the order list and the syncer as synced.
func (syncer *Syncer) Sync() (<-chan orderbook.Notification, error) {
	var notificationChan <-chan orderbook.Notification
	if !syncer.hasSynced {
		// changes := make(orderbook.ChangeSet, len(syncer.orders))
		// i := 0
		// for _, ord := range syncer.orders {
		// changes[i] = orderbook.Change{
		// 	OrderID:       ord.ID,
		// 	OrderParity:   ord.Parity,
		// 	OrderPriority: orderbook.Priority(i),
		// 	OrderStatus:   order.Open,
		// }
		// i++
		// }
		// syncer.hasSynced = true
		return notificationChan, nil
	}

	return notificationChan, nil
}

// ConfirmOrderMatch confirms the two orders are a match and increment the
// match counter by 1.
func (syncer *Syncer) ConfirmOrderMatch(order.ID, order.ID) error {
	syncer.numberOfMatches++
	return nil
}

func (syncer *Syncer) HasSynced() bool {
	return syncer.hasSynced
}
