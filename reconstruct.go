package compute

import (
	"math/big"
	"sync"

	"github.com/republicprotocol/go-sss"

	"github.com/ethereum/go-ethereum/crypto"
)

type ReconstructionID []byte

type Reconstruction struct {
	ID        ReconstructionID
	FstCode   *big.Int
	SndCode   *big.Int
	Price     *big.Int
	MaxVolume *big.Int
	MinVolume *big.Int
}

func (reconstruction *Reconstruction) IsMatch() bool {
	// TODO: Do something sensible here like comparing prices and max volumes and min volumes
	return true
}

type ReconstructionMatrix struct {
	reconstructionsMu *sync.Mutex
	reconstructions   map[string][]*OrderFragment
}

func NewReconstructionMatrix() *ReconstructionMatrix {
	return &ReconstructionMatrix{
		reconstructionsMu: new(sync.Mutex),
		reconstructions:   map[string][]*OrderFragment{},
	}
}

func (matrix *ReconstructionMatrix) AddOrderFragment(orderFragment *OrderFragment, k int, prime *big.Int) (*Reconstruction, error) {
	matrix.reconstructionsMu.Lock()
	defer matrix.reconstructionsMu.Unlock()

	orderIDs := make([][]byte, len(orderFragment.OrderIDs))
	for i := range orderFragment.OrderIDs {
		orderIDs[i] = orderFragment.OrderIDs[i][:]
	}
	reconstructionID := ReconstructionID(crypto.Keccak256(orderIDs...))
	matrix.reconstructions[string(reconstructionID)] = append(matrix.reconstructions[string(reconstructionID)], orderFragment)
	if len(matrix.reconstructions[string(reconstructionID)]) < k {
		return nil, nil
	}

	var err error
	reconstruction := &Reconstruction{
		ID: reconstructionID,
	}

	fstCodeShares := make(sss.Shares, k)
	sndCodeShares := make(sss.Shares, k)
	priceShares := make(sss.Shares, k)
	maxVolumeShares := make(sss.Shares, k)
	minVolumeShares := make(sss.Shares, k)
	for i := range matrix.reconstructions[string(reconstructionID)] {
		orderFragment := matrix.reconstructions[string(reconstructionID)][i]
		fstCodeShares[i] = orderFragment.FstCodeShare
		sndCodeShares[i] = orderFragment.SndCodeShare
		priceShares[i] = orderFragment.PriceShare
		maxVolumeShares[i] = orderFragment.MaxVolumeShare
		minVolumeShares[i] = orderFragment.MinVolumeShare
	}
	if reconstruction.FstCode, err = sss.Join(prime, fstCodeShares); err != nil {
		return nil, err
	}
	if reconstruction.SndCode, err = sss.Join(prime, sndCodeShares); err != nil {
		return nil, err
	}
	if reconstruction.Price, err = sss.Join(prime, priceShares); err != nil {
		return nil, err
	}
	if reconstruction.MaxVolume, err = sss.Join(prime, maxVolumeShares); err != nil {
		return nil, err
	}
	if reconstruction.MinVolume, err = sss.Join(prime, minVolumeShares); err != nil {
		return nil, err
	}
	return reconstruction, nil
}
