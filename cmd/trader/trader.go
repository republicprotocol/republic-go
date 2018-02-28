package main

import (
	"encoding/json"
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

type OrderBook struct {
	LastUpdateId int        `lastUpdateId`
	Bids         [][]string `bids`
	Asks         [][]string `asks`
}

func main() {

	// Nodes get nodes/darkpool details
	multiAddress := []string{
		"/ip4/13.125.159.239/tcp/18514/republic/8MKZ8JwCU9m9affPWHZ9rxp2azXNnE",
		"/ip4/13.229.60.122/tcp/18514/republic/8MHarRJdvWd7SsTJE8vRVfj2jb5cWS",
		"/ip4/54.93.234.49/tcp/18514/republic/8MKDGUTgKtkymyKTH28xeMxiCnJ9xy",
		"/ip4/54.89.239.234/tcp/18514/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto",
		"/ip4/35.161.248.181/tcp/18514/republic/8MJ38m8Nzknh3gVj7QiMjuejmHBMSf",
	}
	nodes := make([]identity.MultiAddress, 5)
	for i := 0; i < 5; i++ {
		multi, err := identity.NewMultiAddressFromString(multiAddress[i])
		if err != nil {
			log.Fatal(err)
		}
		nodes[i] = multi
	}

	// Create a trader address
	config, err := node.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Trader Address: ", config.MultiAddress.String())

	for {
		// Get orders details from Binance
		resp, err := http.Get("https://api.binance.com/api/v1/depth?symbol=ETHBTC&limit=5")
		if err != nil {
			log.Fatal("fail to get data from binance")
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
				log.Fatal("fail to parse the price into a big int")
			}

			amount, err := strconv.ParseFloat(j[1], 10)
			amount = amount * 1000000000000
			if err != nil {
				log.Fatal("fail to parse the amount into a big int")
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
				log.Fatal("fail to parse the price into a big int")
			}

			amount, err := strconv.ParseFloat(j[1], 10)
			amount = amount * 1000000000000
			if err != nil {
				log.Fatal("fail to parse the amount into a big int")
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

					shares, err := order.Split(5, 3, node.Prime)
					if err != nil {
						continue
					}
					do.ForAll(shares, func(i int) {
						rpc.SendOrderFragmentToTarget(nodes[i], nodes[i].Address(), config.MultiAddress, shares[i], 5*time.Second)
					})
				}
			}(orders)
		}
		time.Sleep(15 * time.Second)
	}
}
