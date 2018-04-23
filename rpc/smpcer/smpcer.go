package smpcer

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ErrUnauthorized is returned when an unauthorized client makes a call to
// Compute.
var ErrUnauthorized = errors.New("unauthorized")

// A Conn is a collection of channels for sending, and receiving,
// ComputeMessages, and errors, between a client and an Smpc service.
type Conn struct {
	Sender   chan *ComputeMessage
	Receiver chan *ComputeMessage
}

// Delegate processing of OpenOrderRequests and CancelOrderRequests.
type Delegate interface {
	OpenOrder(signature []byte, orderFragment *OrderFragment) error
	CancelOrder(signature []byte, orderID []byte) error
}

type Smpcer struct {
	delegate Delegate

	connsMu *sync.Mutex
	connsRc map[string]int
	conns   map[string]Conn
}

func NewSmpcer(delegate Delegate) Smpcer {
	return Smpcer{
		delegate: delegate,

		connsMu: new(sync.Mutex),
		connsRc: map[string]int{},
		conns:   map[string]Conn{},
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

	conn := smpcer.AcquireConnection(addr)
	defer smpcer.ReleaseConnection(addr)

	// The done channel will signal to the sender goroutine that it should
	// exit
	done := make(chan struct{})
	defer close(done)

	senderErrs := make(chan error, 1)
	go func() {
		defer close(senderErrs)

		for {
			select {
			case <-done:
				// When the receiver exits the done channel will be closed and
				// this goroutine will eventually exit
				return
			case <-stream.Context().Done():
				senderErrs <- stream.Context().Err()
				return
			case computeMessage, ok := <-conn.Sender:
				if !ok {
					return
				}
				if err := stream.Send(computeMessage); err != nil {
					senderErrs <- err
					return
				}
			}
		}
	}()

	// Receive messages from the client until the context is done, or an error
	// is received
	for {
		message, err := stream.Recv()
		if err != nil {
			return err
		}

		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case err, ok := <-senderErrs:
			if !ok {
				// When the sender exits this error channel will be closed and
				// this goroutine will eventually exit
				return nil
			}
			return err
		case conn.Receiver <- message:
		}
	}
}

// AcquireConnection to an address. If a Connection to that address already
// exits, then it will be returned. Otherwise, a new Connection will be created
// and returned.
func (smpcer *Smpcer) AcquireConnection(addr string) Conn {
	smpcer.connsMu.Lock()
	defer smpcer.connsMu.Unlock()

	if smpcer.connsRc[addr] == 0 {
		// Setup a new connection
		smpcer.conns[addr] = Conn{
			Sender:   make(chan *ComputeMessage),
			Receiver: make(chan *ComputeMessage),
		}
	}

	smpcer.connsRc[addr]++
	return smpcer.conns[addr]
}

// ReleaseConnection to an address. If no other references to the Connection
// exist, then the Connection will be closed and deleted.
func (smpcer *Smpcer) ReleaseConnection(addr string) {
	smpcer.connsMu.Lock()
	defer smpcer.connsMu.Unlock()

	smpcer.connsRc[addr]--
	if smpcer.connsRc[addr] == 0 {
		// Teardown an exited connection
		close(smpcer.conns[addr].Sender)
		close(smpcer.conns[addr].Receiver)
		delete(smpcer.conns, addr)
	}
}
