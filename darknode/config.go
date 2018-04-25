package darknode

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// A Config defines the different settings for a Darknode.
type Config struct {
	EcdsaKey keystore.Key      `json:"ecdsaKey"`
	RsaKey   crypto.RsaKeyPair `json:"rsaKey"`
	Ethereum EthereumConfig    `json:"ethereum"`
	Logs     logger.Options    `json:"logs"`

	Address                 identity.Address        `json:"address"`
	BootstrapMultiAddresses identity.MultiAddresses `json:"bootstrapMultiAddresses"`
	Host                    string                  `json:"host"`
	Port                    string                  `json:"port"`
}

// An EthereumConfig defines the different settings for connecting the Darknode
// to an Ethereum network, and the Republic Protocol smart contracts deployed
// on Ethereum.
type EthereumConfig struct {
	Network                 client.Network `json:"network"` // One of "ganache", "ropsten", or "mainnet" ("mainnet" is not current supported)
	URI                     string         `json:"uri"`
	RepublicTokenAddress    string         `json:"republicTokenAddress"`
	DarknodeRegistryAddress string         `json:"darknodeRegistryAddress"`
	TraderRegistryAddress   string         `json:"traderRegistryAddress"`
	HyperdriveAddress       string         `json:"hyperdriveAddress"`
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
