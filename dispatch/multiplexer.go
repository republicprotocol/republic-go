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
func (multiplexer *Multiplexer) RunMessageQueue(id string, messageQueue MessageQueue) error {

	// Store the MessageQueue until it has finished running
	multiplexer.messageQueuesMu.Lock()
	if _, ok := multiplexer.messageQueues[id]; !ok {
		multiplexer.messageQueues[id] = messageQueue
	}
	multiplexer.messageQueuesMu.Unlock()

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
			multiplexer.messages <- message
		}
	}()

	// Run the MessageQueue until an error is encountered, or the MessageQueue
	// receives a signal to shutdown gracefully, and then return any error
	err := messageQueue.Run()

	// Remove the MessageQueue now that it has finished running
	multiplexer.messageQueuesMu.Lock()
	delete(multiplexer.messageQueues, id)
	multiplexer.messageQueuesMu.Unlock()

	return err
}

// Shutdown gracefully by sending a quit command to all of the MessageQueues
// running in the Multiplexer. This method must only be called exactly once,
// when the Multiplexer is no longer needed.
func (multiplexer *Multiplexer) Shutdown() {

	// Stop letting workers receive message, and stop accepting new
	// MessageQueues
	close(multiplexer.messages)

	multiplexer.messageQueuesMu.Lock()
	defer multiplexer.messageQueuesMu.Unlock()

	// While the mutex is locked, gracefully shutdown all MessageQueues
	for _, messageQueue := range multiplexer.messageQueues {
		// Ignore errors returned during shutdown
		_ = messageQueue.Shutdown()
	}
	multiplexer.messageQueues = map[string]MessageQueue{}
}

// ShutdownMessageQueue by giving its associated ID. If the MessageQueue is
// running on this Multiplexer, it will be gracefully shutdown.
func (multiplexer *Multiplexer) ShutdownMessageQueue(id string) {
	multiplexer.messageQueuesMu.Lock()
	defer multiplexer.messageQueuesMu.Unlock()

	// While the mutex is locked, gracefully shutdown the MessageQueue
	if messageQueue, ok := multiplexer.messageQueues[id]; ok {

		// Ignore errors returned during shutdown
		_ = messageQueue.Shutdown()
		delete(multiplexer.messageQueues, id)
	}
}

// Send a Message directly to the Multiplexer unified channel. If the channel
// is full, then this function will block.
func (multiplexer *Multiplexer) Send(message Message) {
	multiplexer.messages <- message
}

// SendToMessageQueue by giving its associated ID. If no MessageQueue is
// available, then an error is returned. If the MessageQueue is full, then this
// function will block.
func (multiplexer *Multiplexer) SendToMessageQueue(id string, message Message) error {

	// Find the MessageQueue associated with the ID and forward the Message to
	// it
	multiplexer.messageQueuesMu.RLock()
	messageQueue := multiplexer.messageQueues[id]
	multiplexer.messageQueuesMu.RUnlock()
	if messageQueue == nil {
		return fmt.Errorf("cannot send message to %s: no message queue available", id)
	}

	return messageQueue.Send(message)
}

// Recv a Message from the Multiplexer. This function blocks until at least one
// MessageQueue has a Message.
func (multiplexer *Multiplexer) Recv() (Message, bool) {
	message, ok := <-multiplexer.messages
	return message, ok
}
