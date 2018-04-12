package darknode

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/ethereum/client"
)

type DarkNodes []DarkNode

type DarkNode struct {
	config           Config
	darkNodeRegistry dnr.DarkNodeRegistry
	Ocean            darkocean.Ocean
}

func NewDarknode(config Config) DarkNode {

	// Connect to Ethereum
	transactOpts := bind.NewKeyedTransactor(config.Key.PrivateKey)
	client := client.Connect(
		config.EthereumConfig.URI,
		config.EthereumConfig.Network,
		config.EthereumConfig.RepublicTokenAddress,
		config.EthereumConfig.DarkNodeRegistryAddress,
	)

	return DarkNode{
		config:           config,
		darkNodeRegistry: dnr.NewDarkNodeRegistry(context.Background(), &client, transactOpts, &bind.CallOpts{}),
	}
}



func (node *DarkNode) Run(ctx context.Context) {
}
