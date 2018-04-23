package smpcer

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ErrUnauthorized is returned when an unauthorized client makes a call to
// Compute.
var ErrUnauthorized = errors.New("unauthorized")

// Delegate processing of OpenOrderRequests and CancelOrderRequests.
type Delegate interface {
	OpenOrder(signature []byte, orderFragment *OrderFragment) error
	CancelOrder(signature []byte, orderID []byte) error
}

type Smpcer struct {
	delegate Delegate
	client   Client
}

func NewSmpcer(delegate Delegate, client Client) Smpcer {
	return Smpcer{
		delegate: delegate,
		client:   client,
	}
}

// Register the gRPC service to a grpc.Server.
func (smpcer *Smpcer) Register(server *grpc.Server) {
	RegisterSmpcServer(server, smpcer)
}

func (smpcer *Smpcer) OpenOrder(ctx context.Context, request *OpenOrderRequest) (*OpenOrderResponse, error) {

	// Delegate processing of the request in the background
	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		if err := smpcer.delegate.OpenOrder(request.GetSignature(), request.GetOrderFragment()); err != nil {
			errs <- err
		}
	}()

	select {
	case <-ctx.Done():
		return &OpenOrderResponse{}, ctx.Err()
	case err, ok := <-errs:
		if !ok {
			return &OpenOrderResponse{}, nil
		}
		return &OpenOrderResponse{}, err
	}
}

func (smpcer *Smpcer) CancelOrder(ctx context.Context, request *CancelOrderRequest) (*CancelOrderResponse, error) {

	// Delegate processing of the request in the background
	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		if err := smpcer.delegate.CancelOrder(request.GetSignature(), request.GetOrderId()); err != nil {
			errs <- err
		}
	}()

	select {
	case <-ctx.Done():
		return &CancelOrderResponse{}, ctx.Err()
	case err, ok := <-errs:
		if !ok {
			return &CancelOrderResponse{}, nil
		}
		return &CancelOrderResponse{}, err
	}
}

func (smpcer *Smpcer) Compute(stream Smpc_ComputeServer) error {
	auth, err := stream.Recv()
	if err != nil {
		return err
	}
	addr := auth.GetAddress()
	if addr == "" {
		return ErrUnauthorized
	}
	// TODO: Verify the client signature matches the address

	rendezvous := smpcer.client.router.Acquire(addr)
	defer smpcer.client.router.Release(addr)

	return smpcer.client.mergeStreamAndRendezvous(stream, rendezvous)
}
