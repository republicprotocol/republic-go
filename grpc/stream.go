package grpc

import (
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stream"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

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
	addr, err := service.verifyStreamAddress(message)
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
func (service *StreamService) Listen(ctx context.Context, addr identity.Address) (stream.Stream, error) {
	streams := service.setupConn(addr)
	defer service.teardownConn(addr)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case stream := <-streams:
		return stream, nil
	}
}

func (service *StreamService) verifyStreamAddress(message *StreamMessage) (identity.Address, error) {
	var addr string
	if message.GetStreamAddress() != nil && message.GetStreamAddress().GetAddress() != "" {
		addr = message.GetStreamAddress().GetAddress()
	}

	// FIXME: Verify that this message was signed by the sender.
	data := []byte("Republic Protocol: connect: from " + addr + " to " + service.addr.String())
	data = crypto.Keccak256(data)

	return identity.Address(addr), nil
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
	addr     identity.Address
	connPool *ConnPool

	streamsMu *sync.Mutex
	streams   map[identity.Address]StreamService_ConnectClient
}

func NewStreamClient(addr identity.Address, connPool *ConnPool) stream.Client {
	return &streamClient{
		addr:     addr,
		connPool: connPool,

		streamsMu: new(sync.Mutex),
		streams:   map[identity.Address]StreamService_ConnectClient{},
	}
}

func (client *streamClient) Connect(ctx context.Context, multiAddr identity.MultiAddress) (stream.Stream, error) {
	conn, err := client.connPool.Dial(ctx, multiAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot dial %v: %v", multiAddr, err)
	}
	defer conn.Close()

	stream, err := NewStreamServiceClient(conn.ClientConn).Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot open stream: %v", err)
	}
	if err := stream.Send(&StreamMessage{
		StreamAddress: &StreamAddress{
			Signature: []byte{},
			Address:   client.addr.String(),
		},
		Data: []byte{},
	}); err != nil {
		return nil, fmt.Errorf("cannot send stream address: %v", err)
	}

	client.streamsMu.Lock()
	defer client.streamsMu.Unlock()

	client.streams[multiAddr.Address()] = stream
	return newSafeStream(stream), err
}

func (client *streamClient) Disconnect(multiAddr identity.MultiAddress) error {
	client.streamsMu.Lock()
	defer client.streamsMu.Unlock()

	if stream, ok := client.streams[multiAddr.Address()]; ok {
		err := stream.CloseSend()
		delete(client.streams, multiAddr.Address())
		return err
	}
	return nil
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
func (stream safeStream) Close() error {
	stream.sendMu.Lock()
	stream.recvMu.Lock()
	defer stream.sendMu.Unlock()
	defer stream.recvMu.Unlock()

	close(stream.done)
	return nil
}

// Send implements the stream.Stream interface.
func (stream safeStream) Send(message stream.Message) error {
	stream.sendMu.Lock()
	defer stream.sendMu.Unlock()

	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	return stream.stream.SendMsg(&StreamMessage{
		Data: data,
	})
}

// Recv implements the stream.Stream interface.
func (stream safeStream) Recv(message stream.Message) error {
	stream.recvMu.Lock()
	defer stream.recvMu.Unlock()

	streamMessage := StreamMessage{}
	if err := stream.stream.RecvMsg(&streamMessage); err != nil {
		return err
	}
	return message.UnmarshalBinary(streamMessage.Data)
}
