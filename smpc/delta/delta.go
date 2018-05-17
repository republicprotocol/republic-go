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
	Volumn      order.CoExp
	MinVolume   order.CoExp
}

func NewDeltaFromShares(buyOrderID, sellOrderID order.ID, tokenShares shamir.Shares, priceShares, volumeShares, minVolumeShares []order.CoExpShare) Delta {
	delta := Delta{
		BuyOrderID:  buyOrderID,
		SellOrderID: sellOrderID,
	}
	delta.Tokens = shamir.Join(tokenShares)

	priceCos := make ([]shamir.Share, len(priceShares))
	priceExps := make ([]shamir.Share, len(priceShares))
	for i := range priceShares{
		priceCos[i] = priceShares[i].Co
		priceExps[i] = priceShares[i].Exp
	}
	delta.Price.Co = uint32(shamir.Join(priceCos))
	delta.Price.Exp = uint32(shamir.Join(priceExps))

	volumeCos := make ([]shamir.Share, len(volumeShares))
	volumeExps := make ([]shamir.Share, len(volumeShares))
	for i := range volumeShares{
		volumeCos[i] = volumeShares[i].Co
		volumeExps[i] = volumeShares[i].Exp
	}
	delta.Volumn.Co = uint32(shamir.Join(volumeCos))
	delta.Volumn.Exp = uint32(shamir.Join(volumeExps))

	minVolumeCos := make ([]shamir.Share, len(minVolumeShares))
	minVolumeExps := make ([]shamir.Share, len(minVolumeShares))
	for i := range minVolumeShares{
		minVolumeCos[i] = minVolumeShares[i].Co
		minVolumeExps[i] = minVolumeShares[i].Exp
	}
	delta.MinVolume.Co = uint32(shamir.Join(minVolumeCos))
	delta.MinVolume.Exp = uint32(shamir.Join(minVolumeExps))

	var ID [32]byte
	// Compute the ResultID and return the Result.
	copy(ID[:], crypto.Keccak256(delta.BuyOrderID[:], delta.SellOrderID[:]))
	delta.ID = ID

	return delta
}

func (delta *Delta) IsMatch() bool {
	zeroThreshold := shamir.Prime / 2

	if delta.Tokens != 0 {
		return false
	}
	if uint64(delta.Price.Exp) >= zeroThreshold {
		return false
	}
	if uint64(delta.Price.Co) >= zeroThreshold{
		return false
	}
	if uint64(delta.Volumn.Exp) >= zeroThreshold {
		return false
	}
	if uint64(delta.Volumn.Co) >= zeroThreshold{
		return false
	}
	if uint64(delta.MinVolume.Exp) >= zeroThreshold {
		return false
	}
	if uint64(delta.MinVolume.Co) >= zeroThreshold{
		return false
	}

	return true
}
