package compute

type ResidueFragment struct {
	ID ResidueFragmentID
}

type ResidueFragmentID []byte

func NewResidueFragment() (*ResidueFragment, error) {
	return &ResidueFragment{}, nil
}
