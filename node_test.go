package swarm_test

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-dht"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-swarm-network"
	"google.golang.org/grpc"
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

	Context("negative tests", func() {
		var node *swarm.Node

		BeforeEach(func() {
			keyPair, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddress, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/" + keyPair.Address().String())
			Ω(err).ShouldNot(HaveOccurred())
			node = swarm.NewNode(grpc.NewServer(),
				delegate,
				swarm.Options{
					MultiAddress:    multiAddress,
					Debug:           swarm.DebugHigh,
					Alpha:           DefaultOptionsAlpha,
					MaxBucketLength: DefaultOptionsMaxBucketLength,
					Timeout:         DefaultOptionsTimeout,
					TimeoutStep:     DefaultOptionsTimeoutStep,
					TimeoutRetries:  DefaultOptionsTimeoutRetries,
					Concurrent:      DefaultOptionsConcurrent,
				},
			)
		})

		Context("debug logs", func() {
			It("should print logs when we set the debug option to high", func() {
				// Setup testing configuration.
				setupBootstrapNodes(TopologyFull, 8)
				bootstrapNodes[0].Options.Debug = swarm.DebugHigh
				bootstrapNodes[0].Options.Concurrent = true
				setupSwarmNodes(16)
				swarmNodes[0].Options.Debug = swarm.DebugHigh
				swarmNodes[0].Options.Concurrent = true

				// Bootstrap all swarm nodes.
				for _, i := range swarmNodes {
					i.Bootstrap()
				}

				numberOfPings := 0
				for i := 0; i < 20; i++ {
					to, from := PickRandomNodes(swarmNodes)
					if err := Ping(to, from); err == nil {
						numberOfPings++
					} else {
						log.Println(err)
					}
				}
				log.Printf("%v/%v successful pings", numberOfPings, 20)
				Ω(numberOfPings).Should(BeNumerically(">=", 2*20/3))
			})
		})

		Context("context", func() {
			It("should return an error when calling with a canceled context", func() {
				contextWithTimeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
				cancel()
				_, err = node.Ping(contextWithTimeout, &rpc.MultiAddress{})
				Ω(err).Should(HaveOccurred())

				_, err = node.QueryCloserPeers(contextWithTimeout, &rpc.Query{})
				Ω(err).Should(HaveOccurred())

			})
		})

		Context("pruning the dht", func() {
			It("should return an error when pruning with bad dht", func() {
				keyPair, err := identity.NewKeyPair()
				Ω(err).ShouldNot(HaveOccurred())

				correctAddress := keyPair.Address()
				worngFormatTarget := identity.Address("wrongAddress")
				_, err = node.Prune(worngFormatTarget)
				Ω(err).Should(HaveOccurred())

				pruned, err := node.Prune(correctAddress)
				Ω(pruned).Should(BeFalse())
				Ω(err).ShouldNot(HaveOccurred())

				for i := 0; i < dht.IDLengthInBits*node.Options.MaxBucketLength; i++ {
					keyPair, err = identity.NewKeyPair()
					Ω(err).ShouldNot(HaveOccurred())
					multi, err := identity.NewMultiAddressFromString("/ip4/192.168.0.0/tcp/80/republic/" + keyPair.Address().String())
					Ω(err).ShouldNot(HaveOccurred())
					err = node.DHT.UpdateMultiAddress(multi)
					if err == dht.ErrFullBucket {
						pruned, err := node.Prune(keyPair.Address())
						Ω(pruned).Should(BeTrue())
						Ω(err).ShouldNot(HaveOccurred())
						break
					}
				}
			})
		})

		Context("bad rpc input", func() {
			It("should return an error with a bad ping input", func() {
				_, err := node.Ping(context.Background(), &rpc.MultiAddress{})
				Ω(err).Should(HaveOccurred())
			})

			It("should return an error with a bad query input", func() {
				_, err := node.QueryCloserPeers(context.Background(), &rpc.Query{Query: &rpc.Address{}, From: &rpc.MultiAddress{}})
				Ω(err).Should(HaveOccurred())

				err = node.QueryCloserPeersOnFrontier(&rpc.Query{Query: &rpc.Address{}, From: &rpc.MultiAddress{}}, *new(rpc.SwarmNode_QueryCloserPeersOnFrontierServer))
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("bootstrapping", func() {
			It("should return error with offlane bootstrap nodes", func() {
				bootstrapNode, err := identity.NewMultiAddressFromString("/ip4/192.168.0.0/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
				Ω(err).ShouldNot(HaveOccurred())
				node.Options.BootstrapMultiAddresses = identity.MultiAddresses{bootstrapNode}
				node.Bootstrap()
			})
		})

	})
})
