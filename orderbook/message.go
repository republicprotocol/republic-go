package orderbook

import (
	"github.com/republicprotocol/republic-go/order"
)

// Message define status of the order. Along with error happens.
type Message struct {
	Ord    order.Order
	Status order.Status
	Err    error
}

// NewMessage returns a new orderbook message.
func NewMessage(ord order.Order, status order.Status, err error) Message {
	return Message{
		Ord:    ord,
		Status: status,
		Err:    err,
	}
}

