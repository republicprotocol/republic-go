package swarm

import (
	"fmt"

	"net"

	"github.com/republicprotocol/go-dht"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Node implements the gRPC Node service.
type Node struct {
	KeyPair      identity.KeyPair
	MultiAddress identity.MultiAddress
	DHT          dht.DHT
}

// NewNode returns a Node with the given Config, a new DHT, and a new set of grpc.Connections.
func NewNode(config *Config) *Node {
	dht := dht.NewDHT(config.KeyPair.PublicAddress())
	for _, peer := range config.Peers {
		dht.Update(peer)
	}
	return &Node{
		KeyPair:      config.KeyPair,
		MultiAddress: config.MultiAddress,
		DHT:          dht,
	}
}

// StartListening starts a gRPC server.
func (node *Node) StartListening() error {
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
	server := grpc.NewServer()
	rpc.RegisterNodeServer(server, node)
	return server.Serve(listener)
}

// Ping is used to test the connection to the Node and exchange MultiAddresses.
// If the Node does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, peer *rpc.MultiAddress) (*rpc.MultiAddress, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		node.ping(peer)
	}()

	select {
	case <-wait:
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, nil

	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, ctx.Err()
	}
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
		wait <- node.peers()
	}()

	select {
	case ret := <-wait:
		return ret, nil

	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddresses{}, ctx.Err()
	}
}

// Send a Payload.
func (node *Node) Send(ctx context.Context, payload *rpc.Payload) (*rpc.Nothing, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.Nothing{}, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		wait <- node.send(payload)
	}()

	select {
	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.Nothing{}, ctx.Err()

	// Select the value passed by the goroutine.
	case ret := <-wait:
		return &rpc.Nothing{}, nil
	}
}

func (node *Node) ping(peer *rpc.MultiAddress) {
	node.DHT.Update(peer.Multi)
}

func (node *Node) peers() *rpc.MultiAddresses {
	return &rpc.MultiAddresses{}
}

func (node *Node) send() {
}
