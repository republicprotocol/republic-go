package hyper

import (
	"bytes"
	"context"
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

type Block struct {
	Tuples
	Signature
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
			case _, ok := <-blockChIn:
				if !ok {
					return
				}
			}
		}
	}()
	return errCh
}

func getBlockHash(b Block) BlockHash {
	var blockBuffer bytes.Buffer
	binary.Write(&blockBuffer, binary.BigEndian, b.Tuples)
	return sha3.Sum256(blockBuffer.Bytes())
}
