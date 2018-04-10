package smpc

import (
	"context"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/stackint"
)

// Prime used to define the finite field for all computations.
var Prime = func() stackint.Int1024 {
	prime, err := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	if err != nil {
		panic(err)
	}
	return prime
}()

// ComputerID uniquely identifies the Computer.
type ComputerID [32]byte

// ObscureComputeInput accepted by obscure computations. It stores a set of
// read-write channels.
type ObscureComputeInput struct {
	Rng              chan ObscureRng
	RngShares        chan ObscureRngShares
	RngSharesIndexed chan ObscureRngSharesIndexed
	MulShares        chan ObscureMulShares
	MulSharesIndexed chan ObscureMulSharesIndexed
}

// ObscureComputeOutput returned by obscure computations. It stores a set of
// read-only channels.
type ObscureComputeOutput struct {
	Rng              <-chan ObscureRng
	RngShares        <-chan ObscureRngShares
	RngSharesIndexed <-chan ObscureRngSharesIndexed
	MulShares        <-chan ObscureMulShares
	MulSharesIndexed <-chan ObscureMulSharesIndexed
}

// Computer of sMPC messages.
type Computer struct {
	ComputerID

	n, k                      int64
	sharedOrderTable          SharedOrderTable
	sharedObscureResidueTable SharedObscureResidueTable
}

// NewComputer returns a new Computer with the given ComputerID and N-K
// threshold.
func NewComputer(computerID ComputerID, n, k int64) Computer {
	return Computer{
		ComputerID:                computerID,
		sharedOrderTable:          NewSharedOrderTable(),
		sharedObscureResidueTable: NewSharedObscureResidueTable(computerID),
	}
}

// ComputeObscure residues that will be used during order matching to obscure
// Deltas.
func (computer *Computer) ComputeObscure(
	ctx context.Context,
	obscureComputeChs ObscureComputeInput,
) (
	ObscureComputeOutput,
	<-chan error,
) {
	errChs := make([]<-chan error, 6)

	// ProduceObscureRngs to initiate the creation of an obscure residue that
	// will be owned by this Computer
	obscureRngCh, errCh := ProduceObscureRngs(ctx, computer.n, computer.k, &computer.sharedObscureResidueTable)
	errChs[0] = errCh

	// ProcessObscureRngs that were initiated by other Computers and broadcast
	// to this Computer
	obscureRngSharesCh, errCh := ProcessObscureRngs(ctx, obscureComputeChs.Rng, &computer.sharedObscureResidueTable)
	errChs[1] = errCh

	// ProcessObscureRngShares broadcast to this Computer in response to an
	// ObscureRng that was produced by this Computer
	obscureRngSharesIndexedCh, errCh := ProcessObscureRngShares(ctx, obscureComputeChs.RngShares)
	errChs[2] = errCh

	// ProcessObscureRngSharesIndexed broadcast to this Computer by another
	// Computer that is progressing through the creation of its obscure residue
	obscureMulShares, errCh := ProcessObscureRngSharesIndexed(ctx, obscureComputeChs.RngSharesIndexed)
	errChs[3] = errCh

	// ProcessObscureMulShares broadcast to this Computer in response to
	// ObscureRngSharesIndexed that were produced by this Computer
	obscureMulSharesIndexedCh, errCh := ProcessObscureMulShares(ctx, obscureComputeChs.MulShares)
	errChs[4] = errCh

	// ProcessObscureMulSharesIndexed broadcast to this Computer by another
	// Computer that is progressing through the creation of its obscure residue
	obscureResidueFragmentCh, errCh := ProcessObscureMulSharesIndexed(ctx, obscureComputeChs.MulSharesIndexed)
	errChs[5] = errCh

	// Consume all ObscureResidueFragments and store them in the
	// SharedObscureResidueTable, assuming that the owner has already been
	// stored
	go func() {
		for obscureResidueFragment := range obscureResidueFragmentCh {
			computer.sharedObscureResidueTable.InsertObscureResidue(obscureResidueFragment)
		}
	}()

	allowErrCanceled := true
	errCh = dispatch.FilterErrors(dispatch.MergeErrors(errChs...), func(err error) bool {
		if err == context.Canceled && allowErrCanceled {
			allowErrCanceled = false
			return true
		}
		return false
	})

	return ObscureComputeOutput{
		obscureRngCh,
		obscureRngSharesCh,
		obscureRngSharesIndexedCh,
		obscureMulShares,
		obscureMulSharesIndexedCh,
	}, errCh
}

func (computer *Computer) ComputeOrderMatches() {
}

func (computer *Computer) SharedOrderTable() *SharedOrderTable {
	return &computer.sharedOrderTable
}

func (computer *Computer) SharedObscureResidueTable() *SharedObscureResidueTable {
	return &computer.sharedObscureResidueTable
}
