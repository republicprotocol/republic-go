package hyper

import (
	"context"
	"log"
)

type Prepare struct {
	Signature
	Block
	Rank
	Height uint64
}

func ProcessPreparation(ctx context.Context, prepareChIn <-chan Prepare, validator Validator) (<-chan Commit, <-chan Fault, <-chan error) {
	threshold := validator.Threshold()
	commitCh := make(chan Commit, threshold)
	faultCh := make(chan Fault, threshold)
	errCh := make(chan error, threshold)
	prepares := map[[32]byte]uint8{}
	commited := map[[32]byte]bool{}
	counter := uint64(0)

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
				// log.Println("Successfully reading prepares from channels")
				if !ok {
					return
				}
				counter++
				h := BlockHash(prepare.Block)
				if commited[h] {
					continue
				}

				if !validator.ValidatePrepare(prepare) {
					faultCh <- Fault{
						Rank:      prepare.Rank,
						Height:    prepare.Height,
						Signature: validator.Sign(),
					}
					continue
				}
				// log.Println("Counting prepares on", validator.Sign(), prepares[h], "with threshold", threshold)
				if prepares[h] >= threshold-1 {
					commitCh <- Commit{
						Rank:               prepare.Rank,
						Height:             prepare.Height,
						Block:              prepare.Block,
						ThresholdSignature: ThresholdSignature("Threshold_BLS"),
						Signature:          validator.Sign(),
					}
					log.Println("Successfully writing commits to internal channels")
					commited[h] = true
				} else {
					prepares[h]++
				}
			}
		}
	}()

	return commitCh, faultCh, errCh
}
