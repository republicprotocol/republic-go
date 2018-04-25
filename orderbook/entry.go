package orderbook

import (
	"github.com/republicprotocol/republic-go/order"
)

// Entry of an order into the Orderbook, including the epoch hash at which this
// order was discovered.
type Entry struct {
	order.Order
	order.Status
}

// NewEntry returns a new Entry.
func NewEntry(order order.Order, status order.Status) Entry {
	return Entry{
		Order:  order,
		Status: status,
	}
}
