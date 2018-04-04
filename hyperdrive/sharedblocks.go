package hyper

import (
	"bytes"
	"sync"
)

type Element struct {
	ID [32]byte
}

type Elements []Element

type Tuples []Tuple

type Tuple struct {
	ID [32]byte
	Elements
}

type SharedBlocks struct {
	mu      *sync.RWMutex
	history map[[32]byte][32]byte
}

func NewSharedBlocks() SharedBlocks {
	return SharedBlocks{}
}

func (blocks *SharedBlocks) AddBlock(block Block) {
	blocks.mu.Lock()
	defer blocks.mu.Unlock()
}

func (blocks *SharedBlocks) ValidateTuple(tuple Tuple) bool {
	for _, element := range tuple.Elements {
		if tupleID, ok := blocks.history[element.ID]; ok && bytes.Compare(tuple.ID[:], tupleID[:]) != 0 {
			return false
		}
	}
	return true
}
