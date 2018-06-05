package grpc

import (
	"errors"
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
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

// safeStream wraps a grpc.Stream and ensures that it is safe for concurrent
// use. It prevents multiple concurrent writes, and multiple concurrent reads.
// It allows a concurrent writer and reader.
type safeStream struct {
	grpc.Stream

	done   chan struct{}
	sendMu *sync.Mutex
	recvMu *sync.Mutex
}

func newSafeStream(stream grpc.Stream) safeStream {
	return safeStream{
		Stream: stream,

		done:   make(chan struct{}),
		sendMu: new(sync.Mutex),
		recvMu: new(sync.Mutex),
	}
}

// Send implements the stream.Stream interface.
func (str safeStream) Send(message stream.Message) error {
	str.sendMu.Lock()
	defer str.sendMu.Unlock()

	if str.done == nil {
		return stream.ErrSendOnClosedStream
	}

	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	return str.SendMsg(&StreamMessage{
		Data: data,
	})
}

// Recv implements the stream.Stream interface.
func (str safeStream) Recv(message stream.Message) error {
	str.recvMu.Lock()
	defer str.recvMu.Unlock()

	if str.done == nil {
		return stream.ErrRecvOnClosedStream
	}

	data := StreamMessage{}
	if err := str.RecvMsg(&data); err != nil {
		return err
	}
	return message.UnmarshalBinary(data.Data)
}

// Close the stream.
func (str safeStream) Close() error {
	str.sendMu.Lock()
	str.recvMu.Lock()
	defer str.sendMu.Unlock()
	defer str.recvMu.Unlock()

	if str.done == nil {
		return nil
	}
	close(str.done)
	str.done = nil

	if grpcClientStream, ok := str.Stream.(grpc.ClientStream); ok {
		return grpcClientStream.CloseSend()
	}
	return nil
}

// clientStream wraps a safeStream and exposes a method to close the associated
// grpc.ClientConn that is needed when creating a grpc.ClientStream.
type clientStream struct {
	safeStream

	conn *Conn
}

func newClientStream(stream grpc.ClientStream, conn *Conn) clientStream {
	return clientStream{
		safeStream: newSafeStream(stream),

		conn: conn,
	}
}

// Close the stream.
func (str clientStream) Close() error {
	if err := str.safeStream.Close(); err != nil {
		return err
	}
	return str.conn.Close()
}

type streamClient struct {
	signer crypto.Signer
	addr   identity.Address
}

// NewStreamClient implements the stream.Client interface using a gRPC
// bidirectional stream. It accepts a crypto.Signer used to authenticate
// connections requests with the StreamService, and an identity.Address to
// identify itself.
func NewStreamClient(signer crypto.Signer, addr identity.Address) stream.Client {
	return &streamClient{
		signer: signer,
		addr:   addr,
	}
}

func (client *streamClient) Connect(ctx context.Context, multiAddr identity.MultiAddress) (stream.Stream, error) {
	// Establish a connection to the identity.MultiAddress
	conn, err := Dial(ctx, multiAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot dial %v: %v", multiAddr, err)
	}

	// Open a bidirectional stream and continue to backoff the connection
	// until the context.Context is canceled
	var stream StreamService_ConnectClient
	if err := Backoff(ctx, func() error {
		stream, err = NewStreamServiceClient(conn.ClientConn).Connect(ctx)
		return err
	}); err != nil {
		return nil, fmt.Errorf("cannot open stream: %v", err)
	}

	// Sign an authentication message so that the StreamService can verify that
	// the identity.Address of the StreamClient
	data := []byte(fmt.Sprintf("Republic Protocol: connect: from %v to %v", client.addr, multiAddr.Address()))
	data = crypto.Keccak256(data)
	dataSignature, err := client.signer.Sign(data)
	if err != nil {
		return nil, fmt.Errorf("cannot sign stream authentication: %v", err)
	}

	// Send the authentication message
	if err := stream.Send(&StreamMessage{
		Authentication: &StreamAuthentication{
			Signature: dataSignature,
			Address:   client.addr.String(),
		},
		Data: []byte{},
	}); err != nil {
		return nil, fmt.Errorf("cannot send stream address: %v", err)
	}

	// Return a grpc.ClientStream that implements the stream.Stream interface
	// and is safe for concurrent use and will clean the grpc.ClientConn when
	// it is no longer needed
	clientStream := newClientStream(stream, conn)
	go func() {
		<-ctx.Done()
		clientStream.Close()
	}()
	return clientStream, err
}

// StreamService implements the gRPC StreamService. It implements the
// stream.Server interface by forwarding connections from the
// StreamService.Connect RPC when calls to StreamService.Listen happen.
type StreamService struct {
	verifier crypto.Verifier
	addr     identity.Address

	connsMu *sync.Mutex
	conns   map[identity.Address]chan safeStream
}

// NewStreamService returns an implementation of the stream.Server interface
// that uses gRPC for bidirectional streaming.
func NewStreamService(verifier crypto.Verifier, addr identity.Address) StreamService {
	return StreamService{
		verifier: verifier,
		addr:     addr,

		connsMu: new(sync.Mutex),
		conns:   map[identity.Address]chan safeStream{},
	}
}

// Register the StreamService to a Server.
func (service *StreamService) Register(server *Server) {
	RegisterStreamServiceServer(server.Server, service)
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
		logger.Network(logger.LevelDebugLow, "grpc connection closed by client")
		return stream.Context().Err()
	case <-s.done:
		logger.Network(logger.LevelDebugLow, "grpc connection closed by service")
		return nil
	}
}

// Listen implements the stream.Server interface.
func (service *StreamService) Listen(ctx context.Context, addr identity.Address) (stream.Stream, error) {
	streams := service.setupConn(addr)
	defer service.teardownConn(addr)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case stream := <-streams:
		go func() {
			<-ctx.Done()
			stream.Close()
		}()
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

	return identity.Address(addr), service.verifier.Verify(data, signature)
}

func (service *StreamService) setupConn(addr identity.Address) chan safeStream {
	service.connsMu.Lock()
	defer service.connsMu.Unlock()

	if _, ok := service.conns[addr]; !ok {
		service.conns[addr] = make(chan safeStream, 1)
	}
	return service.conns[addr]
}

func (service *StreamService) teardownConn(addr identity.Address) {
	service.connsMu.Lock()
	defer service.connsMu.Unlock()

	delete(service.conns, addr)
}
