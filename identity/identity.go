package identity

import (
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multihash"
)

// IDLength is the number of bytes in an ID.
const IDLength = 20

// An ID is a slice of 20 bytes that can be converted into an Address.
// It must always be exact 20 bytes.
type ID []byte

// String returns the ID as a string.
func (id ID) String() string {
	return id.Address().String()
}

// Address returns the Address of the ID
func (id ID) Address() Address {
	hash := make([]byte, 2, IDLength+2)
	hash[0], hash[1] = multihash.KECCAK_256, IDLength
	hash = append(hash, id...)
	return Address(base58.EncodeAlphabet(hash, base58.BTCAlphabet))
}
