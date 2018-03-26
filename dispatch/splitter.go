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

func (splitter *Splitter) RunMessageQueue(id string, messageQueue MessageQueue) error {
	splitter.outputMu.Lock()
	if _, ok := splitter.output[id]; !ok {
		splitter.output[id] = messageQueue
	}
	splitter.outputMu.Unlock()

	err := messageQueue.Run()

	splitter.outputMu.Lock()
	delete(splitter.output, id)
	splitter.outputMu.Unlock()

	return err
}

func (splitter *Splitter) Shutdown() {
	splitter.outputMu.Lock()
	defer splitter.outputMu.Unlock()

	// While the mutex is locked, gracefully shutdown all MessageQueues
	for _, messageQueue := range splitter.output {
		_ = messageQueue.Shutdown()
	}
	splitter.output = map[string]MessageQueue{}
}

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

func (splitter *Splitter) Recv(id string) (Message, bool) {
	splitter.outputMu.RLock()
	defer splitter.outputMu.RUnlock()

	if _, ok := splitter.output[id]; !ok {
		return nil, false
	}
	return splitter.output[id].Recv()
}
