package node

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
)

// Config information for Miners
type Config struct {
	network.Options `json:"network"`

	Host string `json:"host"`
	Port string `json:"port"`

	EthereumPrivateKey string                `json:"ethereum_private_key"`
	RepublicKeyPair    identity.KeyPair      `json:"republic_key_pair"`
	RSAKeyPair         identity.KeyPair      `json:"rsa_key_pair"`

	Logger                  *logger.Logger          `json:"logger"`
	Dev                     bool                    `json:"dev"`

	EthereumRPC string `json:"ethereum_rpc"`

	Prime *big.Int `json:"prime"`
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

// EthereumKeyPair returns the Ethereum private
func (config Config) EthereumKeyPair() (identity.KeyPair, error) {
	key, err := hex.DecodeString(config.EthereumPrivateKey)
	if err != nil {
		return identity.KeyPair{}, err
	}
	ecdsa, err := crypto.ToECDSA(key)
	if err != nil {
		return identity.KeyPair{}, err
	}
	return identity.NewKeyPairFromPrivateKey(ecdsa)
}
