package dispatch

import (
	"errors"
	"sync"

	do "github.com/republicprotocol/go-do"
)

// The Message interface is type-only interface.
type Message interface{}

// MessageQueues is a slice of MessageQueue interfaces.
type MessageQueues []MessageQueue

// The MessageQueue interface defines a set of expected functionality for a
// queue to be integrated with the Dispatcher.
type MessageQueue interface {

	// Run the MessageQueue, processing all messages that are sent and
	// received. Run must only be called once.
	Run() error

	// Shutdown the MessageQueue gracefully. Shutdown implementations should be
	// idempotent, and handle multiple calls without panicking, or returning an
	// error.
	Shutdown() error

	// Send a Message to the MessageQueue. The implementation shooruld throw a
	// type error if it receives a concrete type that it does not recognize.
	// This method should block if the MessageQueue is full.
	Send(Message) error

	// Recv a message from the MessageQueue. Return nil, and false, if the
	// MessageQueue has been shutdown, otherwise return a Message, and true.
	// This method should block if the MessageQueue is empty.
	Recv() (Message, bool)
}

// The ChannelQueue component is a MessageQueue backed by a Go channel.
type ChannelQueue struct {
	chMu   *sync.RWMutex
	chOpen bool
	ch     chan Message
}

// NewChannelQueue returns a MessageQueue interface that is backed by a Go
// channel. The underlying channel is buffered with a size equal to the message
// queue limit.
func NewChannelQueue(messageQueueLimit int) *ChannelQueue {
	return &ChannelQueue{
		chMu:   new(sync.RWMutex),
		chOpen: true,
		ch:     make(chan Message, messageQueueLimit),
	}
}

// Run the ChannelQueue. The ChannelQueue is an abstraction over a channel of Delta
// components and does not need to be run. This method does nothing.
func (queue *ChannelQueue) Run() error {
	return nil
}

// Shutdown the ChannelQueue. If it has already been Shutdown, an error will be
// returned.
func (queue *ChannelQueue) Shutdown() error {
	queue.chMu.Lock()
	defer queue.chMu.Unlock()

	if !queue.chOpen {
		return errors.New("cannot shutdown channel queue: already shutdown")
	}
	queue.chOpen = false
	close(queue.ch)

	return nil
}

// Send a message to the ChannelQueue.
func (queue *ChannelQueue) Send(message Message) error {
	queue.chMu.RLock()
	defer queue.chMu.RUnlock()

	if !queue.chOpen {
		return nil
	}
	queue.ch <- message

	return nil
}

// Recv a message from the ChannelQueue.
func (queue *ChannelQueue) Recv() (Message, bool) {
	message, ok := <-queue.ch
	return message, ok
}

// The UnboundedQueue component is a MessageQueue backed by a Go slice that is
// grown as necessary. Writing to an UnboundedQueue will never block, because
// the queue can never be filled.
type UnboundedQueue struct {
	do.GuardedObject

	open          bool
	items         []Message
	itemsNotEmpty *do.Guard
}

// NewUnboundedQueue returns a MessageQueue interface that is backed by a Go
// slice. The underlying slice is initially allocated using the buffer hint.
func NewUnboundedQueue(bufferHint int) *UnboundedQueue {
	queue := new(UnboundedQueue)
	queue.GuardedObject = do.NewGuardedObject()
	queue.open = true
	queue.items = make([]Message, bufferHint)
	queue.itemsNotEmpty = queue.Guard(func() bool { return len(queue.items) > 0 })
	return queue
}

// Run the ChannelQueue. The ChannelQueue is an abstraction over a channel of Delta
// components and does not need to be run. This method does nothing.
func (queue *UnboundedQueue) Run() error {
	return nil
}

// Shutdown the UnboundedQueue. If it has already been Shutdown, an error will
// be returned.
func (queue *UnboundedQueue) Shutdown() error {
	queue.Enter(nil)
	defer queue.Exit()

	if !queue.open {
		return errors.New("cannot shutdown unbounded queue: already shutdown")
	}
	queue.open = false
	queue.items = []Message{}
	return nil
}

// Send a message to the UnboundedQueue.
func (queue *UnboundedQueue) Send(message Message) error {
	queue.Enter(nil)
	defer queue.Exit()

	queue.items = append(queue.items, message)
	return nil
}

// Recv a message from the UnboundedQueue.
func (queue *UnboundedQueue) Recv() (Message, bool) {
	queue.Enter(queue.itemsNotEmpty)
	defer queue.Exit()

	if !queue.open {
		return nil, false
	}

	item := queue.items[0]
	queue.items = queue.items[1:]
	return item, true
}
