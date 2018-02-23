package main

import (
	"errors"
	"flag"
	"log"

	"github.com/republicprotocol/go-dark-node"
)

var config *node.Config

func main() {
	// Parse command line arguments and fill the node.Config.
	if err := parseCommandLineFlags(); err != nil {
		log.Println(err)
		flag.Usage()
		return
	}

	// Create a new node.node.
	node, err := node.NewDarkNode(config)
	if err != nil {
		log.Fatal(err)
	}

	// Star the node.
	node.Start()
}

func parseCommandLineFlags() error {
	confFilename := flag.String("config", "", "Path to the JSON configuration file")

	flag.Parse()

	if *confFilename == "" {
		return errors.New("no config file given")
	}

	conf, err := node.LoadConfig(*confFilename)
	if err != nil {
		return err
	}
	config = conf

	return nil
}
