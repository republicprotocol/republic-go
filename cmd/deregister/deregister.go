package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	node "github.com/republicprotocol/republic-go/dark-node"
)

// The Secret private key to use for ethereum transactions
// If it is encrypted, a password must be provided
type Secret struct {
	PrivateKey string `json:"privateKey"`
	Password   string `json:"password"`
}

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

		keypair := configs[i].KeyPair

		clientDetails, err := connection.FromURI("https://ropsten.infura.io/", "ropsten")
		if err != nil {
			// TODO: Handler err
			panic(err)
		}

		raw, err := ioutil.ReadFile("../secrets/secrets.json")
		if err != nil {
			panic(err)
		}

		var s Secret
		json.Unmarshal(raw, &s)

		key := s.PrivateKey
		auth, err := bind.NewTransactor(strings.NewReader(key), s.Password)
		if err != nil {
			panic(err)
		}

		registrar, err := dnr.NewDarkNodeRegistry(context.Background(), &clientDetails, auth, &bind.CallOpts{})
		if err != nil {
			log.Printf("[%v] %sCouldn't connect to registrar%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			return
		}

		isRegistered, err := registrar.IsRegistered(keypair.ID())
		if err != nil {
			log.Printf("[%v] %sCouldn't check node's registration%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			return
		}

		// isPendingRegistration, err := registrar.IsDarkNodePendingRegistration(keypair.ID())
		// if err != nil {
		// 	log.Printf("[%v] %sCouldn't check node's registration%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
		// 	return
		// }

		if isRegistered { // || isPendingRegistration {
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
