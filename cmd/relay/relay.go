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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/ethereum/ganache"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/relay"
	"google.golang.org/grpc"
)

func main() {
	// Parse the command-line arguments
	keystore := flag.String("keystore", "", "path of keystore file")
	passphrase := flag.String("passphrase", "", "passphrase to decrypt keystore")
	bindAddress := flag.String("bind", "", "bind address")
	port := flag.Int("port", 80, "port to bind API") // Defaults to 80
	token := flag.String("token", "", "optional token")
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

	config := relay.Config{
		KeyPair:        keyPair,
		MultiAddress:   multiAddr,
		Token:          *token,
		BootstrapNodes: getBootstrapNodes(),
		BindAddress:    *bindAddress,
	}
	darknodeIDs, err := registrar.GetAllNodes()
	if err != nil {
		fmt.Println(fmt.Errorf("could not get dark nodes: %s", err))
		return
	}

	epoch, err := registrar.CurrentEpoch()
	if err != nil {
		fmt.Println(fmt.Errorf("could not obtain epoch: %s", err))
		return
	}

	darkOcean := darkocean.NewDarkOcean(epoch.Blockhash, darknodeIDs)

	// return darkOcean.Pools()
	pools := darkOcean.Pools()
	book := orderbook.NewOrderbook(100) // TODO: Check max connections
	relay.RunRelay(config, pools, book, registrar)

	server := grpc.NewServer()
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", *bindAddress, *port))
	if err != nil {
		log.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		<-signals
		log.Println("shutting down...")
		server.Stop()
		os.Exit(0)
	}()

	go server.Serve(listener)
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
