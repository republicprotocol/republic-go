package ome

import (
	"github.com/republicprotocol/republic-go/order"
)

const StatusUndefined = 0
const StatusOpen = 1
const StatusConfirmed = 2
const StatusCanceled = 3

// ContractsBinder will define all interactions that the orderbook will
// have with the smart contracts
type ContractsBinder interface {

	// ConfirmOrder match on the Ren Ledger.
	ConfirmOrder(buy order.ID, sell order.ID) error

	// Depth will return depth of confirmation blocks
	Depth(orderID order.ID) (uint, error)

	// Status will return the status of the order
	Status(orderID order.ID) (order.Status, error)

	// Settle the order pair which gets confirmed by the RenLedger
	Settle(buy order.Order, sell order.Order) error

	// OrderMatch of an order, if any.
	OrderMatch(order order.ID) (order.ID, error)
}
