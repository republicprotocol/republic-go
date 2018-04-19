package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/ethereum/ganache"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/rpc"
)

func main() {
	// Parse the command-line arguments
	keystore := flag.String("keystore", "", "path of keystore file")
	passphrase := flag.String("passphrase", "", "passphrase to decrypt keystore")
	bindAddress := flag.String("bind", "", "bind address")
	port := flag.Int("port", 80, "port to bind API") // Defaults to 80
	token := flag.String("token", "", "optional token")
	flag.Parse()

	// Get list of nodes
	nodeAddresses := getBootstrapNodes()

	if flag.Parsed() {
		if *keystore == "" || *passphrase == "" || *bindAddress == "" {
			flag.Usage()
			return
		}

		key, err := getKey(*keystore, *passphrase)
		if err != nil {
			fmt.Println(fmt.Errorf("cannot get key: %s", err))
			return
		}

		keyPair, err := getKeyPair(key)
		if err != nil {
			fmt.Println(fmt.Errorf("cannot get keypair: %s", err))
			return
		}

		relayAddress, err := getRelayMultiaddress(keyPair, *port)
		if err != nil {
			fmt.Println(fmt.Errorf("cannot get multiaddress: %s", err))
			return
		}

		registrar, err := getRegistrar(key)
		if err != nil {
			fmt.Println(fmt.Errorf("cannot create registrar: %s", err))
			return
		}

		pools, err := getDarkPools(key, registrar)
		if err != nil {
			fmt.Println(fmt.Errorf("cannot obtain address and pools: %s", err))
			return
		}

		// Create relay node
		book := orderbook.NewOrderbook(100) // TODO: Check max connections
		relayNode := relay.NewRelay(keyPair, relayAddress, pools, *token, nodeAddresses, book)
		r := relay.NewRouter(relayNode)
		if err := http.ListenAndServe(*bindAddress, r); err != nil {
			fmt.Println(fmt.Errorf("could not start router: %s", err))
			return
		}

		// Handle orderbook synchronization
		multi, err := identity.NewMultiAddressFromString("/ip4/0.0.0.0/tcp/18415/republic/8MGzNX7M1ucyvtumVj7QYbb7wQ8UTx")
		if err != nil {
			fmt.Println(fmt.Errorf("could not generate multiaddress: %s", err))
			return
		}
		clientPool := rpc.NewClientPool(multi, keyPair).WithTimeout(10 * time.Second).WithTimeoutBackoff(0)
		go synchronizeOrderbook(&book, clientPool, registrar)
	}
}

// Synchronize orderbook using 3 randomly selected nodes
func synchronizeOrderbook(book *orderbook.Orderbook, clientPool *rpc.ClientPool, registrar contracts.DarkNodeRegistry) {
	nodes := getBootstrapNodes() // TODO: Select these randomly
	index, connections := 0, 0
	context, cancel := context.WithCancel(context.Background())
	defer cancel() // TODO: Check this
	for {
		// If there are at least 3 connections, try again in 10 seconds
		if connections >= 3 {
			time.Sleep(10 * time.Second)
			break
		}
		// TODO: Handle disconnected nodes
		multiaddressString := nodes[index%len(nodes)]
		index, connections = index+1, connections+1
		multi, err := identity.NewMultiAddressFromString(multiaddressString)
		if err != nil {
			fmt.Println(fmt.Errorf("unable to convert string %s to multiaddress: %s", multiaddressString, err))
		}
		blocks, errs := clientPool.Sync(context, multi)
		select {
		case err, ok := <-errs:
			if !ok {
				break
			}
			if err != nil {
				fmt.Println(fmt.Errorf("error when trying to sync client pool: %s", err))
			}
		case block, ok := <-blocks:
			if !ok {
				break
			}
			var epochHash [32]byte
			if len(block.EpochHash) == 32 {
				copy(epochHash[:], block.EpochHash[:32])
			} else {
				fmt.Println(fmt.Errorf("epoch hash is required to be exactly 32 bytes (%d)", len(block.EpochHash)))
				break
			}
			switch block.OrderBlock.(type) {
			case *rpc.SyncBlock_Open:
				ord := rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Open).Open)
				entry := orderbook.NewEntry(ord, order.Open, epochHash)
				book.Open(entry)
			case *rpc.SyncBlock_Confirmed:
				ord := rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Confirmed).Confirmed)
				entry := orderbook.NewEntry(ord, order.Confirmed, epochHash)
				book.Confirm(entry)
			case *rpc.SyncBlock_Unconfirmed:
				ord := rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Unconfirmed).Unconfirmed)
				entry := orderbook.NewEntry(ord, order.Unconfirmed, epochHash)
				book.Match(entry)
			case *rpc.SyncBlock_Canceled:
				ord := rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Canceled).Canceled)
				entry := orderbook.NewEntry(ord, order.Canceled, epochHash)
				book.Release(entry)
			case *rpc.SyncBlock_Settled:
				ord := rpc.UnmarshalOrder(block.OrderBlock.(*rpc.SyncBlock_Settled).Settled)
				entry := orderbook.NewEntry(ord, order.Settled, epochHash)
				book.Settle(entry)
			default:
				log.Printf("unknown order status, %t", block.OrderBlock)
			}
		}
	}
}

func getDarkPools(key *keystore.Key, registrar contracts.DarkNodeRegistry) (darknode.Pools, error) {
	ocean := darknode.NewOcean(registrar)

	// Return the dark pools
	return ocean.GetPools(), nil
}

func getKey(filename, passphrase string) (*keystore.Key, error) {
	// Read data from keystore file and generate the key
	encryptedKey, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read keystore file: %v", err)
	}

	key, err := keystore.DecryptKey(encryptedKey, passphrase)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt key with provided passphrase: %v", err)
	}

	return key, nil
}

func getKeyPair(key *keystore.Key) (identity.KeyPair, error) {
	id, err := identity.NewKeyPairFromPrivateKey(key.PrivateKey)
	if err != nil {
		return identity.KeyPair{}, fmt.Errorf("cannot generate id from key %v", err)
	}

	return id, nil
}

func getRelayMultiaddress(id identity.KeyPair, port int) (identity.MultiAddress, error) {
	// Get our IP address
	ipInfoOut, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	ipAddress := strings.Trim(string(ipInfoOut), "\n ")

	relayMultiaddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ipAddress, strconv.Itoa(port), id.Address().String()))
	if err != nil {
		return identity.MultiAddress{}, fmt.Errorf("cannot obtain trader multi address %v", err)
	}

	return relayMultiaddress, nil
}

func getRegistrar(key *keystore.Key) (contracts.DarkNodeRegistry, error) {
	// Handle the creation of the Darknode registrar
	conn, err := ganache.Connect("http://localhost:8545")
	auth := bind.NewKeyedTransactor(key.PrivateKey)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch dark node registry: %s", err))
		return contracts.DarkNodeRegistry{}, err
	}
	auth.GasPrice = big.NewInt(6000000000)
	registrar, err := contracts.NewDarkNodeRegistry(context.Background(), conn, auth, &bind.CallOpts{})
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch dark node registry: %s", err))
		return contracts.DarkNodeRegistry{}, err
	}

	return registrar, nil
}

// TODO: (temporary hard-coded bootstrap nodes) Fetch from a config file.
func getBootstrapNodes() []string {
	return []string{
		"/ip4/52.77.88.84/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
		"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
		"/ip4/52.59.176.141/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
		"/ip4/52.21.44.236/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
		"/ip4/52.41.118.171/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
	}
}
