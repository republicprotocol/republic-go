package main

import (
	"log"
	"os"

	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"

	"github.com/republicprotocol/go-dark-node"
)

var config *node.Config

func main() {

	configFiles := os.Args[1:]
	configs := make([]*node.Config, len(configFiles))

	for file := range configFiles {
		fileName := configFiles[file]
		config, err := node.LoadConfig(fileName)
		if err != nil {
			panic(err)
		}
		configs[file] = config
	}

	DeregisterAll(configs)
}

// DeregisterAll takes a slice of republic private keys and deregisters them
func DeregisterAll(configs []*node.Config) {
	const reset = "\x1b[0m"
	const yellow = "\x1b[33;1m"
	const green = "\x1b[32;1m"
	const red = "\x1b[31;1m"

	do.ForAll(configs, func(i int) {

		keypair := configs[i].RepublicKeyPair

		ethereumKeyPair, err := configs[i].EthereumKeyPair()
		if err != nil {
			log.Printf("[%v] %sCouldn't load Ethereum key pair%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			return
		}

		registrar, err := node.ConnectToRegistrar(ethereumKeyPair)
		if err != nil {
			log.Printf("[%v] %sCouldn't connect to registrar%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			return
		}

		isRegistered, err := registrar.IsDarkNodeRegistered(keypair.ID())
		if err != nil {
			log.Printf("[%v] %sCouldn't check node's registration%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			return
		}

		isPendingRegistration, err := registrar.IsDarkNodePendingRegistration(keypair.ID())
		if err != nil {
			log.Printf("[%v] %sCouldn't check node's registration%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			return
		}

		if isRegistered || isPendingRegistration {
			_, err = registrar.Deregister(keypair.ID())
			if err != nil {
				log.Printf("[%v] %sCouldn't deregister node%s: %v\n", base58.Encode(keypair.ID()), yellow, reset, err)
				return
			}
			if isRegistered {
				log.Printf("[%v] %sNode will be deregistered next epoch%s\n", base58.Encode(keypair.ID()), green, reset)
			} else {
				log.Printf("[%v] %sNode deregistered (registration cancelled)%s\n", base58.Encode(keypair.ID()), green, reset)
				registrar.Refund(keypair.ID())
			}
		} else {
			log.Printf("[%v] %sNode already deregistered%s\n", base58.Encode(keypair.ID()), yellow, reset)
		}
	})

}
