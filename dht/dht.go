package dht

import (
	"encoding/base64"
	"fmt"
	"github.com/republicprotocol/go-swarm/crypto"
)

// IDLength is the number of bytes needed to store an ID.
const (
	IDLength       = crypto.PublicAddressLength
	IDLengthBase64 = 28
	IDLengthInBits = IDLength * 8
	Alpha          = 3
	Republic_Code= crypto.Republic_Code
)

// ID is the public address used to identify Nodes, and other entities, in the
// overlay network. It is generated from the public key of a key pair.
// ID is a string in a Base64 encoding
type ID string

// NewID creates a new set of public/private SECP256K1 key pair and
// returns the public address as the ID string
func NewID() (ID, error) {
	secp, err := crypto.NewSECP256K1()
	if err != nil {
		return "", err
	}
	return ID(secp.PublicAddress()), nil
}

// Using xor to calculate distance between two IDs
func (id ID) Xor(other ID) ([]byte, error) {
	// Decode both the IDs into bytes
	idByte, err := base64.StdEncoding.DecodeString(string(id))
	if err != nil {
		return nil, err
	}
	otherByte, err := base64.StdEncoding.DecodeString(string(other))
	if err != nil {
		return nil, err
	}

	xor := make([]byte, IDLength)
	for i := 0; i < IDLength; i++ {
		xor[i] = idByte[i] ^ otherByte[i]
	}
	return xor, nil
}

// Same prefix bits length with another ID
func (id ID) SamePrefixLen(other ID) (int, error) {
	diff, err := id.Xor(other)
	if err != nil {
		return 0, err
	}

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