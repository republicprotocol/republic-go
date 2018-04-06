package orderbook

import (
	"github.com/republicprotocol/republic-go/order"
)

// Message define status of the order. Along with error happens.
type Message struct {
	EpochHash [32]byte
	Ord       order.Order
	Status    order.Status
}

// NewMessage returns a new orderbook message.
func NewMessage(ord order.Order, status order.Status, hash [32]byte) *Message {
	//var epochHash [32]byte
	//if len(hash) != 32 {
	//	log.Println("wrong epoch hash length")
	//	return Message{}
	//}
	//copy(epochHash[:], hash)
	return &Message{
		EpochHash: hash,
		Ord:       ord,
		Status:    status,
	}
}
