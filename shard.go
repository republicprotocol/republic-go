package compute

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

func (shard Shard) Compute() DeltaShard {
	return NewDeltaShard()
}
