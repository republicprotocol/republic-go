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

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

var Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

const reset = "\x1b[0m"
const yellow = "\x1b[33;1m"
const green = "\x1b[32;1m"
const red = "\x1b[31;1m"

type OrderBook struct {
	LastUpdateId int        `json:"lastUpdateId"`
	Bids         [][]interface{} `json:"bids"`
	Asks         [][]interface{} `json:"asks"`
}

func main() {
	// Parse the option parameters
	numberOfOrders := flag.Int("order", 5, "number of orders")
	timeInterval := flag.Int("time", 30, "time interval in second")

	// Get nodes/darkPool details
	multiAddresses := getNodesDetails()
	nodes := make([]identity.MultiAddress, len(multiAddresses))
	for i := 0; i < len(multiAddresses); i++ {
		multi, err := identity.NewMultiAddressFromString(multiAddresses[i])
		if err != nil {
			log.Fatal(err)
		}
		nodes[i] = multi
	}

	// Create a trader address
	address, _, err := identity.NewAddress()
	if err != nil {
		log.Fatal(err)
	}
	multi, err := identity.NewMultiAddressFromString("/ip4/0.0.0.0/tcp/80/republic/" + address.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Trader Address: ", address)

	// Keep sending order fragment
	for {
		// Get orders details from Binance
		resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v1/depth?symbol=ETHBTC&limit=%v", *numberOfOrders))
		if err != nil {
			log.Fatal("fail to get data from binance")
		}

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		orderBook := new(OrderBook)
		err = json.Unmarshal(response, orderBook)
		if err != nil {
			log.Fatal(err)
		}

		// Generate order from the Binance data
		buyOrders := make([]*order.Order, len(orderBook.Asks))
		sellOrders := make([]*order.Order, len(orderBook.Asks))

		for i, j := range orderBook.Asks {
			price, err := strconv.ParseFloat(j[0].(string), 10)
			price = price * 1000000000000
			if err != nil {
				log.Fatal(err)
			}

			amount, err := strconv.ParseFloat(j[1].(string), 10)
			amount = amount * 1000000000000
			if err != nil {
				log.Fatal(err)
			}
			sellOrder := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour),
				order.CurrencyCodeETH, order.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
				big.NewInt(int64(amount)), big.NewInt(1))
			sellOrders[i] = sellOrder

			buyOrder := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour),
				order.CurrencyCodeETH, order.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
				big.NewInt(int64(amount)), big.NewInt(1))
			buyOrders[i] = buyOrder
		}

		// Send order fragment to the nodes
		totalNodes := len(multiAddresses)
		pool := rpc.NewClientPool(multi)
		for i := range buyOrders {
			buyOrder, sellOrder := buyOrders[i], sellOrders[i]
			log.Printf("Sending matched order. [BUY] %s <---> [SELL] %s", buyOrder.ID, sellOrder.ID)

			buyShares, err := buyOrder.Split(int64(totalNodes), int64(totalNodes*2/3+1 ), Prime)
			if err != nil {
			}
			sellShares, err := sellOrder.Split(int64(totalNodes), int64(totalNodes*2/3 +1), Prime)
			if err != nil {
				log.Println(err)
				continue
			}

			do.CoForAll(buyShares, func(j int) {
				err := pool.OpenOrder(nodes[j], &rpc.OrderSignature{}, rpc.SerializeOrderFragment(buyShares[j]))
				if err != nil {
					log.Printf("Coudln't send order fragment to %s\n %s", nodes[j].ID(), err )
				}
			})
			log.Println("finish sending buy order",  buyOrder.ID)

			do.CoForAll(sellShares, func(j int) {
				err := pool.OpenOrder(nodes[j], &rpc.OrderSignature{}, rpc.SerializeOrderFragment(sellShares[j]))
				if err != nil {
					log.Printf("Coudln't send order fragment to %s\n %s", nodes[j].ID() ,err)
				}
			})

			log.Println("finish sending sell order",  sellOrder.ID)

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

	// bootstrap nodes
	return []string{
		"/ip4/52.77.88.84/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
		//"/ip4/52.59.176.141/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
		//"/ip4/52.21.44.236/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
		//"/ip4/52.41.118.171/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
		"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
		"/ip4/13.250.34.9/tcp/18514/republic/8MH9hcbekxW8yUo9ADBhA213PnZ4do",
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
