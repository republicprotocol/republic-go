package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/republicprotocol/republic-go/order"

	. "github.com/republicprotocol/republic-go/relay"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/blockchain/bitcoin/arc"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/arc"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/dht"
	"github.com/republicprotocol/republic-go/rpc/relayer"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/rpc/swarmer"
	"google.golang.org/grpc"
)

func main() {
	keystore := flag.String("keystore", "", "Encrypted keystore file")
	passphrase := flag.String("passphrase", "", "Passphrase for the encrypted keystore file")
	bind := flag.String("bind", "127.0.0.1", "Binding address for the gRPC and HTTP API")
	port := flag.String("port", "18515", "Binding port for the HTTP API")
	token := flag.String("token", "", "Bearer token for restricting access")
	maxConnections := flag.Int("maxConnections", 4, "Maximum number of connections to peers during synchronization")
	flag.Parse()

	fmt.Println("Decrypting keystore...")
	key, err := getKey(*keystore, *passphrase)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot obtain key: %s", err))
		return
	}

	keyPair, err := getKeyPair(key)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot obtain keypair: %s", err))
		return
	}

	multiAddr, err := getMultiaddress(keyPair, *port)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot obtain multiaddress: %s", err))
		return
	}

	registrar, err := getRegistry(key)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot obtain registrar: %s", err))
		return
	}

	// Create gRPC server and TCP listener always using port 18514
	server := grpc.NewServer()
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:18514", *bind))
	if err != nil {
		log.Fatal(err)
	}

	// Create Relay
	config := Config{
		KeyPair:      keyPair,
		MultiAddress: multiAddr,
		Token:        *token,
	}
	book := orderbook.NewOrderbook(100)
	crypter := crypto.NewWeakCrypter()
	dht := dht.NewDHT(multiAddr.Address(), 100)
	connPool := client.NewConnPool(100)
	relayerClient := relayer.NewClient(&crypter, &dht, &connPool)
	smpcerClient := smpcer.NewClient(&crypter, multiAddr, &connPool)
	swarmerClient := swarmer.NewClient(&crypter, multiAddr, &dht, &connPool)
	relay := NewRelay(config, registrar, &book, &relayerClient, &smpcerClient, &swarmerClient)

	entries := make(chan orderbook.Entry)
	defer close(entries)
	go func() {
		defer book.Unsubscribe(entries)
		if err := book.Subscribe(entries); err != nil {
			log.Fatalf("cannot subscribe to orderbook: %v", err)
		}
	}()

	// Server gRPC and RESTful API
	fmt.Println(fmt.Sprintf("Relay API available at %s:%s", *bind, *port))
	dispatch.CoBegin(func() {
		if err := relay.ListenAndServe(*bind, *port); err != nil {
			log.Fatalf("error serving http: %v", err)
		}
	}, func() {
		relay.Register(server)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("error serving grpc: %v", err)
		}
	}, func() {
		relay.Sync(context.Background(), *maxConnections)
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		<-sig
		server.Stop()
	}()

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func processOrderbookEntries(entryInCh <-chan orderbook.Entry) <-chan orderbook.Entry {
	unconfirmedOrders := make(chan orderbook.Entry, 100)
	confirmedEntries := make(chan orderbook.Entry)
	go func() {
		defer close(confirmedEntries)
		for {
			select {
			case entry, ok := <-entryInCh:
				if !ok {
					return
				}
				if !orderConfirmed(entry.Order.ID) {
					unconfirmedOrders <- entry
				} else {
					entry.Status = order.Confirmed
					confirmedEntries <- entry
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case entry, ok := <-unconfirmedOrders:
				if !ok {
					return
				}
				if !orderConfirmed(entry.Order.ID) {
					unconfirmedOrders <- entry
					time.Sleep(time.Second)
				} else {
					entry.Status = order.Confirmed
					confirmedEntries <- entry
				}
			}
		}
	}()
	return confirmedEntries
}

func atomicSwap(entries <-chan orderbook.Entry, privateKey *ecdsa.PrivateKey) error {
	conn, err := client.Connect(uri, network, republicTokenAddress, darknodeRegistryAddr)
	if err != nil {
		return err
	}
	transOps := bind.NewKeyedTransactor(privateKey)
	arc.NewArc(context.Background(), conn, transOps)
	return nil
}

func orderConfirmed(orderID order.ID) bool {
	return false
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

func getMultiaddress(id identity.KeyPair, port string) (identity.MultiAddress, error) {
	// Get our IP address
	ipInfoOut, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	ipAddress := strings.Trim(string(ipInfoOut), "\n ")

	relayMultiaddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ipAddress, port, id.Address().String()))
	if err != nil {
		return identity.MultiAddress{}, fmt.Errorf("cannot obtain trader multi address %v", err)
	}

	return relayMultiaddress, nil
}

func getRegistry(key *keystore.Key) (dnr.DarknodeRegistry, error) {
	conn, err := ganache.Connect("http://localhost:8545")
	auth := bind.NewKeyedTransactor(key.PrivateKey)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch dark node registry: %s", err))
		return dnr.DarknodeRegistry{}, err
	}
	auth.GasPrice = big.NewInt(6000000000)
	registrar, err := dnr.NewDarknodeRegistry(context.Background(), conn, auth, &bind.CallOpts{})
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch dark node registry: %s", err))
		return dnr.DarknodeRegistry{}, err
	}

	return registrar, nil
}
