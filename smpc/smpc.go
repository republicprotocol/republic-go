package smpc

import (
	"context"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
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

// Close all channels. Calling this function more than once will cause a panic.
func (chs *ObscureComputeInput) Close() {
	close(chs.Rng)
	close(chs.RngShares)
	close(chs.RngSharesIndexed)
	close(chs.MulShares)
	close(chs.MulSharesIndexed)
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

// OrderMatchComputeInput accepted by order matching computations. It stores a
// set of read-write channels.
type OrderMatchComputeInput struct {
	OrderFragments chan order.Fragment
	DeltaFragments chan DeltaFragment
}

// Close all channels. Calling this function more than once will cause a panic.
func (chs *OrderMatchComputeInput) Close() {
	close(chs.OrderFragments)
	close(chs.DeltaFragments)
}

// OrderMatchComputeOutput returned by obscure computations. It stores a set of
// read-only channels.
type OrderMatchComputeOutput struct {
	DeltaFragments <-chan DeltaFragment
	Deltas         <-chan Delta
}

// Computer of sMPC messages.
type Computer struct {
	id identity.ID

	n, k                      int64
	sharedOrderTable          SharedOrderTable
	sharedObscureResidueTable SharedObscureResidueTable
	sharedDeltaBuilder        SharedDeltaBuilder
}

// NewComputer returns a new Computer with the given ComputerID and N-K
// threshold.
func NewComputer(id identity.ID, n, k int64) Computer {
	return Computer{
		id:                        id,
		n:                         n,
		k:                         k,
		sharedOrderTable:          NewSharedOrderTable(),
		sharedObscureResidueTable: NewSharedObscureResidueTable(id),
		sharedDeltaBuilder:        NewSharedDeltaBuilder(k, Prime),
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
	obscureRngCh, errCh := ProduceObscureRngs(ctx, computer.n, computer.k, &computer.sharedObscureResidueTable, int(computer.n))
	errChs[0] = errCh

	// ProcessObscureRngs that were initiated by other Computers and broadcast
	// to this Computer
	obscureRngSharesCh, errCh := ProcessObscureRngs(ctx, obscureComputeChs.Rng, &computer.sharedObscureResidueTable, int(computer.n))
	errChs[1] = errCh

	// ProcessObscureRngShares broadcast to this Computer in response to an
	// ObscureRng that was produced by this Computer
	obscureRngSharesIndexedCh, errCh := ProcessObscureRngShares(ctx, obscureComputeChs.RngShares, int(computer.n))
	errChs[2] = errCh

	// ProcessObscureRngSharesIndexed broadcast to this Computer by another
	// Computer that is progressing through the creation of its obscure residue
	obscureMulShares, errCh := ProcessObscureRngSharesIndexed(ctx, obscureComputeChs.RngSharesIndexed, int(computer.n))
	errChs[3] = errCh

	// ProcessObscureMulShares broadcast to this Computer in response to
	// ObscureRngSharesIndexed that were produced by this Computer
	obscureMulSharesIndexedCh, errCh := ProcessObscureMulShares(ctx, obscureComputeChs.MulShares, int(computer.n))
	errChs[4] = errCh

	// ProcessObscureMulSharesIndexed broadcast to this Computer by another
	// Computer that is progressing through the creation of its obscure residue
	obscureResidueFragmentCh, errCh := ProcessObscureMulSharesIndexed(ctx, obscureComputeChs.MulSharesIndexed, int(computer.n))
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

// ComputeOrderMatches using order fragments and delta fragments.
func (computer *Computer) ComputeOrderMatches(done <-chan struct{}, orderFragmentsIn <-chan order.Fragment, deltaFragmentsIn <-chan DeltaFragment) (<-chan DeltaFragment, <-chan Delta) {
	deltaFragmentsOut := make(chan DeltaFragment)
	deltasOut := make(chan Delta)

	go func() {
		defer close(deltaFragmentsOut)
		defer close(deltasOut)

		orderTuples := OrderFragmentsToOrderTuples(done, orderFragmentsIn, &computer.sharedOrderTable, 100)
		deltaFragmentsFromOrderTuples := OrderTuplesToDeltaFragments(done, orderTuples, 100)

		deltaFragments := make(chan DeltaFragment)
		go func() {
			defer close(deltaFragments)
			dispatch.CoBegin(func() {
				dispatch.Split(deltaFragmentsFromOrderTuples, deltaFragments, deltaFragmentsOut)
			}, func() {
				dispatch.Pipe(done, deltaFragmentsIn, deltaFragments)
			})
		}()

		deltas := BuildDeltas(done, deltaFragments, &computer.sharedDeltaBuilder, 100)
		dispatch.Pipe(done, deltas, deltasOut)
	}()

	return deltaFragmentsOut, deltasOut
}

func (computer *Computer) SharedOrderTable() *SharedOrderTable {
	return &computer.sharedOrderTable
}

func (computer *Computer) SharedObscureResidueTable() *SharedObscureResidueTable {
	return &computer.sharedObscureResidueTable
}
