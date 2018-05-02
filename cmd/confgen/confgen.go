package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

func main() {
	keystore, err := crypto.RandomKeystore()
	if err != nil {
		log.Fatalf("cannot create keystore: %v", err)
	}

	conf := darknode.Config{
		Keystore:                keystore,
		Host:                    "0.0.0.0",
		Port:                    "18514",
		Address:                 identity.Address(keystore.Address()),
		BootstrapMultiAddresses: []identity.MultiAddress{},
		Logs: logger.Options{
			Plugins: []logger.PluginOptions{
				logger.PluginOptions{
					File: &logger.FilePluginOptions{
						Path: "stdout",
					},
				},
			},
		},
		Ethereum: ethereum.Config{
			Network:                 ethereum.NetworkRopsten,
			URI:                     "https://ropsten.infura.io",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnGanache.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnGanache.String(),
		},
	}

	file, err := os.Create("config.json")
	if err != nil {
		log.Fatalf("cannot create file: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(conf); err != nil {
		log.Fatalf("cannot write conf to file: %v", err)
	}
}
