package rpc

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// SendOrderFragmentToTarget using a new grpc.ClientConn to make a
// SendOrderFragment RPC to a targetMultiAddress.
func SendOrderFragmentToTarget(to *MultiAddress, orderFragment *OrderFragment) (*Nothing, error) {
	conn, err := Dial(to)
	if err != nil {
		return &Nothing{}, err
	}
	defer conn.Close()
	client := NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.SendOrderFragment(ctx, orderFragment, grpc.FailFast(false))
}

// SendResultFragmentToTarget using a new grpc.ClientConn to make a
// SendResultFragment RPC to a targetMultiAddress.
func SendResultFragmentToTarget(to *MultiAddress, resultFragment *ResultFragment) (*Nothing, error) {
	conn, err := Dial(to)
	if err != nil {
		return &Nothing{}, err
	}
	defer conn.Close()
	client := NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return client.SendResultFragment(ctx, resultFragment, grpc.FailFast(false))
}
