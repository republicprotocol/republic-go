package orderbook

import (
	"github.com/republicprotocol/republic-go/order"
)

// Update reflects an status update for an order.
type Update struct {
	order.ID
	order.Status
}

// NewUpdate creates a new Update which includes the orderID and
// its new status.
func NewUpdate(id order.ID, status order.Status) Update{
	return Update{
		ID : id ,
		Status : status,
	}
}