package dht

import (
	"fmt"
	"github.com/republicprotocol/go-identity"
	"github.com/jbenet/go-base58"
)

// IDLength is the number of bytes needed to store an ID.
const (
	IDLength       = identity.IDLength
	IDLengthInBits = IDLength * 8
	Alpha          = 3
	Republic_Code  = identity.RepublicCode
)

// ID is the public address used to identify Nodes, and other entities, in the
// overlay network. It is generated from the public key of a key pair.
// ID is a string in a Base64 encoding
type ID string

// NewID creates a new set of public/private SECP256K1 key pair and
// returns the public address as the ID string
func NewID() (ID, error) {
	secp, err := identity.NewKeyPair()
	if err != nil {
		return "", err
	}
	return ID(secp.PublicID()), nil
}

// Using xor to calculate distance between two IDs
func (id ID) Xor(other ID) []byte {
	// Decode both the IDs into bytes
	idByte := base58.Decode(string(id))
	otherByte := base58.Decode(string(other))

	xor := make([]byte, IDLength)
	for i := 0; i < IDLength; i++ {
		xor[i] = idByte[i] ^ otherByte[i]
	}
	return xor
}

// Same prefix bits length with another ID
func (id ID) SamePrefixLen(other ID) int {
	diff:= id.Xor(other)
	ret := 0
	for i := 0; i < IDLength; i++ {
		if diff[i] == uint8(0) {
			ret += 8
		} else {
			bits := fmt.Sprintf("%08b", diff[i])
			for j := 0; j < len(bits); j++ {
				if bits[j] == '1' {
					return ret
				}
				ret++
			}
		}
	}

	return ret
}
