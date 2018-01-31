package compute

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-sss"
)

// A ResultID is the Keccak256 hash of the OrderIDs that were used to compute
// the respective Result.
type ResultID []byte

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

func NewResult(resultFragments []*ResultFragment, prime *big.Int) (*Result, error) {
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
	var err error
	result := &Result{
		BuyOrderID:  resultFragments[0].BuyOrderID,
		SellOrderID: resultFragments[0].SellOrderID,
	}
	if result.FstCode, err = sss.Join(prime, fstCodeShares); err != nil {
		return nil, err
	}
	if result.SndCode, err = sss.Join(prime, sndCodeShares); err != nil {
		return nil, err
	}
	if result.Price, err = sss.Join(prime, priceShares); err != nil {
		return nil, err
	}
	if result.MaxVolume, err = sss.Join(prime, maxVolumeShares); err != nil {
		return nil, err
	}
	if result.MinVolume, err = sss.Join(prime, minVolumeShares); err != nil {
		return nil, err
	}

	// Compute the ResultID and return the Result.
	result.ID = ResultID(crypto.Keccak256(result.BuyOrderID[:], result.SellOrderID[:]))
	return result, nil
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

// Bytes returns a ResultFragment serialized into a bytes.
func (resultFragment *ResultFragment) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, resultFragment.BuyOrderID)
	binary.Write(buf, binary.LittleEndian, resultFragment.SellOrderID)
	binary.Write(buf, binary.LittleEndian, resultFragment.BuyOrderFragmentID)
	binary.Write(buf, binary.LittleEndian, resultFragment.SellOrderFragmentID)

	binary.Write(buf, binary.LittleEndian, resultFragment.FstCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, resultFragment.FstCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, resultFragment.SndCodeShare.Key)
	binary.Write(buf, binary.LittleEndian, resultFragment.SndCodeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, resultFragment.PriceShare.Key)
	binary.Write(buf, binary.LittleEndian, resultFragment.PriceShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, resultFragment.MaxVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, resultFragment.MaxVolumeShare.Value.Bytes())
	binary.Write(buf, binary.LittleEndian, resultFragment.MinVolumeShare.Key)
	binary.Write(buf, binary.LittleEndian, resultFragment.MinVolumeShare.Value.Bytes())

	return buf.Bytes()
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
