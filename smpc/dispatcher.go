package smpc

import (
	"fmt"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/network/rpc"
)

// A Dispatcher runs MessageQueues. It aggregates all messages into a unified
// channel, allowing a dynamic number of workers to select messages from a
// dyanmic number of MessageQueues.
type Dispatcher struct {
	messageQueuesMu   *sync.RWMutex
	messageQueues     map[string]*MessageQueue
	messageQueuesQuit map[string]chan struct{}

	messages                 chan *rpc.TauMessage
	messageQueuesMultiplexer chan *MessageQueue
}

// NewDispatcher returns a Dispatcher that uses a buffer hint to estimate the
// number of MessageQueues that it will be handling.
func NewDispatcher(bufferHint int) Dispatcher {
	return Dispatcher{
		messageQueuesMu:   new(sync.RWMutex),
		messageQueues:     make(map[string]*MessageQueue),
		messageQueuesQuit: make(map[string]chan struct{}),

		messages:                 make(chan *rpc.TauMessage, bufferHint*MessageQueueLimit),
		messageQueuesMultiplexer: make(chan *MessageQueue, bufferHint),
	}
}

// Run the Dispatcher by multiplexing all MessageQueues into a unified channel
// that can be used by workers.
func (dispatcher *Dispatcher) Run() {
	for messageQueue := range dispatcher.messageQueuesMultiplexer {
		go func(messageQueue *MessageQueue) {
			for {
				message, ok := messageQueue.Recv()
				if !ok {
					break
				}
				dispatcher.messages <- message
			}
		}(messageQueue)
	}
}

// RunMessageQueue in the Dispatcher. The MessageQueue will run until it
// encounters an error writing to, or reading from, the gRPC stream, or until
// the Dispatcher is shutdown.
func (dispatcher *Dispatcher) RunMessageQueue(multiAddress identity.MultiAddress, messageQueue *MessageQueue) error {
	address := multiAddress.Address().String()
	quit := make(chan struct{})
	dispatcher.messageQueuesMultiplexer <- messageQueue

	// Store the MessageQueue until it has finished running
	dispatcher.messageQueuesMu.Lock()
	dispatcher.messageQueues[address] = messageQueue
	dispatcher.messageQueuesQuit[address] = quit
	dispatcher.messageQueuesMu.Unlock()

	// Run the MessageQueue until an error is encountered, or the MessageQueue
	// receives a signal to shutdown gracefully, and then return any error
	err := messageQueue.Run(quit)

	// Remove the MessageQueue now that it has finished running
	dispatcher.messageQueuesMu.Lock()
	delete(dispatcher.messageQueues, address)
	delete(dispatcher.messageQueuesQuit, address)
	dispatcher.messageQueuesMu.Unlock()

	return err
}

// Shutdown gracefully by sending a quit command to all of the MessageQueues
// running in the Dispatcher.
func (dispatcher *Dispatcher) Shutdown() {
	dispatcher.messageQueuesMu.RLock()
	defer dispatcher.messageQueuesMu.RUnlock()

	// Stop letting workers receive message, and stop accepting new
	// MessageQueues
	close(dispatcher.messages)
	close(dispatcher.messageQueuesMultiplexer)

	dispatcher.messageQueuesMu.Lock()
	defer dispatcher.messageQueuesMu.Unlock()

	// While the mutex is locked, close all quit channels and let the
	// MessageQueues gracefully shutdown on their own
	for _, quit := range dispatcher.messageQueuesQuit {
		close(quit)
	}
	dispatcher.messageQueuesQuit = map[string]chan struct{}{}
}

// ShutdownMessageQueue by giving its associated multi-address. If the
// MessageQueue is running on this Dispatcher, it will be gracefully shutdown.
func (dispatcher *Dispatcher) ShutdownMessageQueue(multiAddress identity.MultiAddress) {
	address := multiAddress.Address().String()

	dispatcher.messageQueuesMu.Lock()
	defer dispatcher.messageQueuesMu.Unlock()

	// While the mutex is locked, close the associated quite channel and let
	// MessageQueue gracefully shutdown on its own
	if quit, ok := dispatcher.messageQueuesQuit[address]; ok {
		close(quit)
		delete(dispatcher.messageQueuesQuit, address)
	}
}

// Send a message to a multi-address by finding an available MessageQueue. If
// no MessageQueue is available, then an error is returned. If the MessageQueue
// is full, then this function will block.
func (dispatcher *Dispatcher) Send(multiAddress identity.MultiAddress, message *rpc.TauMessage) error {
	address := multiAddress.Address().String()

	// Find the MessageQueue associated with the multi-address and forward the
	// message to it
	dispatcher.messageQueuesMu.RLock()
	messageQueue := dispatcher.messageQueues[address]
	dispatcher.messageQueuesMu.RUnlock()
	if messageQueue == nil {
		return fmt.Errorf("cannot send message to %s: no message queue available", address)
	}

	messageQueue.Send(message)
	return nil
}

// Recv a message from the Dispatcher. This function blocks until at least one
// MessageQueue has a message.
func (dispatcher *Dispatcher) Recv() (*rpc.TauMessage, bool) {
	message, ok := <-dispatcher.messages
	return message, ok
}
