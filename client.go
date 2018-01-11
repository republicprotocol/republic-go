package x

import (
	"context"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x/rpc"
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

// RPCSendOrderFragment sends a SendOrderFragment RPC request to the target
// using a new grpc.ClientConn and a new rpc.NodeClient. It returns the result
// of the RPC call, or an error.
func (node *Node) RPCSendOrderFragment(target identity.MultiAddress, fragment *rpc.OrderFragment) (*identity.MultiAddress, error) {
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

	// Call the SendOrderFragment RPC on the target.
	multi, err := client.SendOrderFragment(ctx, fragment, grpc.FailFast(false))
	if err != nil {
		return nil, err
	}

	ret, err := identity.NewMultiAddressFromString(multi.Multi)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
