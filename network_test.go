package network_test

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-rpc"
)

func generateNodes(numberOfNodes int, delegate network.Delegate) ([]*network.Node, error) {
	nodes := make([]*network.Node, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nil, err
		}
		multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000+i, keyPair.Address()))
		if err != nil {
			return nil, err
		}
		node := network.NewNode(
			delegate,
			network.Options{
				MultiAddress:    multi,
				Debug:           network.DebugOff,
				Alpha:           3,
				MaxBucketLength: 100,
			},
		)
		nodes[i] = node
	}
	return nodes, nil
}

func generateFullyConnectedTopology(numberOfNodes int, delegate network.Delegate) ([]*network.Node, map[identity.Address][]*network.Node, error) {
	nodes, err := generateNodes(numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	topology := map[identity.Address][]*network.Node{}
	for i, node := range nodes {
		topology[node.DHT.Address] = []*network.Node{}
		for j, peer := range nodes {
			if i == j {
				continue
			}
			topology[node.DHT.Address] = append(topology[node.DHT.Address], peer)
		}
	}
	return nodes, topology, nil
}

func generateStarTopology(numberOfNodes int, delegate network.Delegate) ([]*network.Node, map[identity.Address][]*network.Node, error) {
	nodes, err := generateNodes(numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	topology := map[identity.Address][]*network.Node{}
	for i, node := range nodes {
		topology[node.DHT.Address] = []*network.Node{}
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

func generateLineTopology(numberOfNodes int, delegate network.Delegate) ([]*network.Node, map[identity.Address][]*network.Node, error) {
	nodes, err := generateNodes(numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	topology := map[identity.Address][]*network.Node{}
	for i, node := range nodes {
		topology[node.DHT.Address] = []*network.Node{}
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

func generateRingTopology(numberOfNodes int, delegate network.Delegate) ([]*network.Node, map[identity.Address][]*network.Node, error) {
	nodes, err := generateNodes(numberOfNodes, delegate)
	if err != nil {
		return nil, nil, err
	}
	topology := map[identity.Address][]*network.Node{}
	for i, node := range nodes {
		topology[node.DHT.Address] = []*network.Node{}
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

func randomNodes(nodes []*network.Node) (*network.Node, *network.Node) {
	left := rand.Intn(len(nodes))
	right := rand.Intn(len(nodes))
	for left == right {
		right = rand.Intn(len(nodes))
	}
	return nodes[left], nodes[right]
}

func ping(nodes []*network.Node, topology map[identity.Address][]*network.Node) error {
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	var muError *sync.Mutex
	var globalError error = nil

	for _, node := range nodes {
		go func(node *network.Node) {
			defer wg.Done()
			peers := topology[node.DHT.Address]
			for _, peer := range peers {
				err := rpc.PingTarget(peer.MultiAddress(), node.MultiAddress(), time.Second)
				if err != nil {
					muError.Lock()
					defer muError.Unlock()
					globalError = err
				}
			}
		}(node)
	}

	wg.Wait()
	return globalError
}

func peers(nodes []*network.Node, topology map[identity.Address][]*network.Node) error {
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	var muError *sync.Mutex
	var globalError error = nil

	for _, node := range nodes {
		go func(node *network.Node) {
			defer wg.Done()
			peers := topology[node.DHT.Address]
			connectedPeers, err := rpc.GetPeersFromTarget(node.MultiAddress(), identity.MultiAddress{}, time.Second)
			if err != nil {
				muError.Lock()
				defer muError.Unlock()
				globalError = err
			}
			for _, peer := range peers {
				connected := false
				for _, connectedPeer := range connectedPeers {
					if peer.MultiAddress().String() == connectedPeer.String() {
						connected = true
					}
				}
				if !connected {
					if err != nil {
						muError.Lock()
						defer muError.Unlock()
						globalError = fmt.Errorf("%s should be connected to %s", node.MultiAddress().String(), peer.MultiAddress().String())
					}
					return
				}
			}
		}(node)
	}

	wg.Wait()
	return globalError
}
