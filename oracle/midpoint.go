package oracle

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/republicprotocol/republic-go/crypto"
)

// A MidpointPrice is a signed message contains the the mid-prices of
// token pairs.
type MidpointPrice struct {
	Signature []byte
	Prices    map[uint64]uint64
	Nonce     uint64
}

// Equals checks if two MidpointPrice objects have equivalent fields.
func (midpointPrice MidpointPrice) Equals(other MidpointPrice) bool {
	if !reflect.DeepEqual(midpointPrice.Prices, other.Prices) {
		return false
	}
	if bytes.Compare(midpointPrice.Signature, other.Signature) != 0 {
		return false
	}

	return midpointPrice.Nonce == other.Nonce
}

// Hash returns the Keccak256 hash of the MidpointPrice.
func (midpointPrice MidpointPrice) Hash() []byte {
	data := make([]byte, 0)
	for token, price := range midpointPrice.Prices {
		tokensBytes := [8]byte{}
		binary.LittleEndian.PutUint64(tokensBytes[:], token)
		data = append(data, tokensBytes[:]...)
		priceBytes := [8]byte{}
		binary.LittleEndian.PutUint64(priceBytes[:], price)
		data = append(data, priceBytes[:]...)
	}

	nonceBytes := [8]byte{}
	binary.LittleEndian.PutUint64(nonceBytes[:], midpointPrice.Nonce)
	data = append(data, nonceBytes[:]...)
	return crypto.Keccak256(data)
}
