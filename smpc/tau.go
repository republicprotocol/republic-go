package smpc

import (
	"fmt"
	"io"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
)

// Tau implements the Tau gRPC service. Tau creates TauMessageQueues for each
// gRPC stream, and runs them on a Multiplexer. The closure of a gRPC stream,
// by the client or by the server, will prompt Tau to shutdown the
// TauMessageQueue.
type Tau struct {
	multiplexer       *dispatch.Multiplexer
	messageQueueLimit int
}

// NewTau returns a new Tau service that will run TauMessageQueues on the given
// Dispatcher. The message queue limit is used used to buffer the size of the
// TauMessageQueues that are created by the Tau.
func NewTau(multiplexer *dispatch.Multiplexer, messageQueueLimit int) Tau {
	return Tau{
		multiplexer:       multiplexer,
		messageQueueLimit: messageQueueLimit,
	}
}

// Connect to the Tau service and begin streaming requests and responses for
// the Tau sMPC protocol.
func (τ *Tau) Connect(stream rpc.TauService_ConnectServer) error {

	// Use a background MessageQueue to handle the connection until an error
	// is returned by the MessageQueue
	ch := make(chan error, 1)
	quit := make(chan struct{}, 1)
	go func() { ch <- τ.connect(stream, quit) }()
	defer close(quit)

	// Select between the context finishing and the background worker
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case err := <-ch:
		return err
	}
}

func (τ *Tau) connect(stream rpc.TauService_ConnectServer, quit chan struct{}) error {
	multiAddress, err := identity.NewMultiAddressFromString("unimplemented")
	if err != nil {
		return err
	}

	// Create a MessageQueue that owns this gRPC stream and run it on the
	// Dispatcher
	messageQueue := NewTauMessageQueue(stream, τ.messageQueueLimit)

	// Shutdown the MessageQueue
	go func() {
		<-quit
		τ.multiplexer.ShutdownMessageQueue(multiAddress.Address().String())
	}()
	return τ.multiplexer.RunMessageQueue(multiAddress.Address().String(), messageQueue)
}

// TauMessageQueue workers own a gRPC stream to a client. All Message
// writing to, and reading from, a gRPC stream must go through a
// TauMessageQueue. The writing and reading channels used by the
// StreaMessageQueue are abstracted away to prevent incorrectly reading from
// the write channel, and writing to the read channel.
type TauMessageQueue struct {
	stream rpc.TauService_ConnectServer
	write  chan *rpc.TauMessage
	read   chan *rpc.TauMessage
	quit   chan struct{}
}

// NewTauMessageQueue returns a MessageQueue interface. It accepts a gRPC
// stream, that will be owned by the MessageQueue and should not be used by any
// other component
func NewTauMessageQueue(stream rpc.TauService_ConnectServer, messageQueueLimit int) dispatch.MessageQueue {
	return &TauMessageQueue{
		stream: stream,
		write:  make(chan *rpc.TauMessage, messageQueueLimit),
		read:   make(chan *rpc.TauMessage, messageQueueLimit),
		quit:   make(chan struct{}),
	}
}

// Run the TauMessageQueue. This will concurrently process all writes to the
// underlying gRPC stream, and buffer all messages received over the gRPC
// stream.
func (queue *TauMessageQueue) Run() error {
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

// Shutdown the TauMessageQueue. If it has already been Shutdown, and error
// will be returned.
func (queue *TauMessageQueue) Shutdown() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on shutdown: %v", r)
		}
	}()

	close(queue.quit)

	return err
}

// Send a message to the TauMessageQueue. The Message must be a pointer to
// a gRPC TauMessage, otherwise an error is returned.
func (queue *TauMessageQueue) Send(message dispatch.Message) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on send: %v", r)
		}
	}()

	switch message := message.(type) {
	case *rpc.TauMessage:
		queue.write <- message
	default:
		return fmt.Errorf("cannot send message: unrecognized type %T", message)
	}

	return err
}

// Recv a message from the TauMessageQueue. All Messages returned will be a
// pointer to a gRPC TauMessage.
func (queue *TauMessageQueue) Recv() (dispatch.Message, bool) {
	message, ok := <-queue.read
	return message, ok
}

// writeAll messages from the messaging queue to the stream.
func (queue *TauMessageQueue) writeAll() error {
	for {
		select {
		case <-queue.quit:
			return nil
		case message := <-queue.write:
			if err := queue.stream.Send(message); err != nil {
				return err
			}
		}
	}
}

// readAll messages from the stream and write them to the output queue. If the
// output queue is full, the TauMessageQueue will stop reading messages
// until the output is read from.
func (queue *TauMessageQueue) readAll() error {
	defer close(queue.read)
	for {
		message, err := queue.stream.Recv()
		if err != nil {
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
