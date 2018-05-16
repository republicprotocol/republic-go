package cal

import (
	"math/big"

	"github.com/republicprotocol/republic-go/order"
)

type RenLedger interface {

	// OpenOrder on the Ren Ledger. The signature will be used to identify the
	// trader that owns the order. The order must be in an undefined state to
	// be opened.
	OpenOrder(signature [65]byte, orderID order.ID) error

	// CancelOrder on the Ren Ledger. The signature will be used to verify that
	// the request was created by the trader that owns the order. The order
	// must be in the opened state to be canceled.
	CancelOrder(signature [65]byte, orderID order.ID) error

	// ConfirmOrder match on the Ren Ledger. Both the id and the matches should
	// be in the opened state to be confirmed.
	ConfirmOrder(id order.ID, matches []order.ID) error

	// Fee required to open an order.
	Fee() (*big.Int, error)

	// Status will return the status of the order
	Status(orderID order.ID) (order.Status, error)

	// Depth will return depth of confirmation blocks
	Depth(orderID order.ID) (uint, error)

	// Orders in the Ren Ledger starting at an offset and returning limited
	// numbers of orders.
	Orders(offset, limit int) ([]order.ID, error)
}
