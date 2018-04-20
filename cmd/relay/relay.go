package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"
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
			fmt.Println(fmt.Errorf("could not obtain key: %s", err))
			return
		}

		keyPair, err := getKeyPair(key)
		if err != nil {
			fmt.Println(fmt.Errorf("could not obtain keypair: %s", err))
			return
		}

		relayAddress, err := getRelayMultiaddress(keyPair, *port)
		if err != nil {
			fmt.Println(fmt.Errorf("could not obtain multiaddress: %s", err))
			return
		}

		registrar, err := getRegistrar(key)
		if err != nil {
			fmt.Println(fmt.Errorf("could not create registrar: %s", err))
			return
		}

		pools, err := getDarkPools(key, registrar)
		if err != nil {
			fmt.Println(fmt.Errorf("could not obtain address and pools: %s", err))
			return
		}

		// Handle orderbook synchronization
		book := orderbook.NewOrderbook(100) // TODO: Check max connections
		multi, err := identity.NewMultiAddressFromString("/ip4/0.0.0.0/tcp/18415/republic/8MGzNX7M1ucyvtumVj7QYbb7wQ8UTx")
		if err != nil {
			fmt.Println(fmt.Errorf("could not generate multiaddress: %s", err))
			return
		}
		clientPool := rpc.NewClientPool(multi, keyPair).WithTimeout(10 * time.Second).WithTimeoutBackoff(0)
		go synchronizeOrderbook(&book, clientPool, registrar)

		// Create relay node
		relayNode := relay.NewRelay(keyPair, relayAddress, pools, *token, nodeAddresses, book)
		r := relay.NewRouter(relayNode)
		if err := http.ListenAndServe(*bindAddress, r); err != nil {
			fmt.Println(fmt.Errorf("could not start router: %s", err))
			return
		}
	}
}

// Synchronize orderbook using 3 randomly selected nodes.
func synchronizeOrderbook(book *orderbook.Orderbook, clientPool *rpc.ClientPool, registrar contracts.DarkNodeRegistry) {
	nodes, err := registrar.GetAllNodes()
	if err != nil {
		fmt.Println(fmt.Errorf("could not retrieve nodes: %s", err))
	}
	connections := int32(0)
	context, cancel := context.WithCancel(context.Background())
	defer cancel() // TODO: Check this
	for {
		// If there are at least 3 connections, we try again in 10 seconds.
		if atomic.LoadInt32(&connections) >= 3 {
			time.Sleep(10 * time.Second)
			break
		}

		// Select a node in a random position and increment the number of
		// connected nodes.
		randIndex := rand.Intn(len(nodes))
		multiaddressString := nodes[randIndex]
		atomic.AddInt32(&connections, 1)

		// Retrieve the multiaddress of the selected node.
		multi, err := identity.NewMultiAddressFromString(string(multiaddressString))
		if err != nil {
			fmt.Println(fmt.Errorf("unable to convert \"%s\" to multiaddress: %s", multiaddressString, err))
		}

		// Check for any messages received from this node and forward them to
		// orderbook stored in the relay.
		blocks, errs := clientPool.Sync(context, multi)
		go func() {
			if err := forwardMessages(blocks, errs, &connections, book); err != nil {
				fmt.Println(err)
			}
		}()
	}
}

// Forward any messages we receive from the channels and store them in the
// orderbook.
func forwardMessages(blocks <-chan *rpc.SyncBlock, errs <-chan error, connections *int32, book *orderbook.Orderbook) error {
	// When the function ends we decrement the total number of connections.
	defer atomic.AddInt32(connections, -1)
	for {
		select {
		case err, ok := <-errs:
			// Output an error and end the connection.
			if !ok || err != nil {
				return fmt.Errorf("error when trying to sync client pool: %s", err)
			}
		case block, ok := <-blocks:
			if !ok {
				return fmt.Errorf("error when trying to sync client pool")
			}

			// The epoch hash we retrieve is stored in a dynamic sized byte
			// array, so we must copy this to one of a fixed length in order
			// to include it in the order entry.
			var epochHash [32]byte
			if len(block.EpochHash) == 32 {
				copy(epochHash[:], block.EpochHash[:32])
			} else {
				return fmt.Errorf("epoch hash is required to be exactly 32 bytes (%d)", len(block.EpochHash))
			}

			// Store this entry in the relay orderbook.
			if err := storeEntry(block, epochHash, book); err != nil {
				return err
			}
		}
	}
}

func storeEntry(block *rpc.SyncBlock, epochHash [32]byte, book *orderbook.Orderbook) error {
	// Check the status of the order message received and call the
	// corresponding function from the orderbook.
	switch block.OrderBlock.(type) {
	case *rpc.SyncBlock_Open:
		ord := rpc.UnmarshalOrder(block.GetOpen())
		entry := orderbook.NewEntry(ord, order.Open, epochHash)
		if err := book.Open(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	case *rpc.SyncBlock_Confirmed:
		ord := rpc.UnmarshalOrder(block.GetConfirmed())
		entry := orderbook.NewEntry(ord, order.Confirmed, epochHash)
		if err := book.Confirm(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	case *rpc.SyncBlock_Unconfirmed:
		ord := rpc.UnmarshalOrder(block.GetUnconfirmed())
		entry := orderbook.NewEntry(ord, order.Unconfirmed, epochHash)
		if err := book.Match(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	case *rpc.SyncBlock_Canceled:
		ord := rpc.UnmarshalOrder(block.GetCanceled())
		entry := orderbook.NewEntry(ord, order.Canceled, epochHash)
		if err := book.Release(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	case *rpc.SyncBlock_Settled:
		ord := rpc.UnmarshalOrder(block.GetSettled())
		entry := orderbook.NewEntry(ord, order.Settled, epochHash)
		if err := book.Settle(entry); err != nil {
			return fmt.Errorf("error when synchronizing order: %s", err)
		}
	default:
		return fmt.Errorf("unknown order status, %t", block.OrderBlock)
	}

	return nil
}

func getDarkPools(key *keystore.Key, registrar contracts.DarkNodeRegistry) (darknode.Pools, error) {
	ocean := darknode.NewOcean(registrar)
	return ocean.GetPools(), nil
}

func getKey(filename, passphrase string) (*keystore.Key, error) {
	// Read data from the keystore file and generate the key
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
