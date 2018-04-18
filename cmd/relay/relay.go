package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"net/http"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/relay"
)

func main() {
	// Parse the command-line arguments
	keystore := flag.String("keystore", "", "path of keystore file")
	passphrase := flag.String("passphrase", "", "passphrase to decrypt keystore")
	bindAddress := flag.String("bind", "", "bind address")
	port := flag.Int("port", 80, "port to bind API")       // defaults to 80
	token := flag.String("token", "", "optional token")
	flag.Parse()

	if flag.Parsed() {
		if *keystore == "" || *passphrase == "" || *bindAddress == "" {
			flag.Usage()
			return
		}
		keyPair, relayAddress, pools, err := getRelayAddressAndDarkPools(*keystore, *passphrase, *port)
		if err != nil {
			fmt.Println(fmt.Errorf("cannot obtain address and pools: %v", err))
			return
		}
		relayNode := relay.NewRelay(keyPair, relayAddress, pools, *token, getBootstrapNodes())
		r := relay.NewRouter(relayNode)
		if err := http.ListenAndServe(*bindAddress, r); err != nil {
			fmt.Printf("could not start router: %v", err)
			return
		}
	}
}

func getRelayAddressAndDarkPools(keyFile, passphrase string, port int) (identity.KeyPair, identity.MultiAddress, dark.Pools, error) {
	// Get key and traderMultiAddress
	keyPair, key, multiAddress, err := getKeyPairAndAddress(keyFile, passphrase, port)
	if err != nil {
		return (identity.KeyPair{}), (identity.MultiAddress{}), nil, err
	}

	// Get dark pools
	pools, err := getDarkPools(key)
	if err != nil {
		return (identity.KeyPair{}), (identity.MultiAddress{}), nil, err
	}

	return keyPair, multiAddress, pools, nil
}

func getLogger() (*logger.Logger, error) {
	return logger.NewLogger(logger.Options{
		Plugins: []logger.PluginOptions{
			logger.PluginOptions{
				File: &logger.FilePluginOptions{
					Path: "stdout",
				},
			},
		}})
}

// Returns key and multi address
func getKeyPairAndAddress(filename, passphrase string, port int) (identity.KeyPair, *keystore.Key, identity.MultiAddress, error) {

	// Get our IP address
	ipInfoOut, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	if err != nil {
		return (identity.KeyPair{}), nil, (identity.MultiAddress{}), err
	}
	ipAddress := strings.Trim(string(ipInfoOut), "\n ")

	// Read data from keystore file and generate the key
	encryptedKey, err := ioutil.ReadFile(filename)
	if err != nil {
		return (identity.KeyPair{}), nil, (identity.MultiAddress{}), fmt.Errorf("cannot read keystore file: %v", err)
	}

	key, err := keystore.DecryptKey(encryptedKey, passphrase)
	if err != nil {
		return (identity.KeyPair{}), nil, (identity.MultiAddress{}), fmt.Errorf("cannot decrypt key with provided passphrase: %v", err)
	}

	id, err := identity.NewKeyPairFromPrivateKey(key.PrivateKey)
	if err != nil {
		return (identity.KeyPair{}), nil, (identity.MultiAddress{}), fmt.Errorf("cannot generate id from key %v", err)
	}

	traderMultiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ipAddress, strconv.Itoa(port), id.Address().String()))
	if err != nil {
		return (identity.KeyPair{}), nil, (identity.MultiAddress{}), fmt.Errorf("cannot obtain trader multi address %v", err)
	}

	return id, key, traderMultiAddress, nil
}

func getDarkPools(key *keystore.Key) (dark.Pools, error) {
	// Logger
	logs, err := getLogger()
	if err != nil {
		return nil, fmt.Errorf("cannot get logger: %v", err)
	}

	// Get Dark Node Registry
	clientDetails, err := connection.FromURI("https://ropsten.infura.io/", connection.ChainRopsten)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch dark node registry: %v", err)
	}

	auth := bind.NewKeyedTransactor(key.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch dark node registry: %v", err)
	}

	// Gas Price
	auth.GasPrice = big.NewInt(6000000000)

	registrar, err := dnr.NewDarkNodeRegistry(context.Background(), &clientDetails, auth, &bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("cannot fetch dark node registry: %v", err)
	}

	// print(registrar.client)
	// Get the Dark Ocean
	ocean, err := dark.NewOcean(logs, 5, registrar)
	if err != nil {
		return nil, fmt.Errorf("cannot read dark ocean: %v", err)
	}

	// return the dark pools
	return ocean.GetPools(), nil
}

//TODO: (temporary hard-coded bootstrap nodes) Fetch from a config file.
func getBootstrapNodes() []string {
	return []string{
		"/ip4/52.77.88.84/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
		"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
		"/ip4/52.59.176.141/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
		"/ip4/52.21.44.236/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
		"/ip4/52.41.118.171/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
	}
}
