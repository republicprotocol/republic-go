package compute

type Shard struct {
	Signature []byte
	Deltas    []DeltaFragment
	Residues  []ResidueFragment
}

func NewShard(deltas []DeltaFragment, residues []ResidueFragment) Shard {
	return Shard{
		Deltas:   deltas,
		Residues: residues,
	}
}

func (shard Shard) Compute() FinalShard {
	return FinalShard{
		Finals: []FinalFragment{},
	}
}

type FinalShard struct {
	Signature []byte
	Finals    []FinalFragment
}

type DeltaFragment struct {
}

type ResidueFragment struct {
}

type FinalFragment struct {
}
