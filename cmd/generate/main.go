package main

import (
	"encoding/json"
	"log"
	"fmt"
	"io/ioutil"
	"github.com/republicprotocol/go-identity"
)

type Config struct {
	KeyPair         *identity.KeyPair        `json:"key"`
	Multi           *identity.MultiAddress   `json:"multi"`
	BootstrapMultis *identity.MultiAddresses `json:"bootstrap_multis"`
}

func main() {
	generateMiners(3, "../miner");
}

func generateMiners(numberOfMiners int,  location string) {
	port := 3000
	var configs []Config
	var multiAddresses identity.MultiAddresses
	for i := 0; i < numberOfMiners; i++ {
		address, keyPair, err := identity.NewAddress()
		if err != nil {
			log.Fatal(err)
		}
	
		multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s",port,address.String()))
		if err != nil {
			log.Fatal(err)
		}
	
		config := Config {
			KeyPair: &keyPair,
			Multi: &multi,
			BootstrapMultis: &identity.MultiAddresses{},
		}
		configs = append(configs,config)
		multiAddresses = append(multiAddresses,multi);
		port++
	}

	for i := 0; i < numberOfMiners ; i++ {
    data, err := json.Marshal(configs[i])
    if err != nil {
      log.Fatal(err)
    }
		d1 := []byte(data)
		err = ioutil.WriteFile(fmt.Sprintf("%s/config-miner-%d.json",location, i), d1, 0644)
		if  err != nil {
			log.Fatal(err)
		}
	}	
	
}