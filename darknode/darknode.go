package darknode

import (
	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/order"
)

type Darknode struct {
	darkpool         cal.Darkpool
	darkpoolAccounts cal.DarkpoolAccounts
	darkpoolFees     cal.DarkpoolFees
}

func NewDarknode(darkpool cal.Darkpool, darkpoolAccounts cal.DarkpoolAccounts, darkpoolFees cal.DarkpoolFees) Darknode {
	return Darknode{
		darkpool:         darkpool,
		darkpoolAccounts: darkpoolAccounts,
		darkpoolFees:     darkpoolFees,
	}
}

// ConfirmOrderMatch implements the ome.Observer interface. A call to
// Darknode.ConfirmOrderMatch happens whenever orders have been matched and
// consensus has been reached on the finality of the match.
func (node *Darknode) ConfirmOrderMatch(orders []order.Order) {
}
