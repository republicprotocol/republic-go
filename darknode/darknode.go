package darknode

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
)

type DarkNodes []DarkNode

type DarkNode struct {
	config           Config
	darkNodeRegistry contracts.DarkNodeRegistry
	darkOcean        *darkocean.Ocean
}

func NewDarkNode(config Config) (DarkNode, error) {

	// Connect to Ethereum
	transactOpts := bind.NewKeyedTransactor(config.Key.PrivateKey)
	//log.Println("config is : ", config.Ethereum.URI, config.Ethereum.Network, config.Ethereum.RepublicTokenAddress, config.Ethereum.DarkNodeRegistryAddress)
	client, err := client.Connect(
		config.Ethereum.URI,
		client.Network(config.Ethereum.Network),
		config.Ethereum.RepublicTokenAddress,
		config.Ethereum.DarkNodeRegistryAddress,
	)
	if err != nil {
		return DarkNode{}, err
	}

	darkNodeRegistry, err := contracts.NewDarkNodeRegistry(context.Background(), &client, transactOpts, &bind.CallOpts{})
	if err != nil {
		return DarkNode{}, err
	}
	darkOcean, err := darkocean.NewOcean(darkNodeRegistry)
	if err != nil {
		return DarkNode{}, err
	}
	return DarkNode{
		config:           config,
		darkNodeRegistry: darkNodeRegistry,
		darkOcean:        darkOcean,
	}, nil
}

func (node *DarkNode) Ocean() *darkocean.Ocean {
	return node.darkOcean
}

func (node *DarkNode) Run(ctx context.Context) {
	node.UpdateDarkOcean(ctx)
}

func (node *DarkNode) UpdateDarkOcean(ctx context.Context) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		minimumEpochIntervalBig, err := node.darkNodeRegistry.MinimumEpochInterval()
		if err != nil {
			errCh <- err
			return
		}
		minimumEpochInterval, err := minimumEpochIntervalBig.ToUint()
		if err != nil {
			errCh <- err
			return
		}

		t := time.NewTicker((time.Duration(minimumEpochInterval) / 24) * time.Second)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-t.C:
				if err := node.darkOcean.Update(); err != nil {
					errCh <- fmt.Errorf("cannot update dark ocean: %v", err)
				}
			}
		}
	}()

	return errCh
}
