package node

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
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

	EthereumKey     *keystore.Key    `json:"ethereum_key"`
	RepublicKeyPair identity.KeyPair `json:"republic_key_pair"`
	RSAKeyPair      identity.KeyPair `json:"rsa_key_pair"`

	Dev bool `json:"dev"`

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
	key, err := hex.DecodeString(config.EthereumKey.Address.String())
	if err != nil {
		return identity.KeyPair{}, err
	}
	ecdsa, err := crypto.ToECDSA(key)
	if err != nil {
		return identity.KeyPair{}, err
	}
	return identity.NewKeyPairFromPrivateKey(ecdsa)
}
