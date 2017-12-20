package swarm

import (
	"container/list"
	"errors"
	"github.com/multiformats/go-multiaddr"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm/dht"
	"github.com/republicprotocol/go-swarm/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

// Node implements the gRPC Node service.
type Node struct {
	DHT *dht.RoutingTable

	ip   string
	port string
}

// NewNode returns a new node with no connections.
func NewNode(ip, port string, address identity.Address) *Node {
	return &Node{DHT: dht.NewRoutingTable(address), ip: ip, port: port}
}

// StartListen starts the node as a grpc server and listens for rpc calls
func (node *Node) StartListen() error {
	// listen to the tcp port
	lis, err := net.Listen("tcp", node.ip+":"+node.port)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	rpc.RegisterDHTServer(s, node)

	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}

// Close stops the node grpc server
func (node *Node) Close() {

}

// MultiAddress returns the multiAddress of the node
func (node *Node) MultiAddress() (multiaddr.Multiaddr, error) {
	multi, err := identity.NewMultiaddr("/ip4/" + node.ip + "/tcp/" + node.port + "/republic/" + string(node.DHT.Address))
	if err != nil {
		return nil, err
	}
	return multi, nil
}

// Ping is used to test connections to a Node. The Node will respond with its
// address. If the Node does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, id *rpc.Node) (*rpc.Node, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.Node{Address: string(node.DHT.Address)}, err
	}

	// Update the sending node in the routing table
	if err := node.updateNode(id); err != nil {
		return nil, err
	}

	return &rpc.Node{Address: string(node.DHT.Address), Ip: node.ip, Port: node.port}, nil
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
	if err := node.updateNode(target); err != nil {
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
	if err := node.updateNode(path.From); err != nil {
		return nil, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		peers, _ := node.closerPeers(identity.Address(path.To))
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
	peers, err := node.DHT.FindNode(address)
	if err != nil {
		return nil, err
	}
	return &rpc.MultiAddresses{Multis: peers.MultiAddresses()}, nil
}

// Update the address in the routing table if there is enough space
func (node *Node) updateNode(newNode *rpc.Node) error {
	rAddress := identity.Address(newNode.Address)

	// Check if there still has place for the new address
	lastNode, err := node.DHT.CheckAvailability(rAddress)
	if err != nil {
		return err
	}

	// Generate multiaddress for the new node
	multiAddress, err := identity.NewMultiaddr("/ip4/" + newNode.Ip +
		"/tcp/" + newNode.Port + "/republic/" + newNode.Address)
	if err != nil {
		return err
	}

	// If the bucket is full
	if lastNode != "" {
		// Try to ping the last node in the bucket
		wait := make(chan interface{})
		go func() {
			defer close(wait)
			pong, err := node.PingNode(lastNode)
			if err != nil {
				wait <- err
			} else {
				wait <- pong
			}
		}()

		// Wait for the response
		resp := <-wait
		// If the last node is active, we do nothing with the new node
		switch resp.(type) {
		case error:
			return resp.(error)
		case *rpc.Node:
			return nil
		}
	}

	return node.DHT.Update(multiAddress)
}

// Connect to other Node and return the grpc client
func connectNode(ctx context.Context, address string) rpc.DHTClient {
	// Set up a connection to the server.
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	return rpc.NewDHTClient(conn)
}

// Ping a node
func (node *Node) PingNode(address string) (*rpc.Node, error) {
	client := connectNode(context.Background(), address)
	pong, err := client.Ping(context.Background(), &rpc.Node{Address: string(node.DHT.Address), Ip: node.ip, Port: node.port})
	if err != nil {
		return nil, err
	}
	multiAddress, err := identity.NewMultiaddr("/ip4/" + pong.Ip +
		"/tcp/" + pong.Port + "/republic/" + pong.Address)
	if err != nil {
		return nil, err
	}

	return pong, node.DHT.Update(multiAddress)
}

// Request all peers of a node
func (node *Node) PeersNode(address string) (*rpc.MultiAddresses, error) {
	client := connectNode(context.Background(), address)
	multiAddresses, err := client.Peers(context.Background(), &rpc.Node{Address: string(node.DHT.Address),
		Ip: node.ip, Port: node.port})
	if err != nil {
		return nil, err
	}
	return multiAddresses, nil
}

// Find a certain node by its address through the p2p network
// Return its multi-adress
func (node *Node) FindNode(target string) (*rpc.MultiAddresses, error) {
	// Find closest peers we know from the routing table
	peers, err := node.DHT.FindNode(identity.Address(target))
	if err != nil {
		return nil, err
	}

	for {
		if peers.Front() == nil {
			return nil, errors.New("empty routing table")
		}
		// Check if the first element is the target
		mostCloserNode, err := peers.Front().Value.(multiaddr.Multiaddr).ValueForProtocol(dht.RepublicCode)
		if err != nil {
			return nil, nil
		}
		if mostCloserNode == string(target) {
			return &rpc.MultiAddresses{Multis: []string{peers.Front().Value.(multiaddr.Multiaddr).String()}}, nil
		}
		// Ping each node in the bucket to find closer nodes
		wait := make(chan []string)
		for e := peers.Front(); e != nil; e = e.Next() {
			// todo: ignore the error for now
			ipAddress, err := e.Value.(multiaddr.Multiaddr).ValueForProtocol(identity.IP4Code)
			if err != nil {
				return nil, nil
			}
			port, err := e.Value.(multiaddr.Multiaddr).ValueForProtocol(identity.TCPCode)
			if err != nil {
				return nil, nil
			}

			client := connectNode(context.Background(), ipAddress+":"+port)
			path := &rpc.Path{From: &rpc.Node{Address: string(node.DHT.Address), Ip: node.ip, Port: node.port}, To: target}

			if e.Next() == nil {
				go func() {
					multiAddresses, _ := client.CloserPeers(context.Background(), path)
					wait <- multiAddresses.Multis
					close(wait)
				}()
			} else {
				go func() {
					multiAddresses, _ := client.CloserPeers(context.Background(), path)
					wait <- multiAddresses.Multis

				}()
			}

		}
		nPeers := dht.RoutingBucket{list.List{}}

		// Check if the node we get is closer than the first node
		for closerNodes := range wait {
			for _, address := range closerNodes {
				multi, err := identity.NewMultiaddr(address)
				if err != nil {
					return nil, err
				}
				nPeers.PushFront(multi)
			}
		}

		nPeers, err = dht.SortBucket(nPeers, identity.Address(target))

		if err != nil {
			return nil, err
		}

		newcloser, err := nPeers.Front().Value.(multiaddr.Multiaddr).ValueForProtocol(dht.RepublicCode)
		if err != nil {
			return nil, err
		}
		if newcloser == mostCloserNode {
			return nil, errors.New(" we can't find the target from the peers i know ")
		}

		peers = nPeers

	}

	return nil, nil
}
