package oracle

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
)

// ErrCursorOutOfRange is returned when an iterator cursor is used to read a
// value outside the range of the iterator.
var ErrCursorOutOfRange = errors.New("cursor out of range")

// A MidpointPrice is a signed message contains the mid-point price
// of a specific token.
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

// Hash returns the Keccak256 hash of the MidpointPrice.
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

// MidpointPriceStorer is used for retrieving/storing MidpointPrice.
type MidpointPriceStorer interface {

	// PutMidpointPrice stores the given midPointPrice into the storer.
	// It doesn't do any nonce check and will always overwrite previous record.
	PutMidpointPrice(tokens, price uint64) error

	// MidpointPrice returns the latest mid-point price of the given token.
	MidpointPrice(tokens uint64) (uint64, error)

	// MidpointPrices returns an iterator of mid-point prices of all tokens
	MidpointPrices() (MidpointPriceIterator, error)
}

// midpointPriceStorer implements MidpointPriceStorer interface with an
// in-memory map implementation. Data will be lost every time it restarts.
type midpointPriceStorer struct {
	mu     *sync.Mutex
	prices map[uint64]uint64
}

// todo
// NewMidpointPriceStorer returns a new MidpointPriceStorer.
func NewMidpointPriceStorer() MidpointPriceStorer {
	return &midpointPriceStorer{
		mu:     new(sync.Mutex),
		prices: map[uint64]uint64{},
	}
}

// PutMidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) PutMidpointPrice(token, price uint64) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.prices[midpointPrice.Tokens] = midpointPrice

	return nil
}

// MidpointPrice implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrice(tokens uint64) (MidpointPrice, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	return storer.prices[tokens], nil
}

// MidpointPrices implements the MidpointPriceStorer interface.
func (storer *midpointPriceStorer) MidpointPrices() (MidpointPriceIterator, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	midpointPrices := make([]MidpointPrice, 0)
	for _, value := range storer.prices {
		midpointPrices = append(midpointPrices, value)
	}

	return NewMidpointPriceIterator(midpointPrices), nil
}

// MidpointPriceIterator is used to iterate over an MidpointPrice collection.
type MidpointPriceIterator interface {
	Next() bool
	Cursor() (MidpointPrice, error)
	Collect() ([]MidpointPrice, error)
	Release()
}

// midpointPriceIterator implements the MidpointPriceIterator interface.
type midpointPriceIterator struct {
	midpointPrices []MidpointPrice
	cursor         int
}

// NewMidpointPriceIterator returns a new MidpointPriceIterator.
func NewMidpointPriceIterator(midpointPrices []MidpointPrice) MidpointPriceIterator {
	return &midpointPriceIterator{
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
		return MidpointPrice{}, ErrCursorOutOfRange
	}
	price := iter.midpointPrices[iter.cursor]
	iter.cursor++
	return price, nil
}

// Collect implements the MidpointPriceIterator interface.
func (iter *midpointPriceIterator) Collect() ([]MidpointPrice, error) {
	return iter.midpointPrices, nil
}

// Release implements the MidpointPriceIterator interface.
func (iter *midpointPriceIterator) Release() {
	return
}
