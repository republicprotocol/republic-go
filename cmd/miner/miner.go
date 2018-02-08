package main

import (
	"errors"
	"flag"
	"log"

	"github.com/republicprotocol/go-miner"
)

var config *miner.Config

func main() {
	// Parse command line arguments and fill the miner.Config.
	if err := parseCommandLineFlags(); err != nil {
		log.Println(err)
		flag.Usage()
		return
	}

	// Create a new miner.Miner.
	miner, err := miner.NewMiner(config)
	if err != nil {
		log.Fatal(err)
	}

	// Star the miner.
	miner.Start()
}

func parseCommandLineFlags() error {
	confFilename := flag.String("config", "", "Path to the JSON configuration file")

	flag.Parse()

	if *confFilename == "" {
		return errors.New("no config file given")
	}

	conf, err := miner.LoadConfig(*confFilename)
	if err != nil {
		return err
	}
	config = conf

	return nil
}
