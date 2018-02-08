package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/republicprotocol/go-order-compute"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-miner"
	"github.com/republicprotocol/go-rpc"
)

var config map[string]interface{}
var order *compute.Order

func main() {

	// Parse command line arguments and fill the miner.Config.
	if err := parseCommandLineFlags(); err != nil {
		log.Println(err)
		flag.Usage()
		return
	}

	traderMultiAddress, err := identity.NewMultiAddressFromString(config["multi_address"].(string))
	if err != nil {
		log.Fatal(err)
	}

	nodeMultiAddressStrings := config["node_multi_addresses"].([]string)
	nodeMultiAddresses := make(identity.MultiAddresses, len(nodeMultiAddressStrings))
	for i := range nodeMultiAddresses {
		multiAddress, err := identity.NewMultiAddressFromString(nodeMultiAddressStrings[i])
		if err != nil {
			log.Fatal(err)
		}
		nodeMultiAddresses[i] = multiAddress
	}

	shares, err := order.Split(miner.N, miner.K, miner.Prime)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v fragmented and distributed\n", base58.Encode(order.ID))

	for i := range shares {
		log.Printf("  %v sent to %v\n", base58.Encode(shares[i].ID), nodeMultiAddresses[i].Address())
		if err := rpc.SendOrderFragmentToTarget(nodeMultiAddresses[i], nodeMultiAddresses[i].Address(), traderMultiAddress, shares[i], 5*time.Second); err != nil {
			log.Fatal(err)
		}
	}
}

func parseCommandLineFlags() error {
	confFilename := flag.String("config", "", "Path to the JSON configuration file")
	orderFilename := flag.String("order", "", "Path to the JSON order file")

	flag.Parse()

	if *confFilename == "" {
		return errors.New("no config file given")
	}

	file, err := os.Open(*confFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return err
	}

	file, err = os.Open(*orderFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	order = new(compute.Order)
	return json.NewDecoder(file).Decode(order)
}
