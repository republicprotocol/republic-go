package hyper

import "context"

type Fault struct {
	Rank
	Height
	Signature
}

func ProcessFault(ctx context.Context, faultChIn chan Fault, threshold uint8) (chan Fault, chan error) {
	faultCh := make(chan Fault)
	errCh := make(chan error)
	faults := map[[32]byte]uint8{}

	go func() {
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case fault := <-faultChIn:
				h := FaultHash(fault)
				if faults[h] >= threshold-1 {
					// updateRank(fault)
				}
				faults[h]++
			}
		}
	}()

	return faultCh, errCh
}
