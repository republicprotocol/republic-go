package smpc

import (
	"sync"

	"github.com/republicprotocol/republic-go/shamir"
)

type ShareBuilder struct {
	k int64

	sharesMu    *sync.Mutex
	shares      map[[32]byte]map[uint64]shamir.Share
	sharesCache shamir.Shares
}

func NewShareBuilder(k int64) *ShareBuilder {
	return &ShareBuilder{
		k:           k,
		sharesMu:    new(sync.Mutex),
		shares:      map[[32]byte]map[uint64]shamir.Share{},
		sharesCache: make(shamir.Shares, 0, k),
	}
}

func (builder *ShareBuilder) InsertShare(id [32]byte, share shamir.Share) (uint64, error) {
	builder.sharesMu.Lock()
	defer builder.sharesMu.Unlock()

	if _, ok := builder.shares[id]; !ok {
		builder.shares[id] = map[uint64]shamir.Share{}
	}

	builder.shares[id][share.Index] = share
	if int64(len(builder.shares[id])) >= builder.k {
		builder.sharesCache = builder.sharesCache[0:0]
		k := int64(0)
		for _, share := range builder.shares[id] {
			builder.sharesCache = append(builder.sharesCache, share)
			if k++; k >= builder.k {
				break
			}
		}
		val := shamir.Join(builder.sharesCache)
		return val, nil
	}

	return 0, ErrInsufficientSharesToJoin
}
