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

func ProcessPreparation(ctx context.Context, prepareChIn <-chan Prepare, signer Signer, validator Validator, threshold uint8) (<-chan Commit, <-chan Fault, <-chan error) {
	commitCh := make(chan Commit)
	faultCh := make(chan Fault)
	errCh := make(chan error)
	prepares := map[[32]byte]uint8{}
	commited := map[[32]byte]bool{}
	counter := uint64(0)
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
				counter++
				h := PrepareHash(prepare)
				if commited[h] {
					continue
				}
				if prepares[h] >= threshold-1 {
					commitCh <- Commit{
						Rank:               prepare.Rank,
						Height:             prepare.Height,
						Block:              prepare.Block,
						ThresholdSignature: ThresholdSignature("Threshold_BLS"),
						Signature:          signer.Sign(),
					}
					commited[h] = true
				} else {
					if validator.Prepare(prepare) {
						prepares[h]++
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
