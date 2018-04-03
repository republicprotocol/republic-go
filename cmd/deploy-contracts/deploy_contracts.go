package main

import (
	"github.com/republicprotocol/republic-go/dark-node"
)

var config *node.Config

const reset = "\x1b[0m"
const green = "\x1b[32;1m"

func main() {
	//	err := parseCommandLineFlags()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	ethereumKeyPair := config.EthereumKey
	//	auth := bind.NewKeyedTransactor(ethereumKeyPair.PrivateKey)
	//
	//	client, err := connection.FromURI("https://ropsten.infura.io/", connection.ChainRopsten)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	// REPLACE REN ADDRESS HERE
	//	renContract := common.HexToAddress("0x889debfe1478971bcff387f652559ae1e0b6d34a")
	//	address, _ , _, err := bind.DeployContract(auth, client, renContract, big.NewInt(1000), big.NewInt(60))
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	//_, err = dnr.PatchedWaitDeployed(context.Background(), client, tx)
	//	//if err != nil {
	//	//	log.Fatalln(err)
	//	//}
	//	fmt.Printf("[%v] Contract deployed at %s%v%s\n", base58.Encode(config.KeyPair.ID()), green, address.Hex(), reset)
}
