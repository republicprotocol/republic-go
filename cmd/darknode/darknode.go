package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/republicprotocol/go-dark-node"
	identity "github.com/republicprotocol/go-identity"
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

	// Start listening.
	go func() {
		if err := node.StartListening(); err != nil {
			log.Fatal(err)
		}
	}()

	// Star the node.
	time.Sleep(time.Second)
	node.Start()
	node.StartListening()
}

func parseCommandLineFlags() error {
	confFilename := flag.String("config", "./default-config.json", "Path to the JSON configuration file")

	flag.Parse()

	conf, err := node.LoadConfig(*confFilename)
	if err != nil {
		conf, err = LoadDefaultConfig()
		if err != nil {
			return err
		}
		config = conf
		return nil
	}
	config = conf
	return nil
}

func LoadDefaultConfig() (*node.Config, error) {
	address, keypair, err := identity.NewAddress()
	if err != nil {
		return &node.Config{}, err
	}
	out, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	out = []byte(strings.Trim(string(out), "\n "))
	if err != nil {
		return &node.Config{}, err
	}
	multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/18514/republic/%s", out, address))
	if err != nil {
		return &node.Config{}, err
	}
	bootstrap1, err := identity.NewMultiAddressFromString("/ip4/52.78.129.188/tcp/18514/republic/8MKZ8JwCU9m9affPWHZ9rxp2azXNnE")
	if err != nil {
		return &node.Config{}, err
	}
	bootstrap2, err := identity.NewMultiAddressFromString("/ip4/54.179.189.157/tcp/18514/republic/8MHarRJdvWd7SsTJE8vRVfj2jb5cWS")
	if err != nil {
		return &node.Config{}, err
	}
	bootstrap3, err := identity.NewMultiAddressFromString("/ip4/52.59.192.207/tcp/18514/republic/8MKDGUTgKtkymyKTH28xeMxiCnJ9xy")
	if err != nil {
		return &node.Config{}, err
	}
	bootstrap4, err := identity.NewMultiAddressFromString("/ip4/54.88.24.57/tcp/18514/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto")
	if err != nil {
		return &node.Config{}, err
	}
	bootstrap5, err := identity.NewMultiAddressFromString("/ip4/34.217.114.249/tcp/18514/republic/8MJ38m8Nzknh3gVj7QiMjuejmHBMSf")
	if err != nil {
		return &node.Config{}, err
	}

	return &node.Config{
		Host:            "0.0.0.0",
		Port:            "18514",
		RepublicKeyPair: keypair,
		MultiAddress:    multiAddress,
		BootstrapMultiAddresses: identity.MultiAddresses{
			bootstrap1, bootstrap2, bootstrap3, bootstrap4, bootstrap5,
		},
	}, nil

}
