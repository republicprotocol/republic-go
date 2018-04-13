package rpc

import (
	"google.golang.org/grpc"
)

// SmpcService implements the Smpc gRPC service. SmpcService creates
// MessageQueues for each gRPC stream, and runs them on a Multiplexer. The
// closure of a gRPC stream, by the client or by the server, will prompt
// SmpcService to shutdown the respective MessageQueue.
type SmpcService struct {
	Options
	streams chan SmpcStream
}

// NewSmpcService returns a new SmpcService that will run MessageQueues on the
// given Dispatcher. The message queue limit is used used to buffer the size of
// the MessageQueues that are created by the SmpcService.
func NewSmpcService(options Options) SmpcService {
	return SmpcService{
		Options: options,
		streams: make(chan SmpcStream),
	}
}

// Streams published by the service.
func (service *SmpcService) Streams() <-chan SmpcStream {
	return service.streams
}

// CloseStreams publications. No new connections can be established after
// calling this method.
func (service *SmpcService) CloseStreams() {
	close(service.streams)
}

// Register the SmpcService with a gRPC server.
func (service *SmpcService) Register(server *grpc.Server) {
	RegisterSmpcServer(server, service)
}

// Compute opens a gRPC stream for streaming computation commands and results
// to the SmpcService.
func (service *SmpcService) Compute(stream Smpc_ComputeServer) error {
	//addr := "" // TODO: Get the address from a signed authentication message
	//writeCh := make(chan *SmpcMessage)
	//readCh := make(chan *SmpcMessage)
	//defer close(writeCh)
	//defer close(readCh)
	//
	//// Publish the message channel so that an external component can read
	//// messages from the stream
	//select {
	//case <-stream.Context().Done():
	//	return stream.Context().Err()
	//case service.messageChs <- messageCh:
	//}
	//
	//for {
	//	message, err := stream.Recv()
	//	if err != nil {
	//		return err
	//	}
	//
	//	select {
	//	case <-stream.Context().Done():
	//		return stream.Context().Err()
	//	case messageCh <- message:
	//	}
	//}
	return nil
}

type SmpcStream struct {
	Address string
	Write   chan<- *SmpcMessage
	Read    <-chan *SmpcMessage
}
