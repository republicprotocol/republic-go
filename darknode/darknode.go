package darknode

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/ethereum/client"
)

type Darknodes []Darknode

type Darknode struct {
	config           Config
	darkNodeRegistry dnr.DarkNodeRegistry
}

func NewDarknode(config Config) Darknode {

	// Connect to Ethereum
	transactOpts := bind.NewKeyedTransactor(config.Key.PrivateKey)
	client := client.Connect(
		config.EthereumConfig.URI,
		config.EthereumConfig.Network,
		config.EthereumConfig.RepublicTokenAddress,
		config.EthereumConfig.DarkNodeRegistryAddress,
	)

	return Darknode{
		config:           config,
		darkNodeRegistry: dnr.NewDarkNodeRegistry(context.Background(), &client, transactOpts, &bind.CallOpts{}),
	}
}

func (node *Darknode) Run(ctx context.Context) {
}
