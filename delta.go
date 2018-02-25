package compute

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-sss"
)

type DeltaShard struct {
	Signature []byte
	Finals    []*DeltaFragment
}

func NewDeltaShard() DeltaShard {
	return DeltaShard{
		Finals: []*DeltaFragment{},
	}
}

// A DeltaID is the Keccak256 hash of the OrderIDs that were used to compute
// the respective Result.
type DeltaID []byte

// Equals checks if two ResultIDs are equal in value.
func (id DeltaID) Equals(other DeltaID) bool {
	return bytes.Equal(id, other)
}

// String returns the ResultID as a string.
func (id DeltaID) String() string {
	return string(id)
}

// A Delta is the publicly computed value of comparing two Orders.
type Delta struct {
	ID          DeltaID
	BuyOrderID  OrderID
	SellOrderID OrderID
	FstCode     *big.Int
	SndCode     *big.Int
	Price       *big.Int
	MaxVolume   *big.Int
	MinVolume   *big.Int
}

func NewDelta(deltaFragments []*DeltaFragment, prime *big.Int) (*Delta, error) {
	// Check that all ResultFragments are compatible with each other.
	err := isCompatible(deltaFragments)
	if err != nil {
		return nil, err
	}

	// Collect sss.Shares across all ResultFragments.
	k := len(deltaFragments)
	fstCodeShares := make(sss.Shares, k)
	sndCodeShares := make(sss.Shares, k)
	priceShares := make(sss.Shares, k)
	maxVolumeShares := make(sss.Shares, k)
	minVolumeShares := make(sss.Shares, k)
	for i, resultFragment := range deltaFragments {
		fstCodeShares[i] = resultFragment.FstCodeShare
		sndCodeShares[i] = resultFragment.SndCodeShare
		priceShares[i] = resultFragment.PriceShare
		maxVolumeShares[i] = resultFragment.MaxVolumeShare
		minVolumeShares[i] = resultFragment.MinVolumeShare
	}

	// Join the sss.Shares into a Result.
	delta := &Delta{
		BuyOrderID:  deltaFragments[0].BuyOrderID,
		SellOrderID: deltaFragments[0].SellOrderID,
	}
	delta.FstCode = sss.Join(prime, fstCodeShares)
	delta.SndCode = sss.Join(prime, sndCodeShares)
	delta.Price = sss.Join(prime, priceShares)
	delta.MaxVolume = sss.Join(prime, maxVolumeShares)
	delta.MinVolume = sss.Join(prime, minVolumeShares)

	// Compute the ResultID and return the Result.
	delta.ID = DeltaID(crypto.Keccak256(delta.BuyOrderID[:], delta.SellOrderID[:]))
	return delta, nil
}

func (delta *Delta) IsMatch(prime *big.Int) bool {
	zeroThreshold := big.NewInt(0).Div(prime, big.NewInt(2))
	if delta.FstCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if delta.SndCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if delta.Price.Cmp(zeroThreshold) == 1 {
		return false
	}
	if delta.MaxVolume.Cmp(zeroThreshold) == 1 {
		return false
	}
	if delta.MinVolume.Cmp(zeroThreshold) == 1 {
		return false
	}
	return true
}

// A DeltaFragmentID is the Keccak256 hash of its OrderFragmentIDs.
type DeltaFragmentID []byte

// Equals checks if two ResultFragmentIDs are equal in value.
func (id DeltaFragmentID) Equals(other DeltaFragmentID) bool {
	return bytes.Equal(id, other)
}

// String returns the ResultFragmentID as a string.
func (id DeltaFragmentID) String() string {
	return string(id)
}

// A DeltaFragment is a secret share of a Result. Is is performing a
// computation over two OrderFragments.
type DeltaFragment struct {
	// Public data.
	ID                  DeltaFragmentID
	BuyOrderID          OrderID
	SellOrderID         OrderID
	BuyOrderFragmentID  OrderFragmentID
	SellOrderFragmentID OrderFragmentID

	// Private data.
	FstCodeShare   sss.Share
	SndCodeShare   sss.Share
	PriceShare     sss.Share
	MaxVolumeShare sss.Share
	MinVolumeShare sss.Share
}

