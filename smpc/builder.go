package smpc

import (
	"sync"

	"github.com/republicprotocol/republic-go/shamir"
)

type shareBuilder struct {
	k int64

	sharesMu *sync.Mutex
	shares   map[[32]byte]shamir.Shares
}

func newShareBuilder(k int64) *shareBuilder {
	return &shareBuilder{
		k:        k,
		sharesMu: new(sync.Mutex),
		shares:   map[[32]byte]shamir.Shares{},
	}
}

func (builder *shareBuilder) insertShare(id [32]byte, share shamir.Share) (uint64, error) {
	builder.sharesMu.Lock()
	defer builder.sharesMu.Unlock()

	if _, ok := builder.shares[id]; !ok {
		builder.shares[id] = shamir.Shares{}
	}

	builder.shares[id] = append(builder.shares[id], share)
	if int64(len(builder.shares[id])) >= builder.k {
		val := shamir.Join(builder.shares[id])
		delete(builder.shares, id)
		return val, nil
	}

	return 0, ErrInsufficientSharesToJoin
}
