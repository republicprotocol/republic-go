package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/republicprotocol/go-dark-node"
	"github.com/republicprotocol/go-identity"
)

func main() {

	numberOfNodes := flag.Int("nodes", 8, "number of nodes")
	numberOfBoostrapNodes := flag.Int("boostrapNodes", 4 , "number of boostrap nodes")

	nodes, err := generateNodes(*numberOfBoostrapNodes, *numberOfNodes)
	if err !=nil {
		log.Fatal("fail to generate nodes-->", err)
	}

	err = deployNodes(nodes)
	if err != nil {
		log.Fatal("fail to deploy nodes-->", err)
	}

	sendingOrders()

	collectLogs()

	log.Printf("finish")
}

func generateNodes(numberOfBootsrapNodes, numberOfTestNodes int) ([]*node.DarkNode, error) {
	// Generate nodes from the config files
	numberOfNodes := numberOfBootsrapNodes + numberOfTestNodes
	nodes := make([]*node.DarkNode, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		config, err := node.LoadConfig(fmt.Sprintf("./configs/config-%d.json", i))
		if err != nil {
			return nil, err
		}
		node, err := node.NewDarkNode(config)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return nodes, nil
}

func generateConfigFiles(numberOfBootsrapNodes, numberOfTestNodes int) error {
	// Generate configs
	numberOfNodes := numberOfBootsrapNodes + numberOfTestNodes
	port := 3000
	configs := make ([]*node.Config, numberOfNodes)
	bootstraps := make([]identity.MultiAddress, numberOfBootsrapNodes)
	for i := 0; i < numberOfNodes; i++ {
		address, keyPair, err := identity.NewAddress()
		if err != nil {
			log.Fatal(err)
		}

		multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", port +i , address.String()))
		if err != nil {
			log.Fatal(err)
		}

		config := node.Config{
			RepublicKeyPair:         keyPair,
			RSAKeyPair: keyPair,
			Host :  "127.0.0.1",
			Port:   fmt.Sprintf("%d",port+i ),
			MultiAddress:   multi,
		}
		configs[i] = &config
		if i < numberOfBootsrapNodes{
			bootstraps[i] = multi
		}
	}

	// Write configs to files
	for i , config := range configs {
		for _, bootstrap := range bootstraps{
			if config.MultiAddress.String() != bootstrap.String() {
				config.BootstrapMultiAddresses = append(config.BootstrapMultiAddresses, bootstrap)
			}
		}

		data ,err  := json.Marshal(config)
		if err != nil {
			return err
		}
		d1 := []byte(data)
		err = ioutil.WriteFile(fmt.Sprintf("./configs/config-%d.json", i), d1, 0644)
		if err != nil {
			return err
		}

	}
	return nil
}

func deployNodes(nodes []*node.DarkNode) error {
	return nil
}

func sendingOrders() {

}

func collectLogs() {

}
