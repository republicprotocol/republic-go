package x

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/republicprotocol/go-identity"
)

// Config information for booting a Node.
type Config struct {
	KeyPair      identity.KeyPair
	MultiAddress identity.MultiAddress
	Peers        identity.MultiAddresses
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

// MarshalJSON implements the json.Marshaler interface.
func (config *Config) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{
		"key":   config.KeyPair,
		"multi": config.MultiAddress.String(),
		"peers": make([]string, len(config.Peers)),
	}
	for i, peer := range config.Peers {
		obj["peers"].([]string)[i] = peer.String()
	}
	return json.Marshal(obj)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (config *Config) UnmarshalJSON(data []byte) error {
	obj := map[string]interface{}{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	multiAddress, ok := obj["multi"]
	if _, okType := multiAddress.(string); !ok || !okType {
		return fmt.Errorf("cannot unmarshal %t into Config struct field .MultiAddress of type string", multiAddress)
	}
	peers, ok := obj["peers"]
	if _, okType := peers.([]string); !ok || !okType {
		return fmt.Errorf("cannot unmarshal %t into Config struct field .Peers of type identity.MultiAddresses", peers)
	}
	var err error
	config.MultiAddress, err = identity.NewMultiAddress(multiAddress.(string))
	if err != nil {
		return err
	}
	config.Peers = make(identity.MultiAddresses, len(peers.([]string)))
	for i, peer := range peers.([]string) {
		config.Peers[i], err = identity.NewMultiAddress(peer)
		if err != nil {
			return err
		}
	}
	return nil
}
