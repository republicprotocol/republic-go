package swarm

import (
	"context"

	"github.com/republicprotocol/republic-go/identity"
)

type Swarmer interface {
	Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses, depth int) <-chan error
	Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error)
}
