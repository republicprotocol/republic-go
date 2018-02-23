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

func (shard Shard) Compute(prime *big.Int) FinalShard {
	return NewFinalShard()
}
