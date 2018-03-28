package dispatch

import (
	"fmt"
	"sync"
)

// A Multiplexer runs MessageQueues. It aggregates all messages into a unified
// channel, allowing a dynamic number of workers to select messages from a
// dynamic number of MessageQueues.
type Multiplexer struct {
	messageQueuesMu *sync.RWMutex
	messageQueues   map[string]MessageQueue

	messagesMu   *sync.RWMutex
	messagesOpen bool
	messages     chan Message
}

// NewMultiplexer returns a Multiplexer that uses a queue limit to buffer the
// unified output channel.
func NewMultiplexer(queueLimit int) Multiplexer {
	return Multiplexer{
		messageQueuesMu: new(sync.RWMutex),
		messageQueues:   make(map[string]MessageQueue),

		messagesMu:   new(sync.RWMutex),
		messagesOpen: true,
		messages:     make(chan Message, queueLimit),
	}
}


// RunMessageQueue in the Multiplexer. All Messages written to the MessageQueue
// will be aggregated into the unified channel. The MessageQueue will run until
// it encounters an error, or until the Multiplexer is shutdown. A MessageQueue
// run using a Multiplexer must not be run anywhere else.
func (multiplexer *Multiplexer) RunMessageQueue(id string, messageQueue MessageQueue) error {

	// Store the MessageQueue until it has finished running
	multiplexer.messageQueuesMu.Lock()
	if _, ok := multiplexer.messageQueues[id]; !ok {
		multiplexer.messageQueues[id] = messageQueue
	} else {
		multiplexer.messageQueuesMu.Unlock()
		return fmt.Errorf("cannot run message queue %s: message queue is already running", id)
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


// Shutdown gracefully by shutting down all MessageQueues running in the
// Multiplexer.
func (multiplexer *Multiplexer) Shutdown() {

	// Stop letting workers receive message, and stop accepting new
	// MessageQueues
	func() {
		multiplexer.messagesMu.Lock()
		defer multiplexer.messagesMu.Unlock()

		// While the mutex is locked, close the channel and prevent further
		// writes
		multiplexer.messagesOpen = false
		close(multiplexer.messages)
	}()

	// While the mutex is locked, gracefully shutdown all MessageQueues
	multiplexer.messageQueuesMu.Lock()
	defer multiplexer.messageQueuesMu.Unlock()

	for _, messageQueue := range multiplexer.messageQueues {
		// Ignore errors returned during shutdown
		_ = messageQueue.Shutdown()
	}
	multiplexer.messageQueues = map[string]MessageQueue{}
}


// Send a Message directly to the Multiplexer unified channel. If the channel
// is full, then this function will block.
func (multiplexer *Multiplexer) Send(message Message) {
	multiplexer.messagesMu.RLock()
	defer multiplexer.messagesMu.RUnlock()
	if multiplexer.messagesOpen {
		multiplexer.messages <- message
	}
}

// Recv a Message from the Multiplexer. This function blocks until at least one
// MessageQueue has a Message.
func (multiplexer *Multiplexer) Recv() (Message, bool) {
	message, ok := <-multiplexer.messages
	return message, ok
}
