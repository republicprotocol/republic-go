package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/urfave/cli"
)

const (
	reset  = "\x1b[0m"
	yellow = "\x1b[33;1m"
	green  = "\x1b[32;1m"
	red    = "\x1b[31;1m"
)

// Registrar command-line tool for interacting with the darknodeRegister contract
// on Ropsten testnet.
// Set up ren contract address:
//   $ registrar --ren 0xContractAddress
// Set up dnr contract address:
//   $ registrar --dnr 0xContractAddress
// Register nodes:
//   $ registrar register 0xaddress1 0xaddress2 0xaddress3
// Deregister nodes:
//   $ registrar deregister 0xaddress1 0xaddress2 0xaddress3
// Calling epoch:
//   $ registrar epoch

func main() {

	// Load ethereum key
	key, err := LoadKey()
	if err != nil {
		log.Fatal("failt to load key from file", err)
	}

	// Create new cli application
	app := cli.NewApp()

	// Define flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ren",
			Value: "0x65d54eda5f032f2275caa557e50c029cfbccbb54",
			Usage: "republic token contract address",
		},
		cli.StringFlag{
			Name:  "dnr",
			Value: "0x69eb8d26157b9e12f959ea9f189A5D75991b59e3",
			Usage: "dark node registry address",
		},
	}

	// Define subcommands
	app.Commands = []cli.Command{
		{
			Name:    "epoch",
			Aliases: []string{"e"},
			Usage:   "calling epoch",
			Action: func(c *cli.Context) error {
				registrar, err := NewRegistrar(c, key)
				if err != nil {
					return err
				}
				_, err = registrar.Epoch()
				return err
			},
		},
		{
			Name:    "register",
			Aliases: []string{"r"},
			Usage:   "register nodes in the dark node registry",
			Action: func(c *cli.Context) error {
				registrar, err := NewRegistrar(c, key)
				if err != nil {
					return err
				}
				return RegisterAll(c.Args(), registrar)
			},
		},
		{
			Name:    "deregister",
			Aliases: []string{"d"},
			Usage:   "deregister nodes in the dark node registry",
			Action: func(c *cli.Context) error {
				registrar, err := NewRegistrar(c, key)
				if err != nil {
					return err
				}
				return DeregisterAll(c.Args(), registrar)
			},
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadKey() (*keystore.Key, error) {
	var keyJSON string = `{"address":"0066ed1af055a568e49e5a20f5b63e9741c81967",
  					"privatekey":"f5ec4b010d7cc3fbf75d71f3ea4700a64bb091e05531ef65424415207c661a39",
  					"id":"06a15d42-1d52-42dd-b714-82582a13782a",
  					"version":3}`
	key := new(keystore.Key)
	err := key.UnmarshalJSON([]byte(keyJSON))

	return key, err
}

func NewRegistrar(c *cli.Context, key *keystore.Key) (dnr.DarknodeRegistry, error) {
	config := ethereum.Config{
		Network:                 ethereum.NetworkRopsten,
		URI:                     "https://ropsten.infura.io",
		RepublicTokenAddress:    c.String("ren"),
		DarknodeRegistryAddress: c.String("dnr"),
	}

	auth := bind.NewKeyedTransactor(key.PrivateKey)
	auth.GasPrice = big.NewInt(40000000000)
	client, err := ethereum.Connect(config)
	if err != nil {
		log.Fatal("fail to connect to ethereum")
	}

	return dnr.NewDarknodeRegistry(context.Background(), client, auth, &bind.CallOpts{})
}

func RegisterAll(addresses []string, registrar dnr.DarknodeRegistry) error {
	for i := range addresses {
		address := common.HexToAddress(addresses[i])
		// Check if node has already been registered
		isRegistered, err := registrar.IsRegistered(address.Bytes())
		if err != nil {
			return fmt.Errorf("[%v] %sCouldn't check node's registration%s: %v\n", []byte(addresses[i]), red, reset, err)
		}

		// Register the node if not registered
		if !isRegistered {
			minimumBond, err := registrar.MinimumBond()
			if err != nil {
				return err
			}
			_, err = registrar.ApproveRen(&minimumBond)
			if err != nil {
				return err
			}

			_, err = registrar.Register(address.Bytes(), []byte{}, &minimumBond)
			if err != nil {
				return fmt.Errorf("[%v] %sCouldn't register node%s: %v\n", address.Hex(), red, reset, err)
			} else {
				return fmt.Errorf("[%v] %sNode will be registered next epoch%s\n", address.Hex(), green, reset)
			}
		} else if isRegistered {
			log.Printf("[%v] %sNode already registered%s\n", address.Hex(), yellow, reset)
		}
	}
	log.Println("Successfully register all node. Run 'registrar epoch' to trigger epoch")

	return nil
}

// DeregisterAll takes a slice of republic private keys and registers them
func DeregisterAll(addresses []string, registrar dnr.DarknodeRegistry) error {
	for i := range addresses {
		address := common.HexToAddress(addresses[i])
		// Check if node has already been registered
		isRegistered, err := registrar.IsRegistered(address.Bytes())
		if err != nil {
			return fmt.Errorf("[%v] %sCouldn't check node's registration%s: %v\n", address.Hex(), red, reset, err)
		}

		if isRegistered {
			_, err = registrar.Deregister(address.Bytes())
			if err != nil {
				log.Printf("[%v] %sCouldn't deregister node%s: %v\n", address.Hex(), red, reset, err)
			} else {
				log.Printf("[%v] %sNode will be deregistered next epoch%s\n", address.Hex(), green, reset)
			}
		} else {
			log.Printf("[%v] %sNode hasn't been registered yet.%s\n", address.Hex(), red, reset, err)
		}
	}
	log.Println("Successfully deregister all node. Run 'registrar epoch' to trigger epoch")

	return nil
}
