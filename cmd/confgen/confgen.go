package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

func main() {
	host := flag.String("host", "0.0.0.0", "ip address of the node")
	network := flag.String("network", "kovan", "ethereum network")

	flag.Parse()

	keystore, err := crypto.RandomKeystore()
	if err != nil {
		log.Fatalf("cannot create keystore: %v", err)
	}

	var ethereumConfig ethereum.Config

	switch *network {
	case "ropsten":
		ethereumConfig = ethereum.Config{
			Network:                 ethereum.NetworkRopsten,
			URI:                     "https://ropsten.infura.io",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnRopsten.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnRopsten.String(),
			RenLedgerAddress:        ethereum.RenLedgerAddressOnRopsten.String(),
		}
	case "kovan":
		ethereumConfig = ethereum.Config{
			Network:                 ethereum.NetworkKovan,
			URI:                     "https://kovan.infura.io",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnKovan.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnKovan.String(),
			RenLedgerAddress:        ethereum.RenLedgerAddressOnKovan.String(),
			RenExAccountsAddress:    ethereum.RenExAccountsAddressOnGanache.String(),
		}
	default:
		log.Fatal("unrecognized network name")
	}

	conf := config.Config{
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
		Ethereum: ethereumConfig,
	}

	bytes, err := json.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bytes))
}
