package main

// DeployDarkNodeRegistry

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
)

const key = `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`

const reset = "\x1b[0m"
const green = "\x1b[32;1m"
const red = "\x1b[31;1m"

func main() {
	clientDetails, err := connection.FromURI("https://ropsten.infura.io/", "ropsten")
	if err != nil {
		// TODO: Handler err
		panic(err)
	}

	auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
	if err != nil {
		panic(err)
	}

	// Gas Price
	auth.GasPrice = big.NewInt(6000000000)

	registrar, err := dnr.NewDarkNodeRegistry(context.Background(), &clientDetails, auth, &bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	minimumEpochTime, err := registrar.MinimumEpochInterval()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Calling Epoch every %s%v seconds%s\n", green, minimumEpochTime, reset)

	callEpoch(registrar)
	uInt, err := minimumEpochTime.ToUint()
	if err != nil {
		log.Fatalln(err)
	}
	ticker := time.NewTicker(time.Duration(uInt) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		/*go */ callEpoch(registrar)
	}
}

func callEpoch(registrar dnr.DarkNodeRegistry) {
	fmt.Printf("Calling Epoch...")
	_, err := registrar.Epoch()

	if err != nil {
		fmt.Printf("\r")
		log.Printf("%sCouldn't call Epoch%s\n", red, reset)
	} else {
		fmt.Printf("\r")
		log.Printf("%sEpoch called%s", green, reset)
	}
}
