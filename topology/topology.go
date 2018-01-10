package topology

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x"
	"github.com/republicprotocol/go-x/rpc"
)

// μ prevents multiple topology tests running in parallel. This is needed to
// protect overlapping ports during tests.
var μ = new(sync.Mutex)

// The number of nodes that should be included in each topology test.
var numberOfNodes = 100

// The number of messages that will be sent through the topology.
var numberOfMessages = 100

// The duration to wait for peers to start listening for RPCs.
var startTimeDelay = time.Second

func generateNodes(numberOfNodes int) ([]*x.Node, error) {
	nodes := make([]*x.Node, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		keyPair, err := identity.NewKeyPair()
		if err != nil {
			return nil, err
		}
		multi, err := identity.NewMultiAddress(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000+i, keyPair.Address()))
		if err != nil {
			return nil, err
		}
		node, err := x.NewNode(&x.Config{
			KeyPair:      keyPair,
			MultiAddress: multi,
			Peers:        make(identity.MultiAddresses, 0, numberOfNodes-1),
		})
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return nodes, nil
}

func sendMessages(nodes []*x.Node) error {
	for i := 0; i < numberOfMessages; i++ {
		left, right := randomNodes(nodes)
		if err := sendMessage(left.MultiAddress, right.MultiAddress); err != nil {
			return err
		}
	}
	return nil
}

func sendMessage(from identity.MultiAddress, to identity.MultiAddress) error {
	client, conn, err := x.NewNodeClient(from)
	defer conn.Close()
	if err != nil {
		return err
	}
	address, err := to.Address()
	if err != nil {
		return err
	}
	_, err = client.Send(context.Background(), &rpc.Payload{
		To:   string(address),
		Data: "message",
	})
	return err
}

func randomNodes(nodes []*x.Node) (*x.Node, *x.Node) {
	left := rand.Intn(len(nodes))
	right := rand.Intn(len(nodes))
	for left == right {
		right = rand.Intn(len(nodes))
	}
	return nodes[left], nodes[right]
}
