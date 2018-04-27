package main

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
)

func main() {
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
	client, err := ethereum.Connect("https://ropsten.infura.io",ethereum.NetworkRopsten, ethereum.RepublicTokenAddressOnRopsten.Hex(),
		ethereum.DarknodeRegistryAddressOnRopsten.Hex(), ethereum.HyperdriveAddressOnRopsten.Hex())
	registrar, err := dnr.NewDarknodeRegistry(context.Background(), client, auth, &bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("calling the epoch")
	_ , err = registrar.Epoch()
	if err != nil {
		log.Fatal("fail to call epoch")
	}
	log.Println("epoch finished ")
}


