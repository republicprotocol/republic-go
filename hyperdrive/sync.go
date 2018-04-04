package hyper

import "context"

type CertifiedCommit struct{}

func ProduceBlocksFromFileSystem(ctx context.Context, sharedBlocks SharedBlocks, location string) (SharedBlocks, error) {
	return sharedBlocks, nil
}

func ProduceBlocksFromSyncingWithNeighbours(ctx context.Context, sharedBlocks SharedBlocks, certifiedCommitChIn chan CertifiedCommit) (SharedBlocks, error) {
	return sharedBlocks, nil
}
