package hyper

import (
	"context"
)

type ThresholdSignature Signature
type Commit struct {
	Rank
	Height
	Block
	ThresholdSignature
	Signature
}

type Commits map[BlockHash]uint8

func ProcessCommit(ctx context.Context, commitChIn <-chan Commit, signer Signer, validator Validator, threshold uint8) (chan Block, chan Fault, chan error) {
	blockCh := make(chan Block)
	faultCh := make(chan Fault)
	errCh := make(chan error)
	commits := make(Commits)

	go func() {
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case commit, ok := <-commitChIn:
				if !ok {
					return
				}
				b := getBlockHash(commit.Block)
				if commits[b] >= threshold-1 {
					blockCh <- commit.Block
					validator.UpdateHeight()
					commits[b] = 0
				} else {
					if validator.Commit(commit) {
						commits[b]++
					} else {
						faultCh <- Fault{
							Rank:      commit.Rank,
							Height:    commit.Height,
							Signature: signer.Sign(),
						}
					}
				}
			}
		}
	}()

	return blockCh, faultCh, errCh
}
