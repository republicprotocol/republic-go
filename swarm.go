package rpc

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// PingTarget using a new grpc.ClientConn to make a Ping RPC to a target
// MultiAddress.
func PingTarget(to *MultiAddress, from *MultiAddress) (*Nothing, error) {
	conn, err := Dial(to)
	if err != nil {
		return &Nothing{}, err
	}
	defer conn.Close()
	client := NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.Ping(ctx, from, grpc.FailFast(false))
}

// GetPeersFromTarget using a new grpc.ClientConn to make a Peers RPC to a
// targetMultiAddress.
func GetPeersFromTarget(to *MultiAddress, from *MultiAddress) (*MultiAddresses, error) {
	conn, err := Dial(to)
	if err != nil {
		return &MultiAddresses{Multis: []*MultiAddress{}}, nil
	}
	defer conn.Close()
	client := NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.Peers(ctx, from, grpc.FailFast(false))
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
