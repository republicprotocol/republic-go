package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
)

var bootstrapNode = []string{
	"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
	"/ip4/52.21.44.236/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
	"/ip4/52.41.118.171/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
	"/ip4/52.59.176.141/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
	"/ip4/52.77.88.84/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
}

func main() {
	err := generateSingleNode(".")
	if err != nil {
		log.Fatal(err)
	}
}

func generateSingleNode(dir string) error {
	_, keyPair, err := identity.NewAddress()
	if err != nil {
		return err
	}

	// Create default network options
	options := network.Options{
		BootstrapMultiAddresses: make([]identity.MultiAddress, len(bootstrapNode)),
		Debug:                network.DebugHigh,
		Alpha:                3,
		MaxBucketLength:      20,
		ClientPoolCacheLimit: 20,
		Timeout:              30 * time.Second,
		TimeoutBackoff:       30 * time.Second,
		TimeoutRetries:       3,
		Concurrent:           false,
	}
	ethKey := keystore.NewKeyForDirectICAP(rand.Reader)
	for i := range bootstrapNode {
		multi, err := identity.NewMultiAddressFromString(bootstrapNode[i])
		if err != nil {
			return err
		}
		options.BootstrapMultiAddresses[i] = multi
	}

	config := &node.Config{
		NetworkOptions: options,
		LoggerOptions: logger.Options{
			Plugins: []logger.PluginOptions{
				logger.PluginOptions{
					File: &logger.FilePluginOptions{Path: "/home/ubuntu/.darknode/darknode.out"},
				},
				logger.PluginOptions{
					WebSocket: &logger.WebSocketPluginOptions{Host: "0.0.0.0", Port: "18515"},
				},
			},
		},
		Host:        "0.0.0.0",
		Port:        fmt.Sprintf("%d", 18514),
		Path:        "/home/ubuntu/.darknode",
		KeyPair:     keyPair,
		EthereumKey: *ethKey,
		EthereumRPC: "https://ropsten.infura.io",
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	d1 := []byte(data)
	err = ioutil.WriteFile(fmt.Sprintf("%s/config.json", dir), d1, 0644)
	return err
}
