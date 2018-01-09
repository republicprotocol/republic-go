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
	*grpc.Server
	KeyPair      identity.KeyPair
	MultiAddress identity.MultiAddress
	DHT          dht.DHT
}

// NewNode returns a Node with the given Config, a new DHT, and a new set of grpc.Connections.
func NewNode(config *Config) (*Node, error) {
	dht := dht.NewDHT(config.KeyPair.Address())
	for _, peer := range config.Peers {
		if err := dht.Update(peer); err != nil {
			return nil, err
		}
	}
	return &Node{
		Server:       grpc.NewServer(),
		KeyPair:      config.KeyPair,
		MultiAddress: config.MultiAddress,
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
	return node.Server.Serve(listener)
}

// Ping is used to test the connection to the Node and exchange MultiAddresses.
// If the Node does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, peer *rpc.MultiAddress) (*rpc.MultiAddress, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan error)
	go func() {
		defer close(wait)
		wait <- node.ping(peer)
	}()

	select {
	case ret := <-wait:
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, ret

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
	wait := make(chan error)
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
		if ret != nil {
			return &rpc.Nothing{}, ret
		}
		return &rpc.Nothing{}, nil
	}
}

func (node *Node) ping(peer *rpc.MultiAddress) error {
	multi, err := identity.NewMultiAddress(peer.Multi)
	if err != nil {
		return err
	}
	err = node.DHT.Update(multi)
	if err == dht.ErrFullBucket {
		address, err := multi.Address()
		if err != nil {
			return err
		}
		pruned, err := node.pruneUnhealthyPeer(address)
		if err != nil {
			return err
		}
		if pruned {
			return node.DHT.Update(multi)
		}
		return nil
	}
	return err
}

func (node *Node) peers() *rpc.MultiAddresses {
	multis := node.DHT.MultiAddresses()
	ret := &rpc.MultiAddresses{
		Multis: make([]*rpc.MultiAddress, len(multis)),
	}
	for i, multi := range multis {
		ret.Multis[i] = &rpc.MultiAddress{Multi: multi.String()}
	}
	return ret
}

func (node *Node) send(payload *rpc.Payload) error {
	return nil
}

func (node *Node) pruneUnhealthyPeer(target identity.Address) (bool, error) {
	bucket, err := node.DHT.Bucket(target)
	if err != nil {
		return false, err
	}
	for i := len(*bucket) - 1; i >= 0; i-- {
		client, conn, err := NewNodeClient((*bucket)[i].MultiAddress)
		if err != nil {
			if err == context.DeadlineExceeded {
				// TODO: prune
				return true, nil
			}
			return false, err
		}
		defer conn.Close()
		_, err = client.Ping(context.Background(), &rpc.MultiAddress{Multi: node.MultiAddress.String()})
		if err != nil {
			if err == context.DeadlineExceeded {
				// TODO: prune
				return true, nil
			}
			return false, err
		}
	}
	return false, nil
}
