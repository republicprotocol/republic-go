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

func ProcessCommit(ctx context.Context, commitChIn <-chan Commit, signer Signer, validator Validator, blocks *SharedBlocks, threshold uint8) (chan Commit, chan Block, chan Fault, chan error) {
	blockCh := make(chan Block)
	commitCh := make(chan Commit)
	faultCh := make(chan Fault)
	errCh := make(chan error)
	commits := map[[32]byte]uint8{}
	certified := map[[32]byte]bool{}
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
				h := CommitHash(commit)
				if certified[h] {
					continue
				}
				if commits[h] >= threshold-1 {
					certified[h] = true
					blockCh <- commit.Block
					blocks.IncrementHeight()
				} else {
					if validator.Commit(commit) {
						commits[h]++
						commitCh <- Commit{
							Rank:               commit.Rank,
							Height:             commit.Height,
							Block:              commit.Block,
							ThresholdSignature: commit.ThresholdSignature,
							Signature:          signer.Sign(),
						}
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
	return commitCh, blockCh, faultCh, errCh
}
