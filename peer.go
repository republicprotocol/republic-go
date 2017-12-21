package swarm

import (
	"errors"
	"fmt"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm/dht"
	"github.com/republicprotocol/go-swarm/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"log"
)

// Peer implements the gRPC Peer service.
type Peer struct {
	Config 		 *Config
	DHT 		 *dht.RoutingTable
	Connections  map[string]*grpc.ClientConn
}

// NewPeer returns a Peer with the given Config, a new DHT, and a new set of grpc.Connections.
func NewPeer(config *Config) *Peer {
	rt := dht.NewRoutingTable(config.KeyPair.PublicAddress())
	for _,peer := range config.Peers{
		rt.Update(peer)
	}

	return &Peer{
		Config: 	 config,
		DHT: 		 rt,
		Connections: map[string]*grpc.ClientConn{},
	}
}

// StartListen starts the peer as a grpc server and listens for rpc calls
func (peer *Peer) StartListening() error {
	// listen to the tcp port
	host, err := peer.Config.MultiAddress.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return err
	}
	port, err := peer.Config.MultiAddress.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host,port))
	if err != nil {
		return err
	}

	// Register the peer as a grpc server
	s := grpc.NewServer()
	rpc.RegisterPeerServer(s, peer)
	return s.Serve(lis)
}

// Ping is used to test connections to a Peer. The Peer will respond with its
// address. If the Peer does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (peer *Peer) Ping(ctx context.Context, multi *rpc.MultiAddress) (*rpc.MultiAddress, error) {

	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddress{Multi: peer.Config.MultiAddress.String()}, err
	}

	// Update the sending Peer in the routing table
	if err := peer.updatePeer(multi); err != nil {
		return nil, err
	}

	return &rpc.MultiAddress{Multi: peer.Config.MultiAddress.String()}, nil
}

// Peers is used to return the rpc.MultiAddresses to which a Peer is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (peer *Peer) Peers(ctx context.Context, sender *rpc.Nothing) (*rpc.MultiAddresses, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		wait <- peer.peers()
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

// SendOrderFragment is the order fragment handler of the
func (peer *Peer) SendOrderFragment(ctx context.Context, fragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	return &rpc.Nothing{}, nil
}

func (peer *Peer) SendFragment(target identity.MultiAddress) (*rpc.Nothing ,error ){
	// Resolve the target network address from it multiAddress
	host, err := target.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, nil
	}
	port, err := target.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, nil
	}
	address := host+ ":" + port

	// Establish a connection with the address
	err = peer.ConnectNode(address)
	if err != nil {
		return nil,err
	}
	client := rpc.NewPeerClient(peer.Connections[address])
	defer peer.CloseConnection(address)

	return client.SendOrderFragment(context.Background(),&rpc.OrderFragment{To:&rpc.MultiAddress{Multi:target.String()}})
}

// Return all peers in the node routing table
func (peer *Peer) peers() *rpc.MultiAddresses {
	multis := peer.DHT.MultiAddresses()
	ret := make([]*rpc.MultiAddress, len(multis))
	for i, j:= range multis{
		ret[i] = &rpc.MultiAddress{ Multi: j.String() }
	}
	return &rpc.MultiAddresses{Multis:ret}
}

// Update the address in the routing table if there is enough space
func (peer *Peer) updatePeer(newPeer *rpc.MultiAddress) error {

	// Get the republic address of the new peer
	mAddress, err:= identity.NewMultiAddress(newPeer.Multi)
	if err != nil {
		return err
	}
	republicAddress,err  := mAddress.ValueForProtocol(identity.RepublicCode)
	if err != nil {
		return err
	}

	// Check if there still has place for the new address
	lastPeer, err := peer.DHT.CheckAvailability(identity.Address(republicAddress))
	if err != nil {
		return err
	}

	// If the bucket is full
	if lastPeer != identity.EmptyMultiAddress{
		// Try to ping the last peer in the bucket
		wait := make(chan interface{})
		go func() {
			defer close(wait)

			ip, err := lastPeer.ValueForProtocol(identity.IP4Code)
			if err != nil {
				wait <- err
			}
			port, err :=lastPeer.ValueForProtocol(identity.TCPCode)
			if err != nil {
				wait <- err
			}

			pong, err := peer.PingPeer(ip+":"+port)
			if err != nil {
				wait <- err
			} else {
				wait <- pong
			}
		}()

		// Wait for the response
		resp := <-wait
		// If the last peer is active, we do nothing with the new node
		switch resp.(type) {
		case error:
			// todo : remove the last node
			return resp.(error)
		case *rpc.MultiAddress:
			return nil
		}
	}

	return peer.DHT.Update(mAddress)
}

// ConnectNode connects to other Peer via its network address
func (peer *Peer) ConnectNode(address string) error {

	// Set up a connection to the server.
	conn ,err := grpc.Dial(address, grpc.WithInsecure())
	peer.Connections[address] = conn
	if err != nil {
		return err
	}
	return nil
}

// CloseConnection closes a connection a peer via its network address
func (peer *Peer) CloseConnection(address string) error {
	conn := peer.Connections[address]
	err :=  conn.Close()
	if err !=nil {
		return err
	}
	delete(peer.Connections, address)
	return nil
}

