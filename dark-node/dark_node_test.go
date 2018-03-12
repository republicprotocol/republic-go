package node_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

const (
	NumberOfBootstrapNodes = 5
	NumberOfOrders         = 20
)

var Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
var trader, _ = identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
var mockRegistrar = dnr.NewMockDarkNodeRegistrar([][]byte{})

type OrderBook struct {
	LastUpdateId int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

var _ = Describe("Dark nodes", func() {
	var err error
	var mu = new(sync.Mutex)
	var bootstrapNodes, testNodes []*node.DarkNode

	for _, numberOfTestNodes := range []int{4} {
		for _, connectivity := range []int{100} {
			func(numberOfTestNodes, connectivity int) {
				Context("integration test", func() {
					BeforeEach(func() {
						mu.Lock()

						bootstrapNodes, err = generateBootstrapNodes(NumberOfBootstrapNodes)
						Ω(err).ShouldNot(HaveOccurred())
						testNodes, err = generateNodes(numberOfTestNodes)
						Ω(err).ShouldNot(HaveOccurred())
						go func() {
							defer GinkgoRecover()

							startNodes(bootstrapNodes, testNodes)
						}()
						time.Sleep(20 * time.Second)

						err = connectNodes(bootstrapNodes, testNodes, connectivity)
						Ω(err).ShouldNot(HaveOccurred())
					})

					AfterEach(func() {
						stopNodes(bootstrapNodes, testNodes)
						mu.Unlock()
					})

					It("should reach consensus on an order match", func() {
						sendOrders(bootstrapNodes, testNodes, NumberOfOrders)
					})
				})
			}(numberOfTestNodes, connectivity)
		}
	}
})

func generateBootstrapNodes(numberOfBootstrapNodes int) ([]*node.DarkNode, error) {
	// Generate nodes from the config files
	nodes := make([]*node.DarkNode, numberOfBootstrapNodes)
	for i := 0; i < numberOfBootstrapNodes; i++ {
		config, err := node.LoadConfig(fmt.Sprintf("./test-configs/bootstrap-%d.json", i+1))
		if err != nil {
			return nil, err
		}
		node, err := node.NewDarkNode(*config, mockRegistrar)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	// Register all nodes to the registrar
	for _, node := range nodes {
		mockRegistrar.Register(node.Config.RepublicKeyPair.ID(), []byte{})
	}
	return nodes, nil
}

func generateNodes(numberOfTestNodes int) ([]*node.DarkNode, error) {
	// Generate nodes from the config files
	nodes := make([]*node.DarkNode, numberOfTestNodes)
	for i := 0; i < numberOfTestNodes; i++ {
		config, err := node.LoadConfig(fmt.Sprintf("./test-configs/node-%d.json", i+1))
		if err != nil {
			return nil, err
		}
		node, err := node.NewDarkNode(*config, mockRegistrar)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	// Register all nodes to the registrar
	for _, node := range nodes {
		mockRegistrar.Register(node.Config.RepublicKeyPair.ID(), []byte{})
	}
	return nodes, nil
}

func startNodes(bootstrapNodes, testNodes []*node.DarkNode) {
	for i := range bootstrapNodes {
		go func(i int) {
			bootstrapNodes[i].Start()
		}(i)
	}
	time.Sleep(5 * time.Second)
	for i := range testNodes {
		go func(i int) {
			testNodes[i].Start()
		}(i)
	}
}

func stopNodes(bootstrapNodes, testNodes []*node.DarkNode) {
	for i := range bootstrapNodes {
		bootstrapNodes[i].Stop()
	}
	for i := range testNodes {
		testNodes[i].Stop()
	}
}

func connectNodes(bootstrapNodes, testNodes []*node.DarkNode, connectivity int) error {
	// Connect bootstrap nodes in a fully-connected terminology
	for i, bootstrapNode := range bootstrapNodes {
		for j, other := range bootstrapNodes {
			if i == j {
				continue
			}

			client, err := bootstrapNode.ClientPool.FindOrCreateClient(other.NetworkOptions.MultiAddress)
			if err != nil {
				return err
			}
			err = client.Ping()
			if err != nil {
				return err
			}
		}
	}

	// Generate connectivity map
	connectivityMap := map[string]map[string]bool{}
	for i, testNode := range testNodes {
		for j := i + 1; j < len(testNodes); j++ {
			other := testNodes[j]
			isConnected := rand.Intn(100) < connectivity
			if _, ok := connectivityMap[testNode.NetworkOptions.MultiAddress.String()]; !ok {
				connectivityMap[testNode.NetworkOptions.MultiAddress.String()] = map[string]bool{}
			}
			if _, ok := connectivityMap[other.NetworkOptions.MultiAddress.String()]; !ok {
				connectivityMap[other.NetworkOptions.MultiAddress.String()] = map[string]bool{}
			}
			connectivityMap[testNode.NetworkOptions.MultiAddress.String()][other.NetworkOptions.MultiAddress.String()] = isConnected
			connectivityMap[other.NetworkOptions.MultiAddress.String()][testNode.NetworkOptions.MultiAddress.String()] = isConnected
		}
	}

	// Connect test nodes depending on the connectivity
	for i, testNode := range testNodes {
		for j, other := range testNodes {
			if i == j {
				continue
			}
			if !connectivityMap[testNode.NetworkOptions.MultiAddress.String()][other.NetworkOptions.MultiAddress.String()] {
				continue
			}
			client, err := testNode.ClientPool.FindOrCreateClient(other.NetworkOptions.MultiAddress)
			if err != nil {
				return err
			}
			err = client.Ping()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sendOrders(bootstrapNodes, testNodes []*node.DarkNode, numberOfOrders int) {
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v1/depth?symbol=ETHBTC&limit=%v", numberOfOrders))
	if err != nil {
		log.Fatal("fail to get data from binance")
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
	nodes := append(bootstrapNodes, testNodes...)
	totalNodes := len(nodes)
	for _, orders := range [][]*order.Order{buyOrders, sellOrders} {
		go func(orders []*order.Order) {
			for _, ord := range orders {

				if ord.Parity == order.ParityBuy {
					log.Println("sending buy order", ord.ID)
				} else {
					log.Println("sending sell order", ord.ID)
				}

				shares, err := ord.Split(int64(totalNodes), int64((totalNodes*2)/3), Prime)
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
}
