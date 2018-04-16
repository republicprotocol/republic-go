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

// ProcessCommits by collecting Commits. Once a threshold of Commits has been
// reached for a Block, the Block is certified and produced to the Block
// channel. The incrementing of height must be done by reading Blocks produced
// by this process, and comparing it to the current height.
func ProcessCommits(ctx context.Context, commitChIn <-chan Commit, validator Validator, capacity int) (<-chan Commit, <-chan Block, <-chan Fault, <-chan error) {
	commitCh := make(chan Commit, capacity)
	blockCh := make(chan Block, capacity)
	faultCh := make(chan Fault, capacity)
	errCh := make(chan error, capacity)

	go func() {
		defer close(commitCh)
		defer close(blockCh)
		defer close(faultCh)
		defer close(errCh)

		threshold := validator.Threshold()

		commitsReceived := map[[32]byte]int{}
		commits := map[[32]byte]Commit{}
		certifications := map[[32]byte]struct{}{}

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
				if _, ok := certifications[h]; ok {
					continue
				}

				if !validator.ValidateCommit(commit) {
					faultCh <- Fault{
						Rank:      commit.Block.Rank,
						Height:    commit.Height,
						Signature: validator.Sign(),
					}
					continue
				}

				commitsReceived[h]++
				if commitsReceived[h] >= threshold {
					certifications[h] = struct{}{}
					blockCh <- commit.Block
				}

				if _, ok := commits[h]; !ok {
					commits[h] = Commit{
						Rank:               commit.Rank,
						Height:             commit.Height,
						Block:              commit.Block,
						ThresholdSignature: commit.ThresholdSignature,
						Signature:          validator.Sign(),
					}
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case commitCh <- commits[h]:
					}
				}
			}
		}
	}()

	return commitCh, blockCh, faultCh, errCh
}
