package orderbook

import "github.com/republicprotocol/republic-go/order"

// ContractsBinder will define all methods that the orderbook will
// require to communicate with smart contracts. All the methods will
// be implemented in contracts.Binder
type ContractsBinder interface {
	BuyOrders(offset, limit int) ([]order.ID, error)

	SellOrders(offset, limit int) ([]order.ID, error)

	Status(orderID order.ID) (order.Status, error)

	BlockNumber(orderID order.ID) (uint, error)

	Trader(orderID order.ID) (string, error)

	Priority(orderID order.ID) (uint64, error)
}