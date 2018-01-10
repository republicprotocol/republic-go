package x

import (
	"context"
	"fmt"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x/rpc"
	"google.golang.org/grpc"
)

// NewNodeClient returns a new rpc.NodeClient that is connected to the given
// target identity.MultiAddress, or an error. It uses a background
// context.Context.
func NewNodeClient(target identity.MultiAddress) (rpc.NodeClient, *grpc.ClientConn, error) {
	host, err := target.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, nil, err
	}
	port, err := target.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, nil, err
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return rpc.NewNodeClient(conn), conn, nil
}

// Ping the target identity.MultiAddress. Returns nil, or an error.
func Ping(target identity.MultiAddress, from *rpc.MultiAddress) error {
	// Create the client.
	client, conn, err := NewNodeClient(target)
	if err != nil {
		return err
	}
	defer conn.Close()
	// Ping.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err = client.Ping(ctx, from, grpc.FailFast(false))
	if err != nil {
		return err
	}
	return nil
}

// Peers asks for all peers connected to the node. Returns nil, or an error.
func Peers(target identity.MultiAddress) (identity.MultiAddresses,error) {
	// Create the client.
	client, conn, err := NewNodeClient(target)
	if err != nil {
		return nil,err
	}
	defer conn.Close()

	// Make grpc call
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	peers , err := client.Peers(ctx, &rpc.Nothing{}, grpc.FailFast(false))
	res := make([]identity.MultiAddress,len(peers.Multis))
	for index,peer :=range peers.Multis{
		multi,err := identity.NewMultiAddressFromString(peer.Multi)
		if err != nil {
			return nil, err
		}
		res[index] = multi
	}
	return res,nil
}


// Send an order fragment to the target identity.MultiAddress. Returns nil, or an
// error.
func SendOrderFragment(target identity.MultiAddress, fragment *rpc.OrderFragment) (*identity.MultiAddress,error) {
	// Create the client.
	client, conn, err := NewNodeClient(target)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Send the order fragment.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	response , err := client.SendOrderFragment(ctx, fragment, grpc.FailFast(false))
	if err != nil {
		return nil, err
	}
	multi,err := identity.NewMultiAddressFromString(response.Multi)
	if err != nil {
		return nil, err
	}
	return &multi, nil
}

