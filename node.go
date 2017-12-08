package swarm

import (
	"github.com/republicprotocol/republic/dht"
	"github.com/republicprotocol/republic/rpc"
	"golang.org/x/net/context"
	_ "container/list"
)

// Node implements the gRPC Node service.
type Node struct {
	DHT *dht.RoutingTable
}

// NewNode returns a new node with no connections.
func NewNode(id dht.ID) *Node {
	return &Node{DHT: dht.NewRoutingTable(id)}
}

// Ping is used to test connections to a Node. The Node will respond with its
// rpc.ID. If the Node does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, id *rpc.ID) (*rpc.ID, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.ID{Address: string(node.DHT.ID)}, err
	}

	// Update the sender in the node routing table
	if err := node.DHT.Update(dht.ID(id.Address)); err !=nil{
		return &rpc.ID{Address: string(node.DHT.ID)}, err
	}

	return &rpc.ID{Address: string(node.DHT.ID)}, nil
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
		wait <- node.peers()
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
		peers, _:= node.closerPeers(target.Address)
		wait <- peers
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

func (node *Node) peers() *rpc.MultiAddresses {
	peers := node.DHT.All()
	var ret []string
	for e := peers.Front(); e != nil ;e.Next(){
		ret = append(ret, e.Value.(string))
	}
	return &rpc.MultiAddresses{Multis:ret}
}

func (node *Node) closerPeers(id string) (*rpc.MultiAddresses,error){
	peers,err  := node.DHT.FindClosest(dht.ID(id))
	if err != nil {
		return nil, err
	}
	var ret []string
	for e := peers.Front(); e != nil ;e.Next(){
		ret = append(ret, e.Value.(string))
	}
	return &rpc.MultiAddresses{Multis:ret},nil
}