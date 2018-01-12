package x

import (
	"fmt"
	"net"

	"github.com/republicprotocol/go-dht"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x/rpc"
	"golang.org/x/net/context"
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
	OnOrderFragmentReceived()
}

// Node implements the gRPC Node service.
type Node struct {
	*grpc.Server
	Delegate
	MultiAddress identity.MultiAddress
	DHT          *dht.DHT
}

// NewNode returns a Node with the given Config, a new DHT, and a new set of grpc.Connections.
func NewNode(multi identity.MultiAddress, multis identity.MultiAddresses, delegate Delegate) (*Node, error) {
	address, err := multi.Address()
	if err != nil {
		return nil, err
	}
	dht := dht.NewDHT(address)
	for _, multi := range multis {
		if err := dht.Update(multi); err != nil {
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
	return node.Server.Serve(listener)
}

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, sender *rpc.Nothing) (*rpc.MultiAddresses, error) {
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		wait <- node.handlePeers()
	}()

	select {
	case ret := <-wait:
		return ret, nil

	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddresses{}, ctx.Err()
	}
}

func (node *Node) handlePeers() *rpc.MultiAddresses {
	// Get all identity.MultiAddresses in the DHT.
	multis := node.DHT.MultiAddresses()
	ret := &rpc.MultiAddresses{
		Multis: make([]*rpc.MultiAddress, len(multis)),
	}
	// Transform them into rpc.MultiAddresses.
	for i, multi := range multis {
		ret.Multis[i] = &rpc.MultiAddress{Multi: multi.String()}
	}
	return ret
}

func (node *Node) pruneMostRecentPeer(target identity.Address) (bool, error) {
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
