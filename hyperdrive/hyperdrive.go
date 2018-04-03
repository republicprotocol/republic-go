package hyper

import (
	"context"
)

type Block struct {
}

type Prepare struct {
}

type Commit struct {
}

type Fault struct {
}

func ProcessPreparation(ctx context.Context, prepareChIn chan Prepare) (chan Prepare, chan Commit, chan Fault, chan error) {
	prepareCh := make(chan Prepare)
	commitCh := make(chan Commit)
	faultCh := make(chan Fault)
	errCh := make(chan error)

	go func() {
		defer close(prepareCh)
		defer close(commitCh)
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case _ = <-prepareChIn:
				// TODO: process the prepare
			}
		}
	}()

	return prepareCh, commitCh, faultCh, errCh
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
