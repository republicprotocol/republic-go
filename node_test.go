package swarm_test

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm-network"
)

type mockDelegate struct {
	mu                                 *sync.Mutex
	numberOfPings                      int
	numberOfQueryCloserPeers           int
	numberOfQueryCloserPeersOnFrontier int
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

func (delegate *mockDelegate) OnQueryCloserPeersOnFrontierReceived(_ identity.MultiAddress) {
	delegate.mu.Lock()
	defer delegate.mu.Unlock()
	delegate.numberOfQueryCloserPeersOnFrontier++
}

// boostrapping
var _ = Describe("Bootstrapping", func() {

	var err error
	var bootstrapNodes []*swarm.Node
	var bootstrapRoutingTable map[identity.Address][]*swarm.Node
	var swarmNodes []*swarm.Node
	var delegate *mockDelegate

	setupBootstrapNodes := func(topology Topology, numberOfNodes int) {
		bootstrapNodes, bootstrapRoutingTable, err = GenerateBootstrapTopology(topology, numberOfNodes, newMockDelegate())
		Ω(err).ShouldNot(HaveOccurred())
		for i, node := range bootstrapNodes {
			go func(i int, node *swarm.Node) {
				defer GinkgoRecover()
				listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", NodePortBootstrap+i))
				Ω(err).ShouldNot(HaveOccurred())
				node.Register()
				Ω(node.Server.Serve(listener)).ShouldNot(HaveOccurred())
			}(i, node)
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
		for i, node := range swarmNodes {
			go func(i int, node *swarm.Node) {
				defer GinkgoRecover()
				listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", NodePortSwarm+i))
				Ω(err).ShouldNot(HaveOccurred())
				node.Register()
				Ω(node.Server.Serve(listener)).ShouldNot(HaveOccurred())
			}(i, node)
		}
		time.Sleep(time.Second)
	}

	BeforeEach(func() {
		delegate = newMockDelegate()
	})

	AfterEach(func() {
		for _, node := range bootstrapNodes {
			func(node *swarm.Node) {
				node.Server.Stop()
			}(node)
		}
		for _, node := range swarmNodes {
			func(node *swarm.Node) {
				node.Server.Stop()
			}(node)
		}
	})

	for _, topology := range []Topology{TopologyFull, TopologyLine, TopologyRing, TopologyStar} {
		func(topology Topology) {
			Context(fmt.Sprintf("when bootstrap nodes are configured in a %s topology.\n", topology), func() {
				for _, numberOfBootstrapNodes := range []int{4, 6, 8} {
					for _, numberOfNodes := range []int{4, 8, 12, 16, 20} {
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
									Ω(numberOfPings).Should(BeNumerically(">=", 2*numberOfNodes/3))
								})
							})
						}(topology, numberOfBootstrapNodes, numberOfNodes)
					}
				}
			})
		}(topology)
	}
})
