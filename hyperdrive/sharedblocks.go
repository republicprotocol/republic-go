package hyper

import "sync"

type Entry string

type Entrys []Entry

type Tuples []Tuple

type Tuple struct {
	entrys Entrys
}

type SharedBlocks struct {
	mu      *sync.RWMutex
	history map[Entry]Entry
}

func NewSharedBlocks() SharedBlocks {
	return SharedBlocks{}
}

func (blocks *SharedBlocks) AddBlock(block Block) {
	blocks.mu.Lock()
	defer blocks.mu.Unlock()
}

func (blocks *SharedBlocks) ValidateTuple(tuple Tuple) bool {
	for i := 0; i < len(tuple.entrys); i++ {
		if i < len(tuple.entrys)-1 {
			if !(blocks.history[tuple.entrys[i]] == blocks.history[tuple.entrys[i+1]]) {
				return false
			}
		} else {
			if !(blocks.history[tuple.entrys[i]] == blocks.history[tuple.entrys[0]]) {
				return false
			}
		}
	}
	return true
}
