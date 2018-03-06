package swarm_test

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-swarm-network"
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

func GenerateBootstrapTopology(topology Topology, numberOfNodes int, delegate swarm.Delegate) ([]*swarm.Node, map[identity.Address][]*swarm.Node, error) {
	var err error
	var nodes []*swarm.Node
	var routingTable map[identity.Address][]*swarm.Node

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

func GenerateNodes(port, numberOfNodes int, delegate swarm.Delegate) ([]*swarm.Node, error) {
	nodes := make([]*swarm.Node, numberOfNodes)
	for i := range nodes {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nil, err
		}
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", port+i, keyPair.Address()))
		if err != nil {
			return nil, err
		}
		node := swarm.NewNode(grpc.NewServer(),
			delegate,
			swarm.Options{
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

func GenerateFullTopology(port, numberOfNodes int, delegate swarm.Delegate) ([]*swarm.Node, map[identity.Address][]*swarm.Node, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*swarm.Node{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*swarm.Node{}
		for j, peer := range nodes {
			if i == j {
				continue
			}
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], peer)
		}
	}
	return nodes, routingTable, nil
}

func GenerateLineTopology(port, numberOfNodes int, delegate swarm.Delegate) ([]*swarm.Node, map[identity.Address][]*swarm.Node, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*swarm.Node{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*swarm.Node{}
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

func GenerateRingTopology(port, numberOfNodes int, delegate swarm.Delegate) ([]*swarm.Node, map[identity.Address][]*swarm.Node, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*swarm.Node{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*swarm.Node{}
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

func GenerateStarTopology(port, numberOfNodes int, delegate swarm.Delegate) ([]*swarm.Node, map[identity.Address][]*swarm.Node, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*swarm.Node{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*swarm.Node{}
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

func Ping(to *swarm.Node, from *swarm.Node) error {
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

func PickRandomNodes(nodes []*swarm.Node) (*swarm.Node, *swarm.Node) {
	i := rand.Intn(len(nodes))
	j := rand.Intn(len(nodes))
	for i == j {
		j = rand.Intn(len(nodes))
	}
	return nodes[i], nodes[j]
}

func ping(nodes []*swarm.Node, topology map[identity.Address][]*swarm.Node) error {
	var wg sync.WaitGroup
	wg.Add(len(nodes))

	muError := new(sync.Mutex)
	var globalError error = nil

	for _, node := range nodes {
		go func(node *swarm.Node) {
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
