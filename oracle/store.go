package oracle

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
)

type MidpointPrice struct {
	Signature []byte
	Tokens    uint64
	Price     uint64
	Nonce     uint64
}

func (midpointPrice MidpointPrice) Equals(other MidpointPrice) bool {
	return bytes.Compare(midpointPrice.Signature, other.Signature) == 0 &&
		midpointPrice.Tokens == other.Tokens &&
		midpointPrice.Price == other.Price &&
		midpointPrice.Nonce == other.Nonce
}

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

type MidpointPriceStorer interface {
	PutMidpointPrice(midpointPrice MidpointPrice) error
	MidpointPrice(tokens uint64) (MidpointPrice, error)
	MidpointPrices() (MidpointPriceIterator, error)
}

type midpointPriceStorer struct {
	mutex          *sync.Mutex
	midpointPrices map[uint64]MidpointPrice
}

func NewMidpointPriceStorer() *midpointPriceStorer {
	return &midpointPriceStorer{
		mutex:          new(sync.Mutex),
		midpointPrices: map[uint64]MidpointPrice{},
	}
}

func (storer *midpointPriceStorer) PutMidpointPrice(midpointPrice MidpointPrice) error {
	storer.mutex.Lock()
	defer storer.mutex.Unlock()

	storer.midpointPrices[midpointPrice.Tokens] = midpointPrice

	return nil
}

func (storer *midpointPriceStorer) MidpointPrice(tokens uint64) (MidpointPrice, error) {
	storer.mutex.Lock()
	defer storer.mutex.Unlock()

	return storer.midpointPrices[tokens], nil
}

func (storer *midpointPriceStorer) MidpointPrices() (MidpointPriceIterator, error) {
	storer.mutex.Lock()
	defer storer.mutex.Unlock()

	iter := NewMidpointPriceIterator(storer.midpointPrices)
	return &iter, nil
}

type MidpointPriceIterator interface {
	Next() bool
	Cursor() (MidpointPrice, error)
	Collect() ([]MidpointPrice, error)
	Release()
}

type midpointPriceIterator struct {
	midpointPrices []MidpointPrice
	cursor         int
}

func NewMidpointPriceIterator(midpointPricesMap map[uint64]MidpointPrice) midpointPriceIterator {
	midpointPrices := make([]MidpointPrice, len(midpointPricesMap))
	for _, value := range midpointPricesMap {
		midpointPrices = append(midpointPrices, value)
	}
	return midpointPriceIterator{
		midpointPrices: midpointPrices,
		cursor:         0,
	}
}

func (iter *midpointPriceIterator) Next() bool {
	if len(iter.midpointPrices) > iter.cursor {
		return true
	}
	return false
}

func (iter *midpointPriceIterator) Cursor() (MidpointPrice, error) {
	if iter.cursor >= len(iter.midpointPrices) {
		return MidpointPrice{}, fmt.Errorf("index out of range")
	}
	return iter.midpointPrices[iter.cursor], nil
}

func (iter *midpointPriceIterator) Collect() ([]MidpointPrice, error) {
	return iter.midpointPrices, nil
}

func (iter *midpointPriceIterator) Release() {
	return
}
