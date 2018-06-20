package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"google.golang.org/grpc"
)

// Dial creates a client connection to the given multiaddress. A context can be
// used to cancel or expire the pending connection. Once this function returns,
// the cancellation and expiration of the Context will do nothing. Users must
// call grpc.ClientConn.Close to terminate all the pending operations after
// this function returns.
func Dial(ctx context.Context, multiAddress identity.MultiAddress) (*grpc.ClientConn, error) {
	host, err := multiAddress.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := multiAddress.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, err
	}
	clientConn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return clientConn, nil
}

// Backoff a function call until the context.Context is done, or the function
// returns nil.
func Backoff(ctx context.Context, f func() error) error {
	timeoutMs := time.Duration(1000)
	for {
		err := f()
		if err == nil {
			return nil
		}
		timer := time.NewTimer(time.Millisecond * timeoutMs)
		select {
		case <-ctx.Done():
			return fmt.Errorf("backoff timeout = %v: %v", ctx.Err(), err)
		case <-timer.C:
			timeoutMs = time.Duration(float64(timeoutMs) * 1.6)
		}
	}
}

// BackoffMax is the same as Backoff but it will not wait longer than the
// maximum time in millisecond.
func BackoffMax(ctx context.Context, f func() error, maxMs int64) error {
	timeoutMs := time.Duration(1000)
	for {
		err := f()
		if err == nil {
			return nil
		}
		timer := time.NewTimer(time.Millisecond * timeoutMs)
		select {
		case <-ctx.Done():
			return fmt.Errorf("backoff timeout = %v: %v", ctx.Err(), err)
		case <-timer.C:
			timeoutMs = time.Duration(float64(timeoutMs) * 1.6)
			if int64(timeoutMs) > maxMs {
				timeoutMs = time.Duration(maxMs)
			}
		}
	}
}
