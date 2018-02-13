package rpc

import (
	"io"
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// SendOrderFragmentToTarget using a new grpc.ClientConn to make a
// SendOrderFragment RPC to a target identity.MultiAddress.
func SendOrderFragmentToTarget(target identity.MultiAddress, to identity.Address, from identity.MultiAddress, orderFragment *compute.OrderFragment, timeout time.Duration) error {
	conn, err := Dial(target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	serializedOrderFragment := SerializeOrderFragment(orderFragment)
	serializedOrderFragment.To = SerializeAddress(to)
	serializedOrderFragment.From = SerializeMultiAddress(from)
	_, err = client.SendOrderFragment(ctx, serializedOrderFragment, grpc.FailFast(false))
	return err
}

// SendResultFragmentToTarget using a new grpc.ClientConn to make a
// SendResultFragment RPC to a target identity.MultiAddress.
func SendResultFragmentToTarget(target identity.MultiAddress, to identity.Address, from identity.MultiAddress, resultFragment *compute.ResultFragment, timeout time.Duration) error {
	conn, err := Dial(target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	serializedResultFragment := SerializeResultFragment(resultFragment)
	serializedResultFragment.To = SerializeAddress(to)
	serializedResultFragment.From = SerializeMultiAddress(from)
	_, err = client.SendResultFragment(ctx, serializedResultFragment, grpc.FailFast(false))
	return err
}

// NotificationsFromTarget using a new grpc.ClientConn to make a
// NotificationsFromTarget RPC to a target identity.MultiAddress. This function
// returns two channels. The first should be used to read compute.Results until
// it is closed. The second should be closed by the caller when they no longer
// want to receive compute.Results. After closing the quit channel, the caller
// should read one last compute.Results from the first channel to guarantee
// that no memory is leaked.
func NotificationsFromTarget(target identity.MultiAddress, traderAddress identity.Address, timeout time.Duration) (chan do.Option, chan struct{}) {
	ret := make(chan do.Option, 1)
	quit := make(chan struct{}, 1)

	go func() {
		defer close(ret)
		conn, err := Dial(target, timeout)
		if err != nil {
			ret <- do.Err(err)
			return
		}
		defer conn.Close()

		client := NewXingNodeClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		stream, err := client.Notifications(ctx, SerializeAddress(traderAddress), grpc.FailFast(false))
		if err != nil {
			ret <- do.Err(err)
			return
		}

		for {

			// Check if the quit channel is closed, without blocking.
			select {
			case _, ok := <-quit:
				if !ok {
					return
				}
			default:
				// Do nothing.
			}

			result, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				ret <- do.Err(err)
				continue
			}
			ret <- do.Ok(DeserializeResult(result))
		}
	}()

	return ret, quit
}

// GetResultsFromTarget using a new grpc.ClientConn to make a
// GetResults RPC to a target identity.MultiAddress.
func GetResultsFromTarget(target identity.MultiAddress, traderAddress identity.Address, timeout time.Duration) ([]*compute.Result, error) {
	conn, err := Dial(target, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewXingNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	results, err := client.GetResults(ctx, SerializeAddress(traderAddress), grpc.FailFast(false))
	ret := make([]*compute.Result, len(results.Result))
	for i, j := range results.Result {
		ret[i] = DeserializeResult(j)
	}
	return ret, nil
}
