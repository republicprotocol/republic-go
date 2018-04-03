package hyper

import "context"

type Fault struct {
	rank   Rank
	height Height
}

func ProcessFault(ctx context.Context, faultChIn chan Commit) (chan Fault, chan error) {
	faultCh := make(chan Fault)
	errCh := make(chan error)

	go func() {
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case _ = <-faultChIn:
				// TODO: process the commit
			}
		}
	}()

	return faultCh, errCh
}
