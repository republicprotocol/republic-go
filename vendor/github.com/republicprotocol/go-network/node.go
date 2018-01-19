package network

import (
	"fmt"
	"log"
	"net"

	"github.com/republicprotocol/go-dht"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"github.com/republicprotocol/go-order-compute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// α determines the maximum number of concurrent client connections that the
// Node is expected to use when running a distributed Dijkstra search.
const α = 3

// Dial the target identity.MultiAddress using a background context.Context.
// Returns a grpc.ClientConn, or an error. The grpc.ClientConn must be closed
// before it exists scope.
func Dial(target identity.MultiAddress) (*grpc.ClientConn, error) {
	host, err := target.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := target.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// The Delegate is used to inject dependencies into the RPC logic.
type Delegate interface {
	OnPingReceived(peer identity.MultiAddress)
	OnOrderFragmentReceived(orderFragment *compute.OrderFragment)
	OnResultFragmentReceived(resultFragment *compute.ResultFragment)
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

func (node *Node) Prune(target identity.Address) (bool, error) {
	bucket, err := node.DHT.FindBucket(target)
	if err != nil {
		return false, err
	}
	multi := bucket.OldestMultiAddress()
	if multi == nil {
		return false, nil
	}
	newMulti, err := node.RPCPing(*multi)
	if err != nil {
		return true, node.DHT.Remove(*multi)
	}
	node.DHT.Update(*newMulti)
	return false, nil
}
