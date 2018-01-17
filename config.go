package miner

import (
	"encoding/json"
	"os"

	"github.com/republicprotocol/go-identity"
)

// Config information for booting a Node.
type Config struct {
	KeyPair         identity.KeyPair        `json:"key"`
	Multi           identity.MultiAddress   `json:"multi"`
	BootstrapMultis identity.MultiAddresses `json:"bootstrap_multis"`
}

// LoadConfig loads a Config object from the given filename. Returns the Config
// object, or an error.
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
