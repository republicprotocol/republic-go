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

	config := map[string]interface{}{
		"bootstrapMultiAddresses": identity.MultiAddresses{},
		"ethereum": ethereum.Config{
			Network:                 ethereum.NetworkRopsten,
			URI:                     "https://ropsten.infura.io",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnRopsten.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnRopsten.String(),
			RenLedgerAddress:        ethereum.RenLedgerAddressOnKovan.String(),
		},
	}

	if err := writeKeystore(keystore); err != nil {
		log.Fatalf("cannot write keystore: %v", err)
	}
	if err := writeConfig(config); err != nil {
		log.Fatalf("cannot write config: %v", err)
	}
}

func writeKeystore(keystore crypto.Keystore) error {
	file, err := os.Create("keystore.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(keystore)
}

func writeConfig(config map[string]interface{}) error {
	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(config)
}
