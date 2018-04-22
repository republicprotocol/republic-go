package darknode

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/rpc"
)

type Config struct {
	EcdsaKey keystore.Key      `json:"ecdsaKey"`
	RsaKey   crypto.RsaKeyPair `json:"rsaKey"`
	Host     string            `json:"host"`
	Port     string            `json:"port"`
	Ethereum EthereumConfig    `json:"ethereum"`
	Network  rpc.Options       `json:"network"`
	Logs     logger.Options    `json:"logs"`
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

func NewLocalConfig(ecdsaKey keystore.Key, host, port string) (Config, error) {
	keyPair, err := identity.NewKeyPairFromPrivateKey(ecdsaKey.PrivateKey)
	if err != nil {
		return Config{}, err
	}

	rsaKey, err := crypto.NewRsaKeyPair()
	if err != nil {
		return Config{}, err
	}

	multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%v/tcp/%v/republic/%v", host, port, keyPair.Address()))
	if err != nil {
		return Config{}, err
	}
	return Config{
		EcdsaKey: ecdsaKey,
		RsaKey:   rsaKey,
		Host:     host,
		Port:     port,
		Network: rpc.Options{
			Alpha:                5,
			MultiAddress:         multi,
			MaxBucketLength:      100,
			ClientPoolCacheLimit: 100,
			Timeout:              10 * time.Second,
			TimeoutBackoff:       0,
			TimeoutRetries:       1,
		},
		Ethereum: EthereumConfig{
			URI:                     "http://localhost:8545",
			Network:                 client.NetworkGanache,
			RepublicTokenAddress:    client.RepublicTokenAddressOnGanache.String(),
			DarkNodeRegistryAddress: client.DarkNodeRegistryAddressOnGanache.String(),
		},
	}, nil
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
