package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
)

var bootstrapNode = []string{
	"/ip4/0.0.0.0/tcp/3000/republic/8MJxpBsezEGKPZBbhFE26HwDFxMtFu",
	"/ip4/0.0.0.0/tcp/3001/republic/8MGB2cj2HbQFepRVs43Ghct5yCRS9C",
	"/ip4/0.0.0.0/tcp/3002/republic/8MGVBvrQJji8ecEf3zmb8SXFCx1PaR",
	"/ip4/0.0.0.0/tcp/3003/republic/8MJNCQhMrUCHuAk977igrdJk3tSzkT",
	"/ip4/0.0.0.0/tcp/3004/republic/8MK6bq5m7UfE1mzRNunJTFH6zTbyss",
}

func main() {
	for i := 0; i < 67; i++ {
		err := generateSingleNode(i)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("finsih")
}

func generateSingleNode(i int) error {
	address, keyPair, err := identity.NewAddress()
	if err != nil {
		return err
	}
	multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%s", 4000+i, address))
	if err != nil {
		return err
	}

	// Create default network options
	option := network.Options{
		MultiAddress:            multiAddress,
		BootstrapMultiAddresses: make([]identity.MultiAddress, len(bootstrapNode)),

		Debug:                network.DebugHigh,
		Alpha:                3,
		MaxBucketLength:      20,
		ClientPoolCacheLimit: 20,
		Timeout:              30 * time.Second,
		TimeoutBackoff:       30 * time.Second,
		TimeoutRetries:       3,
		Concurrent:           false,
	}
	ethKey := keystore.NewKeyForDirectICAP(rand.Reader)
	for i := range bootstrapNode {
		multi, err := identity.NewMultiAddressFromString(bootstrapNode[i])
		if err != nil {
			return err
		}
		option.BootstrapMultiAddresses[i] = multi
	}

	config := &node.Config{
		NetworkOptions: option,
		LoggerOptions: logger.Options{
			Plugins: []logger.PluginOptions{
				logger.PluginOptions{
					File: &logger.FilePluginOptions{Path: "stdout"},
				},
			},
		},
		Host:            "0.0.0.0",
		Port:            fmt.Sprintf("%d", 4000+i),
		RepublicKeyPair: keyPair,
		RSAKeyPair:      keyPair,
		EthereumKey:     ethKey,
		Dev:             true,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	d1 := []byte(data)
	err = ioutil.WriteFile(fmt.Sprintf("./node-%d.json", i+1), d1, 0644)
	return err
}
