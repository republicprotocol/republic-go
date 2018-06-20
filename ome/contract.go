package ome

import (
	"github.com/republicprotocol/republic-go/order"
)

const StatusUndefined = 0
const StatusOpen = 1
const StatusConfirmed = 2
const StatusCanceled = 3

// ContractBinder will define all methods that the order matching
// engine will require to communicate with smart contracts. All the
// methods will be implemented in contract.Binder
type ContractBinder interface {
	ConfirmOrder(buy order.ID, sell order.ID) error

	Depth(orderID order.ID) (uint, error)

	Status(orderID order.ID) (order.Status, error)

	OrderMatch(order order.ID) (order.ID, error)

	Settle(buy order.Order, sell order.Order) error
}
