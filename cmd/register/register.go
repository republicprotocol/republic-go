package main

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/darknode"
)

const RepublicTokenAddress = "0x65d54eda5f032f2275caa557e50c029cfbccbb54"
const DarknodeRegistryAddress = "0x69eb8d26157b9e12f959ea9f189A5D75991b59e3"
const HyperdriveAddress = "0x348496ad820f2ee256268f9f9d0b9f5bacdc26cd"

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

	configFiles := os.Args[1:]
	configs := make([]*darknode.Config, len(configFiles))

	for file := range configFiles {
		fileName := configFiles[file]
		config, err := darknode.LoadConfig(fileName)
		if err != nil {
			log.Fatal(err)
		}
		configs[file] = config
	}

	RegisterAll(configs)
}

// RegisterAll takes a slice of republic private keys and registers them
func RegisterAll(configs []*darknode.Config) {

	/*
		0x3ccB53DBB5f801C28856b3396B01941ecD21Ac1d
		0xFd99C99825781AD6795025b055e063c8e8863e7c
		0x22846b4d7962806c1B450e354B91f9bF33697244
		0x1629de08ec625d2452a564e5e1990f6890f85a5e
	*/

	do.ForAll(configs, func(i int) {
		key := new(keystore.Key)
		file, err := os.Open("key.json")
		if err != nil {
			log.Fatal("cannot read key file")
		}

		err = json.NewDecoder(file).Decode(key)
		if err != nil {
			log.Fatal("fail to parse the ethereum key")
		}

		auth := bind.NewKeyedTransactor(key.PrivateKey)
		auth.GasPrice = big.NewInt(6000000000)

		// Create the eth-client so we can interact with the Registrar contract
		client, err := ethereum.Connect("https://ropsten.infura.io",
			ethereum.NetworkRopsten, ethereum.RepublicTokenAddressOnRopsten.Hex(),
			ethereum.DarknodeRegistryAddressOnRopsten.Hex(),
			ethereum.HyperdriveAddressOnRopsten.Hex())
		registrar, err := dnr.NewDarknodeRegistry(context.Background(), client, auth, &bind.CallOpts{})
		if err != nil {
			log.Fatal(err)
		}

		isRegistered, err := registrar.IsRegistered(configs[i].EcdsaKey.Address.Bytes())

		if err != nil {
			log.Printf("[%v] %sCouldn't check node's registration%s: %v\n", configs[i].EcdsaKey.Address, red, reset, err)
			return
		}


		if !isRegistered {
			minimumBond, err := registrar.MinimumBond()
			if err != nil {
				log.Fatal(err)
			}

			_, err = registrar.ApproveRen(&minimumBond)
			if err != nil {
				log.Fatal(err)
			}

			_, err = registrar.Register(configs[i].EcdsaKey.Id, []byte{}, &minimumBond)
			if err != nil {
				log.Printf("[%v] %sCouldn't register node%s: %v\n",configs[i].EcdsaKey.Address, red, reset, err)
			} else {
				log.Printf("[%v] %sNode will be registered next epoch%s\n", configs[i].EcdsaKey.Address, green, reset)
			}
		} else if isRegistered {
			log.Printf("[%v] %sNode already registered%s\n", configs[i].EcdsaKey.Address, yellow, reset)
		}
	})
}
