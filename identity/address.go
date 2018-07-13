package identity

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jbenet/go-base58"
)

// AddressLength is the number of bytes in an Address.
const AddressLength = 30

// An Address string is generated from an ID.
type Address string

// Addresses is an alias.
type Addresses []Address

// Distance uses a bitwise XOR to calculate distance between two Addresses.
func (address Address) Distance(other Address) ([]byte, error) {
	// Check the length of both Addresses.
	if len(address) != AddressLength || len(other) != AddressLength {
		return nil, ErrWrongAddressLength
	}
	// Decode both Addresses into bytes.
	idByte := base58.DecodeAlphabet(string(address), base58.BTCAlphabet)
	if len(idByte) == 0 {
		return nil, ErrFailToDecode
	}
	otherByte := base58.DecodeAlphabet(string(other), base58.BTCAlphabet)
	if len(otherByte) == 0 {
		return nil, ErrFailToDecode
	}
	// Get distance by using the XOR operation.
	xor := make([]byte, IDLength)
	for i := 0; i < IDLength; i++ {
		xor[i] = idByte[i+2] ^ otherByte[i+2]
	}
	return xor, nil
}

// SamePrefixLength returns the number of prefix bits that match with those of
// another Address, excluding the first 2 bytes.
func (address Address) SamePrefixLength(other Address) (int, error) {
	// Calculator the distance between two Addresses.
	diff, err := address.Distance(other)
	if err != nil {
		return -1, err
	}
	// Calculate the same prefix bits.
	ret := 0
	for i := 0; i < IDLength; i++ {
		if diff[i] == uint8(0) {
			ret += 8
		} else {
			bits := fmt.Sprintf("%08b", diff[i])
			for j := 0; j < len(bits); j++ {
				if bits[j] == '1' {
					return ret, nil
				}
				ret++
			}
		}
	}
	return ret, nil
}

// MultiAddress returns the Republic multi-address of the Address. It can be
// appended with other MultiAddresses.
func (address Address) MultiAddress() (MultiAddress, error) {
	return NewMultiAddressFromString(fmt.Sprintf("/republic/%s", address))
}

// String returns the Address as a string.
func (address Address) String() string {
	return string(address)
}

// ID returns the ID of the Address
func (address Address) ID() ID {
	bytes := base58.DecodeAlphabet(string(address), base58.BTCAlphabet)
	bytes = bytes[2:]
	return ID(bytes)
}

// Hash implements the crypto.Hasher interface for signing.
func (address Address) Hash() []byte {
	return crypto.Keccak256([]byte(address))
}

// Closer returns true if the left Address is closer to the target than the
// right Address, otherwise it returns false. In the case that both Addresses
// are equal distances from the target, it returns true.
func Closer(left, right, target Address) (bool, error) {
	leftDist, err := left.Distance(target)
	if err != nil {
		return false, err
	}
	rightDist, err := right.Distance(target)
	if err != nil {
		return false, err
	}
	for i := 0; i < IDLength; i++ {
		if leftDist[i] < rightDist[i] {
			return true, nil
		} else if leftDist[i] > rightDist[i] {
			return false, nil
		}
	}

	// If the Addresses are the same, return false
	return false, nil
}
