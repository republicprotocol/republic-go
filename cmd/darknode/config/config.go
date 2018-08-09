package config

import (
	"encoding/json"
	"os"

	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

type Config struct {
	Ethereum contract.Config `json:"ethereum"` // TODO: Darknode package should not be dependent on blockchain/ethereum
	Logs     logger.Options  `json:"logs"`

	Address                 identity.Address        `json:"address"`
	OracleAddress           identity.Address        `json:"oracleAddress"`
	BootstrapMultiAddresses identity.MultiAddresses `json:"bootstrapMultiAddresses"`
	Host                    string                  `json:"host"`
	Port                    string                  `json:"port"`
	Alpha                   int                     `json:"alpha"`
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
