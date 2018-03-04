package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-dark-node"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
)

const DefaultTimeOut = 5 * time.Second

func test() {

	numberOfBootstrapNodes := flag.Int("bootstrapNodes", 4, "number of bootstrap nodes")
	numberOfNodes := flag.Int("nodes", 8, "number of nodes")

	nodes, err := generateNodes(*numberOfBootstrapNodes, *numberOfNodes)
	if err != nil {
		log.Fatal("fail to generate nodes --> ", err)
	}

	errChan := deployNodes(nodes, *numberOfBootstrapNodes)
	if len(errChan) > 0 {
		for err := range errChan {
			log.Println("fail to deploy nodes --> ", err)
		}
		log.Fatal("-------------------------------")
	}

	sendingOrders(nodes, 20, 5)
	log.Printf("finish")
}

func generateNodes(numberOfBootstrapNodes, numberOfTestNodes int) ([]*node.DarkNode, error) {
	// Generate config files for the nodes
	err := generateConfigFiles(numberOfBootstrapNodes, numberOfTestNodes)
	if err != nil {
		return nil, err
	}

	// Generate nodes from the config files
	numberOfNodes := numberOfBootstrapNodes + numberOfTestNodes
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

func generateConfigFiles(numberOfBootstrapNodes, numberOfTestNodes int) error {
	// Generate configs
	numberOfNodes := numberOfBootstrapNodes + numberOfTestNodes
	port := 3000
	configs := make([]*node.Config, numberOfNodes)
	bootstraps := make([]identity.MultiAddress, numberOfBootstrapNodes)
	for i := 0; i < numberOfNodes; i++ {
		address, keyPair, err := identity.NewAddress()
		if err != nil {
			log.Fatal(err)
		}

		multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", port+i, address.String()))
		if err != nil {
			log.Fatal(err)
		}

		config := node.Config{
			RepublicKeyPair: keyPair,
			RSAKeyPair:      keyPair,
			Host:            "127.0.0.1",
			Port:            fmt.Sprintf("%d", port+i),
			MultiAddress:    multi,
		}
		configs[i] = &config
		if i < numberOfBootstrapNodes {
			bootstraps[i] = multi
		}
	}

	// Write configs to files
	for i, config := range configs {
		for _, bootstrap := range bootstraps {
			if config.MultiAddress.String() != bootstrap.String() {
				config.BootstrapMultiAddresses = append(config.BootstrapMultiAddresses, bootstrap)
			}
		}

		data, err := json.Marshal(config)
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

func deployNodes(nodes []*node.DarkNode, numberOfBoostrapNodes int) chan error {
	// Start all nodes
	errChan := make(chan error, len(nodes)*3)
	for i := range nodes {
		go func(i int) {
			err := nodes[i].StartListening()
			if err != nil {
				errChan <- err
			}
		}(i)
	}

	// Fully connect all bootstrap nodes
	for i := 0; i < numberOfBoostrapNodes; i++ {
		for j := 0; j < numberOfBoostrapNodes; j++ {
			if i != j {
				go func(i, j int) {
					err := rpc.PingTarget(nodes[j].Swarm.MultiAddress(), nodes[i].Swarm.MultiAddress(), DefaultTimeOut)
					if err != nil {
						errChan <- err
					}
				}(i, j)
			}
		}
	}

	time.Sleep(1 * time.Second)

	for i := range nodes {
		go func(i int) {
			err := nodes[i].Start()
			if err != nil {
				errChan <- err
			}
		}(i)
	}
	return errChan
}

type OrderBook struct {
	LastUpdateId int        `lastUpdateId`
	Bids         [][]string `bids`
	Asks         [][]string `asks`
}

func sendingOrders(nodes []*node.DarkNode, numberOfOrders, timeInterval int) error {
	start := time.Now()
	// Keep sending order fragment
	for time.Now().Sub(start) < 1*time.Minute {

		// Get orders details from Binance
		resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v1/depth?symbol=ETHBTC&limit=%v", numberOfOrders))
		if err != nil {
			return err
		}

		response, err := ioutil.ReadAll(resp.Body)
		orderBook := new(OrderBook)
		err = json.Unmarshal(response, orderBook)

		// Generate order from the Binance data
		sellOrders := make([]*compute.Order, len(orderBook.Asks))
		for i, j := range orderBook.Asks {
			price, err := strconv.ParseFloat(j[0], 10)
			price = price * 1000000000000
			if err != nil {
				return err
			}

			amount, err := strconv.ParseFloat(j[1], 10)
			amount = amount * 1000000000000
			if err != nil {
				return err
			}
			order := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Time{},
				compute.CurrencyCodeETH, compute.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
				big.NewInt(int64(amount)), big.NewInt(1))
			sellOrders[i] = order
		}

		buyOrders := make([]*compute.Order, len(orderBook.Bids))

		//test cast for match
		for i, j := range orderBook.Asks { //change asks/bids
			price, err := strconv.ParseFloat(j[0], 10)
			price = price * 1000000000000
			if err != nil {
				return err
			}

			amount, err := strconv.ParseFloat(j[1], 10)
			amount = amount * 1000000000000
			if err != nil {
				return err
			}

			order := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Time{},
				compute.CurrencyCodeETH, compute.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
				big.NewInt(int64(amount)), big.NewInt(1))
			buyOrders[i] = order
		}
		// Send order fragment to the nodes
		for _, orders := range [][]*compute.Order{buyOrders, sellOrders} {
			go func(orders []*compute.Order) {
				for _, order := range orders {
					//todo
					if order.Parity == compute.OrderParityBuy {
						log.Println("sending buy order :", base58.Encode(order.ID))
					} else {
						log.Println("sending sell order :", base58.Encode(order.ID))
					}

					shares, err := order.Split(8, 5, node.Prime)
					if err != nil {
						continue
					}

					do.ForAll(shares, func(i int) {
						rpc.SendOrderFragmentToTarget(nodes[i].Swarm.MultiAddress(),
							nodes[i].Swarm.Address(), nodes[0].Swarm.MultiAddress(), shares[i], 5*time.Second)
					})
				}
			}(orders)
		}
		time.Sleep(time.Duration(timeInterval) * time.Second)
	}
	return nil
}


