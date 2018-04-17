package darknode

import (
	"context"

	"github.com/republicprotocol/republic-go/identity"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
)

type Darknodes []Darknode

type Darknode struct {
	Config

	darknodeRegistry contracts.DarkNodeRegistry
	ocean            Ocean
}

// NewDarknode returns a new Darknode.
func NewDarknode(config Config) (Darknode, error) {
	node := Darknode{Config: config}

	// Open a connection to the Ethereum network
	transactOpts := bind.NewKeyedTransactor(config.Key.PrivateKey)
	client, err := client.Connect(
		config.Ethereum.URI,
		client.Network(config.Ethereum.Network),
		config.Ethereum.RepublicTokenAddress,
		config.Ethereum.DarkNodeRegistryAddress,
	)
	if err != nil {
		return Darknode{}, err
	}

	// Create bindings to the Ethereum smart contracts
	darknodeRegistry, err := contracts.NewDarkNodeRegistry(context.Background(), client, transactOpts, &bind.CallOpts{})
	if err != nil {
		return Darknode{}, err
	}
	node.darknodeRegistry = darknodeRegistry

	node.ocean = NewOcean(darknodeRegistry)

	return node, nil
}

// ID returns the ID of the Darknode.
func (node *Darknode) ID() identity.ID {
	key, err := identity.NewKeyPairFromPrivateKey(node.Config.Key.PrivateKey)
	if err != nil {
		panic(err)
	}
	return key.ID()
}

// Ocean returns the Ocean used by this Darknode for computing the Pools and
// its position in them.
func (node *Darknode) Ocean() Ocean {
	return node.ocean
}
