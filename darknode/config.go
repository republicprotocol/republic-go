package darknode

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/rpc"
)

type Config struct {
	Key           *keystore.Key  `json:"key"`
	Host          string         `json:"host"`
	Port          string         `json:"port"`
	Ethereum      EthereumConfig `json:"ethereum"`
	NetworkOption rpc.Options    `json:"network"`
	LoggerOptions logger.Options `json:"logger"`
}

type EthereumConfig struct {
	URI                     string         `json:"uri"`
	Network                 client.Network `json:"network"` // One of "ganache", "ropsten", or "mainnet" ("mainnet" is not current supported)
	RepublicTokenAddress    string         `json:"republicTokenAddress"`
	DarkNodeRegistryAddress string         `json:"darkNodeRegistryAddress"`
}

// LoadConfig loads a Config object from the given filename. Returns the Config
// object, or an error.
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	config := new(Config)
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func NewLocalConfig(key *keystore.Key, host, port string) Config {
	address, _, _ := identity.NewAddress()
	multi, _ := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", host, port, address.String()))

	return Config{
		Key:  key,
		Host: host,
		Port: port,
		Ethereum: EthereumConfig{
			URI:                     "http://localhost:8545",
			Network:                 client.NetworkGanache,
			RepublicTokenAddress:    client.RepublicTokenAddressOnGanache.String(),
			DarkNodeRegistryAddress: client.DarkNodeRegistryAddressOnGanache.String(),
		},
		NetworkOption: rpc.Options{
			MultiAddress:      multi,
			Timeout:           3 * time.Second,
			TimeoutBackoff:    3 * time.Second,
			TimeoutRetries:    3,
			MessageQueueLimit: 100,
		},
		LoggerOptions: logger.Options{
			Plugins: []logger.PluginOptions{
				{File: &logger.FilePluginOptions{Path: "stdout"}},
			},
		},
	}
}

func NewFalconConfig() Config {
	return Config{}
}

var FalconBootstrapMultis = []string{
	"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
	"/ip4/52.21.44.236/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
	"/ip4/52.41.118.171/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
	"/ip4/52.59.176.141/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
	"/ip4/52.77.88.84/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
}
