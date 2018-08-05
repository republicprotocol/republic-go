package oracle

import (
	"bytes"
	"encoding/binary"

	"github.com/republicprotocol/republic-go/crypto"
)

// A MidpointPrice is a signed message contains the the mid-prices of
// token pairs.
type MidpointPrice struct {
	Signature  []byte
	TokenPairs []uint64
	Prices     []uint64
	Nonce      uint64
}

// Equals checks if two MidpointPrice objects have equivalent fields.
func (midpointPrice MidpointPrice) Equals(other MidpointPrice) bool {
	if bytes.Compare(midpointPrice.Signature, other.Signature) != 0 {
		return false
	}

	// Compare the token pairs/prices.
	if len(midpointPrice.TokenPairs) != len(other.TokenPairs) || len(midpointPrice.Prices) != len(other.Prices) {
		return false
	}
	tokenPricePairs := map[uint64]uint64{}
	for i, token := range midpointPrice.TokenPairs {
		tokenPricePairs[token] = midpointPrice.Prices[i]
	}
	for i, token := range other.TokenPairs {
		price, ok := tokenPricePairs[token]
		if !ok || price != other.Prices[i] {
			return false
		}
	}

	return midpointPrice.Nonce == other.Nonce
}

// Hash returns the Keccak256 hash of the MidpointPrice.
func (midpointPrice MidpointPrice) Hash() []byte {
	data := make([]byte, 0)
	for i := range midpointPrice.TokenPairs {
		tokensBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(tokensBytes, midpointPrice.TokenPairs[i])
		data = append(data, tokensBytes...)
		priceBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(priceBytes, midpointPrice.Prices[i])
		data = append(data, priceBytes...)
	}

	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, midpointPrice.Nonce)
	data = append(data, nonceBytes...)
	return crypto.Keccak256(data)
}
