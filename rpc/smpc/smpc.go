package smpc

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ErrUnauthorized is returned when an unauthorized client makes a call to
// Compute.
var ErrUnauthorized = errors.New("unauthorized")

// A Connection is a collection of channels for sending, and receiving,
// messages, and errors, between a client and a service.
type Connection struct {
	Sender   chan *ComputeMessage
	Receiver chan *ComputeMessage
}

// A Delegate is responsible for processing OpenOrderRequests and
// CancelOrderRequests.
type Delegate interface {
	OpenOrder(signature []byte, orderFragment *OrderFragment) error
	CancelOrder(signature []byte, orderID []byte) error
}

type Smpc struct {
	delegate Delegate

	connsMu *sync.Mutex
	connsRc map[string]int
	conns   map[string]Connection
}

func NewSmpc(delegate Delegate) Smpc {
	return Smpc{
		delegate: delegate,

		connsMu: new(sync.Mutex),
		connsRc: map[string]int{},
		conns:   map[string]Connection{},
	}
}

// Register the gRPC service to a grpc.Server.
func (smpc *Smpc) Register(server *grpc.Server) {
	RegisterSmpcServer(server, smpc)
}

func (smpc *Smpc) OpenOrder(ctx context.Context, request *OpenOrderRequest) (*OpenOrderResponse, error) {

	// Delegate processing of the request in the background
	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		if err := smpc.delegate.OpenOrder(request.GetSignature(), request.GetOrderFragment()); err != nil {
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

func (smpc *Smpc) CancelOrder(ctx context.Context, request *CancelOrderRequest) (*CancelOrderResponse, error) {

	// Delegate processing of the request in the background
	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		if err := smpc.delegate.CancelOrder(request.GetSignature(), request.GetOrderId()); err != nil {
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

func (smpc *Smpc) Compute(stream Smpc_ComputeServer) error {
	auth, err := stream.Recv()
	if err != nil {
		return err
	}
	addr := auth.GetAddress()
	if addr == "" {
		return ErrUnauthorized
	}
	// TODO: Verify the client signature matches the address

	conn := smpc.AcquireConnection(addr)
	defer smpc.ReleaseConnection(addr)

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
func (smpc *Smpc) AcquireConnection(addr string) Connection {
	smpc.connsMu.Lock()
	defer smpc.connsMu.Unlock()

	if smpc.connsRc[addr] == 0 {
		// Setup a new connection
		smpc.conns[addr] = Connection{
			Sender:   make(chan *ComputeMessage),
			Receiver: make(chan *ComputeMessage),
		}
	}

	smpc.connsRc[addr]++
	return smpc.conns[addr]
}

// ReleaseConnection to an address. If no other references to the Connection
// exist, then the Connection will be closed and deleted.
func (smpc *Smpc) ReleaseConnection(addr string) {
	smpc.connsMu.Lock()
	defer smpc.connsMu.Unlock()

	smpc.connsRc[addr]--
	if smpc.connsRc[addr] == 0 {
		// Teardown an exited connection
		close(smpc.conns[addr].Sender)
		close(smpc.conns[addr].Receiver)
		delete(smpc.conns, addr)
	}
}
