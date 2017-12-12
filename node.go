package main

import (
	_ "container/list"
	"github.com/republicprotocol/republic/dht"
	"github.com/republicprotocol/republic/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	if err := node.updateSender(id); err != nil {
		return nil, err
	}

	return &rpc.ID{Address: string(node.DHT.ID)}, nil
}

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, id *rpc.ID) (*rpc.MultiAddresses, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	// Update the sender in the node routing table
	if err := node.updateSender(id); err != nil {
		return nil, err
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
func (node *Node) CloserPeers(ctx context.Context, path *rpc.Path) (*rpc.MultiAddresses, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	// Update the sender in the node routing table
	if err := node.updateSender(path.From); err != nil {
		return nil, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		peers, _ := node.closerPeers(path.To.Address)
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

// Return all peers in the node routing table
func (node *Node) peers() *rpc.MultiAddresses {
	peers := node.DHT.All()
	var ret []string
	for e := peers.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(string))
	}
	return &rpc.MultiAddresses{Multis: ret}
}

// Return the closer peers in the node routing table
func (node *Node) closerPeers(id string) (*rpc.MultiAddresses, error) {
	peers, err := node.DHT.FindClosest(dht.ID(id))
	if err != nil {
		return nil, err
	}
	var ret []string
	for e := peers.Front(); e != nil; e = e.Next() {
		if e.Value != nil {
			ret = append(ret, e.Value.(string))
		}
	}
	return &rpc.MultiAddresses{Multis: ret}, nil
}

// Every time we receive a request, update the sender
// as a active peer in the node routing table
func (node *Node) updateSender(id *rpc.ID) error {
	return node.DHT.Update(dht.ID(id.Address))
}

// Connect to other Node and return the grpc client
func connectNode(address string) rpc.NodeClient {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	defer conn.Close()
	return rpc.NewNodeClient(conn)
}

// Ping a node
func (node *Node) PingNode(address string) (*rpc.ID, error) {
	client := connectNode(address)
	pong, err := client.Ping(context.Background(), &rpc.ID{Address: string(node.DHT.ID)})
	if err != nil {
		return nil, err
	}
	return pong, nil
}

// Request all peers of a node
func (node *Node) PeersNode(address string) (*rpc.MultiAddresses, error) {
	client := connectNode(address)
	multiAddresses, err := client.Peers(context.Background(), &rpc.ID{Address: string(node.DHT.ID)})
	if err != nil {
		return nil, err
	}
	return multiAddresses, nil
}

// Find close peers from a node
// Ping a node
func (node *Node) findClosePeers(address string, target dht.ID) (*rpc.MultiAddresses, error) {
	client := connectNode(address)
	multiAddresses, err := client.CloserPeers(context.Background(), &rpc.Path{From: &rpc.ID{Address: string(node.DHT.ID)}, To: &rpc.ID{Address: string(target)}})
	if err != nil {
		return nil, err
	}
	return multiAddresses, nil
}

// Find a certain node by its ID  throught the kademlia network
func (node *Node) FindNode(id dht.ID) (string, error) {
	// todo
	return "", nil
}
