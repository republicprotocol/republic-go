package main

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

type OrderBook struct {
	LastUpdateId int        `lastUpdateId`
	Bids         [][]string `bids`
	Asks         [][]string `asks`
}

func main() {
	// Parse the option parameters
	numberOfOrders := flag.Int("order", 10, "number of orders")
	timeInterval := flag.Int("time", 10, "time interval in second")

	// Get nodes/darkpool details
	multiAddress := getNodesDetails()
	nodes := make([]identity.MultiAddress, len(multiAddress))
	for i := 0; i < len(multiAddress); i++ {
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

	// Keep sending order fragment
	for {
		// Get orders details from Binance
		resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v1/depth?symbol=ETHBTC&limit=%v", *numberOfOrders))
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
		time.Sleep(time.Duration(*timeInterval) * time.Second)
	}
}

func getNodesDetails() []string {
	// AWS nodes
	//return []string{
	//	"/ip4/13.125.159.239/tcp/18514/republic/8MKZ8JwCU9m9affPWHZ9rxp2azXNnE",
	//	"/ip4/13.229.60.122/tcp/18514/republic/8MHarRJdvWd7SsTJE8vRVfj2jb5cWS",
	//	"/ip4/54.93.234.49/tcp/18514/republic/8MKDGUTgKtkymyKTH28xeMxiCnJ9xy",
	//	"/ip4/54.89.239.234/tcp/18514/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto",
	//	"/ip4/35.161.248.181/tcp/18514/republic/8MJ38m8Nzknh3gVj7QiMjuejmHBMSf",
	//}

	// susruth's test nodes
	return []string{
		"/ip4/52.21.44.236/tcp/18514/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto",
		"/ip4/52.41.118.171/tcp/18514/republic/8MJ38m8Nzknh3gVj7QiMjuejmHBMSf",
		"/ip4/52.59.176.141/tcp/18514/republic/8MKDGUTgKtkymyKTH28xeMxiCnJ9xy",
		"/ip4/52.77.88.84/tcp/18514/republic/8MHarRJdvWd7SsTJE8vRVfj2jb5cWS",
		"/ip4/52.79.194.108/tcp/18514/republic/8MKZ8JwCU9m9affPWHZ9rxp2azXNnE",
	}

	// Local nodes
	//return []string{
	//	"/ip4/127.0.0.1/tcp/4000/republic/8MGyTXr6poqfizzdp9fWcLd8UpDC5y",
	//	"/ip4/127.0.0.1/tcp/4001/republic/8MJWTpvNJv2SW7meGFpa8c44zNw63f",
	//	"/ip4/127.0.0.1/tcp/4002/republic/8MJKxAujyofThVCwmYfMnpCRCNwnQe",
	//	"/ip4/127.0.0.1/tcp/4003/republic/8MHJYuWArPDwA8VvXzxrzPyEdwrb4s",
	//	"/ip4/127.0.0.1/tcp/4004/republic/8MGonGTeJ6Kz2gUFYgKpJy4TPCH6q3",
	//	"/ip4/127.0.0.1/tcp/4005/republic/8MGYgkWK26wS4U4EqHcdgRaodhR9AS",
	//	"/ip4/127.0.0.1/tcp/4006/republic/8MKYtPMdocyv2DeKQkobciKsPHMTwC",
	//	"/ip4/127.0.0.1/tcp/4007/republic/8MJSSou1rUmjxrvKcYTxaZms8zSvvD",
	//}
}
