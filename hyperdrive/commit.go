package hyper

import "context"

type Commit struct {
}

func ProcessCommit(ctx context.Context, commitChIn chan Commit) (chan Commit, chan Block, chan Fault, chan error) {
	commitCh := make(chan Commit)
	blockCh := make(chan Block)
	faultCh := make(chan Fault)
	errCh := make(chan error)

	go func() {
		defer close(commitCh)
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case _ = <-commitChIn:
				// TODO: process the commit
			}
		}
	}()

	return commitCh, blockCh, faultCh, errCh
}
