package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

func main() {
	host := flag.String("host", "0.0.0.0", "ip address of the node")

	flag.Parse()

	keystore, err := crypto.RandomKeystore()
	if err != nil {
		log.Fatalf("cannot create keystore: %v", err)
	}

	conf := darknode.Config{
		Keystore:                keystore,
		Host:                    *host,
		Port:                    "18514",
		Address:                 identity.Address(keystore.Address()),
		BootstrapMultiAddresses: []identity.MultiAddress{},
		Logs: logger.Options{
			Plugins: []logger.PluginOptions{
				{
					File: &logger.FilePluginOptions{
						Path: "/home/ubuntu/.darknode/darknode.out",
					},
				},
			},
		},
		Ethereum: ethereum.Config{
			Network:                 ethereum.NetworkRopsten,
			URI:                     "https://ropsten.infura.io",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnRopsten.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnRopsten.String(),
			RenLedgerAddress:        ethereum.LedgerAddressOnRopsten.String(),
		},
	}

	bytes, err := json.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bytes))
}
