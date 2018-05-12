package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	netHttp "net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/http"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/stackint"

	"github.com/republicprotocol/republic-go/crypto"
)

func main() {
	relayParam := flag.String("relay", "http://127.0.0.1:18515", "Binding and port of the relay")
	keystoreParam := flag.String("keystore", "", "Optionally encrypted keystore file")
	configParam := flag.String("config", "", "Ethereum configuration file")
	passphraseParam := flag.String("passphrase", "", "Optional passphrase to decrypt the keystore file")
	flag.Parse()

	keystore, err := loadKeystore(*keystoreParam, *passphraseParam)
	if err != nil {
		log.Fatalf("cannot load keystore: %v", err)
	}

	config, err := loadConfig(*configParam)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	auth, darkpool, err := loadSmartContracts(config, keystore)
	if err != nil {
		log.Fatalf("cannot load smart contracts: %v", err)
	}
	log.Printf("ethereum %v", auth.From.Hex())

	nonce, err := stackint.Random(rand.Reader, &smpc.Prime)
	if err != nil {
		log.Fatalf("cannot generate nonce: %v", err)
	}
	ord := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, stackint.One(), stackint.One(), stackint.One(), nonce)

	signatureData := crypto.Keccak256([]byte("Republic Protocol: open: "), ord.ID)
	signatureData = crypto.Keccak256([]byte("\x19Ethereum Signed Message:\n32"), signatureData)
	signature, err := keystore.Sign(crypto.NewHash32(signatureData))

	log.Printf("id = %v; sig = %v", base64.StdEncoding.EncodeToString(ord.ID), base64.StdEncoding.EncodeToString(signature))

	if err != nil {
		log.Fatalf("cannot sign order: %v", err)
	}

	request := http.OpenOrderRequest{
		Signature:            base64.StdEncoding.EncodeToString(signature),
		OrderFragmentMapping: map[string][]order.Fragment{},
	}
	log.Printf("order signature %v", request.Signature)

	pods, err := darkpool.Pods()
	if err != nil {
		log.Fatalf("cannot get pods from darkpool: %v", err)
	}
	for _, pod := range pods {
		n := int64(len(pod.Darknodes))
		k := int64(2 * (len(pod.Darknodes) + 1) / 3)
		hash := base64.StdEncoding.EncodeToString(pod.Hash[:])
		ordFragments, err := ord.Split(n, k, &smpc.Prime)
		if err != nil {
			log.Fatalf("cannot split order to %v: %v", hash, err)
		}
		request.OrderFragmentMapping[hash] = []order.Fragment{}
		for _, ordFragment := range ordFragments {
			request.OrderFragmentMapping[hash] = append(request.OrderFragmentMapping[hash], *ordFragment)
		}
	}

	data, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		log.Fatalf("cannot marshal request: %v", err)
	}
	fmt.Println(string(data))
	buf := bytes.NewBuffer(data)

	res, err := netHttp.DefaultClient.Post(fmt.Sprintf("%v/orders", *relayParam), "application/json", buf)
	if err != nil {
		log.Fatalf("cannot send request: %v", err)
	}
	defer res.Body.Close()

	log.Printf("status: %v", res.StatusCode)
	resText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("cannot read response body: %v", err)
	}
	log.Printf("body: %v", string(resText))
}

func loadKeystore(keystoreFile, passphrase string) (crypto.Keystore, error) {
	file, err := os.Open(keystoreFile)
	if err != nil {
		return crypto.Keystore{}, err
	}
	defer file.Close()

	if passphrase == "" {
		keystore := crypto.Keystore{}
		if err := json.NewDecoder(file).Decode(&keystore); err != nil {
			return keystore, err
		}
		return keystore, nil
	}

	keystore := crypto.Keystore{}
	keystoreData, err := ioutil.ReadAll(file)
	if err != nil {
		return keystore, err
	}
	if err := keystore.DecryptFromJSON(keystoreData, passphrase); err != nil {
		return keystore, err
	}
	return keystore, nil
}

func loadConfig(configFile string) (ethereum.Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return ethereum.Config{}, err
	}
	defer file.Close()
	config := ethereum.Config{}
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return ethereum.Config{}, err
	}
	return config, nil
}

func loadSmartContracts(ethereumConfig ethereum.Config, keystore crypto.Keystore) (*bind.TransactOpts, cal.Darkpool, error) {
	conn, err := ethereum.Connect(ethereumConfig)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot connect to ethereum: %v", err))
		return nil, nil, err
	}
	auth := bind.NewKeyedTransactor(keystore.EcdsaKey.PrivateKey)
	auth.GasPrice = big.NewInt(10)

	registry, err := dnr.NewDarknodeRegistry(context.Background(), conn, auth, &bind.CallOpts{})
	if err != nil {
		fmt.Println(fmt.Errorf("cannot bind to darkpool: %v", err))
		return auth, nil, err
	}
	return auth, &registry, nil
}
