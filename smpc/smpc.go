package smpc

import (
	"fmt"
	"io"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"google.golang.org/grpc"
)

// Smpc implements the Smpc gRPC service. Smpc creates SmpcMessageQueues for
// each gRPC stream, and runs them on a Multiplexer. The closure of a gRPC
// stream, by the client or by the server, will prompt Smpc to shutdown the
// SmpcMessageQueue.
type Smpc struct {
	multiAddress      identity.MultiAddress
	multiplexer       *dispatch.Multiplexer
	messageQueueLimit int
}

// NewSmpc returns a new Smpc service that will run SmpcMessageQueues on the given
// Dispatcher. The message queue limit is used used to buffer the size of the
// SmpcMessageQueues that are created by the Smpc.
func NewSmpc(multiAddress identity.MultiAddress, multiplexer *dispatch.Multiplexer, messageQueueLimit int) Smpc {
	return Smpc{
		multiAddress:      multiAddress,
		multiplexer:       multiplexer,
		messageQueueLimit: messageQueueLimit,
	}
}

// Connect to the Smpc service and begin streaming requests and responses for
// the Smpc sMPC protocol.
func (λ *Smpc) Connect(stream rpc.Smpc_ConnectServer) error {

	// Use a background MessageQueue to handle the connection until an error
	// is returned by the MessageQueue
	ch := make(chan error, 1)
	quit := make(chan struct{}, 1)
	go func() { ch <- λ.connect(stream, quit) }()
	defer close(quit)

	// Select between the context finishing and the background worker
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case err := <-ch:
		return err
	}
}

func (λ *Smpc) connect(stream rpc.Smpc_ConnectServer, quit chan struct{}) error {

	// Send an identification message to the client
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- stream.Send(&rpc.SmpcMessage{
			MultiAddress: rpc.SerializeMultiAddress(λ.multiAddress),
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
	multiAddress, err := rpc.DeserializeMultiAddress(message.MultiAddress)
	if err != nil {
		return err
	}

	// Wait for our identity to be sent, and their identity to be received
	if err, ok := <-ch; ok && err != nil {
		return err
	}

	// Create a MessageQueue that owns this gRPC stream and run it on the
	// Dispatcher
	messageQueue := NewServerStreamQueue(stream, λ.messageQueueLimit)

	// Shutdown the MessageQueue when the quit signal is received
	go func() {
		<-quit
		messageQueue.Shutdown()
	}()
	return λ.multiplexer.RunMessageQueue(multiAddress.Address().String(), messageQueue)
}

// StreamQueue workers own a gRPC stream. All Message writing to, and reading
// from, a gRPC stream must go through a StreamQueue. The writing and reading
// channels used by the StreamQueue are abstracted away to prevent incorrectly
// reading from the write channel, and writing to the read channel.
type StreamQueue struct {
	stream grpc.Stream
	write  chan *rpc.SmpcMessage
	read   chan *rpc.SmpcMessage
	quit   chan struct{}
}

// NewClientStreamQueue returns a MessageQueue interface that is connected to a
// gRPC server. It accepts a gRPC stream, that will be owned by the
// MessageQueue and should not be used by any other component
func NewClientStreamQueue(stream rpc.Smpc_ConnectClient, messageQueueLimit int) dispatch.MessageQueue {
	return &StreamQueue{
		stream: stream,
		write:  make(chan *rpc.SmpcMessage, messageQueueLimit),
		read:   make(chan *rpc.SmpcMessage, messageQueueLimit),
		quit:   make(chan struct{}),
	}
}

// NewServerStreamQueue returns a MessageQueue interface that is connected to a
// gRPC client. It accepts a gRPC stream, that will be owned by the
// MessageQueue and should not be used by any other component
func NewServerStreamQueue(stream rpc.Smpc_ConnectServer, messageQueueLimit int) dispatch.MessageQueue {
	return &StreamQueue{
		stream: stream,
		write:  make(chan *rpc.SmpcMessage, messageQueueLimit),
		read:   make(chan *rpc.SmpcMessage, messageQueueLimit),
		quit:   make(chan struct{}),
	}
}

// Run the StreamQueue. This will concurrently process all writes to the
// underlying gRPC stream, and buffer all messages received over the gRPC
// stream.
func (queue *StreamQueue) Run() error {
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

// Shutdown the StreamQueue. If it has already been Shutdown, and error will be
// returned.
func (queue *StreamQueue) Shutdown() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on shutdown: %v", r)
		}
	}()

	// If the stream is a client stream, then close the sending channel
	if stream, ok := queue.stream.(rpc.Smpc_ConnectClient); ok {
		// I hate that Go makes me do this :(
		err = stream.CloseSend()
	}
	close(queue.quit)

	return err
}

// Send a message to the StreamQueue. The Message must be a WorkerTask that
// wraps a gRPC SmpcMessage, otherwise an error is returned.
func (queue *StreamQueue) Send(message dispatch.Message) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on send: %v", r)
		}
	}()

	switch message := message.(type) {
	case WorkerTask:
		if message.SmpcMessage != nil {
			queue.write <- message.SmpcMessage
		}
	default:
		return fmt.Errorf("cannot send message: unrecognized type %T", message)
	}

	return err
}

// Recv a message from the StreamQueue. All Messages returned will be
// WorkerTasks that wrap a gRPC SmpcMessage.
func (queue *StreamQueue) Recv() (dispatch.Message, bool) {
	message, ok := <-queue.read
	return WorkerTask{SmpcMessage: message}, ok
}

// writeAll messages from the messaging queue to the stream.
func (queue *StreamQueue) writeAll() error {
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
func (queue *StreamQueue) readAll() error {
	defer close(queue.read)
	for {
		message := new(rpc.SmpcMessage)
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
