package network

import (
	"fmt"
	"log"
	"net"

	"github.com/republicprotocol/go-dht"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"github.com/republicprotocol/go-order-compute"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// The Delegate is used as a callback interface to inject logic into the
// different RPCs.
type Delegate interface {
	OnPingReceived(from identity.MultiAddress)
	OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment)
	OnResultFragmentReceived(from identity.MultiAddress, resultFragment *compute.ResultFragment)
}

// Node implements the gRPC Node service.
type Node struct {
	Delegate
	Server  *grpc.Server
	DHT     *dht.DHT
	Options Options
}

// NewNode returns a Node with the given its own identity.MultiAddress, a list
// of boostrap node identity.MultiAddresses, and a delegate that defines
// callbacks for each RPC.
func NewNode(delegate Delegate, options Options) *Node {
	return &Node{
		Delegate: delegate,
		Server:   grpc.NewServer(),
		DHT:      dht.NewDHT(options.MultiAddress.Address(), options.MaxBucketLength),
		Options:  options,
	}
}

// Serve starts the gRPC server.
func (node *Node) Serve() error {
	rpc.RegisterNodeServer(node.Server, node)

	host, err := node.MultiAddress.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return err
	}
	port, err := node.MultiAddress.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return err
	}

	if node.Options.Debug >= DebugLow {
		log.Printf("Listening at %s:%s", host, port)
	}
	return node.Server.Serve(listener)
}

// Stop the gRPC server.
func (node *Node) Stop() {
	node.Server.Stop()
}

// Prune an identity.Address from the dht.DHT. Returns a boolean indicating
// whether or not an identity.Address was pruned.
func (node *Node) Prune(target identity.Address) (bool, error) {
	bucket, err := node.DHT.FindBucket(target)
	if err != nil {
		return false, err
	}
	if bucket == nil || bucket.Length() == 0 {
		return false, nil
	}
	multiAddress := bucket.MultiAddresses[0]
	if _, err := PingTarget(SerializeMultiAddress(multiAddress), SerializeMultiAddress(node.MultiAddress())); err != nil {
		return true, node.DHT.RemoveMultiAddress(multiAddress)
	}
	return false, node.DHT.UpdateMultiAddress(multiAddress)
}

// Address returns the identity.Address of the Node.
func (node *Node) Address() identity.Address {
	return node.Options.MultiAddress.Address()
}

// MultiAddress returns the identity.MultiAddress of the Node.
func (node *Node) MultiAddress() identity.MultiAddress {
	return node.Options.MultiAddress
}

// Ping is used to test the connection to the Node and exchange
// identity.MultiAddresses. If the Node does not respond, or it responds with
// an error, then the connection should be considered unhealthy.
func (node *Node) Ping(ctx context.Context, from *rpc.MultiAddress) (*rpc.Nothing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.ping(from)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (node *Node) ping(from *rpc.MultiAddress) (*rpc.Nothing, error) {
	// Update the DHT.
	fromMultiAddress, err := DeserializeMultiAddress(from)
	if err != nil {
		return &rpc.Nothing{}, err
	}
	if err := node.updatePeer(fromMultiAddress); err != nil {
		return &rpc.Nothing{}, err
	}

	// Notify the delegate of the ping.
	node.Delegate.OnPingReceived(fromMultiAddress)
	return &rpc.Nothing{}, nil
}

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, from *rpc.MultiAddress) (*rpc.MultiAddresses, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		peers, err := node.peers(from)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(peers)
	})

	select {
	case val := <-wait:
		if multiAddresses, ok := val.Ok.(*rpc.MultiAddresses); ok {
			return multiAddresses, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, ctx.Err()
	}
}

func (node *Node) peers(from *rpc.MultiAddress) (*rpc.MultiAddresses, error) {
	// Update the DHT.
	fromMultiAddress, err := DeserializeMultiAddress(from)
	if err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}
	if err := node.updatePeer(fromMultiAddress); err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}

	// Return all peers in the DHT.
	peers := node.DHT.MultiAddresses()
	return SerializeMultiAddresses(peers), nil
}

// FindPeer is used to return the closest rpc.MultiAddresses to a peer with the
// given target rpc.Address. It will not return rpc.MultiAddresses that are
// further away from the target than the Node itself. The rpc.MultiAddresses
// returned are not guaranteed to provide healthy connections and should be
// pinged.
func (node *Node) FindPeer(ctx context.Context, finder *rpc.Finder) (*rpc.MultiAddresses, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		peers, err := node.findPeer(finder)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(peers)
	})

	select {
	case val := <-wait:
		if multiAddresses, ok := val.Ok.(*rpc.MultiAddresses); ok {
			return multiAddresses, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, ctx.Err()
	}
}

