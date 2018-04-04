package dispatch

import (
	"fmt"
	"sync"
)

// A Splitter runs MessageQueues. It reads messages from a unified channel
// and splits them into dynamic number of MessagesQueues
type Splitter struct {
	maxConnections int

	outputMu *sync.RWMutex
	output   map[string]MessageQueue
}

// NewSplitter creates a new splitter with the giving max connections limit.
func NewSplitter(maxConnections int) Splitter {
	return Splitter{
		maxConnections: maxConnections,
		outputMu:       new(sync.RWMutex),
		output:         map[string]MessageQueue{},
	}
}

// RunMessageQueue in the Splitter. All Messages written to the Splitter will
// be written to the MessageQueue.  The MessageQueue will run until it
// encounters an error, or until the Splitter is shutdown. A MessageQueue run
// using a Splitter must not be run anywhere else.
func (splitter *Splitter) RunMessageQueue(id string, messageQueue MessageQueue) error {
	// Check number of connections
	splitter.outputMu.Lock()
	if len(splitter.output) >= splitter.maxConnections {
		splitter.outputMu.Unlock()
		return fmt.Errorf("cannot run message queue %s: max connections reached", id)
	}

	// Register the message queue as a output queue
	if _, ok := splitter.output[id]; !ok {
		splitter.output[id] = messageQueue
	} else {
		splitter.outputMu.Unlock()
		return fmt.Errorf("cannot run message queue %s: message queue is already running", id)
	}
	splitter.outputMu.Unlock()

	// Start streaming message to the message queue
	err := messageQueue.Run()

	// Remove the message queue when finished
	splitter.outputMu.Lock()
	delete(splitter.output, id)
	splitter.outputMu.Unlock()

	return err
}

// Shutdown gracefully by shutting down all MessageQueues running in the
// Splitter.
func (splitter *Splitter) Shutdown() {
	splitter.outputMu.Lock()
	defer splitter.outputMu.Unlock()

	// While the mutex is locked, gracefully shutdown all MessageQueues
	for _, messageQueue := range splitter.output {
		_ = messageQueue.Shutdown()
	}
	splitter.output = map[string]MessageQueue{}
}

// Send a Message to the Splitter. The Message will be forwarded to every
// MessageQueue running in the Splitter. If a MessageQueue is full, this
// function will block.
func (splitter *Splitter) Send(message Message) error {
	splitter.outputMu.RLock()
	defer splitter.outputMu.RUnlock()


	for _, messageQueue := range splitter.output {
		if err := messageQueue.Send(message); err != nil {
			return err
		}
	}
	return nil
}

// ShutdownMessageQueue will shut down a specific queue by its id.
func (splitter *Splitter) ShutdownMessageQueue(id string) {
	splitter.outputMu.Lock()
	defer splitter.outputMu.Unlock()

	// While the mutex is locked, gracefully shutdown the MessageQueue
	if messageQueue, ok := splitter.output[id]; ok {
		// Ignore errors returned during shutdown
		_ = messageQueue.Shutdown()
		delete(splitter.output, id)
	}
}

// SendByID will send the message to specific queue by its id.
func (splitter *Splitter) SendByID(id string, message Message) error {
	splitter.outputMu.RLock()
	defer splitter.outputMu.RUnlock()

	if queue, ok := splitter.output[id]; !ok {
		return fmt.Errorf("cannot send meesage, %s doesn't exist", id)
	} else {
		return queue.Send(message)
	}
}
