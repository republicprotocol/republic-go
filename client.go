package swarm

import (
	"context"
	"fmt"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm/rpc"
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
	client, conn, err := NewNodeClient(target)
	if err != nil {
		return err
	}
	_, err = client.Ping(context.Background(), from)
	if err != nil {
		conn.Close()
		return err
	}
	return conn.Close()
}

// Send an rpc.Payload to the target identity.MultiAddress. Returns nil, or an
// error.
func Send(target identity.MultiAddress, payload *rpc.Payload) error {
	client, conn, err := NewNodeClient(target)
	if err != nil {
		return err
	}
	_, err = client.Send(context.Background(), payload)
	if err != nil {
		conn.Close()
		return err
	}
	return conn.Close()
}
