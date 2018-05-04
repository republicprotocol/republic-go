package smpcer

import (
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
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
	OpenOrder(signature []byte, orderFragment order.Fragment) error
	CancelOrder(signature []byte, orderID order.ID) error
}

type Smpcer struct {
	delegate Delegate
	client   *Client

	connsMu *sync.Mutex
	conns   map[string]struct{}
}

func NewSmpcer(client *Client, delegate Delegate) Smpcer {
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

		orderFragmentSignature := request.GetSignature()
		orderFragment, err := UnmarshalOrderFragment(smpcer.client.crypter, request.GetOrderFragment())
		if err != nil {
			errs <- err
			return
		}
		if err := smpcer.client.crypter.Verify(&orderFragment, orderFragmentSignature); err != nil {
			errs <- err
			return
		}
		if err := smpcer.delegate.OpenOrder(orderFragmentSignature, orderFragment); err != nil {
			errs <- err
			return
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

		orderIDSignature := request.GetSignature()
		orderID := order.ID(request.GetOrderId())
		if err := smpcer.client.crypter.Verify(orderID, orderIDSignature); err != nil {
			errs <- err
			return
		}
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

	smpcer.connsMu.Lock()
	smpcer.conns[addr] = struct{}{}
	smpcer.connsMu.Unlock()
	defer func() {
		smpcer.connsMu.Lock()
		delete(smpcer.conns, addr)
		smpcer.connsMu.Unlock()
	}()

	done := make(chan struct{})
	defer close(done)

	receiver := make(chan interface{})
	defer close(receiver)

	sender := smpcer.client.rendezvous.connect(identity.Address(addr), done, receiver)

	// Read all messages from the sender channel and write them to the gRPC
	// stream
	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		for {
			select {
			case <-done:
				// The receiver has terminated
				return
			case <-stream.Context().Done():
				// Writing to the error channel will cause the receiver to
				// terminate
				errs <- stream.Context().Err()
				return
			case message, ok := <-sender:
				if !ok {
					return
				}
				if err := stream.Send(message); err != nil {
					errs <- err
					return
				}
			}
		}
	}()

	// Read all messages from the gRPC stream and write them to the receiver
	// channel
	for {
		message, err := stream.Recv()
		if err != nil {
			return err
		}
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case err, ok := <-errs:
			// The sender has terminated, possibly with an error that should be
			// returned
			if !ok {
				return nil
			}
			return err
		case receiver <- message:
		}
	}
}
