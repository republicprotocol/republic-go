package hyper

import (
	"context"
)

type Prepare struct {
	Signature
	Block
	Rank
	Height
}

type Prepares map[BlockHash]uint8

func ProcessPreparation(ctx context.Context, prepareChIn <-chan Prepare, signer Signer, validator Validator, threshold uint8) (<-chan Commit, <-chan Fault, <-chan error) {
	commitCh := make(chan Commit)
	faultCh := make(chan Fault)
	errCh := make(chan error)
	prepares := make(Prepares)

	// prepareChIn = FilterPrepare(prepareChIn, func(*Prepare) bool { return false })

	go func() {
		defer close(commitCh)
		defer close(faultCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case prepare, ok := <-prepareChIn:
				if !ok {
					return
				}
				b := getBlockHash(prepare.Block)
				if prepares[b] >= threshold-1 {
					commitCh <- Commit{
						Rank:               prepare.Rank,
						Height:             prepare.Height,
						Block:              prepare.Block,
						ThresholdSignature: ThresholdSignature("Threshold_BLS"),
						Signature:          signer.Sign(),
					}
					prepares[b] = 0
				} else {
					if validator.Prepare(prepare) {
						prepares[b]++
					} else {
						faultCh <- Fault{
							Rank:      prepare.Rank,
							Height:    prepare.Height,
							Signature: signer.Sign(),
						}
					}
				}
			}
		}
	}()

	return commitCh, faultCh, errCh
}
