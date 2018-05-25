package main

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/identity"
)

var addresses = []string{
	"8MGxw5z4Fm8TS4h3nJ535yynHgTw4H",
	"8MJjW5KTkoDwVpqenjBHGQBhww86qT",
	"8MHkxMUmDAP5Y5fKDGDAp24sT2qY9y",
	"8MJSxE3VXcpnWdTR4Q688jAgoYRSDU",
	"8MGFa6sJDmHAztawDai19SoYgQyDZe",
	"8MJpfZfu32Ss3mgRUjshTSe3UM1mF4",
	"8MGuXYKJVVrgTrcKvrCN7ar5LLkeHi",
	"8MGjykjQ2qrosJ3mQkxDPhAmAubDdh",
}

const (
	reset  = "\x1b[0m"
	yellow = "\x1b[33;1m"
	green  = "\x1b[32;1m"
	red    = "\x1b[31;1m"
)

func main() {
	key, err := LoadKey()
	if err != nil {
		log.Fatal(err)
	}
	registry, err := NewRegistry(key)
	if err != nil {
		log.Fatal(err)
	}

	for _, address := range addresses {
		addr := identity.Address(address)

		// Check if node has already been registered
		isRegistered, err := registry.IsRegistered(addr)
		if err != nil {
			log.Fatalf("[%v] %sCouldn't check node's registration%s: %v\n", address, red, reset, err)
		}

		// Register the node if not registered
		if isRegistered {
			_, err = registry.Deregister(addr.ID())
			if err != nil {
				log.Fatalf("[%v] %sCouldn't deregister node%s: %v\n", address, red, reset, err)
			} else {
				log.Printf("[%v] %sNode will be deregistered next epoch%s\n", address, green, reset)
			}
		} else if isRegistered {
			log.Printf("[%v] %sNode already registered%s\n", address, yellow, reset)
		}
	}

	log.Println("triggering epoch ")
	_, err = registry.TriggerEpoch()
	if err != nil {
		log.Fatalf("failt to trigger epoch ")
	}
	log.Println("epoch  called ")
}

func LoadKey() (*keystore.Key, error) {
	var keyJSON string = `{"address":"90e6572ef66a11690b09dd594a18f36cf76055c8",
  					"privatekey":"dc3f937b4aa1fc7bbf7643f1dead1faf37594ad2f1edcd6b56bf6719f85fa406",
  					"id":"ddd54c1c-6c2e-42a9-a224-6532a90fd4e9", "version":3}`
	key := new(keystore.Key)
	err := key.UnmarshalJSON([]byte(keyJSON))

	return key, err
}

func NewRegistry(key *keystore.Key) (dnr.DarknodeRegistry, error) {
	config := ethereum.Config{
		Network:                 ethereum.NetworkRopsten,
		URI:                     "https://ropsten.infura.io",
		RepublicTokenAddress:    "0x65d54eda5f032f2275caa557e50c029cfbccbb54",
		DarknodeRegistryAddress: "0x69eb8d26157b9e12f959ea9f189A5D75991b59e3",
	}

	auth := bind.NewKeyedTransactor(key.PrivateKey)
	auth.GasPrice = big.NewInt(40000000000)
	//auth.GasLimit = 150000

	client, err := ethereum.Connect(config)
	if err != nil {
		log.Fatal("fail to connect to ethereum")
	}

	return dnr.NewDarknodeRegistry(context.Background(), client, auth, &bind.CallOpts{})
}
