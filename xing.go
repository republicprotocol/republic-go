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
func SendOrderFragmentToTarget(to identity.MultiAddress, orderFragment compute.OrderFragment, timeout time.Duration) error {
	conn, err := Dial(to, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = client.SendOrderFragment(ctx, SerializeOrderFragment(&orderFragment), grpc.FailFast(false))
	return err
}

// SendResultFragmentToTarget using a new grpc.ClientConn to make a
// SendResultFragment RPC to a targetMultiAddress.
func SendResultFragmentToTarget(to identity.MultiAddress, resultFragment compute.ResultFragment, timeout time.Duration) error {
	conn, err := Dial(to, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = client.SendResultFragment(ctx, SerializeResultFragment(&resultFragment), grpc.FailFast(false))
	return err
}

// SendTradingToTarget using a new grpc.ClientConn to make a
// SendTradingToTarget RPC to a targetMultiAddress.
func SendTradingToTarget(to identity.MultiAddress, tradingAtom struct{}, timeout time.Duration) error {
	conn, err := Dial(to, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = client.SendTradingAtom(ctx, SerializeTradingAtom(&struct{}{}), grpc.FailFast(false))
	return err
}
