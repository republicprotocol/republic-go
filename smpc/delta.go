package smpc

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

func ProduceRhoSigmaFragments(ctx context.Context) (<-chan RhoSigmaFragment, <-chan error) {
	rhoSigmaFragmentCh := make(chan RhoSigmaFragment)
	errCh := make(chan error)

	go func() {
		defer close(rhoSigmaFragmentCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return

			}
		}
	}()

	return rhoSigmaFragmentCh, errCh
}

type RhoSigmaFragment struct {
}

type DeltaBuilder struct {
	k                    int64
	prime                stackint.Int1024
	fstCodeSharesCache   shamir.Shares
	sndCodeSharesCache   shamir.Shares
	priceSharesCache     shamir.Shares
	minVolumeSharesCache shamir.Shares
	maxVolumeSharesCache shamir.Shares

	mu                     *sync.Mutex
	deltas                 map[string]Delta
	deltaFragments         map[string]DeltaFragment
	deltasToDeltaFragments map[string]DeltaFragments
	deltasQueue            Deltas
}

func NewDeltaBuilder(k int64, prime stackint.Int1024) DeltaBuilder {
	return DeltaBuilder{
		k:                      k,
		prime:                  prime,
		fstCodeSharesCache:     make(shamir.Shares, k),
		sndCodeSharesCache:     make(shamir.Shares, k),
		priceSharesCache:       make(shamir.Shares, k),
		minVolumeSharesCache:   make(shamir.Shares, k),
		maxVolumeSharesCache:   make(shamir.Shares, k),
		mu:                     new(sync.Mutex),
		deltas:                 map[string]Delta{},
		deltaFragments:         map[string]DeltaFragment{},
		deltasToDeltaFragments: map[string]DeltaFragments{},
		deltasQueue:            Deltas{},
	}
}

func (builder *DeltaBuilder) ComputeDelta(deltaFragments DeltaFragments) {
	builder.mu.Lock()
	defer builder.mu.Unlock()

	for _, deltaFragment := range deltaFragments {
		// Store the DeltaFragment if it has not been seen before
		if builder.hasDeltaFragment(deltaFragment.ID) {
			continue
		}
		builder.deltaFragments[string(deltaFragment.ID)] = deltaFragment

		// Associate the DeltaFragment with its respective Delta if the Delta
		// has not been built yet
		if builder.hasDelta(deltaFragment.DeltaID) {
			return
		}
		if _, ok := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]; ok {
			builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] =
				append(builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)], deltaFragment)
		} else {
			builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] = DeltaFragments{deltaFragment}
		}

		// Build the Delta
		deltaFragments := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]
		if int64(len(deltaFragments)) >= builder.k {
			if !IsCompatible(deltaFragments) {
				continue
			}

			for i := int64(0); i < builder.k; i++ {
				builder.fstCodeSharesCache[i] = deltaFragments[i].FstCodeShare
				builder.sndCodeSharesCache[i] = deltaFragments[i].SndCodeShare
				builder.priceSharesCache[i] = deltaFragments[i].PriceShare
				builder.maxVolumeSharesCache[i] = deltaFragments[i].MaxVolumeShare
				builder.minVolumeSharesCache[i] = deltaFragments[i].MinVolumeShare
			}

			delta := NewDeltaFromShares(
				deltaFragments[0].BuyOrderID,
				deltaFragments[0].SellOrderID,
				builder.fstCodeSharesCache,
				builder.sndCodeSharesCache,
				builder.priceSharesCache,
				builder.minVolumeSharesCache,
				builder.maxVolumeSharesCache,
				builder.k, builder.prime)
			builder.deltas[string(delta.ID)] = delta
			builder.deltasQueue = append(builder.deltasQueue, delta)
		}
	}
}

func (builder *DeltaBuilder) Deltas(deltas Deltas) int {
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

func (builder *DeltaBuilder) hasDelta(deltaID DeltaID) bool {
	_, ok := builder.deltas[string(deltaID)]
	return ok
}

func (builder *DeltaBuilder) hasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	_, ok := builder.deltaFragments[string(deltaFragmentID)]
	return ok
}

// A DeltaID is the Keccak256 hash of the order IDs that were used to compute
// the associated Delta.
type DeltaID []byte

// Equal returns an equality check between two DeltaIDs.
func (id DeltaID) Equal(other DeltaID) bool {
	return bytes.Equal(id, other)
}

// String returns a DeltaID as a Base58 encoded string.
func (id DeltaID) String() string {
	return base58.Encode(id)
}

type Deltas []Delta

