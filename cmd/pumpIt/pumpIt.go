package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/relay"
	"github.com/republicprotocol/republic-go/stackint"
)

var numberOfOrders = 10

func main() {
	m := stackint.FromUint(1000000000)
	n, err := stackint.Random(rand.Reader, &m)
	if err != nil {
		log.Fatal(err)
	}

	// Generate orders
	ethToRen := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeREN, stackint.FromUint(1), stackint.FromUint(1), stackint.FromUint(1), n)
	renToEth := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeREN, stackint.FromUint(1), stackint.FromUint(1), stackint.FromUint(1), n)
	log.Println("eth to ren order id ", ethToRen.ID.String())
	log.Println("ren to eth order id ", renToEth.ID.String())

	// Send sell order
	sellReq := relay.OpenOrderRequest{Order: *ethToRen, OrderFragments: relay.OrderFragments{}}
	bufferSell := new(bytes.Buffer)
	if err := json.NewEncoder(bufferSell).Encode(sellReq); err != nil {
		log.Fatal(" fail to marshal order", err)
	}
	resp, err := http.Post("http://localhost:18516/orders", "application/json", bufferSell)
	if err != nil {
		log.Fatal("response fail,", err)
	}
	// Read the response status

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("cannot readall: " + err.Error())
		}
		log.Println(string(b))
		log.Fatalf("reponse fail with status code %v", resp.StatusCode)
	}

	// Send buy order
	buyReq := relay.OpenOrderRequest{Order: *renToEth, OrderFragments: relay.OrderFragments{}}

	bufferBuy := new(bytes.Buffer)
	if err := json.NewEncoder(bufferBuy).Encode(buyReq); err != nil {
		log.Fatal(" fail to marshal order", err)
	}

	resp, err = http.Post("http://localhost:18518/orders", "application/json", bufferBuy)
	if err != nil {
		log.Fatal("response fail ", err)
	}
	// Read the response status

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		r := map[string]interface{}{}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			log.Fatal(err)
		}
		log.Println(r)
		log.Fatalf("reponse fail with status code %v", resp.StatusCode)
	}

	//for {
	//	// Generate buy-sell order pairs
	//	buyOrders, sellOrders := make([]*order.Order, numberOfOrders), make([]*order.Order, numberOfOrders)
	//	for i := 0; i < numberOfOrders; i++ {
	//		price := i * 1000000000000
	//		amount := i * 1000000000000
	//
	//		nonce, err := stackint.Random(rand.Reader, &smpc.Prime)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		sellOrder := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour),
	//			order.CurrencyCodeETH, order.CurrencyCodeBTC, stackint.FromUint(uint(price)), stackint.FromUint(uint(amount)),
	//			stackint.FromUint(uint(amount)), nonce)
	//		sellOrders[i] = sellOrder
	//
	//		buyOrder := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour),
	//			order.CurrencyCodeETH, order.CurrencyCodeBTC, stackint.FromUint(uint(price)), stackint.FromUint(uint(amount)),
	//			stackint.FromUint(uint(amount)), nonce)
	//		buyOrders[i] = buyOrder
	//	}
	//
	//	// Send order fragment to the nodes
	//	nodes := Multiaddresses()
	//	totalNodes := len(nodes)
	//	traderAddr, _, err := identity.NewAddress()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	trader, _ := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%v", traderAddr))
	//	prime, _ := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	//
	//	crypter := crypto.NewWeakCrypter()
	//	connPool := client.NewConnPool(256)
	//	defer connPool.Close()
	//	smpcerClient := smpcer.NewClient(&crypter, trader, &connPool)
	//
	//	for i := 0; i < numberOfOrders; i++ {
	//		buyOrder, sellOrder := buyOrders[i], sellOrders[i]
	//		log.Printf("sending buy/sell pair (%s, %s)", buyOrder.ID, sellOrder.ID)
	//		buyShares, err := buyOrder.Split(int64(totalNodes), int64((totalNodes+1)*2/3), &prime)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		sellShares, err := sellOrder.Split(int64(totalNodes), int64((totalNodes+1)*2/3), &prime)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		for _, shares := range [][]*order.Fragment{buyShares, sellShares} {
	//			do.CoForAll(shares, func(j int) {
	//				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//				defer cancel()
	//
	//				if err := smpcerClient.OpenOrder(ctx, nodes[j], *shares[j]); err != nil {
	//					log.Printf("cannot send order fragment to %s: %v", nodes[j].Address(), err)
	//				}
	//			})
	//		}
	//	}
	//
	//	time.Sleep(5 * time.Second)
	//}
}

func Multiaddresses() []identity.MultiAddress {
	multiaddresses := []string{
		"/ip4/0.0.0.0/tcp/80/republic/123",
	}
	multis := make([]identity.MultiAddress, len(multiaddresses))

	for i, address := range multiaddresses {
		multi, err := identity.NewMultiAddressFromString(address)
		if err != nil {
			log.Fatal(err)
		}
		multis[i] = multi
	}

	return multis
}
