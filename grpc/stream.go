package grpc

import (
	"errors"
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stream"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ErrUnverifiedConnection is returned when a StreamClient does not produce a
// verifiable connection signature as its first StreamMessage to a
// StreamServer.
var ErrUnverifiedConnection = errors.New("unverified connection")

// ErrNilAuthentication is returned when no authentication message is provided
// for a connection.
var ErrNilAuthentication = errors.New("nil authentication")

// StreamService implements the rpc.SmpcServer interface using a gRPC service.
type StreamService struct {
	verifier crypto.Verifier
	addr     identity.Address

	connsMu     *sync.Mutex
	connsStream map[identity.Address]chan safeStream
}

// NewStreamService returns an implementation of the stream.Server interface
// that uses gRPC for bidirectional streaming.
func NewStreamService(verifier crypto.Verifier, addr identity.Address) StreamService {
	return StreamService{
		verifier: verifier,
		addr:     addr,

		connsMu:     new(sync.Mutex),
		connsStream: map[identity.Address]chan safeStream{},
	}
}

// Register the StreamService to a grpc.Server.
func (service *StreamService) Register(server *grpc.Server) {
	RegisterStreamServiceServer(server, service)
}

// Connect implements the gRPC service for an abstract bidirectional stream of
// messages.
func (service *StreamService) Connect(stream StreamService_ConnectServer) error {

	// Verify the stream address of this connection
	message, err := stream.Recv()
	if err != nil {
		return err
	}
	addr, err := service.verifyAuthentication(message.GetAuthentication())
	if err != nil {
		return err
	}

	streams := service.setupConn(addr)
	defer service.teardownConn(addr)

	// Send the stream to the listener or exit when the context is done
	s := newSafeStream(stream)
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case streams <- s:
	}

	// Wait for an error response from the listener or exit when the context is
	// done
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case <-s.done:
		return nil
	}
}

// Listen implements the stream.Server interface.
func (service *StreamService) Listen(ctx context.Context, addr identity.Address) (stream.CloseStream, error) {
	streams := service.setupConn(addr)
	defer service.teardownConn(addr)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case stream := <-streams:
		return stream, nil
	}
}

func (service *StreamService) verifyAuthentication(auth *StreamAuthentication) (identity.Address, error) {
	if auth == nil || auth.GetAddress() == "" || auth.GetSignature() == nil {
		return identity.Address(""), ErrNilAuthentication
	}

	addr := auth.GetAddress()
	data := []byte(fmt.Sprintf("Republic Protocol: connect: from %v to %v", addr, service.addr))
	data = crypto.Keccak256(data)
	signature := auth.GetSignature()

	return identity.Address(addr), crypto.NewEcdsaVerifier(addr).Verify(data, signature)
}

func (service *StreamService) setupConn(addr identity.Address) chan safeStream {
	service.connsMu.Lock()
	defer service.connsMu.Unlock()

	if _, ok := service.connsStream[addr]; !ok {
		service.connsStream[addr] = make(chan safeStream, 1)
	}
	return service.connsStream[addr]
}

func (service *StreamService) teardownConn(addr identity.Address) {
	service.connsMu.Lock()
	defer service.connsMu.Unlock()

	delete(service.connsStream, addr)
}

type streamClient struct {
	signer crypto.Signer
	addr   identity.Address

	connPool *ConnPool
}

func NewStreamClient(signer crypto.Signer, addr identity.Address, connPool *ConnPool) stream.Client {
	return &streamClient{
		signer: signer,
		addr:   addr,

		connPool: connPool,
	}
}

func (client *streamClient) Connect(ctx context.Context, multiAddr identity.MultiAddress) (stream.CloseStream, error) {
	conn, err := client.connPool.Dial(ctx, multiAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot dial %v: %v", multiAddr, err)
	}
	defer conn.Close()

	data := []byte(fmt.Sprintf("Republic Protocol: connect: from %v to %v", client.addr, multiAddr.Address()))
	dataSignature, err := client.signer.Sign(data)
	if err != nil {
		return nil, fmt.Errorf("cannot sign stream authentication: %v", err)
	}

	stream, err := NewStreamServiceClient(conn.ClientConn).Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot open stream: %v", err)
	}
	if err := stream.Send(&StreamMessage{
		Authentication: &StreamAuthentication{
			Signature: dataSignature,
			Address:   client.addr.String(),
		},
		Data: []byte{},
	}); err != nil {
		return nil, fmt.Errorf("cannot send stream address: %v", err)
	}
	return newSafeStream(stream), err
}

// safeStream wraps a gRPC stream and ensures that it is safe for concurrent
// use. It prevents multiple goroutines from concurrent writing, and from
// concurrent reading, but it allows one goroutine to write while another
// goroutine is reading.
type safeStream struct {
	done   chan struct{}
	sendMu *sync.Mutex
	recvMu *sync.Mutex
	stream grpc.Stream
}

func newSafeStream(stream grpc.Stream) safeStream {
	return safeStream{
		done:   make(chan struct{}),
		sendMu: new(sync.Mutex),
		recvMu: new(sync.Mutex),
		stream: stream,
	}
}

// Close implements the stream.Stream interface.
func (safeStream safeStream) Close() error {
	safeStream.sendMu.Lock()
	safeStream.recvMu.Lock()
	defer safeStream.sendMu.Unlock()
	defer safeStream.recvMu.Unlock()

	if safeStream.done == nil {
		return stream.ErrCloseOnClosedStream
	}

	close(safeStream.done)
	safeStream.done = nil

	if stream, ok := safeStream.stream.(StreamService_ConnectClient); ok {
		return stream.CloseSend()
	}
	return nil
}

// Send implements the stream.Stream interface.
func (safeStream safeStream) Send(message stream.Message) error {
	safeStream.sendMu.Lock()
	defer safeStream.sendMu.Unlock()

	if safeStream.done == nil {
		return stream.ErrSendOnClosedStream
	}

	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	return safeStream.stream.SendMsg(&StreamMessage{
		Data: data,
	})
}

// Recv implements the stream.Stream interface.
func (safeStream safeStream) Recv(message stream.Message) error {
	safeStream.recvMu.Lock()
	defer safeStream.recvMu.Unlock()

	if safeStream.done == nil {
		return stream.ErrRecvOnClosedStream
	}

	data := StreamMessage{}
	if err := safeStream.stream.RecvMsg(&data); err != nil {
		return err
	}
	return message.UnmarshalBinary(data.Data)
}