func NewDeltaFragment(left *OrderFragment, right *OrderFragment, prime *big.Int) (*DeltaFragment, error) {
	if err := left.IsCompatible(right); err != nil {
		return nil, err
	}
	var buyOrderFragment *OrderFragment
	var sellOrderFragment *OrderFragment
	if left.OrderParity == OrderParityBuy {
		buyOrderFragment = left
		sellOrderFragment = right
	} else {
		buyOrderFragment = right
		sellOrderFragment = left
	}

	deltaFragment, err := buyOrderFragment.Sub(sellOrderFragment, prime)
	if err != nil {
		return nil, err
	}
	return deltaFragment, nil
}

// Equals checks if two ResultFragments are equal in value.
func (deltaFragment *DeltaFragment) Equals(other *DeltaFragment) bool {
	return deltaFragment.ID.Equals(other.ID) &&
		deltaFragment.BuyOrderID.Equals(other.BuyOrderID) &&
		deltaFragment.SellOrderID.Equals(other.SellOrderID) &&
		deltaFragment.BuyOrderFragmentID.Equals(other.BuyOrderFragmentID) &&
		deltaFragment.SellOrderFragmentID.Equals(other.SellOrderFragmentID) &&
		deltaFragment.FstCodeShare.Key == other.FstCodeShare.Key &&
		deltaFragment.FstCodeShare.Value.Cmp(other.FstCodeShare.Value) == 0 &&
		deltaFragment.SndCodeShare.Key == other.SndCodeShare.Key &&
		deltaFragment.SndCodeShare.Value.Cmp(other.SndCodeShare.Value) == 0 &&
		deltaFragment.PriceShare.Key == other.PriceShare.Key &&
		deltaFragment.PriceShare.Value.Cmp(other.PriceShare.Value) == 0 &&
		deltaFragment.MaxVolumeShare.Key == other.MaxVolumeShare.Key &&
		deltaFragment.MaxVolumeShare.Value.Cmp(other.MaxVolumeShare.Value) == 0 &&
		deltaFragment.MinVolumeShare.Key == other.MinVolumeShare.Key &&
		deltaFragment.MinVolumeShare.Value.Cmp(other.MinVolumeShare.Value) == 0
}

// IsCompatible returns an error when the two ResultFragment do not have
// the same share indices.
func isCompatible(deltaFragments []*DeltaFragment) error {
	if len(deltaFragments) == 0 {
		return NewEmptySliceError("result fragments")
	}
	buyOrderID := deltaFragments[0].BuyOrderID
	sellOrderID := deltaFragments[0].SellOrderID
	for i := range deltaFragments {
		if !deltaFragments[i].BuyOrderID.Equals(buyOrderID) {
			return NewOrderFragmentationError(0, int64(i))
		}
		if !deltaFragments[i].SellOrderID.Equals(sellOrderID) {
			return NewOrderFragmentationError(0, int64(i))
		}
	}
	return nil
}

type DeltaEngine struct {
	deltaFragmentMap map[string][]*DeltaFragment
	deltaMap         map[string]*Delta
}

func NewDeltaEngine() *DeltaEngine {
	return &DeltaEngine{
		deltaFragmentMap: map[string][]*DeltaFragment{},
		deltaMap:         map[string]*Delta{},
	}
}

func (engine DeltaEngine) AddDeltaFragments(deltaFragment *DeltaFragment, k int64, prime *big.Int) (*Delta, error) {
	deltaID := DeltaID(crypto.Keccak256(deltaFragment.BuyOrderID[:], deltaFragment.SellOrderID[:]))

	// If the delta for this delta fragment has already been reconstructed then
	// return nothing, the engine must have already noted the deltas
	// reconstruction.
	if delta, ok := engine.deltaMap[deltaID.String()]; ok && delta != nil {
		return nil, nil
	}

	// Check that this delta fragment has not been collected yet.
	deltaFragmentIsUnique := true
	for _, candidate := range engine.deltaFragmentMap[deltaID.String()] {
		if candidate.ID.Equals(deltaFragment.ID) {
			deltaFragmentIsUnique = false
			break
		}
	}
	if deltaFragmentIsUnique {
		engine.deltaFragmentMap[deltaID.String()] = append(engine.deltaFragmentMap[deltaID.String()], deltaFragment)
	}

	// Check if we can reconstruct a new delta.
	if int64(len(engine.deltaFragmentMap[deltaID.String()])) >= k {
		delta, err := NewDelta(engine.deltaFragmentMap[deltaID.String()], prime)
		if err != nil {
			return nil, err
		}
		engine.deltaMap[deltaID.String()] = delta
		return delta, nil
	}

	return nil, nil
}
