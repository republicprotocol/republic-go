package node_test

import (
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-dark-node"
	"github.com/republicprotocol/go-rpc"
)

const (
	Number_Of_Bootstrap_Nodes = 4
	Number_Of_Test_NODES      = 4
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
			nodes, err = generateNodes(Number_Of_Bootstrap_Nodes, Number_Of_Test_NODES)
			Ω(err).ShouldNot(HaveOccurred())
			startListening(nodes, Number_Of_Bootstrap_Nodes)
		})

		AfterEach(func() {
			stopListening(nodes)
			mu.Unlock()
		})

		It("should be able to run startup successfully", func() {
			for _, node := range nodes {
				Ω(node.Start()).ShouldNot(HaveOccurred())
			}
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
