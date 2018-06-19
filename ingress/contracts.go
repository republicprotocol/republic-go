package ingress

import (
	"math/big"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
)

// ContractsBinder will define all interactions that the orderbook will
// have with the smart contracts
type ContractsBinder interface {

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

	// MinimumEpochInterval returns the minimum number of seconds between
	// epochs.
	MinimumEpochInterval() (*big.Int, error)

	// Epoch returns the current Epoch which includes the Pod configuration.
	Epoch() (registry.Epoch, error)

	// NextEpoch will try to turn the Epoch and returns the resulting Epoch. If
	// the turning of the Epoch failed, the current Epoch is returned.
	NextEpoch() (registry.Epoch, error)

	// Pods returns the Pod configuration for the current Epoch.
	Pods() ([]registry.Pod, error)
}
