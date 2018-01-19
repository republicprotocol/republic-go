package network

import (
	"log"
	"time"

	"github.com/republicprotocol/go-dht"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// RPCPing sends a Ping RPC request to the target using a new grpc.ClientConn
// and a new rpc.NodeClient. It returns the result of the RPC call, or an
// error.
func (node *Node) RPCPing(target identity.MultiAddress) (*identity.MultiAddress, error) {

	// Connect to the target.
	conn, err := Dial(target)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	// Create a timeout context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Call the Ping RPC on the target.
	multi, err := client.Ping(ctx, &rpc.MultiAddress{Multi: node.MultiAddress.String()}, grpc.FailFast(false))
	if err != nil {
		return nil, err
	}

	ret, err := identity.NewMultiAddressFromString(multi.Multi)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// Ping is used to test the connection to the Node and exchange MultiAddresses.
// If the Node does not respond, or it responds with an error, then the
// connection is considered unhealthy.
func (node *Node) Ping(ctx context.Context, peer *rpc.MultiAddress) (*rpc.MultiAddress, error) {
	log.Println("received call to Ping")

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan error)
	go func() {
		defer close(wait)
		wait <- node.handlePing(peer)
	}()

	select {
	case ret := <-wait:
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, ret

	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, ctx.Err()
	}
}

func (node *Node) handlePing(peer *rpc.MultiAddress) error {
	multi, err := identity.NewMultiAddressFromString(peer.Multi)
	if err != nil {
		return err
	}
	node.Delegate.OnPingReceived(multi)

	// Attempt to update the DHT.
	err = node.DHT.Update(multi)
	if err == dht.ErrFullBucket {
		// If the DHT is full then try and prune disconnected peers.
		address, err := multi.Address()
		if err != nil {
			return err
		}
		pruned, err := node.Prune(address)
		if err != nil {
			return err
		}
		// If a peer was pruned, then update the DHT again.
		if pruned {
			return node.DHT.Update(multi)
		}
		return nil
	}
	return err
}
