package delta

import (
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

// A ID is the Keccak256 hash of the order IDs that were used to compute
// the associated Delta.
type ID []byte

// Equal returns an equality check between two DeltaIDs.
func (id ID) Equal(other ID) bool {
	return bytes.Equal(id, other)
}

// String returns a ID as a Base58 encoded string.
func (id ID) String() string {
	return base58.Encode(id)
}

type Deltas []Delta

type Delta struct {
	ID          ID
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
	delta.ID = ID(crypto.Keccak256(delta.BuyOrderID[:], delta.SellOrderID[:]))
	return delta
}

func (delta *Delta) IsMatch(prime stackint.Int1024) bool {
	// zero := stackint.Zero()
	two := stackint.Two()
	zeroThreshold := prime.Div(&two)

	// TODO: Use real tokens
	// if delta.FstCode.Cmp(&zero) != 0 {
	// 	return false
	// }
	// if delta.SndCode.Cmp(&zero) != 0 {
	// 	return false
	// }
	if delta.Price.Cmp(&zeroThreshold) == 1 {
		return false
	}
	// TODO: Unify max volume
	// if delta.MaxVolume.Cmp(&zeroThreshold) == 1 {
	// 	return false
	// }
	// if delta.MinVolume.Cmp(&zeroThreshold) == 1 {
	// 	return false
	// }
	return true
}
