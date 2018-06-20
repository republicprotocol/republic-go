package config

import (
	"encoding/json"
	"os"

	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

type Config struct {
	Keystore crypto.Keystore `json:"keystore"`
	Ethereum contract.Config `json:"ethereum"` // TODO: Darknode package should not be dependent on blockchain/ethereum
	Logs     logger.Options  `json:"logs"`

	Address                 identity.Address        `json:"address"`
	BootstrapMultiAddresses identity.MultiAddresses `json:"bootstrapMultiAddresses"`
	Host                    string                  `json:"host"`
	Port                    string                  `json:"port"`
}

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
