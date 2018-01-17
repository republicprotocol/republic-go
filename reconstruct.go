package compute

import (
	"math/big"
	"sync"

	"github.com/republicprotocol/go-sss"
)

type Reconstruction struct {
	ID        OrderFragmentID
	FstCode   *big.Int
	SndCode   *big.Int
	Price     *big.Int
	MaxVolume *big.Int
	MinVolume *big.Int
}

func (reconstruction *Reconstruction) IsMatch() bool {
	if reconstruction.FstCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if reconstruction.SndCode.Cmp(big.NewInt(0)) != 0 {
		return false
	}
	if reconstruction.Price.Cmp(big.NewInt(0)) != 1 {
		return false
	}
	if reconstruction.MaxVolume.Cmp(big.NewInt(0)) != 1 {
		return false
	}
	if reconstruction.MinVolume.Cmp(big.NewInt(0)) != 1 {
		return false
	}
	return true
}

type ReconstructionMatrix struct {
	reconstructionsMu *sync.Mutex
	reconstructions   map[string][]*ComputedOrderFragment
}

func NewReconstructionMatrix() *ReconstructionMatrix {
	return &ReconstructionMatrix{
		reconstructionsMu: new(sync.Mutex),
		reconstructions:   map[string][]*ComputedOrderFragment{},
	}
}

func (matrix *ReconstructionMatrix) AddComputedOrderFragment(computed *ComputedOrderFragment, k int, prime *big.Int) (*Reconstruction, error) {
	matrix.reconstructionsMu.Lock()
	defer matrix.reconstructionsMu.Unlock()

	matrix.reconstructions[string(computed.ID)] = append(matrix.reconstructions[string(computed.ID)], computed)
	if len(matrix.reconstructions[string(computed.ID)]) < k {
		return nil, nil
	}

	var err error
	reconstruction := &Reconstruction{
		ID: computed.ID,
	}

	fstCodeShares := make(sss.Shares, k)
	sndCodeShares := make(sss.Shares, k)
	priceShares := make(sss.Shares, k)
	maxVolumeShares := make(sss.Shares, k)
	minVolumeShares := make(sss.Shares, k)
	for i := range matrix.reconstructions[string(reconstruction.ID)] {
		orderFragment := matrix.reconstructions[string(reconstruction.ID)][i]
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
