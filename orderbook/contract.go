package orderbook

import (
	"math/big"

	"github.com/republicprotocol/republic-go/order"
)

// ContractBinder for interacting with Ethereum contracts.
type ContractBinder interface {

	// Orders returns order.IDs that have been opened. The offset and limit
	// defines the range of order.IDs that can be returned, based on their
	// logical time ordering.
	Orders(offset, limit int) ([]order.ID, []order.Status, []string, error)

	// BlockNumber when the order.ID was opened.
	BlockNumber(orderID order.ID) (*big.Int, error)

	// Status of an order.ID.
	Status(orderID order.ID) (order.Status, error)

	// MinimumEpochInterval returns the minimum number of blocks between
	// epochs.
	MinimumEpochInterval() (*big.Int, error)

	Depth(orderID order.ID) (uint, error)
}
