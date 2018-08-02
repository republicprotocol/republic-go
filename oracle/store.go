package oracle

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
)

// ErrCursorOutOfRange is returned when an iterator cursor is used to read a
// value outside the range of the iterator.
var ErrCursorOutOfRange = errors.New("cursor out of range")

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

// MidpointPriceStorer is used for retrieving/storing MidpointPrice.
type MidpointPriceStorer interface {

	// PutMidpointPrice stores the given midPointPrice into the storer.
	// It doesn't do any nonce check and will always overwrite previous record.
	PutMidpointPrice(MidpointPrice) error

	// MidpointPrice returns the latest mid-point price of the given token.
	MidpointPrice() (MidpointPrice, error)
}

// midpointPriceStorer implements MidpointPriceStorer interface with an
// in-memory map implementation. Data will be lost every time it restarts.
type midpointPriceStorer struct {
	mu     *sync.Mutex
	prices MidpointPrice
}

// NewMidpointPriceStorer returns a new MidpointPriceStorer.
func NewMidpointPriceStorer() MidpointPriceStorer {
	return &midpointPriceStorer{
		mu: new(sync.Mutex),
	}
}

// PutMidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) PutMidpointPrice(midPointPrice MidpointPrice) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.prices = midPointPrice

	return nil
}

// MidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrice() (MidpointPrice, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	return storer.prices, nil
}
