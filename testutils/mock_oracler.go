package testutils

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/oracle"
)

type MockOracleClient struct {
	addr  identity.Address
	store oracle.MidpointPriceStorer
	hub   map[identity.Address]oracle.Server
}

func NewMockOracleClient(addr identity.Address, hub map[identity.Address]oracle.Server) (oracle.Client, oracle.MidpointPriceStorer, error) {
	storer := leveldb.NewMidpointPriceStorer()
	return &MockOracleClient{
		addr:  addr,
		store: storer,
		hub:   hub,
	}, storer, nil
}

func (client *MockOracleClient) UpdateMidpoint(ctx context.Context, to identity.MultiAddress, midpointPrice oracle.MidpointPrice) error {
	server, ok := client.hub[to.Address()]
	if !ok {
		return errors.New("cannot send midPointPrice from client")
	}
	return server.UpdateMidpoint(ctx, midpointPrice)
}

func (client *MockOracleClient) MultiAddress() identity.MultiAddress {
	multi, err := client.addr.MultiAddress()
	if err != nil {
		log.Println("error retrieving multiAddress from store", err)
		return identity.MultiAddress{}
	}
	return multi
}
