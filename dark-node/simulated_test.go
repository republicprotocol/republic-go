package node_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
	"github.com/republicprotocol/republic-go/network/dht"
	"github.com/republicprotocol/republic-go/order"
	"google.golang.org/grpc"
)

// NewDarkNode return a DarkNode that adheres to the given Config. The DarkNode
// will configure all of the components that it needs to operate but will not
// start any of them.
func NewTestDarkNode(config darknode.Config) *DarkNode {
	if config.Prime == nil {
		config.Prime = Prime
	}

	// TODO: This should come from the DNR.
	k := int64(5)

	node := &DarkNode{Config: config}

	node.Logger = logger.NewLogger()
	node.ClientPool = rpc.NewClientPool(node.MultiAddress)
	node.DHT = dht.NewDHT(node.MultiAddress.Address(), node.MaxBucketLength)

	node.DeltaBuilder = compute.NewDeltaBuilder(k, node.Prime)
	node.DeltaFragmentMatrix = compute.NewDeltaFragmentMatrix(node.Prime)
	node.OrderFragmentWorkerQueue = make(chan *order.Fragment, 100)
	node.OrderFragmentWorker = NewOrderFragmentWorker(node.OrderFragmentWorkerQueue, node.DeltaFragmentMatrix)
	node.DeltaFragmentWorkerQueue = make(chan *compute.DeltaFragment, 100)
	node.DeltaFragmentWorker = NewDeltaFragmentWorker(node.DeltaFragmentWorkerQueue, node.DeltaBuilder)

	// options := network.Options{}
	node.Server = grpc.NewServer(grpc.ConnectionTimeout(time.Minute))
	node.Swarm = network.NewSwarmService(node, node.Options, node.Logger, node.ClientPool, node.DHT)
	node.Dark = network.NewDarkService(node, node.Options, node.Logger)

	registrar, err := ConnectToRegistrar(clientDetails, config)
	if err != nil {
		panic(err)
	}
	node.Registrar = registrar

	return node
}

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
			// ...
		})
	})
})
