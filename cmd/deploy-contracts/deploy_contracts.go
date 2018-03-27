package main

// DeployDarkNodeRegistrar

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/contracts/bindings"
	"github.com/republicprotocol/republic-go/contracts/connection"
	node "github.com/republicprotocol/republic-go/dark-node"
)

var config *node.Config

const reset = "\x1b[0m"
const green = "\x1b[32;1m"

func main() {
	err := parseCommandLineFlags()
	if err != nil {
		log.Fatalln(err)
	}

	auth := bind.NewKeyedTransactor(config.EthereumKey.PrivateKey)

	client, err := connection.FromURI("https://ropsten.infura.io/", "ropsten")
	if err != nil {
		log.Fatal(err)
	}

	// REPLACE REN ADDRESS HERE
	renContract := common.HexToAddress("0x889debfe1478971bcff387f652559ae1e0b6d34a")
	address, tx, _, err := bindings.DeployDarkNodeRegistrar(auth, client.Client, renContract, big.NewInt(1000), big.NewInt(60))
	if err != nil {
		log.Fatalln(err)
	}
	_, err = connection.PatchedWaitDeployed(context.Background(), client.Client, tx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[%v] Contract deployed at %s%v%s\n", base58.Encode(config.KeyPair.ID()), green, address.Hex(), reset)
}

func parseCommandLineFlags() error {
	confFilename := flag.String("config", "../darknode/config/ap-northeast-2.json", "Path to the JSON configuration file")

	flag.Parse()

	conf, err := node.LoadConfig(*confFilename)
	if err != nil {
		return err
	}
	config = conf
	return nil
}
