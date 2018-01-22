package network

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/republicprotocol/go-dht"
	do "github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"github.com/republicprotocol/go-order-compute"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// α determines the maximum number of concurrent client connections that the
// Node is expected to use when running a distributed Dijkstra search.
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
	*grpc.Server
	Delegate
	MultiAddress identity.MultiAddress
	DHT          *dht.DHT
}

// NewNode returns a Node with the given Config, a new DHT, and a new set of grpc.Connections.
func NewNode(multi identity.MultiAddress, bootstrapMultis identity.MultiAddresses, delegate Delegate) (*Node, error) {
	address, err := multi.Address()
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
		Server:       grpc.NewServer(),
		Delegate:     delegate,
		MultiAddress: multi,
		DHT:          dht,
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
	reflection.Register(node.Server)
	log.Printf("Listening at %s:%s", host, port)
	return node.Server.Serve(listener)
}

// Prune an identity.Address from the dht.DHT. Returns a boolean indicating
//
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

// Ping is used to test the connection to the Node and exchange MultiAddresses.
// If the Node does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, from *rpc.MultiAddress) (*rpc.Nothing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ch := do.Process(func() do.Option {
		return node.ping(from)
	})
	defer close(ch)

	select {
	case nothing := <-ch:
		if nothing.Err != nil {
			return nil, nothing.Err
		}
		return &rpc.Nothing{}, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, from *rpc.MultiAddress) (*rpc.MultiAddresses, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ch := do.Process(func() do.Option {
		return node.peers(from)
	})
	defer close(ch)

	select {
	case multiAddressOpt := <-ch:
		if peers, ok := multiAddressOpt.Ok.(*rpc.MultiAddresses); ok {
			return peers, multiAddressOpt.Err
		}
		return nil, multiAddressOpt.Err

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (node *Node) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.Nothing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ch := do.Process(func() do.Option {
		return do.Err(errors.New("unimplemented"))
	})
	defer close(ch)

	select {
	case nothing := <-ch:
		if nothing.Err != nil {
			return nil, nothing.Err
		}
		return &rpc.Nothing{}, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (node *Node) SendResultFragment(ctx context.Context, resultFragment *rpc.ResultFragment) (*rpc.Nothing, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ch := do.Process(func() do.Option {
		return do.Err(errors.New("unimplemented"))
	})
	defer close(ch)

	select {
	case nothing := <-ch:
		if nothing.Err != nil {
			return nil, nothing.Err
		}
		return &rpc.Nothing{}, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (node *Node) ping(from *rpc.MultiAddress) do.Option {
	multiAddress, err := DeserializeMultiAddress(from)
	node.Delegate.OnPingReceived(multiAddress)

	if err := node.DHT.Update(multiAddress); err != nil {
		if err == dht.ErrFullBucket {
			address, err := multiAddress.Address()
			if err != nil {
				return do.Err(err)
			}
			pruned, err := node.Prune(address)
			if err != nil {
				return do.Err(err)
			}
			if pruned {
				if err := node.DHT.Update(multiAddress); err != nil {
					return do.Err(err)
				}
				return do.Ok(SerializeMultiAddress(node.MultiAddress))
			}
		}
		return do.Err(err)
	}

	return do.Ok(SerializeMultiAddress(node.MultiAddress))
}

func (node *Node) peers(from *rpc.MultiAddress) do.Option {
	peers := node.DHT.MultiAddresses()
	return do.Ok(SerializeMultiAddresses(peers))
}
