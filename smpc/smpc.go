package smpc

import (
	"context"
	"sync"

	"github.com/republicprotocol/republic-go/delta"
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

// ComputerID uniquely identifies the Smpc.
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
	DeltaFragments chan delta.Fragment
}

// Close all channels. Calling this function more than once will cause a panic.
func (chs *OrderMatchComputeInput) Close() {
	close(chs.OrderFragments)
	close(chs.DeltaFragments)
}

// OrderMatchComputeOutput returned by obscure computations. It stores a set of
// read-only channels.
type OrderMatchComputeOutput struct {
	DeltaFragments <-chan delta.Fragment
	Deltas         <-chan delta.Delta
}

type Smpc struct {
	id identity.ID

	n, k                      int64
	sharedOrderTable          SharedOrderTable
	sharedObscureResidueTable SharedObscureResidueTable
	sharedDeltaBuilder        SharedDeltaBuilder
}

// NewSmpc returns a new Smpc with the given ComputerID and N-K
// threshold.
func NewSmpc(id identity.ID, n, k int64) Smpc {
	return Smpc{
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
func (computer *Smpc) ComputeObscure(
	ctx context.Context,
	obscureComputeChs ObscureComputeInput,
) (
	ObscureComputeOutput,
	<-chan error,
) {
	errChs := make([]<-chan error, 6)

	// ProduceObscureRngs to initiate the creation of an obscure residue that
	// will be owned by this Smpc
	obscureRngCh, errCh := ProduceObscureRngs(ctx, computer.n, computer.k, &computer.sharedObscureResidueTable, int(computer.n))
	errChs[0] = errCh

	// ProcessObscureRngs that were initiated by other Computers and broadcast
	// to this Smpc
	obscureRngSharesCh, errCh := ProcessObscureRngs(ctx, obscureComputeChs.Rng, &computer.sharedObscureResidueTable, int(computer.n))
	errChs[1] = errCh

	// ProcessObscureRngShares broadcast to this Smpc in response to an
	// ObscureRng that was produced by this Smpc
	obscureRngSharesIndexedCh, errCh := ProcessObscureRngShares(ctx, obscureComputeChs.RngShares, int(computer.n))
	errChs[2] = errCh

	// ProcessObscureRngSharesIndexed broadcast to this Smpc by another
	// Smpc that is progressing through the creation of its obscure residue
	obscureMulShares, errCh := ProcessObscureRngSharesIndexed(ctx, obscureComputeChs.RngSharesIndexed, int(computer.n))
	errChs[3] = errCh

	// ProcessObscureMulShares broadcast to this Smpc in response to
	// ObscureRngSharesIndexed that were produced by this Smpc
	obscureMulSharesIndexedCh, errCh := ProcessObscureMulShares(ctx, obscureComputeChs.MulShares, int(computer.n))
	errChs[4] = errCh

	// ProcessObscureMulSharesIndexed broadcast to this Smpc by another
	// Smpc that is progressing through the creation of its obscure residue
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
func (computer *Smpc) ComputeOrderMatches(done <-chan struct{}, orderFragmentsIn <-chan order.Fragment, deltaFragmentsIn <-chan delta.Fragment) (<-chan delta.Fragment, <-chan delta.Delta) {
	deltaFragmentsOut := make(chan delta.Fragment)
	deltasOut := make(chan delta.Delta)

	go func() {
		defer close(deltaFragmentsOut)
		defer close(deltasOut)

		orderTuples := OrderFragmentsToOrderTuples(done, orderFragmentsIn, &computer.sharedOrderTable, 100)
		deltaFragmentsComputed := OrderTuplesToDeltaFragments(done, orderTuples, 100)

		deltaFragments := make(chan delta.Fragment)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(deltaFragments)
			dispatch.CoBegin(func() {
				dispatch.Split(deltaFragmentsComputed, deltaFragments, deltaFragmentsOut)
			}, func() {
				dispatch.Pipe(done, deltaFragmentsIn, deltaFragments)
			})
		}()

		deltas := BuildDeltas(done, deltaFragments, &computer.sharedDeltaBuilder, 100)
		dispatch.Pipe(done, deltas, deltasOut)

		wg.Wait()
	}()

	return deltaFragmentsOut, deltasOut
}

func (computer *Smpc) SharedOrderTable() *SharedOrderTable {
	return &computer.sharedOrderTable
}

func (computer *Smpc) SharedObscureResidueTable() *SharedObscureResidueTable {
	return &computer.sharedObscureResidueTable
}

func (computer *Smpc) Prime() stackint.Int1024 {
	return Prime
}
