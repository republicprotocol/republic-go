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

type Engine struct {
	n, k             int64
	sharedOrderTable SharedOrderTable
}

func NewEngine(n, k int64) Engine {
	return Engine{
		sharedOrderTable: NewSharedOrderTable(),
	}
}

func (engine *Engine) Compute(
	ctx context.Context,
	obscureRngChIn <-chan ObscureRng,
	obscureRngSharesChIn <-chan ObscureRngShares,
	obscureRngSharesIndexedChIn <-chan ObscureRngSharesIndexed,
	obscureMulSharesChIn <-chan ObscureMulShares,
	obscureMulSharesIndexedChIn <-chan ObscureMulSharesIndexed,
) (
	<-chan ObscureRng,
	<-chan ObscureRngShares,
	<-chan ObscureRngSharesIndexed,
	<-chan ObscureMulShares,
	<-chan ObscureMulSharesIndexed,
	<-chan error,
) {
	errChs := make(chan (<-chan error), 6)

	obscureRngCh, errCh := ProduceObscureRngs(ctx, engine.n, engine.k)
	errChs <- errCh

	obscureRngSharesCh, errCh := ProcessObscureRngs(ctx, obscureRngChIn)
	errChs <- errCh

	obscureRngSharesIndexedCh, errCh := ProcessObscureRngShares(ctx, obscureRngSharesChIn)
	errChs <- errCh

	obscureMulShares, errCh := ProcessObscureRngSharesIndexed(ctx, obscureRngSharesIndexedChIn)
	errChs <- errCh

	obscureMulSharesIndexedCh, errCh := ProcessObscureMulShares(ctx, obscureMulSharesChIn)
	errChs <- errCh

	obscureResidueFragmentCh, errCh := ProcessObscureMulSharesIndexed(ctx, obscureMulSharesIndexedChIn)
	errChs <- errCh

	return obscureRngCh, obscureRngSharesCh, obscureRngSharesIndexedCh, obscureMulShares, obscureMulSharesIndexedCh, dispatch.MergeErrors(errChs)
}
