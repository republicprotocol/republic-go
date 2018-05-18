package delta

import (
	"bytes"

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
