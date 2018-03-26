package smpc

import (
	"io"

	"github.com/republicprotocol/republic-go/network/rpc"
)

// MessageQueueLimit defines the number of messages that can be buffered in
// the connection queue.
const MessageQueueLimit = 100

// MessageQueue workers own a gRPC stream to a client. All messages writing and
// reading from a gRPC stream must go through a MessageQueue. The writing and
// reading channels used by the MessageQueue are abstracted away, to prevent
// incorrectly reading from the write channel, or writing to the read channel.
type MessageQueue struct {
	stream rpc.TauService_ConnectServer
	write  chan *rpc.TauMessage
	read   chan *rpc.TauMessage
	quit   chan struct{}
}

// NewMessageQueue returns a MessageQueue worker that owns a gRPC stream and
// will gracefully stop when the quit channel is closed.
func NewMessageQueue(stream rpc.TauService_ConnectServer) MessageQueue {
	return MessageQueue{
		stream: stream,
		write:  make(chan *rpc.TauMessage, MessageQueueLimit),
		read:   make(chan *rpc.TauMessage, MessageQueueLimit),
		quit:   make(chan struct{}),
	}
}

// Run the MessageQueue and return the first error that happens.
func (messages *MessageQueue) Run() error {
	ch := make(chan error, 2)
	go func() { ch <- messages.writeAll() }()
	go func() { ch <- messages.readAll() }()
	for i := 0; i < 2; i++ {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

// Shutdown the MessageQueue.
func (messages *MessageQueue) Shutdown() {
	close(messages.quit)
}

// Send a message to the MessageQueue. If the queue is full, then this
// function will block.
func (messages *MessageQueue) Send(message *rpc.TauMessage) {
	messages.write <- message
}

// Recv a message from the MessageQueue. If the queue is empty, then this
// function will block.
func (messages *MessageQueue) Recv() (*rpc.TauMessage, bool) {
	message, ok := <-messages.read
	return message, ok
}

// TrySend a message to the MessageQueue. If the queue is full, then this
// function will abandon the message and return immediately.
func (messages *MessageQueue) TrySend(message *rpc.TauMessage) {
	select {
	case messages.write <- message:
	default:
	}
}

// TryRecv a message from the MessageQueue. If the queue is empty, then this
// function will return nil immediately.
func (messages *MessageQueue) TryRecv() (*rpc.TauMessage, bool) {
	select {
	case message, ok := <-messages.read:
		return message, ok
	default:
		return nil, true
	}
}

// writeAll messages from the messaging queue to the stream.
func (messages *MessageQueue) writeAll() error {
	for {
		select {
		case <-messages.quit:
			return nil
		case message := <-messages.write:
			if err := messages.stream.Send(message); err != nil {
				return err
			}
		}
	}
}

// readAll messages from the stream and write them to the output queue. If the
// output queue is full, the MessageQueue will stop reading messages until the
// output is read from.
func (messages *MessageQueue) readAll() error {
	defer close(messages.read)
	for {
		message, err := messages.stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if message != nil {
			select {
			case <-messages.quit:
				return nil
			case messages.read <- message:
			}
		}
	}
}
