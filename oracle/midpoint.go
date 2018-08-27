package oracle

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"sort"

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

	// Sort mid-point prices based on token values.
	var tokens []int
	for token := range midpointPrice.Prices {
		tokens = append(tokens, int(token))
	}
	sort.Ints(tokens)

	for _, token := range tokens {
		tokensBytes := [8]byte{}
		binary.LittleEndian.PutUint64(tokensBytes[:], uint64(token))
		data = append(data, tokensBytes[:]...)
		priceBytes := [8]byte{}
		binary.LittleEndian.PutUint64(priceBytes[:], midpointPrice.Prices[uint64(token)])
		data = append(data, priceBytes[:]...)
	}

	nonceBytes := [8]byte{}
	binary.LittleEndian.PutUint64(nonceBytes[:], midpointPrice.Nonce)
	data = append(data, nonceBytes[:]...)
	return crypto.Keccak256(data)
}

// IsEmpty returns true if the MidpointPrice is nil.
func (midpointPrice *MidpointPrice) IsEmpty() bool {
	return midpointPrice == nil || len(midpointPrice.Prices) == 0
}
