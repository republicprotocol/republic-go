package network_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-network/rpc"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-sss"
)

type mockDelegate struct {
	mu                      *sync.Mutex
	numberOfPings           int
	orderFragmentsReceived  map[string]struct{}
	resultFragmentsReceived map[string]struct{}
}

func newMockDelegate() *mockDelegate {
	return &mockDelegate{
		mu:                      new(sync.Mutex),
		numberOfPings:           0,
		orderFragmentsReceived:  map[string]struct{}{},
		resultFragmentsReceived: map[string]struct{}{},
	}
}

func (delegate *mockDelegate) OnPingReceived(_ identity.MultiAddress) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfPings++
}

func (delegate *mockDelegate) OnOrderFragmentReceived(_ identity.MultiAddress, orderFragment *compute.OrderFragment) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.orderFragmentsReceived[string(orderFragment.ID)] = struct{}{}
}

func (delegate *mockDelegate) OnResultFragmentReceived(_ identity.MultiAddress, resultFragment *compute.ResultFragment) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.resultFragmentsReceived[string(resultFragment.ID)] = struct{}{}
}

var _ = Describe("Pinging", func() {

	run := func(name string, numberOfNodes int) int {
		var nodes []*network.Node
		var topology map[identity.Address][]*network.Node
		var err error

		delegate := newMockDelegate()
		switch name {
		case "full":
			nodes, topology, err = generateFullyConnectedTopology(numberOfNodes, delegate)
		case "star":
			nodes, topology, err = generateStarTopology(numberOfNodes, delegate)
		case "line":
			nodes, topology, err = generateLineTopology(numberOfNodes, delegate)
		case "ring":
			nodes, topology, err = generateRingTopology(numberOfNodes, delegate)
		}
		Ω(err).ShouldNot(HaveOccurred())

		for _, node := range nodes {
			go func(node *network.Node) {
				defer GinkgoRecover()
				Ω(node.Serve()).ShouldNot(HaveOccurred())
			}(node)
			defer func(node *network.Node) {
				defer GinkgoRecover()
				node.Stop()
			}(node)
		}
		time.Sleep(time.Second)

		err = ping(nodes, topology)
		Ω(err).ShouldNot(HaveOccurred())

		return delegate.numberOfPings
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("in a fully connected topology with %d nodes", numberOfNodes), func() {
				It("should update the DHT", func() {
					testMu.Lock()
					defer testMu.Unlock()
					numberOfPings := run("full", numberOfNodes)
					Ω(numberOfPings).Should(Equal(numberOfNodes * (numberOfNodes - 1)))
				})
			})
		}(numberOfNodes)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("in a star topology with %d nodes", numberOfNodes), func() {
				It("should update the DHT", func() {
					testMu.Lock()
					defer testMu.Unlock()
					numberOfPings := run("star", numberOfNodes)
					Ω(numberOfPings).Should(Equal(2 * (numberOfNodes - 1)))
				})
			})
		}(numberOfNodes)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("in a line topology with %d nodes", numberOfNodes), func() {
				It("should update the DHT", func() {
					testMu.Lock()
					defer testMu.Unlock()
					numberOfPings := run("line", numberOfNodes)
					Ω(numberOfPings).Should(Equal(2 * (numberOfNodes - 1)))
				})
			})
		}(numberOfNodes)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("in a ring topology with %d nodes", numberOfNodes), func() {
				It("should update the DHT", func() {
					testMu.Lock()
					defer testMu.Unlock()
					numberOfPings := run("ring", numberOfNodes)
					Ω(numberOfPings).Should(Equal(2 * numberOfNodes))
				})
			})
		}(numberOfNodes)
	}
})

var _ = Describe("Sending order fragments", func() {

	send := func(nodes []*network.Node, numberOfFragments int) {
		var wg sync.WaitGroup
		wg.Add(numberOfFragments)
		for i := 0; i < numberOfFragments; i++ {
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				to, from := randomNodes(nodes)
				orderFragment := generateOrderFragment(to.MultiAddress().Address().String())
				orderFragment.To = &rpc.Address{Address: to.MultiAddress().Address().String()}
				orderFragment.From = nil

				_, err := network.SendOrderFragmentToTarget(network.SerializeMultiAddress(from.MultiAddress()), orderFragment)
				Ω(err).ShouldNot(HaveOccurred())
			}()
		}
		wg.Wait()
	}

	run := func(name string, numberOfNodes, numberOfFragment int) int {
		var nodes []*network.Node
		var topology map[identity.Address][]*network.Node
		var err error

		delegate := newMockDelegate()
		switch name {
		case "full":
			nodes, topology, err = generateFullyConnectedTopology(numberOfNodes, delegate)
		case "star":
			nodes, topology, err = generateStarTopology(numberOfNodes, delegate)
		case "line":
			nodes, topology, err = generateLineTopology(numberOfNodes, delegate)
		case "ring":
			nodes, topology, err = generateRingTopology(numberOfNodes, delegate)
		}
		Ω(err).ShouldNot(HaveOccurred())

		for _, node := range nodes {
			go func(node *network.Node) {
				defer GinkgoRecover()
				Ω(node.Serve()).ShouldNot(HaveOccurred())
			}(node)
			defer func(node *network.Node) {
				defer GinkgoRecover()
				node.Stop()
			}(node)
		}
		time.Sleep(time.Second)

		err = ping(nodes, topology)
		Ω(err).ShouldNot(HaveOccurred())
		send(nodes, numberOfFragment)

		return len(delegate.orderFragmentsReceived)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("in a fully connected topology with %d nodes", numberOfNodes), func() {
				It("should send the order fragments to the right target", func() {
					numberOfMessages := numberOfNodes
					testMu.Lock()
					defer testMu.Unlock()
					numberOfFragments := run("full", numberOfNodes, numberOfMessages)
					Ω(numberOfFragments).Should(Equal(numberOfMessages))
				})
			})
		}(numberOfNodes)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("in a star topology with %d nodes", numberOfNodes), func() {
				It("should send the order fragment to the right target", func() {
					numberOfMessages := numberOfNodes
					testMu.Lock()
					defer testMu.Unlock()
					numberOfFragments := run("star", numberOfNodes, numberOfMessages)
					Ω(numberOfFragments).Should(Equal(numberOfMessages))
				})
			})
		}(numberOfNodes)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("in a line topology with %d nodes", numberOfNodes), func() {
				It("should send the order fragment to the right target", func() {
					numberOfMessages := numberOfNodes
					testMu.Lock()
					defer testMu.Unlock()
					numberOfFragments := run("line", numberOfNodes, numberOfMessages)
					Ω(numberOfFragments).Should(Equal(numberOfMessages))
				})
			})
		}(numberOfNodes)
	}

	for _, numberOfNodes := range []int{10, 20, 40, 80} {
		func(numberOfNodes int) {
			Context(fmt.Sprintf("in a ring topology with %d nodes", numberOfNodes), func() {
				It("should send the order fragment to the right target", func() {
					numberOfMessages := numberOfNodes
					testMu.Lock()
					defer testMu.Unlock()
					numberOfFragments := run("ring", numberOfNodes, numberOfMessages)
					Ω(numberOfFragments).Should(Equal(numberOfMessages))
				})
			})
		}(numberOfNodes)
	}

})

