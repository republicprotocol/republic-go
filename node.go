package swarm

import (
	"github.com/republicprotocol/republic/dht"
	"github.com/republicprotocol/republic/rpc"
	"golang.org/x/net/context"
)

const IdLength = 20

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
func (node *Node) Ping(ctx context.Context, _ *rpc.Nothing) (*rpc.ID, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return node.ID, err
	}
	return node.ID, nil
}

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, _ *rpc.Nothing) (*rpc.MultiAddresses, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		// TODO: implement Peers.
	}()

	select {
	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddresses{}, ctx.Err()

	// Select the value passed by the goroutine.
	case ret := <-wait:
		return ret, nil
	}
}

// CloserPeers returns the peers of an rpc.Node that are closer to a target
// than the rpc.Node itself. Distance is calculated by evaluating a XOR with
// the target address and each peer address.
func (node *Node) CloserPeers(ctx context.Context, target *rpc.ID) (*rpc.MultiAddresses, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		// TODO: implement CloserPeers
	}()

	select {
	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddresses{}, ctx.Err()

	// Select the value passed by the goroutine.
	case ret := <-wait:
		return ret, nil
	}
}



