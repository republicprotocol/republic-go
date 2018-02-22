package compute

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type DeltaID []byte

type Delta struct {
	ID            ComputationID
	DeltaFragment *DeltaFragment
}

func NewComputation(left *OrderFragment, right *OrderFragment) (*Computation, error) {
	if err := left.IsCompatible(right); err != nil {
		return nil, err
	}
	computation := &Computation{}
	if left.OrderParity == OrderParityBuy {
		computation.BuyOrderFragment = left
		computation.SellOrderFragment = right
	} else {
		computation.BuyOrderFragment = right
		computation.SellOrderFragment = left
	}
	computation.ID = ComputationID(crypto.Keccak256(computation.BuyOrderFragment.ID[:], computation.SellOrderFragment.ID[:]))
	return computation, nil
}

func (computation *Computation) Sub(prime *big.Int) (*ResultFragment, error) {
	return computation.BuyOrderFragment.Sub(computation.SellOrderFragment, prime)
}

type ComputationShardID []byte

type ComputationShard struct {
	ID           ComputationShardID
	Computations []*Computation
}

func NewComputationShard(computations []*Computation) ComputationShard {
	computationIDs := make([]byte, 0, len(computations)*32)
	for _, computation := range computations {
		computationIDs = append(computationIDs, []byte(computation.ID)...)
	}
	return ComputationShard{
		ID:           ComputationShardID(crypto.Keccak256(computationIDs)),
		Computations: computations,
	}
}

func (shard ComputationShard) Compute(prime *big.Int) []*ResultFragment {
	resultFragments := make([]*ResultFragment, len(shard.Computations))
	for i := range resultFragments {
		// FIXME: We are processing computations in bulk with the expectation
		// that some of them will fail (hopefully 2/3rds of participiants will
		// succeed). Errors are dropped here.
		resultFragments[i], _ = shard.Computations[i].Sub(prime)
	}
	return resultFragments
}

type ComputationBid int64

const (
	ComputationBidYes = 1
	ComputationBidNo  = 2
)

type ComputationShardBid struct {
	ID   ComputationShardID
	Bids map[string]ComputationBid
}
