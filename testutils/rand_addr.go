package testutils

import (
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

// RandomMultiAddress generates a random crypto.EcdsaKey and uses it to build
// an identity.Address.
func RandomAddress() (identity.Address, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return "", err
	}
	return identity.Address(ecdsaKey.Address()), nil
}

// RandomMultiAddress generates a random crypto.EcdsaKey and uses it to build
// an identity.MultiAddress.
func RandomMultiAddress() (identity.MultiAddress, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	return identity.Address(ecdsaKey.Address()).MultiAddress()
}

// RandomAddressAndMultiAddress generates a random crypto.EcdsaKey and uses it
// to build an identity.Address and identity.MultiAddress.
func RandomAddressAndMultiAddress() (*identity.Address, identity.MultiAddress, error) {
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
