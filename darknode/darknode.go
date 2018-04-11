package darknode

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/contracts/dnr"
)

type Darknodes []Darknode

type Darknode struct {
	config           Config
	darkNodeRegistry dnr.DarkNodeRegistry
}

func NewDarknode(config Config) Darknode {

	transactOpts := bind.NewKeyedTransactor(config.Key.PrivateKey)
	client := nil // TODO: Figure this out
	darkNodeRegistry := dnr.NewDarkNodeRegistry(context.Background(), &client, transactOpts, &bind.CallOpts{})

	return Darknode{
		config:           config,
		darkNodeRegistry: darkNodeRegistry,
	}
}

func (node *Darknode) Run(ctx context.Context) {
}
