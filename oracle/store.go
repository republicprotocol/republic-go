package oracle

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
)

type MidpointPrice struct {
	Signature []byte
	Tokens    []uint64
	Prices    []uint64
	Nonce     uint64
}

// Equals checks if two MidpointPrice objects have equivalent fields.
func (midpointPrice MidpointPrice) Equals(other MidpointPrice) bool {
	// TODO: Replace reflect.DeepEqual() with traditional loop
	return bytes.Compare(midpointPrice.Signature, other.Signature) == 0 &&
		reflect.DeepEqual(midpointPrice.Tokens, other.Tokens) &&
		reflect.DeepEqual(midpointPrice.Prices, other.Prices) &&
		midpointPrice.Nonce == other.Nonce
}

// Hash returns a keccak256 hash of a MidpointPrice object.
func (midpointPrice MidpointPrice) Hash() []byte {
	tokensBytes := make([]byte, 8)
	// TODO:
	binary.LittleEndian.PutUint64(tokensBytes, midpointPrice.Tokens)
	priceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(priceBytes, midpointPrice.Price)
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, midpointPrice.Nonce)
	midpointPriceBytes := append(tokensBytes, append(priceBytes, nonceBytes...)...)
	return crypto.Keccak256(midpointPriceBytes)
}

// MidpointPriceStorer for the latest MidpointPrice objects from the oracle.
type MidpointPriceStorer interface {
	// PutMidpointPrice stores the given MidpointPrice object.
	PutMidpointPrice(midpointPrice MidpointPrice) error

	// MidpointPrice returns the MidpointPrice stored for the given tokens.
	MidpointPrice(tokens uint64) (MidpointPrice, error)
}

type midpointPriceStorer struct {
	mutex         *sync.Mutex
	midpointPrice MidpointPrice
}

// NewMidpointPriceStorer returns an object that implements the
// MidpointPriceStorer interface.
func NewMidpointPriceStorer() *midpointPriceStorer {
	return &midpointPriceStorer{
		mutex:         new(sync.Mutex),
		midpointPrice: MidpointPrice{},
	}
}

// PutMidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) PutMidpointPrice(midpointPrice MidpointPrice) error {
	storer.mutex.Lock()
	defer storer.mutex.Unlock()

	storer.midpointPrice = midpointPrice

	return nil
}

// MidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrice(tokens uint64) (MidpointPrice, error) {
	storer.mutex.Lock()
	defer storer.mutex.Unlock()

	return storer.midpointPrice, nil
}
