package network_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/republic-go/network/dht"
	"google.golang.org/grpc"
)

type Topology string

const (
	TopologyFull = "full"
	TopologyLine = "line"
	TopologyRing = "ring"
	TopologyStar = "star"
)

const (
	NodePortBootstrap = 3000
	NodePortSwarm     = 4000
)

func GenerateBootstrapTopology(topology Topology, numberOfNodes int, delegate SwarmDelegate) ([]*Swarm, map[identity.Address][]*Swarm, error) {
	var err error
	var nodes []*Swarm
	var routingTable map[identity.Address][]*Swarm

	switch topology {
	case TopologyFull:
		nodes, routingTable, err = GenerateFullTopology(NodePortBootstrap, numberOfNodes, delegate)
	case TopologyStar:
		nodes, routingTable, err = GenerateStarTopology(NodePortBootstrap, numberOfNodes, delegate)
	case TopologyLine:
		nodes, routingTable, err = GenerateLineTopology(NodePortBootstrap, numberOfNodes, delegate)
	case TopologyRing:
		nodes, routingTable, err = GenerateRingTopology(NodePortBootstrap, numberOfNodes, delegate)
	}
	return nodes, routingTable, err
}

func GenerateNodes(port, numberOfNodes int, delegate SwarmDelegate) ([]*Swarm, error) {
	nodes := make([]*Swarm, numberOfNodes)
	for i := range nodes {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nil, err
		}
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", port+i, keyPair.Address()))
		if err != nil {
			return nil, err
		}
		node := NewSwarm(grpc.NewServer(),
			delegate,
			Options{
				MultiAddress:    multiAddress,
				Debug:           DefaultOptionsDebug,
				Alpha:           DefaultOptionsAlpha,
				MaxBucketLength: DefaultOptionsMaxBucketLength,
				Timeout:         DefaultOptionsTimeout,
				TimeoutStep:     DefaultOptionsTimeoutStep,
				TimeoutRetries:  DefaultOptionsTimeoutRetries,
				Concurrent:      DefaultOptionsConcurrent,
			},
		)
		nodes[i] = node
	}
	return nodes, nil
}

func GenerateFullTopology(port, numberOfNodes int, delegate SwarmDelegate) ([]*Swarm, map[identity.Address][]*Swarm, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*Swarm{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*Swarm{}
		for j, peer := range nodes {
			if i == j {
				continue
			}
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], peer)
		}
	}
	return nodes, routingTable, nil
}

func GenerateLineTopology(port, numberOfNodes int, delegate SwarmDelegate) ([]*Swarm, map[identity.Address][]*Swarm, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*Swarm{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*Swarm{}
		if i == 0 {
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[i+1])
		} else if i == len(nodes)-1 {
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[i-1])
		} else {
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[i+1])
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[i-1])
		}
	}
	return nodes, routingTable, nil
}

func GenerateRingTopology(port, numberOfNodes int, delegate SwarmDelegate) ([]*Swarm, map[identity.Address][]*Swarm, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*Swarm{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*Swarm{}
		if i == 0 {
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[i+1])
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[len(nodes)-1])
		} else if i == len(nodes)-1 {
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[i-1])
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[0])
		} else {
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[i+1])
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[i-1])
		}
	}
	return nodes, routingTable, nil
}

func GenerateStarTopology(port, numberOfNodes int, delegate SwarmDelegate) ([]*Swarm, map[identity.Address][]*Swarm, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*Swarm{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*Swarm{}
		if i == 0 {
			for j, peer := range nodes {
				if i == j {
					continue
				}
				routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], peer)
			}
		} else {
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], nodes[0])
		}
	}
	return nodes, routingTable, nil
}

