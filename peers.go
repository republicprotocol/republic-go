package network

import (
	"time"

	identity "github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-network/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// RPCPeers sends a Peers RPC request to the target using a new grpc.ClientConn
// and a new rpc.NodeClient. It returns the result of the RPC call, or an error.
func (node *Node) RPCPeers(target identity.MultiAddress) (identity.MultiAddresses, error) {

	// Connect to the target.
	conn, err := Dial(target)
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	// Create a timeout context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Call the Peers RPC on the target.
	multis, err := client.Peers(ctx, &rpc.Nothing{}, grpc.FailFast(false))
	if err != nil {
		return identity.MultiAddresses{}, err
	}

	peers := make(identity.MultiAddresses, 0, len(multis.Multis))
	for _, multi := range multis.Multis {
		peer, err := identity.NewMultiAddressFromString(multi.Multi)
		if err != nil {
			return peers, err
		}
		peers = append(peers, peer)
	}
	return peers, nil
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
