package compute

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-sss"
)

// A ResultID is the Keccak256 hash of the OrderIDs that were used to compute
// the respective Result.
type ResultID []byte

// Equals checks if two ResultIDs are equal in value.
func (id ResultID) Equals(other ResultID) bool {
	return bytes.Equal(id, other)
}

// String returns the ResultID as a string.
func (id ResultID) String() string {
	return string(id)
}

// A Result is the publicly computed value of comparing two Orders.
type Result struct {
	ID          ResultID
	BuyOrderID  OrderID
	SellOrderID OrderID
	FstCode     *big.Int
	SndCode     *big.Int
	Price       *big.Int
	MaxVolume   *big.Int
	MinVolume   *big.Int
}

func NewResult(resultFragments []*ResultFragment, prime *big.Int) *Result {
	// Collect sss.Shares across all ResultFragments.
	// TODO: Check that all ResultFragments are compatible with each other.
	k := len(resultFragments)
	fstCodeShares := make(sss.Shares, k)
	sndCodeShares := make(sss.Shares, k)
	priceShares := make(sss.Shares, k)
	maxVolumeShares := make(sss.Shares, k)
	minVolumeShares := make(sss.Shares, k)
	for i, resultFragment := range resultFragments {
		fstCodeShares[i] = resultFragment.FstCodeShare
		sndCodeShares[i] = resultFragment.SndCodeShare
		priceShares[i] = resultFragment.PriceShare
		maxVolumeShares[i] = resultFragment.MaxVolumeShare
		minVolumeShares[i] = resultFragment.MinVolumeShare
	}

	// Join the sss.Shares into a Result.
	// FIXME: This can panic if there are no ResultFragments.
	result := &Result{
		BuyOrderID:  resultFragments[0].BuyOrderID,
		SellOrderID: resultFragments[0].SellOrderID,
	}
	result.FstCode = sss.Join(prime, fstCodeShares)
	result.SndCode = sss.Join(prime, sndCodeShares)
	result.Price = sss.Join(prime, priceShares)
	result.MaxVolume = sss.Join(prime, maxVolumeShares)
	result.MinVolume = sss.Join(prime, minVolumeShares)

	// Compute the ResultID and return the Result.
	result.ID = ResultID(crypto.Keccak256(result.BuyOrderID[:], result.SellOrderID[:]))
	return result
}

func (result *Result) IsMatch(prime *big.Int) bool {
	zeroThreshold := big.NewInt(0).Div(prime, big.NewInt(2))
	if result.FstCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if result.SndCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if result.Price.Cmp(zeroThreshold) == 1 {
		return false
	}
	if result.MaxVolume.Cmp(zeroThreshold) == 1 {
		return false
	}
	if result.MinVolume.Cmp(zeroThreshold) == 1 {
		return false
	}
	return true
}

// Results is an array of Result
type Results []Result

// A ResultFragmentID is the Keccak256 hash of its OrderFragmentIDs.
type ResultFragmentID []byte

// Equals checks if two ResultFragmentIDs are equal in value.
func (id ResultFragmentID) Equals(other ResultFragmentID) bool {
	return bytes.Equal(id, other)
}

// String returns the ResultFragmentID as a string.
func (id ResultFragmentID) String() string {
	return string(id)
}

// A ResultFragment is a secret share of a Result. Is is performing a
// computation over two OrderFragments.
type ResultFragment struct {
	// Public data.
	ID                  ResultFragmentID
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

func NewResultFragment(buyOrderID, sellOrderID OrderID, buyOrderFragmentID, sellOrderFragmentID OrderFragmentID, fstCodeShare, sndCodeShare, priceShare, maxVolumeShare, minVolumeShare sss.Share) *ResultFragment {
	resultFragment := &ResultFragment{
		BuyOrderID:          buyOrderID,
		SellOrderID:         sellOrderID,
		BuyOrderFragmentID:  buyOrderFragmentID,
		SellOrderFragmentID: sellOrderFragmentID,
		FstCodeShare:        fstCodeShare,
		SndCodeShare:        sndCodeShare,
		PriceShare:          priceShare,
		MaxVolumeShare:      maxVolumeShare,
		MinVolumeShare:      minVolumeShare,
	}
	resultFragment.ID = ResultFragmentID(crypto.Keccak256(resultFragment.BuyOrderFragmentID[:], resultFragment.SellOrderFragmentID[:]))
	return resultFragment
}

// Equals checks if two ResultFragments are equal in value.
func (resultFragment *ResultFragment) Equals(other *ResultFragment) bool {
	return resultFragment.ID.Equals(other.ID) &&
		resultFragment.BuyOrderID.Equals(other.BuyOrderID) &&
		resultFragment.SellOrderID.Equals(other.SellOrderID) &&
		resultFragment.BuyOrderFragmentID.Equals(other.BuyOrderFragmentID) &&
		resultFragment.SellOrderFragmentID.Equals(other.SellOrderFragmentID) &&
		resultFragment.FstCodeShare.Key == other.FstCodeShare.Key &&
		resultFragment.FstCodeShare.Value.Cmp(other.FstCodeShare.Value) == 0 &&
		resultFragment.SndCodeShare.Key == other.SndCodeShare.Key &&
		resultFragment.SndCodeShare.Value.Cmp(other.SndCodeShare.Value) == 0 &&
		resultFragment.PriceShare.Key == other.PriceShare.Key &&
		resultFragment.PriceShare.Value.Cmp(other.PriceShare.Value) == 0 &&
		resultFragment.MaxVolumeShare.Key == other.MaxVolumeShare.Key &&
		resultFragment.MaxVolumeShare.Value.Cmp(other.MaxVolumeShare.Value) == 0 &&
		resultFragment.MinVolumeShare.Key == other.MinVolumeShare.Key &&
		resultFragment.MinVolumeShare.Value.Cmp(other.MinVolumeShare.Value) == 0
}
