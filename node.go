package swarm

import (
	"fmt"

	"github.com/republic/swarm/dht"
	"github.com/republic/swarm/rpc"
	"golang.org/x/net/context"
)

// Node implements the gRPC Node service.
type Node struct {
	ID  *rpc.ID
	DHT *dht.RoutingTable
}

// NewNode returns a new node with no connections.
func NewNode(id *rpc.ID) *Node {
	return &Node{id, &dht.RoutingTable{}}
}

// Ping is used to test connections to a Node. The Node will respond with its
// rpc.MultiAddress. If the Node does not respond, or it responds with an
// error, then the connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, _ *rpc.Nothing) (*rpc.MultiAddress, error) {
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddress{}, err
	}
	return &rpc.MultiAddress{
		Multi: fmt.Sprintf("/ip4/0.0.0.0/tcp/5000/republic/%s", node.ID.Address),
	}, nil
}

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, _ *rpc.Nothing) (*rpc.MultiAddresses, error) {
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	return &rpc.MultiAddresses{}, nil
}

// CloserPeers returns the peers of an rpc.Node that are closer to a target
// than the rpc.Node itself. Distance is calculated by evaluating a XOR with
// the target address and each peer address.
func (node *Node) CloserPeers(ctx context.Context, target *rpc.ID) (*rpc.MultiAddresses, error) {
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	return &rpc.MultiAddresses{}, nil
}
