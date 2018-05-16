package grpc

import (
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/stream"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ErrClosedConnectionWhileListening is returned when a channel, notifying
// listeners about new client connections, is closed while a listener is
// listening.
var ErrClosedConnectionWhileListening = errors.New("closed connection while listening")

// StreamService implements the rpc.SmpcServer interface using a gRPC service.
type StreamService struct {
	connsMu   *sync.Mutex
	conns     map[string]chan StreamService_ConnectServer
	connsDone map[string]chan (struct{})
}

// NewStreamService returns a gRP
func NewStreamService(server stream.Server) StreamService {
	return StreamService{
		server: server,
	}
}

// Register the gRPC service to a grpc.Server.
func (service *StreamService) Register(server *grpc.Server) {
	RegisterStreamServer(server, service)
}

// Connect implements the gRPC service for an abstract bidirectional stream of
// messages.
func (service *StreamService) Connect(stream StreamService_ConnectServer) error {
	defer func() {
		service.connsMu.Lock()
		defer service.connsMu.Unlock()
		delete(service.conns, addr.String())
		delete(service.connsDone, addr.String())
	}()

	service.connsMu.Lock()
	if _, ok := service.conns[addr.String()]; !ok {
		service.conns[addr.String()] = make(chan StreamService_ConnectServer, 1)
		service.connsDone[addr.String()] = make(chan (struct{}))
	}
	conn := service.conns[addr.String()]
	done := service.connsDone[addr.String()]
	service.connsMu.Unlock()

	select {
	case <-stream.Contetx().Done():
		return stream.Context().Err()
	case <-done:
		return nil
	case conn <- stream:
	}

	select {
	case <-stream.Contetx().Done():
		return stream.Context().Err()
	case <-done:
		return nil
	}
}

func (service *StreamService) Listen(ctx context.Context, addr identiy.Address) (stream.Stream, error) {
	defer func() {
		service.connsMu.Lock()
		defer service.connsMu.Unlock()
		delete(service.conns, addr.String())
		delete(service.connsDone, addr.String())
	}()

	service.connsMu.Lock()
	if _, ok := service.conns[addr.String()]; !ok {
		service.conns[addr.String()] = make(chan StreamService_ConnectServer, 1)
	}
	conn := service.conns[addr.String()]
	service.connsMu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case stream := <-conn:
		return newSafeStream(<-stream), nil
	}
}

// safeStream wraps a gRPC stream and ensures that it is safe for concurrent
// use. It prevents multiple goroutines from concurrent writing, and from
// concurrent reading, but it allows one goroutine to write while another
// goroutine is reading.
type safeStream struct {
	sendMu *sync.Mutex
	recvMu *sync.Mutex
	stream StreamService_ConnectServer
}

func newSafeStream(stream StreamService_ConnectServer) safeStream {
	return safeStream{
		sendMu: new(sync.Mutex),
		recvMu: new(sync.Mutex),
		stream: stream,
	}
}

// Send implements the stream.Stream interface.
func (stream safeStream) Send(message stream.Message) error {
	stream.sendMu.Lock()
	defer adapter.sendMu.Unlock()

	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	return adapter.stream.Send(&SmpcMessage{
		Data: data,
	})
}

// Recv implements the stream.Stream interface.
func (stream safeStream) Recv(message stream.Message) error {
	stream.recvMu.Lock()
	defer stream.recvMu.Unlock()

	data, err := stream.stream.Recv()
	if err != nil {
		return err
	}
	return message.UnmarshalBinary(data.Data)
}
