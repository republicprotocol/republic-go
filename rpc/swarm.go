package rpc

import (
	"context"

	"github.com/republicprotocol/republic-go/identity"
)

type SwarmClient interface {
	Ping(ctx context.Context, to identity.MultiAddress) (identity.MultiAddress, error)
	Query(ctx context.Context, to identity.MultiAddress, query identity.Address) (identity.MultiAddresses, error)
}
