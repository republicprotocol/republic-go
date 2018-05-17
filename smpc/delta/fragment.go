package delta

import (
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

// A FragmentID is the Keccak256 hash of the order IDs that were used to
// compute the associated Fragment.
type FragmentID [32]byte

// Equal returns an equality check between two DeltaFragmentIDs.
func (id FragmentID) Equal(other FragmentID) bool {
	return bytes.Equal(id[:], other[:])
}

// String returns a FragmentID as a Base58 encoded string.
func (id FragmentID) String() string {
	return base58.Encode(id[:])
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

	TokenShare     shamir.Share
	PriceShare     order.CoExpShare
	VolumeShare    order.CoExpShare
	MinVolumeShare order.CoExpShare
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

	token := buyOrderFragment.Tokens.Sub(&sellOrderFragment.Tokens)
	priceCo := buyOrderFragment.Price.Co.Sub(&sellOrderFragment.Price.Co)
	priceExp := buyOrderFragment.Price.Exp.Sub(&sellOrderFragment.Price.Exp)
	volumeCo := buyOrderFragment.Volume.Co.Sub(&sellOrderFragment.Volume.Co)
	volumeExp := buyOrderFragment.Volume.Exp.Sub(&sellOrderFragment.Volume.Exp)
	minVolumeCo := buyOrderFragment.MinimumVolume.Co.Sub(&sellOrderFragment.MinimumVolume.Co)
	minVolumeExp := buyOrderFragment.MinimumVolume.Exp.Sub(&sellOrderFragment.MinimumVolume.Exp)
	var fragmentID, deltaID [32]byte
	copy(fragmentID[:], crypto.Keccak256(buyOrderFragment.ID[:], sellOrderFragment.ID[:]))
	copy(deltaID[:], crypto.Keccak256(buyOrderFragment.OrderID[:], sellOrderFragment.OrderID[:]))

	return Fragment{
		ID: fragmentID,
		DeltaID:             deltaID,
		BuyOrderID:          buyOrderFragment.OrderID,
		SellOrderID:         sellOrderFragment.OrderID,
		BuyOrderFragmentID:  buyOrderFragment.ID,
		SellOrderFragmentID: sellOrderFragment.ID,
		TokenShare:        token,
		PriceShare :  order.CoExpShare{
			Co:priceCo,
			Exp:priceExp,
		},
		VolumeShare:        order.CoExpShare{
			Co:volumeCo,
			Exp:volumeExp,
		},
		MinVolumeShare:       order.CoExpShare{
			Co:minVolumeCo,
			Exp:minVolumeExp,
		},
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
		deltaFragment.TokenShare.Equal(&other.TokenShare) &&
		deltaFragment.PriceShare.Equal(&other.PriceShare) &&
		deltaFragment.VolumeShare.Equal(&other.VolumeShare) &&
		deltaFragment.MinVolumeShare.Equal(&other.MinVolumeShare)
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
