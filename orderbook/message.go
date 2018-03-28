package orderbook

import (
	"github.com/republicprotocol/republic-go/order"
)

type Message struct {
	Err error

	Ord    order.Order
	Status order.Status
}

func NewMessage(err error, ord order.Order, status order.Status) Message {
	return Message{
		Err:    err,
		Ord:    ord,
		Status: status,
	}
}
