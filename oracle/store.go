package oracle

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	return bytes.Compare(midpointPrice.Signature, other.Signature) == 0 &&
		reflect.DeepEqual(midpointPrice.Tokens, other.Tokens) &&
		reflect.DeepEqual(midpointPrice.Prices, other.Prices) &&
		midpointPrice.Nonce == other.Nonce
}

// Hash returns a keccak256 hash of a MidpointPrice object.
func (midpointPrice MidpointPrice) Hash() []byte {
	tokensBytes := make([]byte, 8)
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

	// MidpointPrices returns a new MidpointPriceIterator object using the
	// MidpointPrice objects in the store.
	MidpointPrices() (MidpointPriceIterator, error)
}

type midpointPriceStorer struct {
	mutex          *sync.Mutex
	midpointPrices map[uint64]MidpointPrice
}

// NewMidpointPriceStorer returns an object that implements the
// MidpointPriceStorer interface.
func NewMidpointPriceStorer() *midpointPriceStorer {
	return &midpointPriceStorer{
		mutex:          new(sync.Mutex),
		midpointPrices: map[uint64]MidpointPrice{},
	}
}

// PutMidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) PutMidpointPrice(midpointPrice MidpointPrice) error {
	storer.mutex.Lock()
	defer storer.mutex.Unlock()

	storer.midpointPrices[midpointPrice.Tokens] = midpointPrice

	return nil
}

// MidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrice(tokens uint64) (MidpointPrice, error) {
	storer.mutex.Lock()
	defer storer.mutex.Unlock()

	return storer.midpointPrices[tokens], nil
}

// MidpointPrices implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrices() (MidpointPriceIterator, error) {
	storer.mutex.Lock()
	defer storer.mutex.Unlock()

	midpointPrices := make([]MidpointPrice, 0)
	for _, value := range storer.midpointPrices {
		midpointPrices = append(midpointPrices, value)
	}
	iter := NewMidpointPriceIterator(midpointPrices)
	return &iter, nil
}

// MidpointPriceIterator is used to iterate over a MidpointPrice collection.
type MidpointPriceIterator interface {
	// Next progresses the cursor. Returns true if the new cursor is still in
	// the range of the MidpointPrice collection, otherwise false.
	Next() bool

	// Cursor returns the MidpointPrice at the current cursor location. Returns
	// an error if the cursor is out of range.
	Cursor() (MidpointPrice, error)

	// Collect all MidpointPrices in the iterator into slices.
	Collect() ([]MidpointPrice, error)

	// Release the resources allocated by the iterator.
	Release()
}

type midpointPriceIterator struct {
	midpointPrices []MidpointPrice
	cursor         int
}

// NewMidpointPriceIterator returns an object that implements the
// MidpointPriceIterator interface.
func NewMidpointPriceIterator(midpointPrices []MidpointPrice) midpointPriceIterator {
	return midpointPriceIterator{
		midpointPrices: midpointPrices,
		cursor:         0,
	}
}

// Next implements the MidpointPriceIterator interface.
func (iter *midpointPriceIterator) Next() bool {
	if len(iter.midpointPrices) > iter.cursor {
		return true
	}
	return false
}

// Cursor implements the MidpointPriceIterator interface.
func (iter *midpointPriceIterator) Cursor() (MidpointPrice, error) {
	if iter.cursor >= len(iter.midpointPrices) {
		return MidpointPrice{}, fmt.Errorf("index out of range")
	}
	return iter.midpointPrices[iter.cursor], nil
}

// Collect implements the MidpointPriceIterator interface.
func (iter *midpointPriceIterator) Collect() ([]MidpointPrice, error) {
	return iter.midpointPrices, nil
}

// Release implements the MidpointPriceIterator interface.
func (iter *midpointPriceIterator) Release() {
	return
}
