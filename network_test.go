package network_test

import (
	"fmt"
	"math/rand"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-rpc"
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

func GenerateBootstrapTopology(topology Topology, numberOfNodes int, delegate network.Delegate) ([]*network.Node, map[identity.Address][]*network.Node, error) {
	var err error
	var nodes []*network.Node
	var routingTable map[identity.Address][]*network.Node

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

func GenerateNodes(port, numberOfNodes int, delegate network.Delegate) ([]*network.Node, error) {
	nodes := make([]*network.Node, numberOfNodes)
	for i := range nodes {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nil, err
		}
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", port+i, keyPair.Address()))
		if err != nil {
			return nil, err
		}
		node := network.NewNode(
			delegate,
			network.Options{
				Host:            "127.0.0.1",
				Port:            fmt.Sprintf("%d", port+i),
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

func GenerateFullTopology(port, numberOfNodes int, delegate network.Delegate) ([]*network.Node, map[identity.Address][]*network.Node, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*network.Node{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*network.Node{}
		for j, peer := range nodes {
			if i == j {
				continue
			}
			routingTable[node.DHT.Address] = append(routingTable[node.DHT.Address], peer)
		}
	}
	return nodes, routingTable, nil
}

func GenerateLineTopology(port, numberOfNodes int, delegate network.Delegate) ([]*network.Node, map[identity.Address][]*network.Node, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*network.Node{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*network.Node{}
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

func GenerateRingTopology(port, numberOfNodes int, delegate network.Delegate) ([]*network.Node, map[identity.Address][]*network.Node, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*network.Node{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*network.Node{}
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

func GenerateStarTopology(port, numberOfNodes int, delegate network.Delegate) ([]*network.Node, map[identity.Address][]*network.Node, error) {
	nodes, err := GenerateNodes(port, numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	routingTable := map[identity.Address][]*network.Node{}
	for i, node := range nodes {
		routingTable[node.DHT.Address] = []*network.Node{}
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

func Ping(to *network.Node, from *network.Node) error {
	var target *identity.MultiAddress

	multiAddress, err := from.DHT.FindMultiAddress(to.Address())
	if err != nil {
		return err
	}
	if multiAddress != nil {
		target = multiAddress
	}

	if target == nil {
		multiAddresses, errs := rpc.QueryCloserPeersOnFrontierFromTarget(
			from.MultiAddress(),
			from.MultiAddress(),
			to.Address(),
			DefaultOptionsTimeout,
		)

		streaming := true
		for streaming {
			select {
			// Peers returned by the query will be added to the DHT.
			case multiAddress, ok := <-multiAddresses:
				if !ok {
					streaming = false
				} else if multiAddress.Address() == to.Address() {
					target = &multiAddress
					streaming = false
				}

			// Errors are not returned because it is reasonable that a
			// bootstrap Node might be unavailable at this time.
			case err, ok := <-errStream:
				if !ok {
					streaming = false
				} else {
					return err
				}
			}
		}
	}
	if target != nil {
		return rpc.PingTarget(*target, from.MultiAddress(), DefaultOptionsTimeout)
	}
	return fmt.Errorf("ping error: %v could not find %v", from.Address(), to.Address())
}

func PickRandomNodes(nodes []*network.Node) (*network.Node, *network.Node) {
	i := rand.Intn(len(nodes))
	j := rand.Intn(len(nodes))
	for i == j {
		j = rand.Intn(len(nodes))
	}
	return nodes[i], nodes[j]
}
