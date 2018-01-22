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

// α determines the maximum number of concurrent client connections that the
// Node is expected to use when routing messages.
const α = 3

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
	Server       *grpc.Server
	DHT          *dht.DHT
	Address      identity.Address
	MultiAddress identity.MultiAddress
}

// NewNode returns a Node with the given its own identity.MultiAddress, a list
// of boostrap node identity.MultiAddresses, and a delegate that defines
// callbacks for each RPC.
func NewNode(multiAddress identity.MultiAddress, bootstrapMultis identity.MultiAddresses, delegate Delegate) (*Node, error) {
	address, err := multiAddress.Address()
	if err != nil {
		return nil, err
	}
	dht := dht.NewDHT(address, 100)
	for _, bootstrapMulti := range bootstrapMultis {
		if err := dht.Update(bootstrapMulti); err != nil {
			return nil, err
		}
	}
	return &Node{
		Delegate:     delegate,
		Server:       grpc.NewServer(),
		DHT:          dht,
		Address:      address,
		MultiAddress: multiAddress,
	}, nil
}

// Serve starts the gRPC server.
func (node *Node) Serve() error {
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
	rpc.RegisterNodeServer(node.Server, node)
	log.Printf("Listening at %s:%s", host, port)
	return node.Server.Serve(listener)
}

// Prune an identity.Address from the dht.DHT. Returns a boolean indicating
// whether or not an identity.Address was pruned.
func (node *Node) Prune(target identity.Address) (bool, error) {
	bucket, err := node.DHT.FindBucket(target)
	if err != nil {
		return false, err
	}
	multi := bucket.OldestMultiAddress()
	if multi == nil {
		return false, nil
	}
	if err := PingTarget(*multi, node.MultiAddress); err != nil {
		return true, node.DHT.Remove(*multi)
	}
	return false, nil
}

// Ping is used to test the connection to the Node and exchange
// identity.MultiAddresses. If the Node does not respond, or it responds with
// an error, then the connection should be considered unhealthy.
func (node *Node) Ping(ctx context.Context, from *rpc.MultiAddress) (*rpc.Nothing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		return do.Err(node.ping(from))
	})

	select {
	case val := <-wait:
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
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
		return val.Ok.(*rpc.MultiAddresses), val.Err

	case <-ctx.Done():
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, ctx.Err()
	}
}

// SendOrderFragment to the Node. If the rpc.OrderFragment is not destined for
// this Node then it will be forwarded on to the correct destination.
func (node *Node) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		return do.Err(node.sendOrderFragment(orderFragment))
	})

	select {
	case val := <-wait:
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

// SendResultFragment to the Node. If the rpc.ResultFragment is not destined
// for this Node then it will be forwarded on to the correct destination.
func (node *Node) SendResultFragment(ctx context.Context, resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	wait := do.Process(func() do.Option {
		return do.Err(node.sendResultFragment(resultFragment))
	})

	select {
	case val := <-wait:
		return &rpc.Nothing{}, val.Err

	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()
	}
}

func (node *Node) ping(from *rpc.MultiAddress) error {
	// Update the DHT.
	fromMultiAddress, err := DeserializeMultiAddress(from)
	if err != nil {
		return err
	}
	if err := node.updateDHT(fromMultiAddress); err != nil {
		return err
	}

	// Notify the delegate of the ping.
	node.Delegate.OnPingReceived(fromMultiAddress)
	return nil
}

func (node *Node) peers(from *rpc.MultiAddress) (*rpc.MultiAddresses, error) {
	// Update the DHT.
	fromMultiAddress, err := DeserializeMultiAddress(from)
	if err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}
	if err := node.updateDHT(fromMultiAddress); err != nil {
		return &rpc.MultiAddresses{Multis: []*rpc.MultiAddress{}}, nil
	}

	// Return all peers in the DHT.
	peers := node.DHT.MultiAddresses()
	return SerializeMultiAddresses(peers), nil
}

func (node *Node) sendOrderFragment(orderFragment *rpc.OrderFragment) error {
	// Update the DHT.
	fromMultiAddress, err := DeserializeMultiAddress(orderFragment.From)
	if err != nil {
		return err
	}
	if err := node.updateDHT(fromMultiAddress); err != nil {
		return err
	}

	// Check if the rpc.OrderFragment has reached its destination.
	if orderFragment.To.String() == node.Address.String() {
		deserializedOrderFragment, err := DeserializeOrderFragment(orderFragment)
		if err != nil {
			return err
		}
		node.OnOrderFragmentReceived(fromMultiAddress, deserializedOrderFragment)
		return nil
	}

	peer, err := node.FindPeer(identity.Address(orderFragment.To.String()))
	if err != nil {
		return nil
	}
	return SendOrderFragmentToTarget(target, node.MultiAddress, orderFragment)
}

func (node *Node) sendResultFragment(resultFragment *rpc.ResultFragment) error {
	// Update the DHT.
	fromMultiAddress, err := DeserializeMultiAddress(resultFragment.From)
	if err != nil {
		return err
	}
	if err := node.updateDHT(fromMultiAddress); err != nil {
		return err
	}

	// Check if the rpc.OrderFragment has reached its destination.
	if resultFragment.To.String() == node.Address.String() {
		deserializedResultFragment, err := DeserializeResultFragment(resultFragment)
		if err != nil {
			return err
		}
		node.OnResultFragmentReceived(fromMultiAddress, deserializedResultFragment)
		return nil
	}

	// TODO: Find the closest target.
	return nil
}

func (node *Node) updateDHT(multiAddress identity.MultiAddress) error {
	if err := node.DHT.Update(multiAddress); err != nil {
		if err == dht.ErrFullBucket {
			address, err := multiAddress.Address()
			if err != nil {
				return err
			}
			pruned, err := node.Prune(address)
			if err != nil {
				return err
			}
			if pruned {
				return node.DHT.Update(multiAddress)
			}
			return nil
		}
		return err
	}
	return nil
}
