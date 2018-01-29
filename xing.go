package rpc

import (
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// SendOrderFragmentToTarget using a new grpc.ClientConn to make a
// SendOrderFragment RPC to a targetMultiAddress.
func SendOrderFragmentToTarget(to identity.MultiAddress, orderFragment compute.OrderFragment) error {
	conn, err := Dial(to)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err = client.SendOrderFragment(ctx, SerializeOrderFragment(&orderFragment), grpc.FailFast(false))
	return err
}

// SendResultFragmentToTarget using a new grpc.ClientConn to make a
// SendResultFragment RPC to a targetMultiAddress.
func SendResultFragmentToTarget(to identity.MultiAddress, resultFragment compute.ResultFragment) error {
	conn, err := Dial(to)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = client.SendResultFragment(ctx, SerializeResultFragment(&resultFragment), grpc.FailFast(false))
	return err
}
