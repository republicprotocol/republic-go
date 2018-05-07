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
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/darkocean"
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
			Name:    "checkreg",
			Aliases: []string{"c"},
			Usage:   "check if the node is registered or not",
			Action: func(c *cli.Context) error {
				registrar, err := NewRegistrar(c, key)
				if err != nil {
					return err
				}
				return CheckRegistration(c.Args(), registrar)
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
		{
			Name:    "pool",
			Aliases: []string{"p"},
			Usage:   "get the index of the pool the node is in, return -1 if no pool found",
			Action: func(c *cli.Context) error {
				registrar, err := NewRegistrar(c, key)
				if err != nil {
					return err
				}
				return GetPool(c.Args(), registrar)
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
		// Convert republic address to ethereum address
		addByte := base58.DecodeAlphabet(addresses[i], base58.BTCAlphabet)[2:]
		if len(addByte) == 0 {
			return errors.New("fail to decode the address")
		}
		address := common.BytesToAddress(addByte)

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
		// Convert republic address to ethereum address
		addByte := base58.DecodeAlphabet(addresses[i], base58.BTCAlphabet)[2:]
		if len(addByte) == 0 {
			return errors.New("fail to decode the address")
		}
		address := common.BytesToAddress(addByte)

		// Check if node has already been registered
		isRegistered, err := registrar.IsRegistered(address.Bytes())
		if err != nil {
			return fmt.Errorf("[%v] %sCouldn't check node's registration%s: %v\n", address.Hex(), red, reset, err)
		}

		if isRegistered {
			registrar.SetGasLimit(4000000)
			_, err = registrar.Refund(address.Bytes())
			if err != nil {
				return fmt.Errorf("[%v] %sCouldn't refund node%s: %v\n", address.Hex(), red, reset, err)
			}

			_, err = registrar.Deregister(address.Bytes())
			if err != nil {
				return fmt.Errorf("[%v] %sCouldn't deregister node%s: %v\n", address.Hex(), red, reset, err)
			} else {
				log.Printf("[%v] %sNode will be deregistered next epoch%s\n", address.Hex(), green, reset)
			}
			registrar.SetGasLimit(0)

		} else {
			return fmt.Errorf("[%v] %sNode hasn't been registered yet.%s\n", address.Hex(), red, reset)
		}
	}

	return nil
}

// GetPool will get the index of the pool the node is in.
// The address should be the ethereum address
func GetPool(addresses []string, registrar dnr.DarknodeRegistry) error {
	if len(addresses) != 1 {
		return fmt.Errorf("%sPlease provide one node address.%s\n", red, reset)
	}
	address := common.HexToAddress(addresses[0])

	currentEpoch, err := registrar.CurrentEpoch()
	if err != nil {
		return err
	}
	nodes, err := registrar.GetAllNodes()
	if err != nil {
		return err
	}

	ocean := darkocean.NewDarkOcean(currentEpoch.Blockhash, nodes)
	poolIndex := ocean.PoolIndex(address.Bytes())
	fmt.Println(poolIndex)

	return nil
}

// CheckRegistration will check if the node with given address is registered with
// the darknode registry. The address will be the ethereum address.
func CheckRegistration(addresses []string, registrar dnr.DarknodeRegistry) error {
	if len(addresses) != 1 {
		return fmt.Errorf("%sPlease provide one node address.%s\n", red, reset)
	}
	address := common.HexToAddress(addresses[0])
	isRegistered, err := registrar.IsRegistered(address.Bytes())
	if err != nil {
		return err
	}
	fmt.Println(isRegistered)

	return nil
}
