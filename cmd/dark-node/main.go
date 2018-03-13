package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
)

func main() {

	// Load the default configuration
	config, err := LoadConfig("/home/.darknode/config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create a dark node registrar.
	darkNodeRegistrar, err := CreateDarkNodeRegistrar(config.EthereumKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new node.node.
	node, err := node.NewDarkNode(*config, darkNodeRegistrar)
	if err != nil {
		log.Fatal(err)
	}

	go node.StartServices()
	node.StartBackgroundWorkers()
	node.Bootstrap()
	node.WatchDarkOcean()
}

// LoadConfig returns a default Config object for the Falcon testnet.
func LoadConfig(filename string) (*node.Config, error) {

	// Load configuration file
	config, err := node.LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	// Generate our ethereum keypair
	if config.EthereumKey == nil {
		config.EthereumKey = keystore.NewKeyForDirectICAP(rand.Reader)
	}

	if config.RepublicKeyPair == nil {
		// Get an address and keypair
		address, keyPair, err := identity.NewAddress()
		if err != nil {
			return nil, err
		}

		// Get our IP address
		out, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
		out = []byte(strings.Trim(string(out), "\n "))
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		// Generate our multiaddress
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%v/tcp/%v/republic/%v", string(out), config.Port, address.String()))
		if err != nil {
			return nil, err
		}

		config.RepublicKeyPair = &keyPair
		config.NetworkOptions.MultiAddress = multiAddress
	}

	return config, err
}

func CreateDarkNodeRegistrar(key *keystore.Key) (dnr.DarkNodeRegistrar, error) {
	auth := bind.NewKeyedTransactor(key.PrivateKey)
	client, err := connection.FromURI("https://ropsten.infura.io/")
	if err != nil {
		return nil, err
	}
	return dnr.NewEthereumDarkNodeRegistrar(context.Background(), &client, auth, &bind.CallOpts{})
}
