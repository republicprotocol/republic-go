package registry

import (
	"crypto/rsa"

	"github.com/republicprotocol/republic-go/identity"
)

// ContractsBinder will define all interactions that the orderbook will
// have with the smart contracts
type ContractsBinder interface {

	// PublicKey returns the RSA public key of the Darknode registered with the
	// given identity.Address.
	PublicKey(addr identity.Address) (rsa.PublicKey, error)

	// IsRegistered returns true if the identity.Address is a current
	// registered Darknode. Otherwise, it returns false.
	IsRegistered(addr identity.Address) (bool, error)
}