type Delta struct {
	ID          DeltaID
	BuyOrderID  order.ID
	SellOrderID order.ID
	FstCode     stackint.Int1024
	SndCode     stackint.Int1024
	Price       stackint.Int1024
	MaxVolume   stackint.Int1024
	MinVolume   stackint.Int1024
}

func NewDeltaFromShares(buyOrderID, sellOrderID order.ID, fstCodeShares, sndCodeShares, priceShares, maxVolumeShares, minVolumeShares shamir.Shares, k int64, prime stackint.Int1024) Delta {
	// Join the Shares into a Result.
	delta := Delta{
		BuyOrderID:  buyOrderID,
		SellOrderID: sellOrderID,
	}
	delta.FstCode = shamir.Join(&prime, fstCodeShares)
	delta.SndCode = shamir.Join(&prime, sndCodeShares)
	delta.Price = shamir.Join(&prime, priceShares)
	delta.MaxVolume = shamir.Join(&prime, maxVolumeShares)
	delta.MinVolume = shamir.Join(&prime, minVolumeShares)

	// Compute the ResultID and return the Result.
	delta.ID = DeltaID(crypto.Keccak256(delta.BuyOrderID[:], delta.SellOrderID[:]))
	return delta
}

func (delta *Delta) IsMatch(prime stackint.Int1024) bool {
	zero := stackint.Zero()
	two := stackint.Two()
	zeroThreshold := prime.Div(&two)
	if delta.FstCode.Cmp(&zero) != 0 {
		return false
	}
	if delta.SndCode.Cmp(&zero) != 0 {
		return false
	}
	if delta.Price.Cmp(&zeroThreshold) == 1 {
		return false
	}
	if delta.MaxVolume.Cmp(&zeroThreshold) == 1 {
		return false
	}
	if delta.MinVolume.Cmp(&zeroThreshold) == 1 {
		return false
	}
	return true
}

// A DeltaFragmentID is the Keccak256 hash of the order IDs that were used to
// compute the associated DeltaFragment.
type DeltaFragmentID []byte

// Equal returns an equality check between two DeltaFragmentIDs.
func (id DeltaFragmentID) Equal(other DeltaFragmentID) bool {
	return bytes.Equal(id, other)
}

// String returns a DeltaFragmentID as a Base58 encoded string.
func (id DeltaFragmentID) String() string {
	return base58.Encode(id)
}

type DeltaFragments []DeltaFragment

// A DeltaFragment is a secret share of a Final. Is is performing a
// computation over two OrderFragments.
type DeltaFragment struct {
	ID                  DeltaFragmentID
	DeltaID             DeltaID
	BuyOrderID          order.ID
	SellOrderID         order.ID
	BuyOrderFragmentID  order.FragmentID
	SellOrderFragmentID order.FragmentID

	FstCodeShare   shamir.Share
	SndCodeShare   shamir.Share
	PriceShare     shamir.Share
	MaxVolumeShare shamir.Share
	MinVolumeShare shamir.Share
}

func NewDeltaFragment(left, right *order.Fragment, prime stackint.Int1024) DeltaFragment {
	var buyOrderFragment, sellOrderFragment *order.Fragment
	if left.OrderParity == order.ParityBuy {
		buyOrderFragment = left
		sellOrderFragment = right
	} else {
		buyOrderFragment = right
		sellOrderFragment = left
	}

	fstCodeShare := shamir.Share{
		Key:   buyOrderFragment.FstCodeShare.Key,
		Value: buyOrderFragment.FstCodeShare.Value.SubModulo(&sellOrderFragment.FstCodeShare.Value, &prime),
	}
	sndCodeShare := shamir.Share{
		Key:   buyOrderFragment.SndCodeShare.Key,
		Value: buyOrderFragment.SndCodeShare.Value.SubModulo(&sellOrderFragment.SndCodeShare.Value, &prime),
	}
	priceShare := shamir.Share{
		Key:   buyOrderFragment.PriceShare.Key,
		Value: buyOrderFragment.PriceShare.Value.SubModulo(&sellOrderFragment.PriceShare.Value, &prime),
	}
	maxVolumeShare := shamir.Share{
		Key:   buyOrderFragment.MaxVolumeShare.Key,
		Value: buyOrderFragment.MaxVolumeShare.Value.SubModulo(&sellOrderFragment.MinVolumeShare.Value, &prime),
	}
	minVolumeShare := shamir.Share{
		Key:   buyOrderFragment.MinVolumeShare.Key,
		Value: sellOrderFragment.MaxVolumeShare.Value.SubModulo(&buyOrderFragment.MinVolumeShare.Value, &prime),
	}

	return DeltaFragment{
		ID:                  DeltaFragmentID(crypto.Keccak256([]byte(buyOrderFragment.ID), []byte(sellOrderFragment.ID))),
		DeltaID:             DeltaID(crypto.Keccak256([]byte(buyOrderFragment.OrderID), []byte(sellOrderFragment.OrderID))),
		BuyOrderID:          buyOrderFragment.OrderID,
		SellOrderID:         sellOrderFragment.OrderID,
		BuyOrderFragmentID:  buyOrderFragment.ID,
		SellOrderFragmentID: sellOrderFragment.ID,
		FstCodeShare:        fstCodeShare,
		SndCodeShare:        sndCodeShare,
		PriceShare:          priceShare,
		MaxVolumeShare:      maxVolumeShare,
		MinVolumeShare:      minVolumeShare,
	}
}

