package smpc

import (
	"context"
	"time"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/stackint"
)

// OrderFragmentReceiver receives order Fragments from an input channel and
// uses a DeltaFragmentMatrix to compute new DeltaFragments with them.
// Cancelling the context will shutdown the DeltaFragmentReader.It returns an
// error, or nil.
func OrderFragmentReceiver(ctx context.Context, orderFragments chan order.Fragment, computationMatrix *ComputationMatrix) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case orderFragment, ok := <-orderFragments:
			if !ok {
				return nil
			}
			if orderFragment.OrderParity == order.ParityBuy {
				computationMatrix.InsertBuyOrder(orderFragment)
			} else {
				computationMatrix.InsertSellOrder(orderFragment)
			}
		}
	}
}

func DeltaFragmentComputer(ctx context.Context, computationMatrix *ComputationMatrix, bufferLimit int, prime stackint.Int1024) (chan DeltaFragment, chan error) {
	deltaFragments := make(chan DeltaFragment, bufferLimit)
	errors := make(chan error, bufferLimit)

	go func() {
		defer close(deltaFragments)
		defer close(errors)

		buffer := make([]Computation, bufferLimit)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			case <-ticker.C:
				for i, n := 0, computationMatrix.Computations(buffer[:]); i < n; i++ {
					// FIXME: NewDeltaFragment blocks, depending on what the
					// shares look like. This is probably to do with our custom
					// stackint library.
					deltaFragment := NewDeltaFragment(buffer[i].BuyOrderFragment, buffer[i].SellOrderFragment, prime)
					select {
					case <-ctx.Done():
						errors <- ctx.Err()
						return
					case deltaFragments <- deltaFragment:
					}
				}
			}
		}
	}()

	return deltaFragments, errors
}

// DeltaFragmentReceiver receives DeltaFragments from an input channel and uses
// a DeltaBuilder to build new Deltas with them. Cancelling the context will
// shutdown the DeltaFragmentReader. It returns an error, or nil.
func DeltaFragmentReceiver(ctx context.Context, deltaFragments chan DeltaFragment, builder *DeltaBuilder) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case deltaFragment, ok := <-deltaFragments:
			if !ok {
				return nil
			}
			builder.ComputeDelta(DeltaFragments{deltaFragment})
		}
	}
}

// DeltaFragmentGarbageCollector receives Deltas from an input channel and
// removes the associated order Fragments from the DeltaFragmentMatrix. This
// improves performance by removing computations on orders that have already
// been matched.
func DeltaFragmentGarbageCollector(ctx context.Context, deltas chan Delta, computationMatrix *ComputationMatrix) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case delta, ok := <-deltas:
			if !ok {
				return nil
			}
			computationMatrix.RemoveOrder(delta.BuyOrderID)
			computationMatrix.RemoveOrder(delta.SellOrderID)
		}
	}
}

// DeltaBroadcaster reads Deltas from the DeltaBuilder and writes them to an
// output channel. Cancelling the context will shutdown the DeltaBroadcaster.
// Errors are written to an error channel
func DeltaBroadcaster(ctx context.Context, builder *DeltaBuilder, bufferLimit int) (chan Delta, chan error) {
	deltas := make(chan Delta, bufferLimit)
	errors := make(chan error, bufferLimit)

	go func() {
		defer close(deltas)
		defer close(errors)

		buffer := make(Deltas, bufferLimit)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			case <-ticker.C:

				for i, n := 0, builder.WaitForDeltas(buffer[:]); i < n; i++ {
					select {
					case <-ctx.Done():
						errors <- ctx.Err()
						return
					case deltas <- buffer[i]:
					}
				}
			}
		}
	}()

	return deltas, errors
}
