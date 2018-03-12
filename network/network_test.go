package network_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/network"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

const (
	numberOfBootstrapNodes = 5
)

func generateNodes(numberOfNodes int) ([]*network.SwarmService, error) {
	// Generate nodes from the config files
	nodes := make([]*node.DarkNode, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		var err error
		var config *node.Config
		if i < numberOfBootstrapNodes {
			config, err = node.LoadConfig(fmt.Sprintf("./test-configs/bootstrap-%d.json", i+1))
		} else {
			config, err = node.LoadConfig(fmt.Sprintf("./test-configs/node-%d.json", i+1))
		}
		if err != nil {
			return nil, err
		}

		config.NetworkOptions.Timeout = 1 * time.Second
		config.NetworkOptions.TimeoutBackoff = 0 * time.Second
		config.NetworkOptions.TimeoutRetries = 1
		node, err := network.NewSwarmService(mockSwarmDelegate, config.NetworkOptions, logge)
		if err != nil {
			return nil, err
		}
		mockRegistrar.Register(node.Config.RepublicKeyPair.ID(), []byte{})
		nodes[i] = node
	}
	return nodes, nil
}

func startNodes(nodes []*node.DarkNode) {
	do.CoForAll(nodes, func(i int) {
		log.Println(nodes[i].Swarm.Address())
		nodes[i].Start()
	})
}

func watchNodes(nodes []*node.DarkNode) {
	for i := range nodes {
		go func(i int) {
			nodes[i].WatchDarkOcean()
		}(i)
	}
}

func stopNodes(nodes []*node.DarkNode) {
	for i := range nodes {
		nodes[i].Stop()
	}
}

func connectNodes(nodes []*node.DarkNode, connectivity int) error {
	for i, from := range nodes {
		for j, to := range nodes {
			if i == j {
				continue
			}
			// Connect bootstrap nodes in a fully connected topology
			if i < NumberOfBootstrapNodes {
				if j < NumberOfBootstrapNodes {
					client, err := from.ClientPool.FindOrCreateClient(to.NetworkOptions.MultiAddress)
					if err != nil {
						return err
					}
					if err := client.Ping(); err != nil {
						return err
					}
				}
				continue
			}
			// Connect standard nodes randomly
			isConnected := rand.Intn(100) < connectivity
			if isConnected {
				client, err := from.ClientPool.FindOrCreateClient(to.NetworkOptions.MultiAddress)
				if err != nil {
					return err
				}
				if err := client.Ping(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func sendOrders(nodes []*node.DarkNode) error {
	// Get order data from Binance
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v1/depth?symbol=ETHBTC&limit=%v", NumberOfOrders))
	if err != nil {
		return err
	}

	response, err := ioutil.ReadAll(resp.Body)
	orderBook := new(OrderBook)
	err = json.Unmarshal(response, orderBook)

	// Generate order from the Binance data
	buyOrders := make([]*order.Order, len(orderBook.Bids))
	sellOrders := make([]*order.Order, len(orderBook.Asks))

	for i, j := range orderBook.Asks {
		price, err := strconv.ParseFloat(j[0], 10)
		price = price * 1000000000000
		if err != nil {
			log.Fatal("fail to parse the price into a big int")
		}

		amount, err := strconv.ParseFloat(j[1], 10)
		amount = amount * 1000000000000
		if err != nil {
			log.Fatal("fail to parse the amount into a big int")
		}
		sellOrder := order.NewOrder(order.TypeLimit, order.ParitySell, time.Time{},
			order.CurrencyCodeETH, order.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
			big.NewInt(int64(amount)), big.NewInt(1))
		sellOrders[i] = sellOrder

		buyOrder := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Time{},
			order.CurrencyCodeETH, order.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
			big.NewInt(int64(amount)), big.NewInt(1))
		buyOrders[i] = buyOrder
	}

	// Send order fragment to the nodes
	totalNodes := len(nodes)
	for _, orders := range [][]*order.Order{buyOrders, sellOrders} {
		go func(orders []*order.Order) {
			for _, ord := range orders {

				if ord.Parity == order.ParityBuy {
					log.Println("sending buy order", ord.ID)
				} else {
					log.Println("sending sell order", ord.ID)
				}

				shares, err := ord.Split(int64(totalNodes), int64((totalNodes+1)*2/3), Prime)
				if err != nil {
					log.Println("cannot split the order", ord.ID)
					continue
				}

				do.ForAll(shares, func(i int) {
					client, err := rpc.NewClient(nodes[i].NetworkOptions.MultiAddress, trader)
					if err != nil {
						log.Fatal(err)
					}
					err = client.OpenOrder(&rpc.OrderSignature{}, rpc.SerializeOrderFragment(shares[i]))
					if err != nil {
						log.Printf("Coudln't send order fragment to %s\n", nodes[i].NetworkOptions.MultiAddress.Address())
						return
					}
				})
			}
		}(orders)
	}

	return nil
}
