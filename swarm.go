package rpc

import (
	"time"

	identity "github.com/republicprotocol/go-identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// PingTarget using a new grpc.ClientConn to make a Ping RPC to a target
// identity.MultiAddress.
func PingTarget(to identity.MultiAddress, from identity.MultiAddress) error {
	conn, err := Dial(to)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := client.Ping(ctx, from, grpc.FailFast(false))
	return err
}

// GetPeersFromTarget using a new grpc.ClientConn to make a Peers RPC to a
// target identity.MultiAddress.
func GetPeersFromTarget(to identity.MultiAddress, from identity.MultiAddress) (identity.MultiAddresses, error) {
	conn, err := Dial(to)
	if err != nil {
		return identity.MultiAddresses{}, nil
	}
	defer conn.Close()
	client := NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	multiAddresses, err := client.Peers(ctx, from, grpc.FailFast(false))
	if err != nil {
		return identity.MultiAddresses{}, nil
	}
	return DeserializeMultiAddresses(multiAddresses)
}

// QueryCloserPeersFromTarget using a new grpc.ClientConn to make a
// QueryCloserPeers RPC to a targetMultiAddress.
func QueryCloserPeersFromTarget(to *MultiAddress, from *MultiAddress, query *Address, deep bool) (*MultiAddresses, error) {
	conn, err := Dial(to)
	if err != nil {
		return &MultiAddresses{Multis: []*MultiAddress{}}, nil
	}
	defer conn.Close()
	client := NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.QueryCloserPeers(ctx, &Query{From: from, Query: query, Deep: deep}, grpc.FailFast(false))
}
