package swarm

import (
	"github.com/republicprotocol/go-identity"
	"os"
	"encoding/json"
)

// Config struct holds configuration details for connecting to peers
// at boot.
type Config struct {
	Address identity.Address `json:"address"`
	Host    string           `json:"host"`
	Port    int              `json:"port"`
	Peers   []Config		 `json:"peers,omitempty"`
}

// LoadConfig loads a Config object from the given filename. Returns the Config
// object, or an error.
func LoadConfig(filename string ) (*Config, error)  {

	file , err := os.Open(filename)
	if err !=nil {
		return nil, err
	}

	config := new(Config)
	if err := json.NewDecoder(file).Decode(config); err !=nil {
		return nil, err
	}

	return config, nil
}