func generateOrderFragment(to string) *rpc.OrderFragment {
	return network.SerializeOrderFragment(compute.NewOrderFragment(
		crypto.Keccak256([]byte(fmt.Sprintf("%v", rand.Int()))),
		compute.OrderTypeIBBO,
		compute.OrderBuy,
		sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		},
		sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		},
		sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		},
		sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		},
		sss.Share{
			Key:   int64(1),
			Value: big.NewInt(1),
		},
	))
}

// var _ = Describe("Peers RPC", func() {

// 	run := func(name string, numberOfNodes int) int {
// 		var nodes []*network.Node
// 		var topology map[identity.Address][]*network.Node
// 		var err error

// 		delegate := newPingDelegate()
// 		switch name {
// 		case "full":
// 			nodes, topology, err = generateFullyConnectedTopology(numberOfNodes, delegate)
// 		case "star":
// 			nodes, topology, err = generateStarTopology(numberOfNodes, delegate)
// 		case "line":
// 			nodes, topology, err = generateLineTopology(numberOfNodes, delegate)
// 		case "ring":
// 			nodes, topology, err = generateRingTopology(numberOfNodes, delegate)
// 		}
// 		Ω(err).ShouldNot(HaveOccurred())

// 		for _, node := range nodes {
// 			go func(node *network.Node) {
// 				defer GinkgoRecover()
// 				Ω(node.Serve()).ShouldNot(HaveOccurred())
// 			}(node)
// 			defer func(node *network.Node) {
// 				defer GinkgoRecover()
// 				node.Stop()
// 			}(node)
// 		}
// 		time.Sleep(time.Second)
// 		// Ping nodes to make sure they are connected.
// 		err = ping(nodes, topology)
// 		Ω(err).ShouldNot(HaveOccurred())
// 		// Check that the nodes have the expected peers.
// 		err = peers(nodes, topology)
// 		Ω(err).ShouldNot(HaveOccurred())

// 		return int(delegate.numberOfPings)
// 	}

// 	for _, numberOfNodes := range []int{10, 20, 40, 80} {
// 		Context(fmt.Sprintf("in a fully connected topology with %d nodes", numberOfNodes), func() {
// 			It("should be connected to the peers described in the topology", func() {
// 				testMu.Lock()
// 				defer testMu.Unlock()
// 				numberOfPings := run("full", numberOfNodes)
// 				Ω(numberOfPings).Should(Equal(numberOfNodes * (numberOfNodes - 1)))
// 			})
// 		})
// 	}

// 	for _, numberOfNodes := range []int{10, 20, 40, 80} {
// 		Context(fmt.Sprintf("in a star topology with %d nodes", numberOfNodes), func() {
// 			It("should be connected to the peers described in the topology", func() {
// 				testMu.Lock()
// 				defer testMu.Unlock()
// 				numberOfPings := run("star", numberOfNodes)
// 				Ω(numberOfPings).Should(Equal(2 * (numberOfNodes - 1)))
// 			})
// 		})
// 	}

// 	for _, numberOfNodes := range []int{10, 20, 40, 80} {
// 		Context(fmt.Sprintf("in a line topology with %d nodes", numberOfNodes), func() {
// 			It("should be connected to the peers described in the topology", func() {
// 				testMu.Lock()
// 				defer testMu.Unlock()
// 				numberOfPings := run("line", numberOfNodes)
// 				Ω(numberOfPings).Should(Equal(2 * (numberOfNodes - 1)))
// 			})
// 		})
// 	}

// 	for _, numberOfNodes := range []int{10, 20, 40, 80} {
// 		Context(fmt.Sprintf("in a ring topology with %d nodes", numberOfNodes), func() {
// 			It("should be connected to the peers described in the topology", func() {
// 				testMu.Lock()
// 				defer testMu.Unlock()
// 				numberOfPings := run("ring", numberOfNodes)
// 				Ω(numberOfPings).Should(Equal(2 * numberOfNodes))
// 			})
// 		})
// 	}

// })
