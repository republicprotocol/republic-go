package swarm

import (
	"context"

	"github.com/republicprotocol/republic-go/identity"
)

type SwarmClient interface {
	Ping(ctx context.Context, to identity.MultiAddress) (identity.MultiAddress, error)
	Query(ctx context.Context, to identity.MultiAddress, query identity.Address) (identity.MultiAddresses, error)
}

type Swarmer interface {
	Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses, depth int) <-chan error
	Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error)
}

type swarmer struct {
	client SwarmClient
}

func NewSwarmer(client SwarmClient) Swarmer {
	return &swarmer{}
}

func (swarmer *swarmer) Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses, depth int) <-chan error {
	panic("unimplemented")
}

func (swarmer *swarmer) Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error) {
	panic("unimplemented")
}
