package rpc

import (
	"context"

	"github.com/republicprotocol/republic-go/identity"
)

type SwarmRPC interface {
	Ping(ctx context.Context, to identity.Address, multiAddr identity.MultiAddress) (identity.MultiAddress, error)
	Query(ctx context.Context, to identity.Address, query identity.Address) ([]identity.MultiAddress, error)
}

type SwarmClient struct {
	swarm SwarmRPC
}

func (client *SwarmClient) Bootstrap() {
}

func (client *SwarmClient) Query(ctx context.Context, addr identity.Address, depth int) (identity.MultiAddress, error) {
	panic("unimplemented")
}
