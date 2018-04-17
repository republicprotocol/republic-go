package hyper

import (
	"bytes"
	"sync/atomic"
)

type SharedBlocks struct {
	history map[[32]byte][32]byte
	height  uint64
	Rank
}

func NewSharedBlocks(h uint64, r Rank) SharedBlocks {
	return SharedBlocks{
		history: map[[32]byte][32]byte{},
		height:  h,
		Rank:    r,
	}
}

func (sb *SharedBlocks) NextHeight() uint64 {
	return atomic.AddUint64(&sb.height, 1)
}

func (blocks *SharedBlocks) ValidateTuple(tuple Tuple) bool {
	for _, element := range tuple.Elements {
		if tupleID, ok := blocks.history[element.ID]; ok && bytes.Compare(tuple.ID[:], tupleID[:]) != 0 {
			return false
		}
	}
	return true
}
