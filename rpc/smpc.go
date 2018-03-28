package rpc

import (
	"fmt"
	"io"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/smpc"
	"google.golang.org/grpc"
)

// SmpcService implements the Smpc gRPC service. SmpcService creates
// MessageQueues for each gRPC stream, and runs them on a Multiplexer. The
// closure of a gRPC stream, by the client or by the server, will prompt
// SmpcService to shutdown the respective MessageQueue.
type SmpcService struct {
	multiAddress      *identity.MultiAddress
	multiplexer       *dispatch.Multiplexer
	messageQueueLimit int
}

// NewSmpcService returns a new SmpcService that will run MessageQueues on the
// given Dispatcher. The message queue limit is used used to buffer the size of
// the MessageQueues that are created by the SmpcService.
func NewSmpcService(multiAddress *identity.MultiAddress, multiplexer *dispatch.Multiplexer, messageQueueLimit int) SmpcService {
	return SmpcService{
		multiAddress:      multiAddress,
		multiplexer:       multiplexer,
		messageQueueLimit: messageQueueLimit,
	}
}

// Register the SmpcService with a gRPC server.
func (service *SmpcService) Register(server *grpc.Server) {
	RegisterSmpcServer(server, service)
}

// Compute opens a gRPC stream for streaming computation commands and results
// to the SmpcService.
func (service *SmpcService) Compute(stream Smpc_ComputeServer) error {

	// Use a background MessageQueue to handle the connection until an error
	// is returned by the MessageQueue
	ch := make(chan error, 1)
	quit := make(chan struct{}, 1)
	go func() { ch <- service.connect(stream, quit) }()
	defer close(quit)

	// Select between the context finishing and the background worker
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case err := <-ch:
		return err
	}
}

func (service *SmpcService) connect(stream Smpc_ComputeServer, quit chan struct{}) error {

	// Send an identification message to the client
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- stream.Send(&SmpcMessage{
			MultiAddress: MarshalMultiAddress(service.multiAddress),
		})
	}()

	// Receive an identity message from the client
	message, err := stream.Recv()
	if err != nil {
		return err
	}
	if message.MultiAddress == nil {
		return fmt.Errorf("unverified identity: no signature available")
	}
	multiAddress, err := UnmarshalMultiAddress(message.MultiAddress)
	if err != nil {
		return err
	}

	// Wait for our identity to be sent, and their identity to be received
	if err, ok := <-ch; ok && err != nil {
		return err
	}

	// Create a MessageQueue that owns this gRPC stream and run it on the
	// Dispatcher
	messageQueue := NewSmpcServerStreamQueue(stream, service.messageQueueLimit)

	// Shutdown the MessageQueue when the quit signal is received
	go func() {
		<-quit
		messageQueue.Shutdown()
	}()
	return service.multiplexer.RunMessageQueue(multiAddress.Address().String(), &messageQueue)
}

// SmpcStreamQueue workers own a gRPC stream. All Message writing to, and
// reading from, a gRPC stream must go through a StreamQueue. The writing and
// reading channels used by the StreamQueue are abstracted away to prevent
// incorrectly reading from the write channel, and writing to the read channel.
type SmpcStreamQueue struct {
	stream grpc.Stream
	write  chan *SmpcMessage
	read   chan *SmpcMessage
	quit   chan struct{}
}

// NewSmpcClientStreamQueue returns a MessageQueue interface that is connected
// to an Smpc server. It accepts a gRPC stream, that will be owned by the
// MessageQueue.
func NewSmpcClientStreamQueue(stream Smpc_ComputeClient, messageQueueLimit int) SmpcStreamQueue {
	return SmpcStreamQueue{
		stream: stream,
		write:  make(chan *SmpcMessage, messageQueueLimit),
		read:   make(chan *SmpcMessage, messageQueueLimit),
		quit:   make(chan struct{}),
	}
}

// NewSmpcServerStreamQueue returns a MessageQueue interface that is connected
// to an Smpc client. It accepts a gRPC stream, that will be owned by the
// MessageQueue.
func NewSmpcServerStreamQueue(stream Smpc_ComputeServer, messageQueueLimit int) SmpcStreamQueue {
	return SmpcStreamQueue{
		stream: stream,
		write:  make(chan *SmpcMessage, messageQueueLimit),
		read:   make(chan *SmpcMessage, messageQueueLimit),
		quit:   make(chan struct{}),
	}
}

// Run the SmpcStreamQueue. This will concurrently process all writes to the
// underlying gRPC stream, and buffer all messages received over the gRPC
// stream.
func (queue *SmpcStreamQueue) Run() error {
	ch := make(chan error, 2)
	go func() { ch <- queue.writeAll() }()
	go func() { ch <- queue.readAll() }()
	for i := 0; i < 2; i++ {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

// Shutdown the SmpcStreamQueue. If it has already been Shutdown, and error
// will be returned.
func (queue *SmpcStreamQueue) Shutdown() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on shutdown: %v", r)
		}
	}()

	// If the stream is a client stream, then close the sending channel
	if stream, ok := queue.stream.(Smpc_ComputeClient); ok {
		// I hate that Go makes me do this :(
		err = stream.CloseSend()
	}
	close(queue.quit)

	return err
}

// Send a message to the SmpcStreamQueue. The Message must be a WorkerTask that
// wraps a gRPC SmpcMessage, otherwise an error is returned.
func (queue *SmpcStreamQueue) Send(message dispatch.Message) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on send: %v", r)
		}
	}()

	switch message := message.(type) {
	case smpc.Message:
		val := new(SmpcMessage)
		if message.DeltaFragments != nil {
			val.DeltaFragments = MarshalDeltaFragments(message.DeltaFragments)
		}
		// TODO: support other message formats
		queue.write <- val
	default:
		return fmt.Errorf("cannot send message: unrecognized type %T", message)
	}

	return err
}

// Recv a message from the SmpcStreamQueue. All Messages returned will be
// WorkerTasks that wrap a gRPC SmpcMessage.
func (queue *SmpcStreamQueue) Recv() (dispatch.Message, bool) {
	message, ok := <-queue.read
	if !ok {
		return message, ok
	}
	if message.DeltaFragments != nil {
		deltaFragments, err := UnmarshalDeltaFragments(message.DeltaFragments)
		if err != nil {
			return smpc.Message{Error: err}, true
		}
		return smpc.Message{DeltaFragments: deltaFragments}, true
	}

	// TODO: Handle other message types

	return smpc.Message{}, true
}

// writeAll messages from the messaging queue to the stream.
func (queue *SmpcStreamQueue) writeAll() error {
	for {
		select {
		case <-queue.quit:
			return nil
		case message := <-queue.write:
			if err := queue.stream.SendMsg(message); err != nil {
				return err
			}
		}
	}
}

// readAll messages from the stream and write them to the output queue. If the
// output queue is full, the SmpcStreamQueue will stop reading messages until
// the output is read from.
func (queue *SmpcStreamQueue) readAll() error {
	defer close(queue.read)
	for {
		message := new(SmpcMessage)
		if err := queue.stream.RecvMsg(message); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if message != nil {
			select {
			case <-queue.quit:
				return nil
			case queue.read <- message:
			}
		}
	}
}
