package hyper

import "context"

const THRESHOLD uint8 = 72

type Prepare struct {
	Signature
	Block
	Rank
	Height
}

type Prepares map[BlockHash]uint8

func ProcessPreparation(ctx context.Context, prepareChIn chan Prepare, signer Signer, sharedBlocks SharedBlocks) (chan Prepare, chan Commit, chan Fault, chan error) {
	prepareCh := make(chan Prepare)
	commitCh := make(chan Commit)
	faultCh := make(chan Fault)
	errCh := make(chan error)
	prepares := make(Prepares)

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
			case prepare := <-prepareChIn:
				b := getBlockHash(prepare.Block)
				if prepares[b] >= THRESHOLD {
					commitCh <- Commit{}
				} else {
					if validatePrepare(prepare, sharedBlocks) {
						prepareCh <- Prepare{}
						prepares[b]++
					} else {
						faultCh <- Fault{}
					}
				}
			}
		}
	}()

	return prepareCh, commitCh, faultCh, errCh
}

func validatePrepare(prepare Prepare, sharedBlocks SharedBlocks) bool {
	valid := validateBlock(prepare.Block, sharedBlocks)
	return valid
}

func signPrepare(prepare Prepare, signer Signer) (Prepare, err) {
	
}

func signCommit(prepare Prepare, signer Signer) (Commit, err) {

}
