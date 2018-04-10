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
	Height
	Rank
}

func NewSharedBlocks(r Rank, h Height) SharedBlocks {
	return SharedBlocks{
		mu:      new(sync.RWMutex),
		history: map[[32]byte][32]byte{},
		Rank:    r,
		Height:  h,
	}
}

func (blocks *SharedBlocks) IncrementHeight() {
	blocks.mu.Lock()
	defer blocks.mu.Unlock()
	blocks.Height++
}

func (blocks *SharedBlocks) ReadHeight() Height {
	blocks.mu.Lock()
	h := blocks.Height
	defer blocks.mu.Unlock()
	return h
}

func (blocks *SharedBlocks) GetSharedBlocks() *SharedBlocks {
	return blocks
}

func (blocks *SharedBlocks) ValidateTuple(tuple Tuple) bool {
	for _, element := range tuple.Elements {
		if tupleID, ok := blocks.history[element.ID]; ok && bytes.Compare(tuple.ID[:], tupleID[:]) != 0 {
			return false
		}
	}
	return true
}
