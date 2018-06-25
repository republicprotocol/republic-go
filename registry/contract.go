package registry

import (
	"crypto/rsa"

	"github.com/republicprotocol/republic-go/identity"
)

// ContractBinder will define all methods that the registry will
// require to communicate with smart contracts. All the methods will
// be implemented in contract.Binder
type ContractBinder interface {
	PublicKey(addr identity.Address) (rsa.PublicKey, error)

	IsRegistered(addr identity.Address) (bool, error)
}
