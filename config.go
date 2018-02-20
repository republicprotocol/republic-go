package node

import (
	"encoding/json"
	"os"

	"math/big"

	"github.com/republicprotocol/go-identity"
)

// Config information for Miners
type Config struct {
	Host                    string                  `json:"host"`
	Port                    string                  `json:"port"`
	EthereumPrivateKey      string                  `json:"ethereum_private_key"`
	RepublicKeyPair         identity.KeyPair        `json:"republic_key_pair"`
	RSAKeyPair              string                  `json:"rsa_key_pair"`
	MultiAddress            identity.MultiAddress   `json:"multi_address"`
	BootstrapMultiAddresses identity.MultiAddresses `json:"bootstrap_multi_addresses"`

	ComputationChunkSize     int      `json:"computation_chunk_size"`
	ComputationChunkInterval int      `json:"computation_chunk_interval"`
	Prime                    *big.Int `json:"prime"`
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
