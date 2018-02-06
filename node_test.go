package network_test

import (
	"fmt"
	"log"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network"
)

type mockDelegate struct {
	mu                       *sync.Mutex
	numberOfPings            int
	numberOfQueryCloserPeers int
}

func newMockDelegate() *mockDelegate {
	return &mockDelegate{
		mu:                       new(sync.Mutex),
		numberOfPings:            0,
		numberOfQueryCloserPeers: 0,
	}
}

func (delegate *mockDelegate) OnPingReceived(_ identity.MultiAddress) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfPings++
}

func (delegate *mockDelegate) OnQueryCloserPeersReceived(_ identity.MultiAddress) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfQueryCloserPeers++
}

// boostrapping
var _ = Describe("Bootstrapping", func() {

	var err error
	var bootstrapNodes []*network.Node
	var bootstrapRoutingTable map[identity.Address][]*network.Node
	var swarmNodes []*network.Node
	var delegate *mockDelegate

	setupBootstrapNodes := func(topology Topology, numberOfNodes int) {
		bootstrapNodes, bootstrapRoutingTable, err = GenerateBootstrapTopology(topology, numberOfNodes, newMockDelegate())
		Ω(err).ShouldNot(HaveOccurred())
		for i, node := range bootstrapNodes {
			By(fmt.Sprintf("%dth bootstrap node is %s", i, node.MultiAddress()))
		}
		for _, node := range bootstrapNodes {
			go func(node *network.Node) {
				defer GinkgoRecover()
				Ω(node.Serve()).ShouldNot(HaveOccurred())
			}(node)
		}
		time.Sleep(time.Second)
		err = ping(bootstrapNodes, bootstrapRoutingTable)
	}

	setupSwarmNodes := func(numberOfNodes int) {
		swarmNodes, err = GenerateNodes(NodePortSwarm, numberOfNodes, newMockDelegate())
		Ω(err).ShouldNot(HaveOccurred())
		for _, swarmNode := range swarmNodes {
			for _, bootstrapNode := range bootstrapNodes {
				swarmNode.Options.BootstrapMultiAddresses = append(swarmNode.Options.BootstrapMultiAddresses, bootstrapNode.MultiAddress())
			}
		}
		for _, node := range swarmNodes {
			go func(node *network.Node) {
				defer GinkgoRecover()
				Ω(node.Serve()).ShouldNot(HaveOccurred())
			}(node)
		}
		time.Sleep(time.Second)
	}

	BeforeEach(func() {
		delegate = newMockDelegate()
	})

	AfterEach(func() {
		for _, node := range bootstrapNodes {
			func(node *network.Node) {
				node.Stop()
			}(node)
		}
		for _, node := range swarmNodes {
			func(node *network.Node) {
				node.Stop()
			}(node)
		}
	})

	for _, topology := range []Topology{TopologyFull, TopologyLine, TopologyRing, TopologyStar} {
		func(topology Topology) {
			Context(fmt.Sprintf("when bootstrap nodes are configured in a %s topology.\n", topology), func() {
				for _, numberOfBootstrapNodes := range []int{2, 4, 6} {
					for _, numberOfNodes := range []int{2, 6, 10, 14, 18} {
						func(topology Topology, numberOfBootstrapNodes, numberOfNodes int) {
							Context(fmt.Sprintf("with %d bootstrap nodes and %d swarm nodes.\n", numberOfBootstrapNodes, numberOfNodes), func() {
								It("should be able to successfully ping between nodes", func() {
									// Tests should be run serially to prevent
									// port overlaps.STEP: 0th bootstrap node is /ip4/127.0.0.1/tcp/3000/republic/8MJtdcgaGFxrBLJ1RhwXS3SQd2DcTG
									testMu.Lock()
									defer testMu.Unlock()

									// Setup testing configuration.
									setupBootstrapNodes(topology, numberOfBootstrapNodes)
									setupSwarmNodes(numberOfNodes)

									// Bootstrap all swarm nodes.
									for _, node := range swarmNodes {
										node.Bootstrap()
									}

									numberOfPings := 0
									for i := 0; i < numberOfNodes; i++ {
										to, from := PickRandomNodes(swarmNodes)
										if err := Ping(to, from); err == nil {
											numberOfPings++
										} else {
											log.Println(err)
										}
									}
									log.Printf("%v/%v successful pings", numberOfPings, numberOfNodes)
									Ω(numberOfPings).Should(BeNumerically(">", 2*numberOfNodes/3))
								})
							})
						}(topology, numberOfBootstrapNodes, numberOfNodes)
					}
				}
			})
		}(topology)
	}
})

