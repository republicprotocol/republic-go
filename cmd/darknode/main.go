package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
)

const PATH = "/home/ubuntu/"
//const PATH = ""

var config *node.Config

func main() {

	// Parse command line arguments and fill the node.Config.
	if err := parseCommandLineFlags(); err != nil {
		log.Println(err)
		flag.Usage()
		return
	}
	// Create a new node.node.
	node, err := node.NewDarkNode(*config)
	if err != nil {
		log.Fatal(err)
	}

	node.Start()
	defer node.Stop()
}

// Parse the config file path and read config from it.
func parseCommandLineFlags() error {
	confFilename := flag.String("config", fmt.Sprintf("%sdefault-config.json", PATH), "Path to the JSON configuration file")

	flag.Parse()

	conf, err := node.LoadConfig(*confFilename)
	if err != nil {
		log.Fatal("error :", err)
		conf, err = LoadDefaultConfig()
		if err != nil {
			return err
		}
		config = conf
		return nil
	}
	config = conf
	// Create plugins for logger.
	stdoutPlugin := logger.NewFilePlugin("stout")
	filePlugin := logger.NewFilePlugin(PATH)
	websocketPlugin := logger.NewWebSocketPlugin("0.0.0.0", "8080", "", "")
	config.Logger = logger.NewLogger(stdoutPlugin, filePlugin, websocketPlugin)

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

	// 5 bootstraps nodes set up by the team.
	bootstrapNodes := []string{
		"/ip4/52.21.44.236/tcp/18514/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto",
		"/ip4/52.41.118.171/tcp/18514/republic/8MJ38m8Nzknh3gVj7QiMjuejmHBMSf",
		"/ip4/52.59.176.141/tcp/18514/republic/8MKDGUTgKtkymyKTH28xeMxiCnJ9xy",
		"/ip4/52.77.88.84/tcp/18514/republic/8MHarRJdvWd7SsTJE8vRVfj2jb5cWS",
		"/ip4/52.79.194.108/tcp/18514/republic/8MKZ8JwCU9m9affPWHZ9rxp2azXNnE",
	}

	// Create default network options
	option := network.Options{
		MultiAddress:            multiAddress,
		BootstrapMultiAddresses: make([]identity.MultiAddress, len(bootstrapNodes)),
		Debug:                network.DebugHigh,
		Alpha:                3,
		MaxBucketLength:      20,
		ClientPoolCacheLimit: 20,
		Timeout:              30 * time.Second,
		TimeoutBackoff:       30 * time.Second,
		TimeoutRetries:       1,
		Concurrent:           false,
	}
	for i, bootstrapNode := range bootstrapNodes {
		multi, err := identity.NewMultiAddressFromString(bootstrapNode)
		if err != nil {
			return &node.Config{}, err
		}
		option.BootstrapMultiAddresses[i] = multi
	}

	// Create plugins for logger.
	stdoutPlugin := logger.NewFilePlugin("stout")
	filePlugin := logger.NewFilePlugin(PATH)
	websocketPlugin := logger.NewWebSocketPlugin("0.0.0.0", "8080", "", "")

	ethKey := keystore.NewKeyForDirectICAP(rand.Reader)
	// Generate default config file
	config = &node.Config{
		Options:         option,
		Host:            "0.0.0.0",
		Port:            "18514",
		RepublicKeyPair: keyPair,
		RSAKeyPair:      keyPair,
		EthereumKey:     ethKey,
		Logger:          logger.NewLogger(stdoutPlugin, filePlugin, websocketPlugin),
		//todo : missing some of the fields
	}
	//err = saveConfigFile(config)
	return config, err
}

func saveConfigFile(config *node.Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	d1 := []byte(data)
	err = ioutil.WriteFile(fmt.Sprintf("%sdefault-config.json", PATH), d1, 0644)
	return err
}
