package orderbook

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
)

// ErrOrderNotFound is return when attempting to load an order that cannot be
// found.
var ErrOrderNotFound = errors.New("order not found")

// ContractBinder will define all methods that the orderbook will
// require to communicate with smart contracts. All the methods will
// be implemented in contract.Binder
type ContractBinder interface {
	BuyOrders(offset, limit int) ([]order.ID, error)

	SellOrders(offset, limit int) ([]order.ID, error)

	Status(orderID order.ID) (order.Status, error)

	BlockNumber(orderID order.ID) (uint, error)

	Trader(orderID order.ID) (string, error)

	Priority(orderID order.ID) (uint64, error)

	Depth(orderID order.ID) (uint64, error)
}
