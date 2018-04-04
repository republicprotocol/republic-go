package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

var Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

const key = `{"version":3,"id":"7844982f-abe7-4690-8c15-34f75f847c66","address":"db205ea9d35d8c01652263d58351af75cfbcbf07","Crypto":{"ciphertext":"378dce3c1279b36b071e1c7e2540ac1271581bff0bbe36b94f919cb73c491d3a","cipherparams":{"iv":"2eb92da55cc2aa62b7ffddba891f5d35"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"80d3341678f83a14024ba9c3edab072e6bd2eea6aa0fbc9e0a33bae27ffa3d6d","n":8192,"r":8,"p":1},"mac":"3d07502ea6cd6b96a508138d8b8cd2e46c3966240ff276ce288059ba4235cb0d"}}`

const reset = "\x1b[0m"
const red = "\x1b[31;1m"

type OrderBook struct {
	LastUpdateId int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

func main() {
	// Parse the option parameters
	numberOfOrders := flag.Int("order", 10, "number of orders")
	timeInterval := flag.Int("time", 15, "time interval in second")
	cancelRequest := flag.Bool("cancel", false, "cancel request")

	// Logger
	logs, err := logger.NewLogger(logger.Options{
		Plugins: []logger.PluginOptions{
			logger.PluginOptions{
				File: &logger.FilePluginOptions{
					Path: "stdout",
				},
			},
		}})

	// Create a test trader address
	address, _, err := identity.NewAddress()
	if err != nil {
		panic(err)
	}
	multi, err := identity.NewMultiAddressFromString("/ip4/0.0.0.0/tcp/80/republic/" + address.String())
	if err != nil {
		panic(err)
	}
	log.Println("Trader Address: ", address)

	// Get the Dark Node Registry
	registrar, err := getDarkNodeRegistrar()
	if err != nil {
		panic(err)
	}

	// Get the Dark Ocean
	ocean, err := dark.NewOcean(logs, 5, registrar)

	// Get the Dark Pools
	pools := ocean.GetPools()

	if *cancelRequest {
		cancelOrder(pools, multi)
	} else {
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
				order := order.NewOrder(order.TypeLimit, order.ParitySell, time.Time{},
					order.CurrencyCodeETH, order.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
					big.NewInt(int64(amount)), big.NewInt(1))
				sellOrders[i] = order
			}

			buyOrders := make([]*order.Order, len(orderBook.Bids))
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

				order := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Time{},
					order.CurrencyCodeETH, order.CurrencyCodeBTC, big.NewInt(int64(price)), big.NewInt(int64(amount)),
					big.NewInt(int64(amount)), big.NewInt(1))
				buyOrders[i] = order
			}

			// For all orders generated
			for _, orders := range [][]*order.Order{buyOrders, sellOrders} {
				go func(orders []*order.Order) {
					for _, ord := range orders {
						//todo
						if ord.Parity == order.ParityBuy {
							log.Println("sending buy order :", base58.Encode(ord.ID))
						} else {
							log.Println("sending sell order :", base58.Encode(ord.ID))
						}

						// For every Dark Pool
						for i := range pools {

							// Split order into (Number of nodes in each pool) * 2/3 fragments
							shares, err := ord.Split(int64(pools[i].Size()), int64(pools[i].Size()*2/3), Prime)
							if err != nil {
								continue
							}

							// Send the shares across all nodes within the Dark Pool
							j := 0
							pools[i].ForAll(func(n *dark.Node) {

								// Get multiaddress
								multiaddress, err := getMultiAddress(n.ID.Address(), multi)

								// Create a client
								client, err := rpc.NewClient(multiaddress, multi)
								if err != nil {
									log.Fatal(err)
								}

								// Send fragment to node
								err = client.OpenOrder(&rpc.OrderSignature{}, rpc.SerializeOrderFragment(shares[j]))
								if err != nil {
									log.Println(err)
									log.Printf("%sCoudln't send order fragment to %v%s\n", red, base58.Encode(n.ID), reset)
									return
								}	
								j++
							})
						}
					}
				}(orders)
			}
			time.Sleep(time.Duration(*timeInterval) * time.Second)
		}
	}
}

func cancelOrder(pools dark.Pools, traderAddress identity.MultiAddress) {
	// For every Dark Pool
	for i := range pools {

		// Cancel orders for all nodes in the pool
		pools[i].ForAll(func(n *dark.Node) {

			// Get multiaddress
			multiaddress, err := getMultiAddress(n.ID.Address(), traderAddress)

			// Create a client
			client, err := rpc.NewClient(multiaddress, traderAddress)
			if err != nil {
				log.Fatal(err)
			}

			// Close order
			err = client.CancelOrder(&rpc.OrderSignature{})
			if err != nil {
				log.Println(err)
				log.Printf("%sCoudln't cancel order to %v%s\n", red, base58.Encode(n.ID), reset)
				return
			}
		})
	}
}

// Function to obtain multiaddress of a node by sending requests to bootstrap nodes
func getMultiAddress(address identity.Address, traderMultiAddress identity.MultiAddress) (identity.MultiAddress, error) {
	BootstrapMultiAddresses := []string{
		"/ip4/52.77.88.84/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
		"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
		"/ip4/52.59.176.141/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
		"/ip4/52.21.44.236/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
		"/ip4/52.41.118.171/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
	}

	serializedTarget := rpc.SerializeAddress(address)
	for _, peer := range BootstrapMultiAddresses {

		bootStrapMultiAddress, err := identity.NewMultiAddressFromString(peer)
		if err != nil {
			log.Fatal(err)
		}

		client, err := rpc.NewClient(bootStrapMultiAddress, traderMultiAddress)
		if err != nil {
			log.Fatal(err)
		}

		candidates, err := client.QueryPeersDeep(serializedTarget)
		if err != nil {
			return traderMultiAddress, err
		}
		
		for candidate := range candidates {
			deserializedCandidate, err := rpc.DeserializeMultiAddress(candidate)
			if err != nil {
				return traderMultiAddress, err
			}
			if address == deserializedCandidate.Address() {
				fmt.Println("Found the target : ", address)
				return deserializedCandidate, nil
			}
		}
	}
	return traderMultiAddress, nil
}

// Get Dark Node Registry
func getDarkNodeRegistrar () (dnr.DarkNodeRegistrar, error) {

	clientDetails, err := connection.FromURI("https://ropsten.infura.io/", "ropsten")
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewTransactor(strings.NewReader(key), "password1")
	if err != nil {
		panic(err)
	}

	// Gas Price
	auth.GasPrice = big.NewInt(6000000000)

	return dnr.NewEthereumDarkNodeRegistrar(context.Background(), &clientDetails, auth, &bind.CallOpts{})
}