package ingress

import (
	"math/big"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
)

// ContractsBinder will define all methods that the ingresser will
// need to communicate with smart contracts. All the methods will
// be implemented in contracts.Binder
type ContractsBinder interface {
	OpenBuyOrder(signature [65]byte, orderID order.ID) error

	OpenSellOrder(signature [65]byte, orderID order.ID) error

	CancelOrder(signature [65]byte, orderID order.ID) error

	MinimumEpochInterval() (*big.Int, error)

	Epoch() (registry.Epoch, error)

	NextEpoch() (registry.Epoch, error)

	Pods() ([]registry.Pod, error)
}
