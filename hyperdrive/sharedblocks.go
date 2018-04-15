package hyper

import (
	"bytes"
	"sync/atomic"
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
	history map[[32]byte][32]byte
	Height  uint64
	Rank
}

func NewSharedBlocks(h uint64, r Rank) SharedBlocks {
	return SharedBlocks{
		history: map[[32]byte][32]byte{},
		Height:  h,
		Rank:    r,
	}
}

func (sb *SharedBlocks) ReadHeight() uint64 {
	return atomic.LoadUint64(&sb.Height)
}

func (sb *SharedBlocks) IncrementHeight() uint64 {
	return atomic.AddUint64(&sb.Height, 1)
}

func (blocks *SharedBlocks) ValidateTuple(tuple Tuple) bool {
	for _, element := range tuple.Elements {
		if tupleID, ok := blocks.history[element.ID]; ok && bytes.Compare(tuple.ID[:], tupleID[:]) != 0 {
			return false
		}
	}
	return true
}
