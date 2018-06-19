package orderbook

import "github.com/republicprotocol/republic-go/order"

// ContractsBinder will define all interactions that the orderbook will
// have with the smart contracts
type ContractsBinder interface {

	// BuyOrders in the Ren Ledger starting at an offset and returning limited
	// numbers of buy orders.
	BuyOrders(offset, limit int) ([]order.ID, error)

	// SellOrders in the Ren Ledger starting at an offset and returning limited
	// numbers of sell orders.
	SellOrders(offset, limit int) ([]order.ID, error)

	// Status will return the status of the order
	Status(orderID order.ID) (order.Status, error)

	// BlockNumber will return the block number when the order status
	// last mode modified
	BlockNumber(orderID order.ID) (uint, error)

	// Trader returns the trader who submits the order
	Trader(orderID order.ID) (string, error)

	// Priority will return the priority of the order
	Priority(orderID order.ID) (uint64, error)
}
