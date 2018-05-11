package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/republicprotocol/republic-go/http/adapter"
	. "github.com/republicprotocol/republic-go/relay"

	abiBind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/grpc/client"
	"github.com/republicprotocol/republic-go/grpc/dht"
	"github.com/republicprotocol/republic-go/grpc/smpcer"
	"github.com/republicprotocol/republic-go/grpc/swarmer"
	"github.com/republicprotocol/republic-go/http"
	"github.com/republicprotocol/republic-go/identity"
)

type Config struct {
	EthereumConfig      ethereum.Config         `json:"ethereum"`
	BootstrapMultiAddrs identity.MultiAddresses `json:"bootstrapMultiAddresses"`
}

func main() {
	bindParam := flag.String("bind", "127.0.0.1", "Binding address for the gRPC and HTTP API")
	portParam := flag.String("port", "18515", "Binding port for the HTTP API")
	configParam := flag.String("config", "", "Ethereum configuration file")
	keystoreParam := flag.String("keystore", "", "Optionally encrypted keystore file")
	passphraseParam := flag.String("passphrase", "", "Optional passphrase to decrypt the keystore file")
	flag.Parse()

	config, err := loadConfig(*configParam)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	keystore, err := loadKeystore(*keystoreParam, *passphraseParam)
	if err != nil {
		log.Fatalf("cannot load keystore: %v", err)
	}

	multiAddr, err := getMultiaddress(keystore, *portParam)
	if err != nil {
		log.Fatalf("cannot get multi-address: %v", err)
	}

	registry, err := getRegistry(config.EthereumConfig, keystore)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot get registry: %s", err))
		return
	}

	dht := dht.NewDHT(multiAddr.Address(), 100)
	connPool := client.NewConnPool(100)
	crypter := darkocean.NewCrypter(keystore, registry, 128, time.Minute)
	smpcerClient := smpcer.NewClient(&crypter, multiAddr, &connPool)
	swarmerClient := swarmer.NewClient(&crypter, multiAddr, &dht, &connPool)
	relay := NewRelay(&registry, nil /*RenLedger*/, &swarmerClient, &smpcerClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for err := range swarmerClient.Bootstrap(ctx, config.BootstrapMultiAddrs, -1) {
		log.Printf("error bootstrapping %v", err)
	}

	relayAdapter := adapter.NewRelayAdapter(&relay)

	log.Printf("Listening at %v:%v", *bindParam, *portParam)
	if err := http.ListenAndServe(*bindParam, *portParam, &relayAdapter, &relayAdapter); err != nil {
		log.Fatalf("error listening and serving: %v", err)
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

func getMultiaddress(keystore crypto.Keystore, port string) (identity.MultiAddress, error) {
	// Get our IP address
	ipInfoOut, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	ipAddress := strings.Trim(string(ipInfoOut), "\n ")
	relayMultiaddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ipAddress, port, keystore.Address()))
	if err != nil {
		return identity.MultiAddress{}, fmt.Errorf("cannot obtain trader multi address %v", err)
	}
	return relayMultiaddress, nil
}

func getRegistry(ethereumConfig ethereum.Config, keystore crypto.Keystore) (dnr.DarknodeRegistry, error) {
	conn, err := ethereum.Connect(ethereumConfig)
	auth := abiBind.NewKeyedTransactor(keystore.EcdsaKey.PrivateKey)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch dark node registry: %s", err))
		return dnr.DarknodeRegistry{}, err
	}
	auth.GasPrice = big.NewInt(6000000000)
	registry, err := dnr.NewDarknodeRegistry(context.Background(), conn, auth, &abiBind.CallOpts{})
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch dark node registry: %s", err))
		return dnr.DarknodeRegistry{}, err
	}
	return registry, nil
}

func loadConfig(configFile string) (Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	config := new(Config)
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return Config{}, err
	}
	return *config, nil
}

func loadKeystore(keystoreFile, passphrase string) (crypto.Keystore, error) {
	file, err := os.Open(keystoreFile)
	if err != nil {
		return crypto.Keystore{}, err
	}
	defer file.Close()
	keystore := new(crypto.Keystore)
	data, err := ioutil.ReadFile(keystoreFile)
	if err := json.NewDecoder(file).Decode(keystore); err != nil {
		return crypto.Keystore{}, err
	}
	err = keystore.DecryptFromJSON(data, passphrase)
	if err := json.NewDecoder(file).Decode(keystore); err != nil {
		return crypto.Keystore{}, err
	}
	return *keystore, nil
}
