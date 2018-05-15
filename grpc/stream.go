package grpc

import (
	"sync"

	"github.com/republicprotocol/republic-go/rpc"
	"google.golang.org/grpc"
)

// StreamService implements the rpc.SmpcServer interface using a gRPC service.
type StreamService struct {
	server rpc.StreamServer
}

func NewStreamService(server rpc.StreamServer) StreamService {
	return StreamService{
		server: server,
	}
}

// Register the gRPC service to a grpc.Server.
func (service *StreamService) Register(server *grpc.Server) {
	RegisterStreamServer(server, service)
}

// Connect implements the gRPC service for Smpc. It delegates control
// immediately to an rpc.SmpcServer interface.
func (service *StreamService) Connect(stream StreamService_ConnectServer) error {
	return service.server.ConnectFrom(stream.Context(), stream)
}

type StreamAdapter struct {
	sendMu *sync.Mutex
	recvMu *sync.Mutex
	stream StreamService_ConnectServer
}

// Send implements the rpc.Stream interface so that the gRPC stream can be used
// by other packages without depending on the underlying implementation of
// gRPC.
func (adapter StreamAdapter) Send(message rpc.StreamMessage) error {
	adapter.sendMu.Lock()
	defer adapter.sendMu.Unlock()

	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	return adapter.stream.Send(&SmpcMessage{
		Data: data,
	})
}

// Recv implements the rpc.SmpcStream interface so that the gRPC stream can be
// used by other packages without depending on the underlying implementation
// provided by gRPC.
func (adapter StreamAdapter) Recv(message rpc.StreamMessage) error {
	adapter.recvMu.Lock()
	defer adapter.recvMu.Unlock()

	data, err := adapter.stream.Recv()
	if err != nil {
		return err
	}
	return message.UnmarshalBinary(data.Data)
}
