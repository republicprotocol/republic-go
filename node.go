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
	"net"
	"google.golang.org/grpc/peer"
	"fmt"
	"log"
)

// Node implements the gRPC Node service.
type Node struct {
	IP   net.IP
	Port int
	DHT  *dht.RoutingTable
}

// NewNode returns a new node with no connections.
func NewNode(ip net.IP, port int, address identity.Address) *Node {
	return &Node{
		IP:    ip,
		Port:  port,
		DHT:   dht.NewRoutingTable(address),
		}
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
	if err := node.updateNode(ctx, id); err != nil {
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
	if err := node.updateNode(ctx, target); err != nil {
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
	if err := node.updateNode(ctx, path.From); err != nil {
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
	peers, err := node.DHT.FindNode(address)
	if err != nil {
		return nil, err
	}
	return &rpc.MultiAddresses{Multis: peers.MultiAddresses()}, nil
}

// Update the address in the routing table if there is enough space
func (node *Node) updateNode(ctx context.Context, address *rpc.Node) error {
	rAddress := identity.Address(address.Address)

	// Check if there still has place for the new address
	lastNode, err := node.DHT.CheckAvailability(rAddress)
	if err != nil {
		return err
	}

	// todo : generate the multiaddress from its republic address and ip infor.
	mAddress, err := rAddress.MultiAddress()
	if err != nil {
		return err
	}
	p, ok := peer.FromContext(ctx)
	if !ok || p==nil {
		return fmt.Errorf("can't get node network infor from the context")
	}
	log.Println("peer is : "+ p.Addr.String())

	tcpAddress,err  := net.ResolveTCPAddr("tcp", p.Addr.String())
	if err != nil {
		return err
	}
	networkAddress, err := multiaddr.NewMultiaddr("/ip4/"+tcpAddress.IP.String()+"/tcp/"+string(tcpAddress.Port))
	if err != nil {
		return err
	}
	mAddress = multiaddr.Join(networkAddress,mAddress)

	// If the bucket is full
	if lastNode != "" {
		// Try to ping the last node in the bucket
		wait := make(chan interface{})
		go func() {
			defer close(wait)
			pong, err := node.PingNode(lastNode)
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

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:node.IP,
			Port:node.Port,
		},
	}
	log.Println("peer who is pinging :"+p.Addr.String())
	ctx := peer.NewContext(context.Background(),p)
	client := connectNode(ctx, address)
	pong, err := client.Ping(ctx, &rpc.Node{Address: string(node.DHT.Address)})
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
	client := connectNode(context.Background(),address)
	multiAddresses, err := client.Peers(context.Background(), &rpc.Node{Address: string(node.DHT.Address)})
	if err != nil {
		return nil, err
	}
	return multiAddresses, nil
}

// Find a certain node by its ID through the kademlia network
// Return its multi-adresses
func (node *Node) FindNode(target identity.Address) (multiaddr.Multiaddr, error) {
	// Find closest peers we know from the routing table
	peers, err := node.DHT.FindNode(target)
	if err != nil {
		return nil, err
	}

	path := &rpc.Path{From: &rpc.Node{Address: string(node.DHT.Address)}, To: &rpc.Node{Address: string(target)}}

	for {
		// Check if the first element is the target
		peerValue, err := peers.Front().Value.(multiaddr.Multiaddr).ValueForProtocol(dht.RepublicCode)
		if err != nil {
			return nil, nil
		}
		if peers.Front() != nil && peerValue == string(target) {
			return peers.Front().Value.(multiaddr.Multiaddr), nil
		}

		// Ping each node and try to get closer nodes
		nPeers := dht.RoutingBucket{list.List{}}

		// todo : parallel  and get address from multi address
		for e := peers.Front(); e != nil; e = e.Next() {
			ipAddress, err := e.Value.(multiaddr.Multiaddr).ValueForProtocol(identity.IP4Code)
			if err != nil {
				return nil,err
			}
			port, err := e.Value.(multiaddr.Multiaddr).ValueForProtocol(identity.TCPCode)
			if err != nil {
				return nil,err
			}

			client := connectNode(context.Background(),ipAddress+":"+port)

			multiAddresses, err := client.CloserPeers(context.Background(), path)

			if err != nil {
				return nil, err
			}
			for _, address := range multiAddresses.Multis {
				nPeers.PushBack(address)
			}
		}
		nPeers, err = dht.SortBucket(nPeers,target)
		if err != nil {
			return nil, err
		}

		// Check if the new peers are unchanged from the previous list
		// which means we can't get closer peers from the node we know
		if dht.CompareList(nPeers, peers) {
			return nil, errors.New("can't find target from peers I know, closet peer is: "+peers.Front().Value.(multiaddr.Multiaddr).String())
		}

		peers = nPeers
	}

	return nil, nil
}

// Find close peers from a node
func (node *Node) findClosePeers(address string, target identity.Address) (*rpc.MultiAddresses, error) {
	client := connectNode(context.Background(), address)
	multiAddresses, err := client.CloserPeers(context.Background(), &rpc.Path{
		From: &rpc.Node{Address: string(node.DHT.Address)}, To: &rpc.Node{Address: string(target)}})
	if err != nil {
		return nil, err
	}
	return multiAddresses, nil
}
