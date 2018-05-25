package darknode

import (
	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

type Darknode struct {
	addr identity.Address

	darkpool         cal.Darkpool
	darkpoolAccounts cal.DarkpoolAccounts
	darkpoolFees     cal.DarkpoolFees

	ξ         cal.Epoch
	ξListener cal.EpochListener
}

func NewDarknode(addr identity.Address, darkpool cal.Darkpool, darkpoolAccounts cal.DarkpoolAccounts, darkpoolFees cal.DarkpoolFees) Darknode {
	return Darknode{
		addr:             addr,
		darkpool:         darkpool,
		darkpoolAccounts: darkpoolAccounts,
		darkpoolFees:     darkpoolFees,
	}
}

// Sync updates the Epoch and the Pod configuration to the Darknode using a
// cal.Darkpool. If the Epoch has changes, an event will be emitted to the
// ome.Omer so that it can update its internal configuration (specifically, the
// smpc.Smpcer configuration).
func (node *Darknode) Sync() error {
	// Sync ξ from the Darkpool
	ξ, err := node.darkpool.Epoch()
	if err != nil {
		return err
	}

	// Check whether or not ξ has changed and emit an event
	if ξ.Equal(&node.ξ) {
		return nil
	}
	node.ξ = ξ
	node.ξListener.OnChangeEpoch(node.ξ)
	return nil
}

// SetEpochListener to be notified when there is a change to the Epoch.
func (node *Darknode) SetEpochListener(ξListener cal.EpochListener) {
	node.ξListener = ξListener
}

// Address of the Darknode.
func (node *Darknode) Address() identity.Address {
	return node.addr
}

// OnConfirmOrder implements the ome.Delegate interface.
func (node *Darknode) OnConfirmOrderMatch(buy order.Order, sell order.Order) error {
	return node.darkpoolAccounts.Settle(buy, sell)
}
