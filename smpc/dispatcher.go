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
	messageQueuesMu *sync.RWMutex
	messageQueues   map[string]*MessageQueue

	messages chan *rpc.TauMessage
}

// NewDispatcher returns a Dispatcher that uses a buffer hint to estimate the
// number of MessageQueues that it will be handling.
func NewDispatcher(bufferHint int) Dispatcher {
	return Dispatcher{
		messageQueuesMu: new(sync.RWMutex),
		messageQueues:   make(map[string]*MessageQueue),

		messages: make(chan *rpc.TauMessage, bufferHint*MessageQueueLimit),
	}
}

// RunMessageQueue in the Dispatcher. The MessageQueue will run until it
// encounters an error writing to, or reading from, the gRPC stream, or until
// the Dispatcher is shutdown.
func (dispatcher *Dispatcher) RunMessageQueue(multiAddress identity.MultiAddress, messageQueue *MessageQueue) error {
	address := multiAddress.Address().String()

	// Store the MessageQueue until it has finished running
	dispatcher.messageQueuesMu.Lock()
	dispatcher.messageQueues[address] = messageQueue
	dispatcher.messageQueuesMu.Unlock()

	// Multiplex messages from this MessageQueue to the unified Dispatcher
	// channel
	go func() {
		defer func() { recover() }() // Recover from writing to a potential closed channel
		for {
			message, ok := messageQueue.Recv()
			if !ok {
				break
			}
			dispatcher.messages <- message
		}
	}()

	// Run the MessageQueue until an error is encountered, or the MessageQueue
	// receives a signal to shutdown gracefully, and then return any error
	err := messageQueue.Run()

	// Remove the MessageQueue now that it has finished running
	dispatcher.messageQueuesMu.Lock()
	delete(dispatcher.messageQueues, address)
	dispatcher.messageQueuesMu.Unlock()

	return err
}

// Shutdown gracefully by sending a quit command to all of the MessageQueues
// running in the Dispatcher. This method must only be called exactly once,
// when the Dispatcher is no longer needed.
func (dispatcher *Dispatcher) Shutdown() {

	// Stop letting workers receive message, and stop accepting new
	// MessageQueues
	close(dispatcher.messages)

	dispatcher.messageQueuesMu.Lock()
	defer dispatcher.messageQueuesMu.Unlock()

	// While the mutex is locked, gracefully shutdown all MessageQueues
	for _, messageQueue := range dispatcher.messageQueues {
		messageQueue.Shutdown()
	}
	dispatcher.messageQueues = map[string]*MessageQueue{}
}

// ShutdownMessageQueue by giving its associated multi-address. If the
// MessageQueue is running on this Dispatcher, it will be gracefully shutdown.
func (dispatcher *Dispatcher) ShutdownMessageQueue(multiAddress identity.MultiAddress) {
	address := multiAddress.Address().String()

	dispatcher.messageQueuesMu.Lock()
	defer dispatcher.messageQueuesMu.Unlock()

	// While the mutex is locked, gracefully shutdown the MessageQueue
	if messageQueue, ok := dispatcher.messageQueues[address]; ok {
		messageQueue.Shutdown()
		delete(dispatcher.messageQueues, address)
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
