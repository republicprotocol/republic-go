package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

func main() {

	keystore, err := crypto.RandomKeystore()
	if err != nil {
		log.Fatalf("cannot create ecdsa key: %v", err)
	}

	conf := map[string]interface{}{
		"keystore":                keystore,
		"bootstrapMultiAddresses": identity.MultiAddresses{},
		"ethereum": ethereum.Config{
			Network:                 ethereum.NetworkRopsten,
			URI:                     "https://ropsten.infura.io",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnRopsten.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnRopsten.String(),
			RenLedgerAddress:       ethereum.RenAddressOnRopsten.String(),
			ArcAddress:              ethereum.ArcAddressOnRopsten.String(),
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
