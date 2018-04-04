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
	"github.com/republicprotocol/republic-go/rpc"
)

const (
	DefaultDebugLevel           = 3
	DefaultAlpha                = 3
	DefaultMaxBucketLength      = 20
	DefaultClientPoolCacheLimit = 20
	DefaultTimeout              = 30
	DefaultTimeBackoff          = 30
	DefaultTimeoutRetries       = 3
	DefaultConcurrent           = false
	DefaultMaxConnections       = 3
	DefaultMessageQueueLimit    = 100
)

var bootstrapNodes = []string{
	"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
	"/ip4/52.21.44.236/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
	"/ip4/52.41.118.171/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
	"/ip4/52.59.176.141/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
	"/ip4/52.77.88.84/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
}

var localBootstrapNodes = []string{
	"/ip4/0.0.0.0/tcp/3000/republic/8MJxpBsezEGKPZBbhFE26HwDFxMtFu",
	"/ip4/0.0.0.0/tcp/3001/republic/8MGB2cj2HbQFepRVs43Ghct5yCRS9C",
	"/ip4/0.0.0.0/tcp/3002/republic/8MGVBvrQJji8ecEf3zmb8SXFCx1PaR",
	"/ip4/0.0.0.0/tcp/3003/republic/8MJNCQhMrUCHuAk977igrdJk3tSzkT",
	"/ip4/0.0.0.0/tcp/3004/republic/8MK6bq5m7UfE1mzRNunJTFH6zTbyss",
}

func main() {
	err := generateLocalBootstrapNodeConfigs()
	if err != nil {
		log.Fatal(err)
	}

	err = generateLocalDarkNodeConfigs(15)
	if err != nil {
		log.Fatal(err)
	}
}

// Generate configs of bootstrap nodes for local testing
func generateLocalBootstrapNodeConfigs() error {
	for i := range localBootstrapNodes {
		// Generate ECDSA key pair and ethereum key
		_, keyPair, err := identity.NewAddress()
		if err != nil {
			return err
		}
		ethKey := keystore.NewKeyForDirectICAP(rand.Reader)

		// Get multiAddress
		multiAddress, err := identity.NewMultiAddressFromString(localBootstrapNodes[i])
		if err != nil {
			return err
		}
		swarmOptions := generateSwarmOptions(localBootstrapNodes)

		// Remove self multi address from the bootstrap nodes
		for j := 0; j < len(swarmOptions.BootstrapMultiAddresses); j++ {
			if swarmOptions.BootstrapMultiAddresses[j].String() == multiAddress.String() {
				numberOfBootstrapNodes := len(localBootstrapNodes)
				swarmOptions.BootstrapMultiAddresses[j] = swarmOptions.BootstrapMultiAddresses[numberOfBootstrapNodes-1]
				swarmOptions.BootstrapMultiAddresses = swarmOptions.BootstrapMultiAddresses[:numberOfBootstrapNodes-1]
				j--
			}
		}
		syncerOptions := rpc.SyncerOptions{
			MaxConnections: DefaultMaxConnections,
		}
		config := generateConfig(multiAddress, keyPair, *ethKey, 3000+i, swarmOptions, syncerOptions)
		writeConfigToFile(fmt.Sprintf("./bootstrap-node-%d.json", i), config)
	}
	return nil
}

func generateLocalDarkNodeConfigs(numberOfNodes int) error {
	for i := 0; i < numberOfNodes; i++ {
		// Generate ECDSA key pair and ethereum key pair
		address, keyPair, err := identity.NewAddress()
		if err != nil {
			return err
		}
		ethKey := keystore.NewKeyForDirectICAP(rand.Reader)

		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 4000+i, address))
		if err != nil {
			return err
		}

		swarmOptions := generateSwarmOptions(localBootstrapNodes)
		syncerOptions := rpc.SyncerOptions{
			MaxConnections: DefaultMaxConnections,
		}
		config := generateConfig(multiAddress, keyPair, *ethKey, 4000+i, swarmOptions, syncerOptions)
		writeConfigToFile(fmt.Sprintf("./node-%d.json", i), config)
	}

	return nil
}

func generateSwarmOptions(nodes []string) rpc.SwarmOptions {
	multis := make([]identity.MultiAddress, len(nodes))
	swarmOptions := rpc.NewSwarmOptions(multis, DefaultDebugLevel, DefaultAlpha, DefaultMaxBucketLength, DefaultClientPoolCacheLimit, DefaultConcurrent)
	for i := range nodes {
		multi, err := identity.NewMultiAddressFromString(localBootstrapNodes[i])
		if err != nil {
			log.Fatal(err)
		}
		swarmOptions.BootstrapMultiAddresses[i] = multi
	}

	return swarmOptions
}

func generateConfig(multiAddress identity.MultiAddress, keyPair identity.KeyPair, ethKey keystore.Key, port int, swarmOptions rpc.SwarmOptions, syncerOptions rpc.SyncerOptions) *node.Config {
	networkOptions := rpc.Options{
		SwarmOptions:      swarmOptions,
		SyncerOptions:     syncerOptions,
		MultiAddress:      multiAddress,
		Timeout:           DefaultTimeout * time.Second,
		TimeoutBackoff:    DefaultTimeBackoff * time.Second,
		TimeoutRetries:    DefaultTimeoutRetries,
		MessageQueueLimit: DefaultMessageQueueLimit,
	}
	config := &node.Config{
		NetworkOptions: networkOptions,
		LoggerOptions: logger.Options{
			Plugins: []logger.PluginOptions{
				{
					File: &logger.FilePluginOptions{Path: "stdout"},
				},
			},
		},
		Host:        "0.0.0.0",
		Port:        fmt.Sprintf("%d", port),
		KeyPair:     keyPair,
		EthereumKey: ethKey,
		EthereumRPC: "https://ropsten.infura.io",
	}

	return config
}

func writeConfigToFile(filePath string, config *node.Config) {
	data, err := json.Marshal(config)
	if err != nil {
		log.Fatal("failt to marshal config", err)
	}
	d1 := []byte(data)
	err = ioutil.WriteFile(filePath, d1, 0644)
	if err != nil {
		log.Fatal("fail to write config to file ", err)
	}
}
