package smpc_test

import (
	"context"
	"math/rand"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/swarm"
)

type mockSwarmer struct {
}

func (swarmer mockSwarmer) Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses) error {
	return nil
}

func (swarmer mockSwarmer) Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error) {
	n := rand.Intn(1000)
	if n <= 200 {
		return identity.MultiAddress{}, swarm.ErrMultiAddressNotFound
	}
	return query.MultiAddress()
}
