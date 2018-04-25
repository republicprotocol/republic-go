package smpc

import (
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/delta"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

// OrderTuplesToDeltaFragments reads OrderTuples from an input channel and
// uses them to build DeltaFragments.
func OrderTuplesToDeltaFragments(done <-chan struct{}, orderTuples <-chan OrderTuple, bufferLimit int) <-chan delta.Fragment {
	deltaFragments := make(chan delta.Fragment, bufferLimit)

	go func() {
		defer close(deltaFragments)

		for {
			select {
			case <-done:
				return
			case orderTuple, ok := <-orderTuples:
				if !ok {
					return
				}
				deltaFragment := delta.NewDeltaFragment(orderTuple.BuyOrderFragment, orderTuple.SellOrderFragment, &Prime)
				select {
				case <-done:
					return
				case deltaFragments <- deltaFragment:
				}
			}
		}
	}()

	return deltaFragments
}

// BuildDeltas by reading DeltaFragments from an input channel and using a
// SharedDeltaBuilder to store them. Deltas can be built from the
// SharedDeltaBuilder after a threshold of DeltaFragments has been reached.
func BuildDeltas(done <-chan struct{}, deltaFragments <-chan delta.Fragment, sharedDeltaBuilder *SharedDeltaBuilder, bufferLimit int) <-chan delta.Delta {
	deltas := make(chan delta.Delta, bufferLimit)

	// Insert DeltaFragments into the SharedDeltaBuilder
	go func() {
		for {
			select {
			case <-done:
				return
			case deltaFragment, ok := <-deltaFragments:
				println("RECEIVED DELTA FRAGMENTs")
				if !ok {
					return
				}
				sharedDeltaBuilder.InsertDeltaFragment(deltaFragment)
			}
		}
	}()

	// Periodically read computed Deltas from the SharedDeltaBuilder into a
	// buffer
	go func() {
		defer close(deltas)

		buffer := make([]delta.Delta, bufferLimit)
		tick := time.NewTicker(time.Millisecond)
		defer tick.Stop()

		for {
			select {
			case <-done:
				return
			case <-tick.C:
				for i, n := 0, sharedDeltaBuilder.Deltas(buffer[:]); i < n; i++ {
					select {
					case <-done:
						return
					case deltas <- buffer[i]:
					}
				}
			}
		}
	}()

	return deltas
}

type SharedDeltaBuilder struct {
	k                    int64
	prime                stackint.Int1024
	fstCodeSharesCache   shamir.Shares
	sndCodeSharesCache   shamir.Shares
	priceSharesCache     shamir.Shares
	minVolumeSharesCache shamir.Shares
	maxVolumeSharesCache shamir.Shares

	mu                     *sync.Mutex
	deltas                 map[string]delta.Delta
	deltaFragments         map[string]delta.Fragment
	deltasToDeltaFragments map[string]delta.Fragments
	deltasQueue            delta.Deltas
}

func NewSharedDeltaBuilder(k int64, prime stackint.Int1024) SharedDeltaBuilder {
	return SharedDeltaBuilder{
		k:                      k,
		prime:                  prime,
		fstCodeSharesCache:     make(shamir.Shares, k),
		sndCodeSharesCache:     make(shamir.Shares, k),
		priceSharesCache:       make(shamir.Shares, k),
		minVolumeSharesCache:   make(shamir.Shares, k),
		maxVolumeSharesCache:   make(shamir.Shares, k),
		mu:                     new(sync.Mutex),
		deltas:                 map[string]delta.Delta{},
		deltaFragments:         map[string]delta.Fragment{},
		deltasToDeltaFragments: map[string]delta.Fragments{},
		deltasQueue:            delta.Deltas{},
	}
}

func (builder *SharedDeltaBuilder) InsertDeltaFragment(deltaFragment delta.Fragment) {
	builder.mu.Lock()
	defer builder.mu.Unlock()

	// Store the DeltaFragment if it has not been seen before
	if builder.hasDeltaFragment(deltaFragment.ID) {
		println("same fragment twice")
		return
	}
	builder.deltaFragments[string(deltaFragment.ID)] = deltaFragment

	// Associate the DeltaFragment with its respective Delta if the Delta
	// has not been built yet
	if builder.hasDelta(deltaFragment.DeltaID) {
		println("same delta twice")
		return
	}
	if _, ok := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]; ok {
		builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] =
			append(builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)], deltaFragment)
	} else {
		builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] = delta.Fragments{deltaFragment}
	}

	// Build the Delta
	deltaFragments := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]
	if int64(len(deltaFragments)) >= builder.k {
		if !delta.IsCompatible(deltaFragments) {
			println("incompatible")
			return
		}

		for i := int64(0); i < builder.k; i++ {
			builder.fstCodeSharesCache[i] = deltaFragments[i].FstCodeShare
			builder.sndCodeSharesCache[i] = deltaFragments[i].SndCodeShare
			builder.priceSharesCache[i] = deltaFragments[i].PriceShare
			builder.maxVolumeSharesCache[i] = deltaFragments[i].MaxVolumeShare
			builder.minVolumeSharesCache[i] = deltaFragments[i].MinVolumeShare
		}

		delta := delta.NewDeltaFromShares(
			deltaFragments[0].BuyOrderID,
			deltaFragments[0].SellOrderID,
			builder.fstCodeSharesCache,
			builder.sndCodeSharesCache,
			builder.priceSharesCache,
			builder.minVolumeSharesCache,
			builder.maxVolumeSharesCache,
			builder.k,
			builder.prime)
		builder.deltas[string(delta.ID)] = delta
		builder.deltasQueue = append(builder.deltasQueue, delta)
	}
}

func (builder *SharedDeltaBuilder) Deltas(deltas delta.Deltas) int {
	builder.mu.Lock()
	defer builder.mu.Unlock()

	n := 0
	for i := 0; i < len(deltas) && i < len(builder.deltasQueue); i++ {
		deltas[i] = builder.deltasQueue[i]
		n++
	}

	if n >= len(builder.deltasQueue) {
		builder.deltasQueue = builder.deltasQueue[0:0]
	} else {
		builder.deltasQueue = builder.deltasQueue[n:]
	}
	return n
}

func (builder *SharedDeltaBuilder) hasDelta(deltaID delta.ID) bool {
	_, ok := builder.deltas[string(deltaID)]
	return ok
}

func (builder *SharedDeltaBuilder) hasDeltaFragment(deltaFragmentID delta.FragmentID) bool {
	_, ok := builder.deltaFragments[string(deltaFragmentID)]
	return ok
}