func Ping(to *Swarm, from *Swarm) error {
	var target *identity.MultiAddress

	multiAddress, err := from.DHT.FindMultiAddress(to.Address())
	if err != nil {
		return err
	}
	if multiAddress != nil {
		target = multiAddress
	}

	if target == nil {
		multiAddresses, err := rpc.QueryCloserPeersOnFrontierFromTarget(
			from.MultiAddress(),
			from.MultiAddress(),
			to.Address(),
			DefaultOptionsTimeout,
		)
		if err != nil {
			return err
		}
		for _, multiAddress := range multiAddresses {
			if to.Address() == multiAddress.Address() {
				target = &multiAddress
				break
			}
		}
	}
	if target != nil {
		return rpc.PingTarget(*target, from.MultiAddress(), DefaultOptionsTimeout)
	}
	return fmt.Errorf("ping error: %v could not find %v", from.Address(), to.Address())
}

func PickRandomNodes(nodes []*Swarm) (*Swarm, *Swarm) {
	i := rand.Intn(len(nodes))
	j := rand.Intn(len(nodes))
	for i == j {
		j = rand.Intn(len(nodes))
	}
	return nodes[i], nodes[j]
}

func ping(nodes []*Swarm, topology map[identity.Address][]*Swarm) error {
	var wg sync.WaitGroup
	wg.Add(len(nodes))

	muError := new(sync.Mutex)
	var globalError error = nil

	for _, node := range nodes {
		go func(node *Swarm) {
			defer wg.Done()
			peers := topology[node.DHT.Address]
			for _, peer := range peers {
				err := rpc.PingTarget(peer.MultiAddress(), node.MultiAddress(), time.Second)
				if err != nil {
					muError.Lock()
					globalError = err
					muError.Unlock()
				}
			}
		}(node)
	}

	wg.Wait()
	return globalError
}

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
	var bootstrapNodes []*Swarm
	var bootstrapRoutingTable map[identity.Address][]*Swarm
	var swarmNodes []*Swarm
	var delegate *mockDelegate

	setupBootstrapNodes := func(topology Topology, numberOfNodes int) {
		bootstrapNodes, bootstrapRoutingTable, err = GenerateBootstrapTopology(topology, numberOfNodes, newMockDelegate())
		Ω(err).ShouldNot(HaveOccurred())
		for i, node := range bootstrapNodes {
			go func(i int, node *Swarm) {
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
			go func(i int, node *Swarm) {
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
			func(node *Swarm) {
				node.Server.Stop()
			}(node)
		}
		for _, node := range swarmNodes {
			func(node *Swarm) {
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

									// test finding node function
									left, right := PickRandomNodes(swarmNodes)
									found, err := left.FindNode(right.Address().ID())
									Ω(err).ShouldNot(HaveOccurred())
									if found != nil {
										Ω(*found).Should(Equal(right.MultiAddress()))
										log.Println("found the target node by its ID ")
									} else {
										log.Println("fail to found the target")
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
		var node *Swarm

		BeforeEach(func() {
			keyPair, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddress, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80/republic/" + keyPair.Address().String())
			Ω(err).ShouldNot(HaveOccurred())
			node = NewSwarm(grpc.NewServer(),
				delegate,
				Options{
					MultiAddress:    multiAddress,
					Debug:           DefaultOptionsDebug,
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
				bootstrapNodes[0].Options.Debug = DebugHigh
				bootstrapNodes[0].Options.Concurrent = true
				setupSwarmNodes(16)
				swarmNodes[0].Options.Debug = DebugHigh
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
			It("should return an error when calling with a cancelled context", func() {
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
				wrongFormatTarget := identity.Address("wrongAddress")
				_, err = node.Prune(wrongFormatTarget)
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

			})
		})

		Context("bootstrapping", func() {
			It("should print error with offline bootstrap nodes", func() {
				bootstrapNode, err := identity.NewMultiAddressFromString("/ip4/192.168.0.0/tcp/80/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f")
				Ω(err).ShouldNot(HaveOccurred())
				node.Options.BootstrapMultiAddresses = identity.MultiAddresses{bootstrapNode}
				node.Options.Debug = DebugHigh

				node.Bootstrap()
			})
		})
	})
})
