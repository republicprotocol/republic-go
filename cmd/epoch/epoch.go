package main

// DeployDarkNodeRegistrar

import (
	"flag"
	"log"

	node "github.com/republicprotocol/go-dark-node"
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

	ethereumKeyPair, err := config.EthereumKeyPair()
	if err != nil {
		log.Fatalln(err)
	}

	registrar, err := node.ConnectToRegistrar(ethereumKeyPair)

	_, err = registrar.Epoch()

	if err != nil {
		log.Fatalf("%sCouldn't call epoch%s\n", red, reset)
	}
	log.Printf("%sEpoch called%s", green, reset)
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
