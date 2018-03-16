package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	node "github.com/republicprotocol/republic-go/dark-node"
)

const reset = "\x1b[0m"
const yellow = "\x1b[33;1m"
const green = "\x1b[32;1m"
const red = "\x1b[31;1m"

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

	RegisterAll(configs)
}

// RegisterAll takes a slice of republic private keys and registers them
func RegisterAll(configs []*node.Config) {

	/*
		0x3ccB53DBB5f801C28856b3396B01941ecD21Ac1d
		0xFd99C99825781AD6795025b055e063c8e8863e7c
		0x22846b4d7962806c1B450e354B91f9bF33697244
		0x1629de08ec625d2452a564e5e1990f6890f85a5e
	*/

	do.ForAll(configs, func(i int) {
		keypair := configs[i].RepublicKeyPair
		var decimalMultiplier = big.NewInt(1000000000000000000)
		var bondTokenCount = big.NewInt(100)
		var bond = decimalMultiplier.Mul(decimalMultiplier, bondTokenCount)
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

		registrar, err := dnr.NewEthereumDarkNodeRegistrar(context.Background(), &clientDetails, auth, &bind.CallOpts{})
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

		if !isRegistered && !isPendingRegistration {
			_, err = registrar.Register(keypair.ID(), append(keypair.PublicKey.X.Bytes(), keypair.PublicKey.Y.Bytes()...), bond)
			if err != nil {
				log.Printf("[%v] %sCouldn't register node%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			} else {
				log.Printf("[%v] %sNode will be registered next epoch%s\n", base58.Encode(keypair.ID()), green, reset)
			}
		} else if isRegistered {
			log.Printf("[%v] %sNode already registered%s\n", base58.Encode(keypair.ID()), yellow, reset)
		} else if isPendingRegistration {
			log.Printf("[%v] %sNode will already be registered next epoch%s\n", base58.Encode(keypair.ID()), yellow, reset)
		}
	})

}
