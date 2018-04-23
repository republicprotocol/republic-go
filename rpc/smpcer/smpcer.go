package smpcer

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// MaxConnections that can be serviced by an Smpc service.
const MaxConnections = 256

// ErrUnauthorized is returned when an unauthorized client makes a call to
// Compute.
var ErrUnauthorized = errors.New("unauthorized")

// ErrMaxConnectionsReached is returned when an Smpc service is servicing too
// many connection, or a client has attempted to connect more than once.
var ErrMaxConnectionsReached = errors.New("max connections reached")

// Delegate processing of OpenOrderRequests and CancelOrderRequests.
type Delegate interface {
	OpenOrder(signature []byte, orderFragment *OrderFragment) error
	CancelOrder(signature []byte, orderID []byte) error
}

type Smpcer struct {
	delegate Delegate
	client   Client

	connsMu *sync.Mutex
	conns   map[string]struct{}
}

func NewSmpcer(delegate Delegate, client Client) Smpcer {
	return Smpcer{
		delegate: delegate,
		client:   client,

		connsMu: new(sync.Mutex),
		conns:   map[string]struct{}{},
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
	if addr == "" || addr > smpcer.client.Address().String() {
		// The address cannot be empty and cannot be greater than the Smpc
		// service address (in such cases, the Smpc service is meant to call
		// the RPC)
		return ErrUnauthorized
	}
	// FIXME: Verify the client signature matches the address

	smpcer.connsMu.Lock()
	_, ok := smpcer.conns[addr]
	maxConnectionsReached := ok || len(smpcer.conns) >= MaxConnections
	smpcer.connsMu.Unlock()

	if maxConnectionsReached {
		return ErrMaxConnectionsReached
	}

	conn := smpcer.client.rendezvous.acquireConn(addr)
	defer smpcer.client.rendezvous.releaseConn(addr)

	return conn.stream(stream)
}
