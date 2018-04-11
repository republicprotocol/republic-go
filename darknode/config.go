package darknode

import (
	"crypto/rand"
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// // Config contains all configuration details for running a DarkNode.
// type Config struct {
// 	NetworkOptions rpc.Options    `json:"network"`
// 	LoggerOptions  logger.Options `json:"logger"`

// 	Path string `json:"path"`
// 	Host string `json:"host"`
// 	Port string `json:"port"`

// 	KeyPair     identity.KeyPair `json:"keyPair"`
// 	EthereumKey keystore.Key     `json:"ethereumKey"`
// 	EthereumRPC string           `json:"ethereumRPC"`
// }

type Config struct {
	Key            keystore.Key `json:"key"`
	Host string `json:"host"`
	Port string `json:"port"`
	Ethereum EthereumConfig `json:"ethereum"`
}

type EthereumConfig struct {
	URI                     string `json:"uri"`
	Network                 string `json:"network"` // One of "ganache", "ropsten", or "mainnet" ("mainnet" is not current supported)
	RepublicTokenAddress    string `json:"republicTokenAddress"`
	DarkNodeRegistryAddress string `json:"darkNodeRegistryAddress"`
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
	return Config{
		Key:         key,
		Host:  host,
		Port: port,
		Ethereum: EthereumConfig{
			URI: "http://localhost:8545",
			Network: client.NetworkGanache,
			RepublicTokenAddress: client.RepublicTokenAddressOnGanache,
			DarkNodeRegistryAddress: client.DarkNodeRegistryAddressOnGanache,
		}
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