func (node *Node) findPeer(finder *rpc.Finder) (*rpc.MultiAddresses, error) {
	// Update the DHT.
	if finder.From != nil {
		fromMultiAddress, err := DeserializeMultiAddress(finder.From)
		if err != nil {
			return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, err
		}
		if err := node.updatePeer(fromMultiAddress); err != nil {
			return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, err
		}
	}

	// Get the target identity.Address for which this Node is searching for
	// peers.
	target := identity.Address(finder.Peer.Address)
	targetPeers := &rpc.MultiAddresses{Multis: make([]*rpc.MultiAddress, 0, node.Options.Alpha)}

	// Create the closed and open data structures for performing the Kademlia
	// search.
	peers, err := node.DHT.FindMultiAddressNeighbors(target, node.Options.Alpha)
	if err != nil {
		return targetPeers, err
	}

	// Filter away peers that are further from the target than this Node.
	for _, peer := range peers {
		peerAddress := peer.Address()
		if peerAddress == target {
			return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{SerializeMultiAddress(peer)}}, nil
		}
		closer, err := identity.Closer(peerAddress, node.Address(), target)
		if err != nil {
			return targetPeers, err
		}
		if closer {
			targetPeers.Multis = append(targetPeers.Multis, SerializeMultiAddress(peer))
		}
	}

	return targetPeers, nil
}

// SendOrderFragment to the Node. If the rpc.OrderFragment is not destined for
// this Node then it will be forwarded on to the correct destination.
func (node *Node) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.sendOrderFragment(orderFragment)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (node *Node) sendOrderFragment(orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	// Update the DHT.
	var err error
	var fromMultiAddress identity.MultiAddress
	if orderFragment.From != nil {
		fromMultiAddress, err := DeserializeMultiAddress(orderFragment.From)
		if err != nil {
			return &rpc.Nothing{}, err
		}
		if err := node.updatePeer(fromMultiAddress); err != nil {
			return &rpc.Nothing{}, err
		}
	}

	// Check if the rpc.OrderFragment has reached its destination.
	if orderFragment.To.Address == node.Address().String() {
		deserializedOrderFragment, err := DeserializeOrderFragment(orderFragment)
		if err != nil {
			return &rpc.Nothing{}, err
		}
		node.OnOrderFragmentReceived(fromMultiAddress, deserializedOrderFragment)
		return &rpc.Nothing{}, nil
	}

	// Forward the rpc.OrderFragment to the closest peers.
	peers, err := node.findPeer(&rpc.Finder{Peer: orderFragment.To, From: nil})
	if err != nil {
		return &rpc.Nothing{}, err
	}
	log.Println(peers.Multis, orderFragment.To)
	do.CoForAll(peers.Multis, func(i int) {
		SendOrderFragmentToTarget(peers.Multis[i], orderFragment)
	})
	return &rpc.Nothing{}, nil
}

// SendResultFragment to the Node. If the rpc.ResultFragment is not destined
// for this Node then it will be forwarded on to the correct destination.
func (node *Node) SendResultFragment(ctx context.Context, resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		nothing, err := node.sendResultFragment(resultFragment)
		if err != nil {
			return do.Err(err)
		}
		return do.Ok(nothing)
	})

	select {
	case val := <-wait:
		if nothing, ok := val.Ok.(*rpc.Nothing); ok {
			return nothing, val.Err
		}
		return nil, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (node *Node) sendResultFragment(resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	// Update the DHT.
	var err error
	var fromMultiAddress identity.MultiAddress
	if resultFragment.From != nil {
		fromMultiAddress, err = DeserializeMultiAddress(resultFragment.From)
		if err != nil {
			return &rpc.Nothing{}, err
		}
		if err := node.updatePeer(fromMultiAddress); err != nil {
			return &rpc.Nothing{}, err
		}
	}

	// Check if the rpc.OrderFragment has reached its destination.
	if resultFragment.To.String() == node.Address().String() {
		deserializedResultFragment, err := DeserializeResultFragment(resultFragment)
		if err != nil {
			return &rpc.Nothing{}, err
		}
		node.OnResultFragmentReceived(fromMultiAddress, deserializedResultFragment)
		return &rpc.Nothing{}, nil
	}

	// Forward the rpc.OrderFragment to the closest peers.
	peers, err := node.findPeer(&rpc.Finder{Peer: resultFragment.To, From: nil})
	if err != nil {
		return &rpc.Nothing{}, err
	}
	for _, peer := range peers.Multis {
		go SendResultFragmentToTarget(peer, resultFragment)
	}
	return &rpc.Nothing{}, nil
}

func (node *Node) updatePeer(multiAddress identity.MultiAddress) error {
	if err := node.DHT.UpdateMultiAddress(multiAddress); err != nil {
		if err == dht.ErrFullBucket {
			pruned, err := node.Prune(multiAddress.Address())
			if err != nil {
				return err
			}
			if pruned {
				return node.DHT.UpdateMultiAddress(multiAddress)
			}
			return nil
		}
		log.Println(err)
		return err
	}
	return nil
}
