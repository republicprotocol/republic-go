package x

import (
	"context"
	"fmt"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x/rpc"
	"google.golang.org/grpc"
)

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

// Peers asks for all peers connected to the node. Returns nil, or an error.
func Peers(target identity.MultiAddress) (identity.MultiAddresses, error) {
	// Create the client.
	client, conn, err := NewNodeClient(target)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Make grpc call
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	peers, err := client.Peers(ctx, &rpc.Nothing{}, grpc.FailFast(false))
	res := make([]identity.MultiAddress, len(peers.Multis))
	for index, peer := range peers.Multis {
		multi, err := identity.NewMultiAddress(peer.Multi)
		if err != nil {
			return nil, err
		}
		res[index] = multi
	}
	return res, nil
}

// Send an order fragment to the target identity.MultiAddress. Returns nil, or an
// error.
func SendOrderFragment(target identity.MultiAddress, fragment *rpc.OrderFragment) (*identity.MultiAddress, error) {
	// Create the client.
	client, conn, err := NewNodeClient(target)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Send the order fragment.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	response, err := client.SendOrderFragment(ctx, fragment, grpc.FailFast(false))
	if err != nil {
		return nil, err
	}
	multi, err := identity.NewMultiAddress(response.Multi)
	if err != nil {
		return nil, err
	}
	return &multi, nil
}

// CallPing sends an RPC request to the target using a new grpc.ClientConn and
// a new rpc.NodeClient. It returns the result of the RPC call, or an error.
func (node *Node) CallPing(target identity.MultiAddress) (*rpc.MultiAddress, error) {

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
	return client.Ping(ctx, &rpc.MultiAddress{Multi: node.MultiAddress.String()}, grpc.FailFast(false))
}
