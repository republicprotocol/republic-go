package node_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-dark-node"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
)

const (
	NumberOfBootstrapNodes = 4
	NumberOfTestNODES      = 4
)

var _ = Describe("Dark nodes", func() {
	var mu = new(sync.Mutex)
	var nodes []*node.DarkNode
	var err error

	startListening := func(nodes []*node.DarkNode, bootstrapNodes int) {
		// Start all the nodes listening for rpc calls
		for _, n := range nodes {
			go func(n *node.DarkNode) {
				defer GinkgoRecover()

				Ω(n.StartListening()).ShouldNot(HaveOccurred())
			}(n)
		}

		time.Sleep(3 * time.Second)

		// Fully connect the bootstrap nodes
		for i := 0; i < bootstrapNodes; i++ {
			for j := 0; j < bootstrapNodes; j++ {
				if i == j {
					continue
				}
				rpc.PingTarget(nodes[j].Configuration.MultiAddress, nodes[i].Configuration.MultiAddress, time.Second*5)
			}
		}
	}

	stopListening := func(nodes []*node.DarkNode) {
		for _, node := range nodes {
			node.StopListening()
		}
	}

	Context("nodes start up", func() {
		BeforeEach(func() {
			mu.Lock()
			nodes, err = generateNodes(NumberOfBootstrapNodes, NumberOfTestNODES)
			Ω(err).ShouldNot(HaveOccurred())
			startListening(nodes, NumberOfBootstrapNodes)
		})

		AfterEach(func() {
			stopListening(nodes)
			mu.Unlock()
		})

		It("should be able to run startup successfully", func() {
			for _, n := range nodes {
				go func(n *node.DarkNode) {
					Ω(n.Start()).ShouldNot(HaveOccurred())
				}(n)
			}

			time.Sleep(3 * time.Second)
			orderFileNames := []string{"./test_orders/btc-eth.json", "./test_orders/eth-btc.json"}
			for i := range orderFileNames {
				order, err := readOrderFromFile(orderFileNames[i])
				Ω(err).ShouldNot(HaveOccurred())
				shares, err := order.Split(NumberOfBootstrapNodes+NumberOfTestNODES, 5, node.Prime)
				Ω(err).ShouldNot(HaveOccurred())
				for i := range shares {
					if err := rpc.SendOrderFragmentToTarget(nodes[i].Configuration.MultiAddress, nodes[i].Configuration.MultiAddress.Address(), nodes[0].Configuration.MultiAddress, shares[i], 5*time.Second); err != nil {
						log.Fatal(err)
					}
				}
			}
			time.Sleep(1 * time.Minute)
		})
	})
})

func generateNodes(numberOfBootsrapNodes, numberOfTestNodes int) ([]*node.DarkNode, error) {
	// Generate nodes from the config files
	numberOfNodes := numberOfBootsrapNodes + numberOfTestNodes
	nodes := make([]*node.DarkNode, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		config, err := node.LoadConfig(fmt.Sprintf("./test_configs/config-%d.json", i))
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

func readOrderFromFile(orderFile string) (*compute.Order, error) {
	file, err := os.Open(orderFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	order := new(compute.Order)
	if err := json.NewDecoder(file).Decode(order); err != nil {
		return nil, err
	}

	rand.Seed(int64(time.Now().Nanosecond()))
	order.Nonce = big.NewInt(rand.Int63())
	order.GenerateID()
	return order, nil
}
