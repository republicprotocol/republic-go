package darknode

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
)

type DarkNodes []DarkNode

type DarkNode struct {
	config           Config
	darkNodeRegistry contracts.DarkNodeRegistry
	ocean            darkocean.Ocean
}

func NewDarkNode(config Config) (DarkNode, error) {

	// Connect to Ethereum
	transactOpts := bind.NewKeyedTransactor(config.Key.PrivateKey)
	client, err := client.Connect(
		config.Ethereum.URI,
		config.Ethereum.Network,
		config.Ethereum.RepublicTokenAddress,
		config.Ethereum.DarkNodeRegistryAddress,
	)
	if err != nil {
		return DarkNode{}, err
	}

	registry, err :=  contracts.NewDarkNodeRegistry(context.Background(), &client, transactOpts, &bind.CallOpts{})
	if err != nil {
		return DarkNode{}, err
	}
	return DarkNode{
		config:           config,
		darkNodeRegistry: registry,
	}, nil
}

func (node *DarkNode) Ocean() darkocean.Ocean  {
	return node.ocean
}


func (node *DarkNode) Run(ctx context.Context) {
}
