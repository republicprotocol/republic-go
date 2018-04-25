package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	. "github.com/republicprotocol/republic-go/relay"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
	"github.com/republicprotocol/republic-go/identity"
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

	multiAddr, err := getMultiaddress(keyPair, *port)
	if err != nil {
		fmt.Println(fmt.Errorf("could not obtain multiaddress: %s", err))
		return
	}

	registrar, err := getRegistrar(key)
	if err != nil {
		fmt.Println(fmt.Errorf("could not obtain registrar: %s", err))
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
		KeyPair:        keyPair,
		MultiAddress:   multiAddr,
		Token:          *token,
		BootstrapNodes: getBootstrapNodes(),
	}
	relay := NewRelay(...)

	// Server gRPC and RESTful API
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
	}),


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

func getMultiaddress(id identity.KeyPair, port int) (identity.MultiAddress, error) {
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

func getRegistry(key *keystore.Key) (dnr.DarknodeRegistry, error) {
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
