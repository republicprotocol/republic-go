package dispatch

import (
	"fmt"
	"sync"
)

// A Multiplexer runs MessageQueues. It aggregates all messages into a unified
// channel, allowing a dynamic number of workers to select messages from a
// dyanmic number of MessageQueues.
type Multiplexer struct {
	messageQueuesMu *sync.RWMutex
	messageQueues   map[string]MessageQueue

	messages chan Message
}

// NewMultiplexer returns a Multiplexer that uses a queue limit to buffer the
// unified output channel.
func NewMultiplexer(queueLimit int) Multiplexer {
	return Multiplexer{
		messageQueuesMu: new(sync.RWMutex),
		messageQueues:   make(map[string]MessageQueue),

		messages: make(chan Message, queueLimit),
	}
}

// RunMessageQueue in the Multiplexer. The MessageQueue will run until it
// encounters an error, or until the Multiplexer is shutdown.
func (dispatcher *Multiplexer) RunMessageQueue(id string, messageQueue MessageQueue) error {

	// Store the MessageQueue until it has finished running
	dispatcher.messageQueuesMu.Lock()
	if _, ok := dispatcher.messageQueues[id]; !ok {
		dispatcher.messageQueues[id] = messageQueue
	}
	dispatcher.messageQueuesMu.Unlock()

	// Multiplex messages from this MessageQueue to the unified Multiplexer
	// channel
	go func() {
		// Recover from writing to a potential closed channel
		defer func() { recover() }()

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
	delete(dispatcher.messageQueues, id)
	dispatcher.messageQueuesMu.Unlock()

	return err
}

// Shutdown gracefully by sending a quit command to all of the MessageQueues
// running in the Multiplexer. This method must only be called exactly once,
// when the Multiplexer is no longer needed.
func (dispatcher *Multiplexer) Shutdown() {

	// Stop letting workers receive message, and stop accepting new
	// MessageQueues
	close(dispatcher.messages)

	dispatcher.messageQueuesMu.Lock()
	defer dispatcher.messageQueuesMu.Unlock()

	// While the mutex is locked, gracefully shutdown all MessageQueues
	for _, messageQueue := range dispatcher.messageQueues {
		// Ignore errors returned during shutdown
		_ = messageQueue.Shutdown()
	}
	dispatcher.messageQueues = map[string]MessageQueue{}
}

// ShutdownMessageQueue by giving its associated ID. If the MessageQueue is
// running on this Multiplexer, it will be gracefully shutdown.
func (dispatcher *Multiplexer) ShutdownMessageQueue(id string) {
	dispatcher.messageQueuesMu.Lock()
	defer dispatcher.messageQueuesMu.Unlock()

	// While the mutex is locked, gracefully shutdown the MessageQueue
	if messageQueue, ok := dispatcher.messageQueues[id]; ok {

		// Ignore errors returned during shutdown
		_ = messageQueue.Shutdown()
		delete(dispatcher.messageQueues, id)
	}
}

// Send a Message to a MessageQueue by giving its associated ID. If no
// MessageQueue is available, then an error is returned. If the MessageQueue is
// full, then this function will block.
func (dispatcher *Multiplexer) Send(id string, message Message) error {

	// Find the MessageQueue associated with the ID and forward the Message to
	// it
	dispatcher.messageQueuesMu.RLock()
	messageQueue := dispatcher.messageQueues[id]
	dispatcher.messageQueuesMu.RUnlock()
	if messageQueue == nil {
		return fmt.Errorf("cannot send message to %s: no message queue available", id)
	}

	return messageQueue.Send(message)
}

// Recv a Message from the Multiplexer. This function blocks until at least one
// MessageQueue has a Message.
func (dispatcher *Multiplexer) Recv() (Message, bool) {
	message, ok := <-dispatcher.messages
	return message, ok
}
