package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

func main() {
	ecdsaKey := keystore.NewKeyForDirectICAP(rand.Reader)

	keyPair, err := identity.NewKeyPairFromPrivateKey(ecdsaKey.PrivateKey)
	if err != nil {
		log.Fatalf("cannot create ecdsa key: %v", err)
	}

	rsaKey, err := crypto.NewRsaKeyPair()
	if err != nil {
		log.Fatalf("cannot create rsa key: %v", err)
	}

	conf := darknode.Config{
		EcdsaKey:                ecdsaKey,
		RsaKey:                  rsaKey,
		Host:                    "0.0.0.0",
		Port:                    "18514",
		Address:                 keyPair.Address(),
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
