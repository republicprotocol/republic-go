package testutils

import (
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

func ComputationID(buy, sell order.ID) [32]byte {
	var id [32]byte
	copy(id[:], crypto.Keccak256(buy[:], sell[:]))
	return id
}

func CreateMultiaddress() (identity.MultiAddress, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	return identity.Address(ecdsaKey.Address()).MultiAddress()
}

func CreateAddressAndMultiaddress() (*identity.Address, identity.MultiAddress, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return nil, identity.MultiAddress{}, err
	}
	address := identity.Address(ecdsaKey.Address())
	multiaddress, err := address.MultiAddress()
	if err != nil {
		return nil, identity.MultiAddress{}, err
	}
	return &address, multiaddress, nil
}
