package node

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
)

// Config contains all configuration details for running a DarkNode.
type Config struct {
	NetworkOptions network.Options `json:"network"`
	LoggerOptions  logger.Options  `json:"logger"`

	Host string `json:"host"`
	Port string `json:"port"`

	PrivateKey         *identity.KeyPair `json:"privateKey"`
	EthereumPrivateKey *keystore.Key     `json:"ethereumPrivateKey"`
	EthereumRPC        string            `json:"ethereumRPC"`
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