// Equals checks if two DeltaFragments are equal in value.
func (deltaFragment *DeltaFragment) Equals(other *DeltaFragment) bool {
	return deltaFragment.ID.Equal(other.ID) &&
		deltaFragment.DeltaID.Equal(other.DeltaID) &&
		deltaFragment.BuyOrderID.Equal(other.BuyOrderID) &&
		deltaFragment.SellOrderID.Equal(other.SellOrderID) &&
		deltaFragment.BuyOrderFragmentID.Equal(other.BuyOrderFragmentID) &&
		deltaFragment.SellOrderFragmentID.Equal(other.SellOrderFragmentID) &&
		deltaFragment.FstCodeShare.Key == other.FstCodeShare.Key &&
		deltaFragment.FstCodeShare.Value.Cmp(&other.FstCodeShare.Value) == 0 &&
		deltaFragment.SndCodeShare.Key == other.SndCodeShare.Key &&
		deltaFragment.SndCodeShare.Value.Cmp(&other.SndCodeShare.Value) == 0 &&
		deltaFragment.PriceShare.Key == other.PriceShare.Key &&
		deltaFragment.PriceShare.Value.Cmp(&other.PriceShare.Value) == 0 &&
		deltaFragment.MaxVolumeShare.Key == other.MaxVolumeShare.Key &&
		deltaFragment.MaxVolumeShare.Value.Cmp(&other.MaxVolumeShare.Value) == 0 &&
		deltaFragment.MinVolumeShare.Key == other.MinVolumeShare.Key &&
		deltaFragment.MinVolumeShare.Value.Cmp(&other.MinVolumeShare.Value) == 0
}

// IsCompatible returns true if all DeltaFragments are fragments of the same
// Delta, otherwise it returns false.
func IsCompatible(deltaFragments DeltaFragments) bool {
	if len(deltaFragments) == 0 {
		return false
	}
	for i := range deltaFragments {
		if !deltaFragments[i].DeltaID.Equal(deltaFragments[0].DeltaID) {
			return false
		}
	}
	return true
}

// DeltaQueues is a slice of DeltaQueue components.
type DeltaQueues []DeltaQueue

// A DeltaQueue owns a channel of Delta components.
type DeltaQueue struct {
	chMu   *sync.RWMutex
	chOpen bool
	ch     chan Delta
}

// NewDeltaQueue returns a MessageQueue interface that channels Delta
//components.
func NewDeltaQueue(messageQueueLimit int) DeltaQueue {
	return DeltaQueue{
		chMu:   new(sync.RWMutex),
		chOpen: true,
		ch:     make(chan Delta, messageQueueLimit),
	}
}

// Run the DeltaQueue. The DeltaQueue is an abstraction over a channel of Delta
// components and does not need to be run. This method does nothing.
func (queue *DeltaQueue) Run() error {
	return nil
}

// Shutdown the DeltaQueue. If it has already been Shutdown, an error will be
// returned.
func (queue *DeltaQueue) Shutdown() error {
	queue.chMu.Lock()
	defer queue.chMu.Unlock()

	queue.chOpen = false
	close(queue.ch)
	return nil
}

// Send a message to the DeltaQueue. The Message must be a Delta component,
// otherwise an error is returned.
func (queue *DeltaQueue) Send(message dispatch.Message) error {
	queue.chMu.RLock()
	defer queue.chMu.RUnlock()

	if !queue.chOpen {
		return nil
	}

	switch message := message.(type) {
	case Delta:
		queue.ch <- message
	default:
		return fmt.Errorf("cannot send message: unrecognized type %T", message)
	}
	return nil
}

// Recv a message from the DeltaQueue. All Messages returned will be Delta
// components.
func (queue *DeltaQueue) Recv() (dispatch.Message, bool) {
	message, ok := <-queue.ch
	return message, ok
}
