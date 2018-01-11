package topology

import (
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
var numberOfNodes = 20

// The number of messages that will be sent through the topology.
var numberOfMessages = 20

// The duration to wait for peers to start listening for RPCs.
var startTimeDelay = time.Second

// generateNodes generates nodes at the beginning of each topology test.
func generateNodes(numberOfNodes int) ([]*x.Node, error) {
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
		delegate := NewMockDelegate()
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

func sendMessages(nodes []*x.Node) error {
	for i := 0; i < numberOfMessages; i++ {
		left, right := randomNodes(nodes)
		if err := sendMessage(left, right); err != nil {
			return err
		}
	}
	return nil
}

func sendMessage(from, to *x.Node) error {
	address, err := to.MultiAddress.Address()
	if err != nil {
		return err
	}
	orderFragment := &rpc.OrderFragment{
		To:              string(address),
		From:            string(from.MultiAddress.String()),
		OrderID:         []byte("orderID"),
		OrderFragmentID: []byte("fragmentID"),
		OrderFragment:   []byte(address),
	}
	return from.ForwardOrderFragment(orderFragment)

}

func randomNodes(nodes []*x.Node) (*x.Node, *x.Node) {
	left := rand.Intn(len(nodes))
	right := rand.Intn(len(nodes))
	for left == right {
		right = rand.Intn(len(nodes))
	}
	return nodes[left], nodes[right]
}
