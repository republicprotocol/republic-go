package main

// DeployDarkNodeRegistrar

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark-node"
)

var config *node.Config

const reset = "\x1b[0m"
const green = "\x1b[32;1m"
const red = "\x1b[31;1m"

func main() {
	err := parseCommandLineFlags()
	if err != nil {
		log.Fatalln(err)
	}

	clientDetails, err := connection.FromURI("https://ropsten.infura.io/")
	if err != nil {
		// TODO: Handler err
		panic(err)
	}
	registrar, err := node.ConnectToRegistrar(clientDetails, *config)
	if err != nil {
		// TODO: Handler err
		panic(err)
	}

	minimumEpochTime, err := registrar.MinimumEpochInterval()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Calling Epoch every %s%v seconds%s\n", green, minimumEpochTime, reset)

	callEpoch(registrar)
	ticker := time.NewTicker(time.Duration(minimumEpochTime.Int64()) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		/*go */ callEpoch(registrar)
	}
}

func callEpoch(registrar *dnr.DarkNodeRegistrar) {
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
