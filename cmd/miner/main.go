package main

import (
	"errors"
	"flag"
	"log"

	"github.com/republicprotocol/go-miner"
)

var config *miner.Config

func main() {
	if err := parseCommandLineFlags(); err != nil {
		log.Println(err)
		flag.Usage()
		return
	}
}

func parseCommandLineFlags() error {
	confFilename := flag.String("config", "", "Path to the JSON configuration file")

	flag.Parse()

	if *confFilename == "" {
		return errors.New("no config file given")
	}

	conf, err := miner.LoadConfig(confFilename)
	if err != nil {
		return err
	}
	config = conf

	return nil
}
