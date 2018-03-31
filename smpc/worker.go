package smpc

import (
	"context"
	"time"

	"github.com/republicprotocol/republic-go/order"
)

func OrderFragmentReceiver(ctx context.Context, orderFragments chan order.Fragment, matrix *DeltaFragmentMatrix) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case orderFragment := <-orderFragments:
			if orderFragment.OrderParity == order.ParityBuy {
				matrix.ComputeBuyOrder(&orderFragment)
			} else {
				matrix.ComputeSellOrder(&orderFragment)
			}
		}
	}
}

func DeltaFragmentReceiver(ctx context.Context, deltaFragments chan DeltaFragment, builder *DeltaBuilder) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case deltaFragment := <-deltaFragments:
			builder.ComputeDelta(DeltaFragments{deltaFragment})
		}
	}
}

func DeltaFragmentWaiter(ctx context.Context, matrix *DeltaFragmentMatrix, builder *DeltaBuilder) (chan DeltaFragment, chan error) {
	deltaFragments := make(chan DeltaFragment, 1)
	errors := make(chan error, 1)

	go func() {
		defer close(deltaFragments)
		defer close(errors)

		buffer := [128]DeltaFragment{}
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			case <-ticker.C:
				if n := matrix.WaitForDeltaFragments(buffer[:]); n > 0 {
					builder.ComputeDelta(buffer[:n])
					for i := 0; i < n; i++ {
						deltaFragments <- buffer[i]
					}
				}
			}
		}
	}()

	return deltaFragments, errors
}

func DeltaBroadcaster(ctx context.Context, builder *DeltaBuilder) (chan Delta, chan error) {
	deltas := make(chan Delta, 1)
	errors := make(chan error, 1)

	go func() {
		defer close(deltas)
		defer close(errors)

		buffer := [128]Delta{}
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			case <-ticker.C:
				n := builder.WaitForDeltas(buffer[:])
				for i := 0; i < n; i++ {
					deltas <- buffer[i]
				}
			}
		}
	}()

	return deltas, errors
}

func DeltaReconstructer(ctx context.Context, chan DeltaFragment) (chan Delta, chan error) {
	
}