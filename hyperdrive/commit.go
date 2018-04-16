package hyper

import (
	"context"
)

type ThresholdSignature Signature

type Commit struct {
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
	commits := map[[32]byte]int{}
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
				h := BlockHash(commit.Block)
				if certified[h] {
					continue
				}

				if !validator.ValidateCommit(commit) {
					faultCh <- Fault{
						Rank:      commit.Rank,
						Height:    commit.Height,
						Signature: validator.Sign(),
					}
					continue
				}

				// log.Println("Counting commits on", validator.Sign(), commits[h], "with threshold", threshold)
				if commits[h] >= threshold-1 {
					certified[h] = true
					blockCh <- commit.Block
					blocks.NextHeight()
				} else {
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
				}
			}
		}
	}()
	return commitCh, blockCh, faultCh, errCh
}