//var _ = Describe("Pinging", func() {
//
//	run := func(name string, numberOfNodes int) int {
//		var nodes []*network.Node
//		var topology map[identity.Address][]*network.Node
//		var err error
//
//		delegate := newMockDelegate()
//		switch name {
//		case "full":
//			nodes, topology, err = generateFullyConnectedTopology(numberOfNodes, delegate)
//		case "star":
//			nodes, topology, err = generateStarTopology(numberOfNodes, delegate)
//		case "line":
//			nodes, topology, err = generateLineTopology(numberOfNodes, delegate)
//		case "ring":
//			nodes, topology, err = generateRingTopology(numberOfNodes, delegate)
//		}
//		Ω(err).ShouldNot(HaveOccurred())
//
//		for _, node := range nodes {
//			go func(node *network.Node) {
//				defer GinkgoRecover()
//				Ω(node.Serve()).ShouldNot(HaveOccurred())
//			}(node)
//			defer func(node *network.Node) {
//				defer GinkgoRecover()
//				node.Stop()
//			}(node)
//		}
//		time.Sleep(time.Second)
//
//		err = ping(nodes, topology)
//		Ω(err).ShouldNot(HaveOccurred())
//
//		return delegate.numberOfPings
//	}
//
//	for _, numberOfNodes := range []int{10, 20, 40, 80} {
//		func(numberOfNodes int) {
//			Context(fmt.Sprintf("in a fully connected topology with %d nodes", numberOfNodes), func() {
//				It("should update the DHT", func() {
//					testMu.Lock()
//					defer testMu.Unlock()
//					numberOfPings := run("full", numberOfNodes)
//					Ω(numberOfPings).Should(Equal(numberOfNodes * (numberOfNodes - 1)))
//				})
//			})
//		}(numberOfNodes)
//	}
//
//	for _, numberOfNodes := range []int{10, 20, 40, 80} {
//		func(numberOfNodes int) {
//			Context(fmt.Sprintf("in a star topology with %d nodes", numberOfNodes), func() {
//				It("should update the DHT", func() {
//					testMu.Lock()
//					defer testMu.Unlock()
//					numberOfPings := run("star", numberOfNodes)
//					Ω(numberOfPings).Should(Equal(2 * (numberOfNodes - 1)))
//				})
//			})
//		}(numberOfNodes)
//	}
//
//	for _, numberOfNodes := range []int{10, 20, 40, 80} {
//		func(numberOfNodes int) {
//			Context(fmt.Sprintf("in a line topology with %d nodes", numberOfNodes), func() {
//				It("should update the DHT", func() {
//					testMu.Lock()
//					defer testMu.Unlock()
//					numberOfPings := run("line", numberOfNodes)
//					Ω(numberOfPings).Should(Equal(2 * (numberOfNodes - 1)))
//				})
//			})
//		}(numberOfNodes)
//	}
//
//	for _, numberOfNodes := range []int{10, 20, 40, 80} {
//		func(numberOfNodes int) {
//			Context(fmt.Sprintf("in a ring topology with %d nodes", numberOfNodes), func() {
//				It("should update the DHT", func() {
//					testMu.Lock()
//					defer testMu.Unlock()
//					numberOfPings := run("ring", numberOfNodes)
//					Ω(numberOfPings).Should(Equal(2 * numberOfNodes))
//				})
//			})
//		}(numberOfNodes)
//	}
//})

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
