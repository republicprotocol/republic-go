package rpc

import (
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// SendOrderFragmentToTarget using a new grpc.ClientConn to make a
// SendOrderFragment RPC to a target identity.MultiAddress.
func SendOrderFragmentToTarget(to identity.MultiAddress, from identity.MultiAddress, target identity.Address, orderFragment *compute.OrderFragment, timeout time.Duration) error {
	conn, err := Dial(to, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	serializedOrderFragment := SerializeOrderFragment(orderFragment)
	serializedOrderFragment.To = SerializeAddress(target)
	serializedOrderFragment.From = SerializeMultiAddress(from)
	_, err = client.SendOrderFragment(ctx, serializedOrderFragment, grpc.FailFast(false))
	return err
}

// SendResultFragmentToTarget using a new grpc.ClientConn to make a
// SendResultFragment RPC to a target identity.MultiAddress.
func SendResultFragmentToTarget(to identity.MultiAddress, from identity.MultiAddress, target identity.Address, resultFragment *compute.ResultFragment, timeout time.Duration) error {
	conn, err := Dial(to, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	serializedResultFragment := SerializeResultFragment(resultFragment)
	serializedResultFragment.To = SerializeAddress(target)
	serializedResultFragment.From = SerializeMultiAddress(from)
	_, err = client.SendResultFragment(ctx, serializedResultFragment, grpc.FailFast(false))
	return err
}
