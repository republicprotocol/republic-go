package delta

import (
	"bytes"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
)

// A ID is the Keccak256 hash of the order IDs that were used to compute
// the associated Delta.
type ID [32]byte

// Equal returns an equality check between two delta ID .
func (id ID) Equal(other ID) bool {
	return bytes.Equal(id[:], other[:])
}

// String returns a ID as a Base58 encoded string.
func (id ID) String() string {
	return base58.Encode(id[:])
}

type Deltas []Delta

type Delta struct {
	ID          ID
	BuyOrderID  order.ID
	SellOrderID order.ID
	Tokens      uint64
	Price       order.CoExp
	Volume      order.CoExp
	MinVolume   order.CoExp
}

func NewDeltaFromShares(buyOrderID, sellOrderID order.ID, tokenShares, priceCoshares, priceExpShares, volumeCoShares, volumeExpShare, minVolumeCoShare, minVolumeExpShare []shamir.Share) *Delta {
	delta := Delta{
		BuyOrderID:  buyOrderID,
		SellOrderID: sellOrderID,
	}
	delta.Tokens = shamir.Join(tokenShares)
	delta.Price.Co = uint32(shamir.Join(priceCoshares))
	delta.Price.Exp = uint32(shamir.Join(priceExpShares))
	delta.Volume.Co = uint32(shamir.Join(volumeCoShares))
	delta.Volume.Exp = uint32(shamir.Join(volumeExpShare))
	delta.MinVolume.Co = uint32(shamir.Join(minVolumeCoShare))
	delta.MinVolume.Exp = uint32(shamir.Join(minVolumeExpShare))

	// Compute the ResultID and return the Result.
	var ID [32]byte
	copy(ID[:], crypto.Keccak256(delta.BuyOrderID[:], delta.SellOrderID[:]))
	delta.ID = ID

	return &delta
}

func (delta *Delta) IsMatch() bool {
	zeroThreshold := shamir.Prime / 2

	if delta.Tokens != 0 {
		return false
	}
	if uint64(delta.Price.Exp) >= zeroThreshold {
		return false
	}
	if uint64(delta.Price.Co) >= zeroThreshold {
		return false
	}
	if uint64(delta.Volume.Exp) >= zeroThreshold {
		return false
	}
	if uint64(delta.Volume.Co) >= zeroThreshold {
		return false
	}
	if uint64(delta.MinVolume.Exp) >= zeroThreshold {
		return false
	}
	if uint64(delta.MinVolume.Co) >= zeroThreshold {
		return false
	}

	return true
}

type deltaBuilder struct {
	PeersID []byte
	k       int
	mu      *sync.Mutex

	deltas                 map[ID][]Fragment
	tokenSharesCache       []shamir.Share
	priceCoShareCache      []shamir.Share
	priceExpShareCache     []shamir.Share
	volumeCoShareCache     []shamir.Share
	volumeExpShareCache    []shamir.Share
	minVolumeCoShareCache  []shamir.Share
	minVolumeExpShareCache []shamir.Share
}

func NewDeltaBuilder(peersID []byte, k int) *deltaBuilder {
	builder := new(deltaBuilder)
	builder.PeersID = peersID
	builder.k = k
	builder.mu = new(sync.Mutex)
	builder.deltas = map[ID][]Fragment{}
	builder.tokenSharesCache = make([]shamir.Share, k)
	builder.priceCoShareCache = make([]shamir.Share, k)
	builder.priceExpShareCache = make([]shamir.Share, k)
	builder.volumeCoShareCache = make([]shamir.Share, k)
	builder.volumeExpShareCache = make([]shamir.Share, k)
	builder.minVolumeCoShareCache = make([]shamir.Share, k)
	builder.minVolumeCoShareCache = make([]shamir.Share, k)

	return builder
}

func (builder *deltaBuilder) InsertDeltaFragment(fragment Fragment) (*Delta, error) {
	builder.mu.Lock()
	defer builder.mu.Unlock()

	builder.deltas[fragment.DeltaID] = append(builder.deltas[fragment.DeltaID], fragment)
	if len(builder.deltas[fragment.DeltaID]) > builder.k {
		// join the shares to a delta
		dlt, err := builder.Join(builder.deltas[fragment.DeltaID]...)
		if err != nil {
			return nil, err
		}
		delete(builder.deltas, fragment.DeltaID)

		return dlt, nil
	}

	return nil, nil
}

func (builder *deltaBuilder) Join(deltaFragments ...Fragment) (*Delta, error) {
	if !IsCompatible(deltaFragments) {
		return nil, errors.New("delta fragment are not compatible with each other ")
	}

	for i := 0; i < builder.k; i++ {
		builder.tokenSharesCache[i] = deltaFragments[i].TokenShare
		builder.priceCoShareCache[i] = deltaFragments[i].PriceShare.Co
		builder.priceExpShareCache[i] = deltaFragments[i].PriceShare.Exp
		builder.volumeCoShareCache[i] = deltaFragments[i].VolumeShare.Co
		builder.volumeExpShareCache[i] = deltaFragments[i].VolumeShare.Exp
		builder.minVolumeCoShareCache[i] = deltaFragments[i].MinVolumeShare.Co
		builder.minVolumeExpShareCache[i] = deltaFragments[i].MinVolumeShare.Exp
	}

	return NewDeltaFromShares(
		deltaFragments[0].BuyOrderID,
		deltaFragments[0].SellOrderID,
		builder.tokenSharesCache,
		builder.priceCoShareCache,
		builder.priceExpShareCache,
		builder.volumeCoShareCache,
		builder.volumeExpShareCache,
		builder.minVolumeCoShareCache,
		builder.minVolumeExpShareCache,
	), nil
}
