package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/stackint"
)

const reset = "\x1b[0m"
const yellow = "\x1b[33;1m"
const green = "\x1b[32;1m"
const red = "\x1b[31;1m"

const key = `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`

// The Secret private key to use for ethereum transactions
// If it is encrypted, a password must be provided
type Secret struct {
	PrivateKey string `json:"privateKey"`
	Password   string `json:"password"`
}

func main() {

	secretFile := os.Args[1]
	configFiles := os.Args[2:]
	configs := make([]*node.Config, len(configFiles))

	for file := range configFiles {
		fileName := configFiles[file]
		config, err := node.LoadConfig(fileName)
		if err != nil {
			log.Fatal(err)
		}
		configs[file] = config
	}

	RegisterAll(secretFile, configs)
}

// RegisterAll takes a slice of republic private keys and registers them
func RegisterAll(secretFile string, configs []*node.Config) {

	/*
		0x3ccB53DBB5f801C28856b3396B01941ecD21Ac1d
		0xFd99C99825781AD6795025b055e063c8e8863e7c
		0x22846b4d7962806c1B450e354B91f9bF33697244
		0x1629de08ec625d2452a564e5e1990f6890f85a5e
	*/

	do.ForAll(configs, func(i int) {
		keypair := configs[i].KeyPair
		var decimalMultiplier = stackint.FromUint(1000000000000000000)
		var bondTokenCount = stackint.Zero()
		var bond = decimalMultiplier.Mul(&bondTokenCount)
		clientDetails, err := connection.FromURI("https://ropsten.infura.io/", "ropsten")
		if err != nil {
			log.Fatal(err)
		}

		auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
		if err != nil {
			log.Fatal(err)
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

		if !isRegistered { // && !isPendingRegistration {
			_, err = registrar.Register(keypair.ID(), append(keypair.PublicKey.X.Bytes(), keypair.PublicKey.Y.Bytes()...), &bond)
			if err != nil {
				log.Printf("[%v] %sCouldn't register node%s: %v\n", base58.Encode(keypair.ID()), red, reset, err)
			} else {
				log.Printf("[%v] %sNode will be registered next epoch%s\n", base58.Encode(keypair.ID()), green, reset)
			}
		} else if isRegistered {
			log.Printf("[%v] %sNode already registered%s\n", base58.Encode(keypair.ID()), yellow, reset)
		} // else if isPendingRegistration {
		// 	log.Printf("[%v] %sNode will already be registered next epoch%s\n", base58.Encode(keypair.ID()), yellow, reset)
		// }
	})

}
