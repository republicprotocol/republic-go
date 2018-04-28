package relay

import (
	"encoding/json"
	"os"

	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/identity"
)

// A Config defines the different settings for a Relay.
type Config struct {
	Token                   string                  `json:"token"`
	EthereumAddress         string                  `json:"ethereumAddress"`
	KeyPair                 identity.KeyPair        `json:"keypair"`
	MultiAddress            identity.MultiAddress   `json:"multiAddress"`
	BootstrapMultiAddresses identity.MultiAddresses `json:"bootstrapMultiAddresses"`
	Ethereum                ethereum.Config         `json:"ethereum"`
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
