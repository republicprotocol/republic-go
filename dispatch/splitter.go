package dispatch

import (
	"sync"
)

type Splitter struct {
	outputMu *sync.RWMutex
	output   map[string]MessageQueue
}

func NewSplitter() Splitter {
	return Splitter{
		outputMu: new(sync.RWMutex),
		output:   map[string]MessageQueue{},
	}
}

// RunMessageQueue in the Splitter. All Messages written to the Splitter will
// be written to the MessageQueue.  The MessageQueue will run until it
// encounters an error, or until the Splitter is shutdown. A MessageQueue run
// using a Splitter must not be run anywhere else.
func (splitter *Splitter) RunMessageQueue(id string, messageQueue MessageQueue) error {
	splitter.outputMu.Lock()
	if _, ok := splitter.output[id]; !ok {
		splitter.output[id] = messageQueue
		splitter.outputMu.Unlock()
		return nil
	}
	splitter.outputMu.Unlock()

	err := messageQueue.Run()

	splitter.outputMu.Lock()
	delete(splitter.output, id)
	splitter.outputMu.Unlock()

	return err
}

// Shutdown gracefully by shuttding down all MessageQueues running in the
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
