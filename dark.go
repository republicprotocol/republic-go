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
	client := NewDarkNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	serializedOrderFragment := SerializeOrderFragment(orderFragment)
	serializedOrderFragment.To = SerializeAddress(to)
	serializedOrderFragment.From = SerializeMultiAddress(from)
	_, err = client.SendOrderFragment(ctx, serializedOrderFragment, grpc.FailFast(false))
	return err
}

// SendOrderFragmentCommitmentToTarget using a new grpc.ClientConn to make a
// SendOrderFragmentCommitment RPC to a target identity.MultiAddress.
func SendOrderFragmentCommitmentToTarget(target identity.MultiAddress, from identity.MultiAddress, timeout time.Duration) error {
	conn, err := Dial(target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewDarkNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	commitment := &OrderFragmentCommitment{
		From:SerializeMultiAddress(from),
		Signature: []byte{}, // todo :
		OrderFragment: []byte{}, // todo :
	}
	_, err = client.SendOrderFragmentCommitment(ctx, commitment, grpc.FailFast(false))
	return err
}

// NotificationsFromTarget using a new grpc.ClientConn to make a
// NotificationsFromTarget RPC to a target identity.MultiAddress. This function
// returns two channels. The first should be used to read compute.Results until
// it is closed. The second should be closed by the caller when they no longer
// want to receive compute.Results. After closing the quit channel, the caller
// should read one last compute.Results from the first channel to guarantee
// that no memory is leaked.
func NotificationsFromTarget(target identity.MultiAddress, traderAddress identity.MultiAddress, timeout time.Duration) (chan do.Option, chan struct{}) {
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

		client := NewDarkNodeClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		stream, err := client.Notifications(ctx, SerializeMultiAddress(traderAddress), grpc.FailFast(false))
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
			ret <- do.Ok(DeserializeFinal(result))
		}
	}()

	return ret, quit
}

// GetResultsFromTarget using a new grpc.ClientConn to make a
// GetResultsFromTarget RPC to a target identity.MultiAddress.
func GetFinalsFromTarget(target identity.MultiAddress, traderAddress identity.MultiAddress, timeout time.Duration) ([]*compute.Final, error) {
	results := make([]*compute.Final, 0)
	conn, err := Dial(target, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := NewDarkNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := client.GetFinals(ctx, SerializeMultiAddress(traderAddress), grpc.FailFast(false))
	if err != nil {
		return nil, err
	}

	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, DeserializeFinal(result))
	}

	return results, nil
}
