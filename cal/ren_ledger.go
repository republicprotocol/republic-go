package cal

import (
	"math/big"

	"github.com/republicprotocol/republic-go/order"
)

const StatusUndefined = 0
const StatusOpen = 1
const StatusConfirmed = 2
const StatusCanceled = 3

type RenLedger interface {

	// OpenBuyOrder on the Ren Ledger. The signature will be used to identify
	// the trader that owns the order. The order must be in an undefined state
	// to be opened.
	OpenBuyOrder(signature [65]byte, orderID order.ID) error

	// OpenSellOrder on the Ren Ledger. The signature will be used to identify
	// the trader that owns the order. The order must be in an undefined state
	// to be opened.
	OpenSellOrder(signature [65]byte, orderID order.ID) error

	// CancelOrder on the Ren Ledger. The signature will be used to verify that
	// the request was created by the trader that owns the order. The order
	// must be in the opened state to be canceled.
	CancelOrder(signature [65]byte, orderID order.ID) error

	// ConfirmOrder match on the Ren Ledger.
	ConfirmOrder(buy order.ID, sell order.ID) error

	// OrderMatch of an order, if any.
	OrderMatch(order order.ID) (order.ID, error)

	// Fee required to open an order.
	Fee() (*big.Int, error)

	// Status will return the status of the order
	Status(orderID order.ID) (order.Status, error)

	// Priority will return the priority of the order
	Priority(orderID order.ID) (uint64, error)

	// Depth will return depth of confirmation blocks
	Depth(orderID order.ID) (uint, error)

	// BuyOrders in the Ren Ledger starting at an offset and returning limited
	// numbers of  buy orders.
	BuyOrders(offset, limit int) ([]order.ID, error)

	// SellOrders in the Ren Ledger starting at an offset and returning limited
	// numbers of  sell orders.
	SellOrders(offset, limit int) ([]order.ID, error)

	// Trader returns the trader who submit the order
	Trader(orderID order.ID) (string, error)
}
