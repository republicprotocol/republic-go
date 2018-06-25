package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

func main() {
	network := flag.String("network", "nightly", "Republic Protocol network")

	flag.Parse()

	keystore, err := crypto.RandomKeystore()
	if err != nil {
		log.Fatalf("cannot create keystore: %v", err)
	}

	var ethereumConfig contract.Config

	switch *network {
	case string(contract.NetworkTestnet):
		ethereumConfig = contract.Config{
			Network: contract.NetworkTestnet,
			URI:     "https://kovan.infura.io",
		}
	case string(contract.NetworkFalcon):
		ethereumConfig = contract.Config{
			Network: contract.NetworkFalcon,
			URI:     "https://kovan.infura.io",
		}
	case string(contract.NetworkNightly):
		ethereumConfig = contract.Config{
			Network: contract.NetworkNightly,
			URI:     "https://kovan.infura.io",
		}
	default:
		log.Fatal("unrecognized network name")
	}

	conf := config.Config{
		Keystore:                keystore,
		Host:                    "0.0.0.0",
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
