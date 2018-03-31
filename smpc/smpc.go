package smpc

type ComputeEngine struct {
	workers             Workers
	broadcasters        Broadcasters
	deltaFragmentMatrix DeltaFragmentMatrix
	deltaBuilder        DeltaBuilder
	deltaOutput         chan Delta
}

func NewComputeEngine() ComputeEngine {
	return ComputeEngine{}
}
