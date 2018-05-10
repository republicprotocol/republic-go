package cal

import (
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

type Pool struct {
	Hash      [32]byte
	Darknodes []identity.Address
}

type Darkpool interface {
	Pools() ([]Pool, error)
}

type RenLedger interface {
	OpenOrder(orderID order.ID, orderType order.Type, orderParity order.Parity, orderExpiry int64, signature [65]byte) error
	WaitForOpenOrder(orderID order.ID) error

	CancelOrder(orderID order.ID, signature [65]byte) error
	WaitForCancelOrder(orderID order.ID) error
}

type RenToken interface {
}
