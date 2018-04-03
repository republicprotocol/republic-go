package hyper

import (
	"context"
)

type Block struct {
	tuples    Tuples
	signature Signature
}

type BlockHash [32]byte

func ConsumeBlock(ctx context.Context, blockChIn chan Block, sharedBlocks *SharedBlocks) chan error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case block, ok := <-blockChIn:
				if !ok {
					return
				}
				sharedBlocks.AddBlock(block)
			}
		}
	}()

	return errCh
}

func validateBlock(b Block, sb SharedBlocks) bool {
	for _, tuple := range b.tuples {
		if !sb.ValidateTuple(tuple) {
			return false
		}
	}
	// TODO: Replace the validate signature with actual logic
	return validateSignature()
}

func validateSignature() bool {
	return true
}

func getBlockHash(b Block) BlockHash {
	return BlockHash{}
}
