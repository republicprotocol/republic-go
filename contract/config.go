package contract

import (
	"encoding/json"
	"os"
)

// Network is used to represent a Republic Protocol network.
type Network string

const (
	// NetworkMainnet represents the mainnet
	NetworkMainnet Network = "mainnet"
	// NetworkTestnet represents the internal F∅ testnet
	NetworkTestnet Network = "testnet"
	// NetworkFalcon represents the internal Falcon testnet
	NetworkFalcon Network = "falcon"
	// NetworkNightly represents the internal Nightly testnet
	NetworkNightly Network = "nightly"
	// NetworkLocal represents a local network
	NetworkLocal Network = "local"
)

// Config defines the different settings for connecting to Ethereum on
// different Republic Protocol networks.
type Config struct {
	Network                 Network `json:"network"`
	URI                     string  `json:"uri"`
	RepublicTokenAddress    string  `json:"republicTokenAddress"`
	DarknodeRegistryAddress string  `json:"darknodeRegistryAddress"`
	OrderbookAddress        string  `json:"orderbookAddress"`
	RewardVaultAddress      string  `json:"rewardVaultAddress"`
	RenExBalancesAddress    string  `json:"renExBalancesAddress"`
	RenExSettlementAddress  string  `json:"renExSettlementAddress"`
}

// NewConfigFromJSONFile unmarshals a JSON into a Config object
func NewConfigFromJSONFile(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	conf := Config{}
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		return Config{}, err
	}
	return conf, nil
}
