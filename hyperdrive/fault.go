package hyper

import (
	"context"
)

type Fault struct {
	Rank
	Height uint64
	Signature
}

func ProcessFault(ctx context.Context, faultChIn chan Fault, validator Validator) (chan Fault, chan error) {
	faultCh := make(chan Fault, validator.Threshold())
	errCh := make(chan error, validator.Threshold())
	faults := map[[32]byte]uint8{}
	certified := map[[32]byte]bool{}

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
				if certified[h] {
					continue
				}
				if faults[h] >= validator.Threshold()-1 {
					faultCh <- Fault{
						fault.Rank,
						fault.Height,
						Signature("Threshold_BLS"),
					}
				} else {
					faultCh <- Fault{
						fault.Rank,
						fault.Height,
						validator.Sign(),
					}
					faults[h]++
				}
			}
		}
	}()

	return faultCh, errCh
}
