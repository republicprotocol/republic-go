package smpc

import (
	"sync"

	"github.com/republicprotocol/republic-go/shamir"
)

type ShareBuilderObserver interface {
	OnNotifyBuild(id, networkID [32]byte, value uint64)
}

type ShareBuilder struct {
	k int64

	sharesMu    *sync.Mutex
	sharesCache shamir.Shares
	shares      map[[32]byte]map[uint64]shamir.Share

	observersMu *sync.Mutex
	observers   map[[32]byte]map[[32]byte]ShareBuilderObserver
}

func NewShareBuilder(k int64) *ShareBuilder {
	return &ShareBuilder{
		k: k,

		sharesMu:    new(sync.Mutex),
		shares:      map[[32]byte]map[uint64]shamir.Share{},
		sharesCache: make(shamir.Shares, 0, k),

		observersMu: new(sync.Mutex),
		observers:   map[[32]byte]map[[32]byte]ShareBuilderObserver{},
	}
}

func (builder *ShareBuilder) Insert(id [32]byte, share shamir.Share) error {
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
		builder.notify(id, val)
		return nil
	}

	return ErrInsufficientSharesToJoin
}

func (builder *ShareBuilder) Observe(id, networkID [32]byte, observer ShareBuilderObserver) {
	builder.observersMu.Lock()
	defer builder.observersMu.Unlock()

	if _, ok := builder.observers[id]; !ok {
		builder.observers[id] = map[[32]byte]ShareBuilderObserver{}
	}

	if observer == nil {
		delete(builder.observers[id], networkID)
		if len(builder.observers[id]) == 0 {
			delete(builder.observers, id)
		}
		return
	}

	builder.observers[id][networkID] = observer
}

func (builder *ShareBuilder) notify(id [32]byte, val uint64) {
	builder.observersMu.Lock()
	defer builder.observersMu.Unlock()

	if observers, ok := builder.observers[id]; ok {
		for networkID, observer := range observers {
			observer.OnNotifyBuild(id, networkID, val)
		}
	}
}
