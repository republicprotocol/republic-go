package x_test

import (
	"fmt"
	"math/rand"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x"
)

func generateNodes(numberOfNodes int, delegate x.Delegate) ([]*x.Node, error) {
	nodes := make([]*x.Node, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nil, err
		}
		multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000+i, keyPair.Address()))
		if err != nil {
			return nil, err
		}
		node, err := x.NewNode(
			multi,
			make(identity.MultiAddresses, 0, numberOfNodes-1),
			delegate,
		)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return nodes, nil
}

func generateFullyConnectedTopology(numberOfNodes int, delegate x.Delegate) ([]*x.Node, map[identity.Address][]*x.Node, error) {
	nodes, err := generateNodes(numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	topology := map[identity.Address][]*x.Node{}
	for i, node := range nodes {
		topology[node.DHT.Address] = []*x.Node{}
		for j, peer := range nodes {
			if i == j {
				continue
			}
			topology[node.DHT.Address] = append(topology[node.DHT.Address], peer)
		}
	}
	return nodes, topology, nil
}

func generateStarTopology(numberOfNodes int, delegate x.Delegate) ([]*x.Node, map[identity.Address][]*x.Node, error) {
	nodes, err := generateNodes(numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	topology := map[identity.Address][]*x.Node{}
	for i, node := range nodes {
		topology[node.DHT.Address] = []*x.Node{}
		if i == 0 {
			for j, peer := range nodes {
				if i == j {
					continue
				}
				topology[node.DHT.Address] = append(topology[node.DHT.Address], peer)
			}
		} else {
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[0])
		}
	}
	return nodes, topology, nil
}

func generateLineTopology(numberOfNodes int, delegate x.Delegate) ([]*x.Node, map[identity.Address][]*x.Node, error) {
	nodes, err := generateNodes(numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	topology := map[identity.Address][]*x.Node{}
	for i, node := range nodes {
		topology[node.DHT.Address] = []*x.Node{}
		if i == 0 {
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[i+1])
		} else if i == len(nodes)-1 {
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[i-1])
		} else {
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[i+1])
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[i-1])
		}
	}
	return nodes, topology, nil
}

func generateRingTopology(numberOfNodes int, delegate x.Delegate) ([]*x.Node, map[identity.Address][]*x.Node, error) {
	nodes, err := generateNodes(numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	topology := map[identity.Address][]*x.Node{}
	for i, node := range nodes {
		topology[node.DHT.Address] = []*x.Node{}
		if i == 0 {
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[i+1])
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[len(nodes)-1])
		} else if i == len(nodes)-1 {
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[i-1])
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[0])
		} else {
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[i+1])
			topology[node.DHT.Address] = append(topology[node.DHT.Address], nodes[i-1])
		}
	}
	return nodes, topology, nil
}

func randomNodes(nodes []*x.Node) (*x.Node, *x.Node) {
	left := rand.Intn(len(nodes))
	right := rand.Intn(len(nodes))
	for left == right {
		right = rand.Intn(len(nodes))
	}
	return nodes[left], nodes[right]
}