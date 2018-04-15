package main

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/republicprotocol/republic-go/darknode"
)

func main() {
}

func NewLocalDarkNodes(numberOfDarkNodes, numberOfBootstrapDarkNodes int) (darknode.DarkNodes, []context.Context, []context.CancelFunc, error) {
	darkNodes := make(darknode.DarkNodes, numberOfDarkNodes)
	ctxs := make([]context.Context, numberOfDarkNodes)
	cancels := make([]context.CancelFunc, numberOfDarkNodes)

	configs := make([]darknode.Config, numberOfDarkNodes)
	for i := 0; i < numberOfDarkNodes; i++ {
		key := keystore.NewKeyForDirectICAP(rand.Reader)
		configs[i] = darknode.NewLocalConfig(key, "127.0.0.1", fmt.Sprintf("%d", 3000+i))
	}
	for i := 0; i < numberOfDarkNodes; i++ {
		for j := 0; j < numberOfDarkNodes; j++ {
			configs[i].NetworkOption.BootstrapMultiAddresses = append(configs[i].NetworkOption.BootstrapMultiAddresses, configs[j].NetworkOption.MultiAddress)
		}
	}
	var err error
	for i := 0; i < numberOfDarkNodes; i++ {
		darkNodes[i], err = darknode.NewDarkNode(configs[i])
		if err != nil {
			return nil, nil, nil, err
		}
		ctxs[i], cancels[i] = context.WithCancel(context.Background())
	}
	return darkNodes, ctxs, cancels, nil
}
