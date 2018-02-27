package compute

import "math/big"

type Shard struct {
	Signature []byte
	Deltas    []*DeltaFragment
	Residues  []*ResidueFragment
}

func NewShard(deltas []*DeltaFragment, residues []*ResidueFragment) Shard {
	return Shard{
		Deltas:   deltas,
		Residues: residues,
	}
}

func (shard Shard) Compute(prime *big.Int) DeltaShard {
	deltaFragments := make([]*DeltaFragment, len(shard.Deltas))
	for i := range deltaFragments {
		deltaFragments[i] = deltaFragments[i]
	}
	return NewDeltaShard(deltaFragments)
}
