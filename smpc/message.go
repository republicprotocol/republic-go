package smpc

import "github.com/republicprotocol/republic-go/order"

// A Message for Workers to receive and interpret.
type Message struct {
	Error error

	OrderFragment  *order.Fragment
	Delta          *Delta
	DeltaFragments DeltaFragments
}
