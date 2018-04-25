package delta

import (
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

// A FragmentID is the Keccak256 hash of the order IDs that were used to
// compute the associated Fragment.
type FragmentID []byte

// Equal returns an equality check between two DeltaFragmentIDs.
func (id FragmentID) Equal(other FragmentID) bool {
	return bytes.Equal(id, other)
}

// String returns a FragmentID as a Base58 encoded string.
func (id FragmentID) String() string {
	return base58.Encode(id)
}

type Fragments []Fragment

// A Fragment is a secret share of a Final. Is is performing a
// computation over two OrderFragments.
type Fragment struct {
	ID                  FragmentID
	DeltaID             ID
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

func NewDeltaFragment(left, right *order.Fragment, prime *stackint.Int1024) Fragment {
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
		Value: buyOrderFragment.FstCodeShare.Value.SubModulo(&sellOrderFragment.FstCodeShare.Value, prime),
	}
	sndCodeShare := shamir.Share{
		Key:   buyOrderFragment.SndCodeShare.Key,
		Value: buyOrderFragment.SndCodeShare.Value.SubModulo(&sellOrderFragment.SndCodeShare.Value, prime),
	}
	priceShare := shamir.Share{
		Key:   buyOrderFragment.PriceShare.Key,
		Value: buyOrderFragment.PriceShare.Value.SubModulo(&sellOrderFragment.PriceShare.Value, prime),
	}
	maxVolumeShare := shamir.Share{
		Key:   buyOrderFragment.MaxVolumeShare.Key,
		Value: buyOrderFragment.MaxVolumeShare.Value.SubModulo(&sellOrderFragment.MinVolumeShare.Value, prime),
	}
	minVolumeShare := shamir.Share{
		Key:   buyOrderFragment.MinVolumeShare.Key,
		Value: sellOrderFragment.MaxVolumeShare.Value.SubModulo(&buyOrderFragment.MinVolumeShare.Value, prime),
	}

	return Fragment{
		ID:                  FragmentID(crypto.Keccak256([]byte(buyOrderFragment.ID), []byte(sellOrderFragment.ID))),
		DeltaID:             ID(crypto.Keccak256([]byte(buyOrderFragment.OrderID), []byte(sellOrderFragment.OrderID))),
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

// Equals checks if two Fragments are equal in value.
func (deltaFragment *Fragment) Equals(other *Fragment) bool {
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

// IsCompatible returns true if all Fragments are fragments of the same
// Delta, otherwise it returns false.
func IsCompatible(deltaFragments Fragments) bool {
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
