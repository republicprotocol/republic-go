package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/republicprotocol/go-dark-node"
	"github.com/republicprotocol/go-do"
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

	// Setup output log file
	f, err := os.OpenFile("/home/ubuntu/darknode.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Create a new node.node.
	node, err := node.NewDarkNode(config)
	if err != nil {
		log.Fatal(err)
	}
	// Start the dark node.
	options := do.CoBegin(func() do.Option {
		return do.Err(node.StartListening())
	}, func() do.Option {
		time.Sleep(time.Second)
		return do.Err(node.Start())
	})
	for _, option := range options {
		if option.Err != nil {
			log.Println(option.Err)
		}
	}
}

func parseCommandLineFlags() error {

	confFilename := flag.String("config", "/home/ubuntu/default-config.json", "Path to the JSON configuration file")
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

// LoadDefaultConfig loads a default config if no config is provided
func LoadDefaultConfig() (*node.Config, error) {
	address, keyPair, err := identity.NewAddress()
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

	bootstrapNodes := []string{
		"/ip4/52.21.44.236/tcp/18514/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto",
		"/ip4/52.41.118.171/tcp/18514/republic/8MJ38m8Nzknh3gVj7QiMjuejmHBMSf",
		"/ip4/52.59.176.141/tcp/18514/republic/8MKDGUTgKtkymyKTH28xeMxiCnJ9xy",
		"/ip4/52.77.88.84/tcp/18514/republic/8MHarRJdvWd7SsTJE8vRVfj2jb5cWS",
		"/ip4/52.79.194.108/tcp/18514/republic/8MKZ8JwCU9m9affPWHZ9rxp2azXNnE",
	}
	config := &node.Config{
		Host:                    "0.0.0.0",
		Port:                    "18514",
		RepublicKeyPair:         keyPair,
		MultiAddress:            multiAddress,
		BootstrapMultiAddresses: make([]identity.MultiAddress, len(bootstrapNodes)),
	}

	for i, bootstrapNode := range bootstrapNodes {
		multi, err := identity.NewMultiAddressFromString(bootstrapNode)
		if err != nil {
			return &node.Config{}, err
		}
		config.BootstrapMultiAddresses[i] = multi
	}

	return config, nil
}