// PingNode sends a ping another Peer via its network address
func (peer *Peer) PingPeer(address string) (*rpc.MultiAddress, error) {

	// Establish a connection with the address
	err := peer.ConnectNode(address)
	if err != nil {
		return nil,err
	}
	client := rpc.NewPeerClient(peer.Connections[address])
	defer peer.CloseConnection(address)

	pong, err := client.Ping(context.Background(), &rpc.MultiAddress{Multi:peer.Config.MultiAddress.String()})
	if err != nil {
		return nil, err
	}

	respAddress, err := identity.NewMultiAddress(pong.Multi)
	if err != nil {
		return nil, err
	}
	return pong, peer.DHT.Update(respAddress)
}

// Request all peers of a peer knows
func (peer *Peer) AskPeers(address string) (*rpc.MultiAddresses, error) {

	// Establish a connection with the address
	err := peer.ConnectNode(address)
	if err != nil {
		return nil,err
	}
	client := rpc.NewPeerClient(peer.Connections[address])
	defer peer.CloseConnection(address)

	multiAddresses, err := client.Peers(context.Background(), &rpc.Nothing{})
	if err != nil {
		return nil, err
	}
	return multiAddresses, nil
}

// Find a certain node by its republic address through the p2p network
// Return its multiAdress
func (peer *Peer) FindPeer(target identity.Address) (*rpc.MultiAddresses, error) {
	log.Println("start finding ")
	// Find closest peers we know from the routing table
	peers, err := peer.DHT.FindPeer(identity.Address(target))
	if err != nil {
		return nil, err
	}
	visited := map[identity.MultiAddress]bool{peer.Config.MultiAddress:true}

	for {
		// Check if we know any peers
		if peers.Front() == nil {
			return nil, errors.New("can't find the target from the known peers")
		}
		log.Println("start finding ")
		// Check if we find the peer
		for e:= peers.Front(); e!= nil ;e = e.Next(){
			rAddress, err := e.Value.(identity.MultiAddress).ValueForProtocol(identity.RepublicCode)
			log.Println("we have "+ rAddress)
			if err != nil {
				return nil, err
			}
			if rAddress == string(target){
				return &rpc.MultiAddresses{Multis:[]*rpc.MultiAddress{{Multi:e.Value.(identity.MultiAddress).String()}}},nil
			}
		}

		// Ping the first node from the bucket
		first := peers.Remove(peers.Front()).(identity.MultiAddress)
		_, in :=  visited[first]
		if in {
			continue
		}

		// Get the network address of the peer
		host, err := first.ValueForProtocol(identity.IP4Code)
		if err != nil {
			return nil, nil
		}
		port, err := first.ValueForProtocol(identity.TCPCode)
		if err != nil {
			return nil, nil
		}

		multiAddresses, err := peer.AskPeers(host+":"+port)
		if err != nil {
			return nil, err
		}
		for _,j := range multiAddresses.Multis{
			mAddress, err := identity.NewMultiAddress(j.Multi)
			if err != nil {
				return nil, err
			}
			peers.PushBack(mAddress)
		}
		visited[first] = true

		// todo : need to figure out how to sort the bucket
		peers,err  = dht.SortBucket(peers,target)
		if err!= nil {
			return nil, err
		}
	}

	return nil, nil
}

// CloserPeers returns the peers of an rpc.Node that are closer to a target
// than the rpc.Node itself. Distance is calculated by evaluating a XOR with
// the target address and each peer ID.
//func (peer *Peer) CloserPeers(ctx context.Context, ) (*rpc.MultiAddresses, error) {
//	// Check for errors in the context.
//	if err := ctx.Err(); err != nil {
//		return &rpc.MultiAddresses{}, err
//	}
//
//	// Update the sender in the node routing table
//	if err := node.updateNode(path.From); err != nil {
//		return nil, err
//	}
//
//	// Spawn a goroutine to evaluate the return value.
//	wait := make(chan *rpc.MultiAddresses)
//	go func() {
//		defer close(wait)
//		peers, _ := node.closerPeers(identity.Address(path.To))
//		wait <- peers
//	}()
//
//	select {
//	// Select the timeout from the context.
//	case <-ctx.Done():
//		return &rpc.MultiAddresses{}, ctx.Err()
//
//	// Select the value passed by the goroutine.
//	case ret := <-wait:
//		return ret, nil
//	}
//}

//// Return the closer peers in the node routing table
//func (node *Node) closerPeers(address identity.Address) (*rpc.MultiAddresses, error) {
//	peers, err := node.DHT.FindNode(address)
//	if err != nil {
//		return nil, err
//	}
//	return &rpc.MultiAddresses{Multis: peers.MultiAddresses()}, nil
//}

//// Ask for
//func (node *Node) AskCloserNode(address ,target string)(*rpc.MultiAddresses, error ){
//	// Establish a connection with the address
//	err := node.ConnectNode(address)
//	if err != nil {
//		return nil,err
//	}
//	client := rpc.NewDHTClient(node.conn)
//	defer node.Close()
//
//	path := &rpc.Path{From: &rpc.Node{Address: string(node.DHT.Address), Ip: node.ip, Port: node.port}, To: target}
//	multiAddresses, err := client.CloserPeers(context.Background(), path)
//	if err != nil {
//		return nil, err
//	}
//	return multiAddresses, nil
//}
