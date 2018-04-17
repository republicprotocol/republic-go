package hyper

import (
	"context"
)

type Fault struct {
	Rank
	Height

	// Signatures of the Replicas that signed this Fault
	Signatures []Signature
}

func ProcessFault(ctx context.Context, faultChIn chan Fault, validator Validator) (chan Fault, chan error) {
	faultCh := make(chan Fault, validator.Threshold())
	errCh := make(chan error, validator.Threshold())
	faults := map[int]int{}
	certified := map[int]bool{}

	go func() {
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case fault := <-faultChIn:
				if certified[fault.Height] {
					continue
				}
				if faults[fault.Height] >= validator.Threshold()-1 {
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
					faults[fault.Height]++
				}
			}
		}
	}()

	return faultCh, errCh
}
