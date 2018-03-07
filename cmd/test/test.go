package main

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
		log.Fatal(err)
	}

	node, err := node.NewDarkNode(config)
	if err != nil {
		log.Fatal(err)
	}

	node.Start()
}

func parseCommandLineFlags() error {
	confFilename := flag.String("config", "./config.json", "Path to the JSON configuration file")

	flag.Parse()

	conf, err := node.LoadConfig(*confFilename)
	if err != nil {
		return err
	}
	config = conf
	return nil
}
