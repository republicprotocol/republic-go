package hyper

import (
	"context"
	"log"
)

type ThresholdSignature Signature
type Commit struct {
	Rank
	Height uint64
	Block
	ThresholdSignature
	Signature
}

func ProcessCommit(ctx context.Context, commitChIn <-chan Commit, validator Validator) (chan Commit, chan Block, chan Fault, chan error) {
	blockCh := make(chan Block, validator.Threshold())
	commitCh := make(chan Commit, validator.Threshold())
	faultCh := make(chan Fault, validator.Threshold())
	errCh := make(chan error, validator.Threshold())
	blocks := validator.SharedBlocks()
	threshold := validator.Threshold()
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
				log.Println("Counting commits on", validator.Sign(), commits[h], "with threshold", threshold)
				if commits[h] >= threshold-1 {
					certified[h] = true
					blockCh <- commit.Block
					blocks.IncrementHeight()
				} else {
					if validator.ValidateCommit(commit) {
						commits[h]++
						if len(commitCh) == int(validator.Threshold()) {
							continue
						}
						commitCh <- Commit{
							Rank:               commit.Rank,
							Height:             commit.Height,
							Block:              commit.Block,
							ThresholdSignature: commit.ThresholdSignature,
							Signature:          validator.Sign(),
						}
					} else {
						faultCh <- Fault{
							Rank:      commit.Rank,
							Height:    commit.Height,
							Signature: validator.Sign(),
						}
					}
				}
			}
		}
	}()
	return commitCh, blockCh, faultCh, errCh
}
