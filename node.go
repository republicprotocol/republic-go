package swarm

import (
	"container/list"

	"github.com/pkg/errors"
	"github.com/republicprotocol/go-swarm/dht"
	"github.com/republicprotocol/go-swarm/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/republicprotocol/go-identity"
	"github.com/multiformats/go-multiaddr"
)

// Node implements the gRPC Node service.
type Node struct {
	DHT *dht.RoutingTable
}

// NewNode returns a new node with no connections.
func NewNode(address identity.Address) *Node {
	return &Node{DHT: dht.NewRoutingTable(address)}
}

// Ping is used to test connections to a Node. The Node will respond with its
// address. If the Node does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, id *rpc.Node) (*rpc.Node, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.Node{Address: string(node.DHT.Address)}, err
	}

	// Update the sender in the node routing table
	if err := node.updateSender(id); err != nil {
		return nil, err
	}

	return &rpc.Node{Address: string(node.DHT.Address)}, nil
}

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, target *rpc.Node) (*rpc.MultiAddresses, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	// Update the sender in the node routing table
	if err := node.updateSender(target); err != nil {
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
// the target address and each peer ID.
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
		peers, _ := node.closerPeers(identity.Address(path.To.Address))
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
	multis := node.DHT.MultiAddresses()
	return &rpc.MultiAddresses{Multis: multis}
}

// Return the closer peers in the node routing table
func (node *Node) closerPeers(address identity.Address) (*rpc.MultiAddresses, error) {
	peers, err := node.DHT.FindClosest(address)
	if err != nil {
		return nil, err
	}
	return &rpc.MultiAddresses{Multis: peers.MultiAddresses()}, nil
}

// Update the address in the routing table if there is enough space
func (node *Node) updateSender(address *rpc.Node) error {
	rAddress := identity.Address(address.Address)
	// Check if there still has place for the new address
	peer, err := node.DHT.CheckAvailability(rAddress)
	if err != nil {
		return err
	}

	// todo : generate the multiaddress from its republic address and ip infor.
	mAddress, err := rAddress.MultiAddress()
	if err != nil {
		return err
	}

	// If the bucket is full
	if peer != "" {
		// Try to ping the last node in the bucket
		wait := make(chan interface{})
		go func() {
			defer close(wait)
			pong, err := node.PingNode(peer)
			if err != nil {
				wait <- err
			}else{
				wait <- pong
			}
		}()

		// Wait for the response
		resp := <- wait
		// If the last node is active, we do nothing with the new node
		switch resp.(type){
		case error:
			return resp.(error)
		case *rpc.Node:
			return nil
		}
	}

	return node.DHT.Update(mAddress)
}

// Connect to other Node and return the grpc client
func connectNode(address string) rpc.DHTClient {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	return rpc.NewDHTClient(conn)
}

// Ping a node
func (node *Node) PingNode(address string) (*rpc.Node, error) {
	client := connectNode(address)
	pong, err := client.Ping(context.Background(), &rpc.Node{Address: string(node.DHT.Address)})
	if err != nil {
		return nil, err
	}
	multiAddress, err := identity.NewMultiaddr("/republic/" + pong.Address)
	if err != nil {
		return nil, err
	}

	return pong, node.DHT.Update(multiAddress)
}

// Request all peers of a node
func (node *Node) PeersNode(address string) (*rpc.MultiAddresses, error) {
	client := connectNode(address)
	multiAddresses, err := client.Peers(context.Background(), &rpc.Node{Address: string(node.DHT.Address)})
	if err != nil {
		return nil, err
	}
	return multiAddresses, nil
}

// Find a certain node by its ID through the kademlia network
// Return its multi-adresses
func (node *Node) FindNode(address identity.Address) (multiaddr.Multiaddr, error) {
	// Find closest peers we know from the routing table
	peers, err := node.DHT.FindClosest(address)
	if err != nil {
		return nil, err
	}

	path := &rpc.Path{From: &rpc.Node{Address: string(node.DHT.Address)}, To: &rpc.Node{Address: string(address)}}

	for {

		peerValue, err := peers.Front().Value.(multiaddr.Multiaddr).ValueForProtocol(dht.RepublicCode)
		if err != nil {
			return nil, nil
		}
		// Return the multiaddress if we find the target
		if peers.Front() != nil && peerValue == string(address) {
			return peers.Front().Value.(multiaddr.Multiaddr), nil
		}

		nPeers := dht.RoutingBucket{list.List{}}
		// todo : parallel  and get address from multi address
		for e := peers.Front(); e != nil; e = e.Next() {
			if e.Value == "/republic/"+node.DHT.Address {
				continue
			}
			address := "localhost:8080"
			client := connectNode(address)

			multiAddresses, err := client.CloserPeers(context.Background(), path)

			if err != nil {
				return nil, err
			}
			for _, address := range multiAddresses.Multis {
				nPeers.PushBack(address)
			}
		}


		// Check if the new peers are unchanged from the previous list
		// which means we can't get closer peers from the node we know
		if dht.CompareList(nPeers, peers, 3) {
			// todo: what to do if we can't find node
			return nil, errors.New("can't find target from peers I know, closet peer is:")
		}
		peers = nPeers
	}

	return nil, nil
}

// Find close peers from a node
func (node *Node) findClosePeers(address string, target identity.Address) (*rpc.MultiAddresses, error) {
	client := connectNode(address)
	multiAddresses, err := client.CloserPeers(context.Background(), &rpc.Path{
		From: &rpc.Node{Address: string(node.DHT.Address)}, To: &rpc.Node{Address: string(target)}})
	if err != nil {
		return nil, err
	}
	return multiAddresses, nil
}
