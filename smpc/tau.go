package smpc

import (
	"fmt"
	"io"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
	"google.golang.org/grpc"
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
func (tau *Tau) Connect(stream rpc.Tau_ConnectServer) error {

	// Use a background MessageQueue to handle the connection until an error
	// is returned by the MessageQueue
	ch := make(chan error, 1)
	quit := make(chan struct{}, 1)
	go func() { ch <- tau.connect(stream, quit) }()
	defer close(quit)

	// Select between the context finishing and the background worker
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case err := <-ch:
		return err
	}
}

func (tau *Tau) connect(stream rpc.Tau_ConnectServer, quit chan struct{}) error {
	multiAddress, err := identity.NewMultiAddressFromString("unimplemented")
	if err != nil {
		return err
	}

	// Create a MessageQueue that owns this gRPC stream and run it on the
	// Dispatcher
	messageQueue := NewTauServerQueue(stream, tau.messageQueueLimit)

	// Shutdown the MessageQueue
	go func() {
		<-quit
		messageQueue.Shutdown()
	}()
	return tau.multiplexer.RunMessageQueue(multiAddress.Address().String(), messageQueue)
}

// TauStreamQueue workers own a gRPC stream. All Message writing to, and
// reading from, a gRPC stream must go through a TauStreamQueue. The writing
// and reading channels used by the TauStreamQueue are abstracted away to
// prevent incorrectly reading from the write channel, and writing to the read
// channel.
type TauStreamQueue struct {
	stream grpc.Stream
	write  chan *rpc.TauMessage
	read   chan *rpc.TauMessage
	quit   chan struct{}
}

// NewTauClientQueue returns a MessageQueue interface that is connected to a
// gRPC server. It accepts a gRPC stream, that will be owned by the
// MessageQueue and should not be used by any other component
func NewTauClientQueue(stream rpc.Tau_ConnectClient, messageQueueLimit int) dispatch.MessageQueue {
	return &TauStreamQueue{
		stream: stream,
		write:  make(chan *rpc.TauMessage, messageQueueLimit),
		read:   make(chan *rpc.TauMessage, messageQueueLimit),
		quit:   make(chan struct{}),
	}
}

// NewTauServerQueue returns a MessageQueue interface that is connected to a
// gRPC client. It accepts a gRPC stream, that will be owned by the
// MessageQueue and should not be used by any other component
func NewTauServerQueue(stream rpc.Tau_ConnectServer, messageQueueLimit int) dispatch.MessageQueue {
	return &TauStreamQueue{
		stream: stream,
		write:  make(chan *rpc.TauMessage, messageQueueLimit),
		read:   make(chan *rpc.TauMessage, messageQueueLimit),
		quit:   make(chan struct{}),
	}
}

// Run the TauStreamQueue. This will concurrently process all writes to the
// underlying gRPC stream, and buffer all messages received over the gRPC
// stream.
func (queue *TauStreamQueue) Run() error {
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

// Shutdown the TauStreamQueue. If it has already been Shutdown, and error will
// be returned.
func (queue *TauStreamQueue) Shutdown() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on shutdown: %v", r)
		}
	}()

	// If the stream is a client stream, then close the sending channel
	if stream, ok := queue.stream.(rpc.Tau_ConnectClient); ok {
		// I hate that Go makes me do this :(
		err = stream.CloseSend()
	}
	close(queue.quit)

	return err
}

// Send a message to the TauStreamQueue. The Message must be a WorkerTask that
// wraps a gRPC TauMessage, otherwise an error is returned.
func (queue *TauStreamQueue) Send(message dispatch.Message) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on send: %v", r)
		}
	}()

	switch message := message.(type) {
	case WorkerTask:
		if message.TauMessage != nil {
			queue.write <- message.TauMessage
		}
	default:
		return fmt.Errorf("cannot send message: unrecognized type %T", message)
	}

	return err
}

// Recv a message from the TauStreamQueue. All Messages returned will be
// WorkerTasks that wrap a gRPC TauMessage.
func (queue *TauStreamQueue) Recv() (dispatch.Message, bool) {
	message, ok := <-queue.read
	return WorkerTask{TauMessage: message}, ok
}

// writeAll messages from the messaging queue to the stream.
func (queue *TauStreamQueue) writeAll() error {
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
// output queue is full, the TauStreamQueue will stop reading messages until
// the output is read from.
func (queue *TauStreamQueue) readAll() error {
	defer close(queue.read)
	for {
		message := new(rpc.TauMessage)
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
