package compute

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-sss"
)

type ResultID []byte

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

func NewResult(prime *big.Int, resultFragments []*ResultFragment) (*Result, error) {
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
	if result.FstCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if result.SndCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if result.Price.Cmp(big.NewInt(0)) == -1 {
		return false
	}
	if result.MaxVolume.Cmp(big.NewInt(0)) == -1 {
		return false
	}
	if result.MinVolume.Cmp(big.NewInt(0)) == -1 {
		return false
	}
	return true
}

// A ResultFragmentID is the Keccak256 hash of its OrderFragmentIDs.
type ResultFragmentID []byte

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
