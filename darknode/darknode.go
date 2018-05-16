package darknode

import (
	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/order"
)

type Darknode struct {
	addr             identity.Address
	engine           ome.Omer
	darkpool         cal.Darkpool
	darkpoolAccounts cal.DarkpoolAccounts
	darkpoolFees     cal.DarkpoolFees
	ξ                cal.Epoch
}

func NewDarknode(addr identity.Address, engine ome.Omer, darkpool cal.Darkpool, darkpoolAccounts cal.DarkpoolAccounts, darkpoolFees cal.DarkpoolFees) Darknode {
	return Darknode{
		addr:             addr,
		engine:           engine,
		darkpool:         darkpool,
		darkpoolAccounts: darkpoolAccounts,
		darkpoolFees:     darkpoolFees,
	}
}

// SyncDarkpool updates the Epoch and the Pod configuration to the Darknode
// using a cal.Darkpool. If the Epoch has changes, an event will be emitted to
// the ome.Omer so that it can update its internal configuration (specifically,
// the smpc.Smpcer configuration).
func (node *Darknode) SyncDarkpool() error {
	// Sync ξ from the Darkpool
	ξ, err := node.darkpool.Epoch()
	if err != nil {
		return err
	}

	// Check whether or not ξ has changed
	if ξ.Equal(&node.ξ) {
		return nil
	}
	node.ξ = ξ
	node.engine.OnChangeEpoch(node.ξ)
	return nil
}

// Address of the Darknode.
func (node *Darknode) Address() identity.Address {
	return node.addr
}

// OnConfirmOrder implements the ome.Delegate interface.
func (node *Darknode) OnConfirmOrder(order order.Order, matches []order.Order) error {
	return node.darkpoolAccounts.Settle(order, matches)
}
