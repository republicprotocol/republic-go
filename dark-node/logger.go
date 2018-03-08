package node

import (
	"time"

	"github.com/republicprotocol/go-do"
)

// LogQueue allows multiple clients to receive logs from a node
type LogQueue struct {
	do.GuardedObject

	channels []chan do.Option
}

// NewLogQueue returns a new LogQueue
func NewLogQueue() *LogQueue {
	logQueue := new(LogQueue)
	logQueue.GuardedObject = do.NewGuardedObject()
	logQueue.channels = nil
	return logQueue
}

// Publish allows a node to push a log to each client
func (logQueue *LogQueue) Publish(val do.Option) {
	logQueue.Enter(nil)
	defer logQueue.Exit()

	var logQueueLength = len(logQueue.channels)
	for i := 0; i < logQueueLength; i++ {
		timer := time.NewTicker(10 * time.Second)
		defer timer.Stop()
		select {
		case logQueue.channels[i] <- val:
		case <-timer.C:
			// Deregister the channel
			logQueue.channels[i] = logQueue.channels[logQueueLength-1]
			logQueue.channels = logQueue.channels[:logQueueLength-1]
			logQueueLength--
			i--
		}
	}
}

// Subscribe allows a new client to listen to events from a node
func (logQueue *LogQueue) Subscribe(channel chan do.Option) {
	logQueue.Enter(nil)
	defer logQueue.Exit()

	logQueue.channels = append(logQueue.channels, channel)
}

// Unsubscribe ...
func (logQueue *LogQueue) Unsubscribe(channel chan do.Option) {
	logQueue.Enter(nil)
	defer logQueue.Exit()
	length := len(logQueue.channels)
	for i := 0; i < length; i++ {
		// https://golang.org/ref/spec#Comparison_operators
		// Two channel values are equal if they were created by the same call to make
		// or if both have value nil.
		if logQueue.channels[i] == channel {
			logQueue.channels[i] = logQueue.channels[length-1]
			logQueue.channels = logQueue.channels[:length-1]
			break
		}
	}
}

// SubscribeToLogs will start sending log events to logChannel
func (node *DarkNode) SubscribeToLogs(logChannel chan do.Option) {
	node.logQueue.Subscribe(logChannel)
}

// UnsubscribeFromLogs will stop sending log events to logChannel
func (node *DarkNode) UnsubscribeFromLogs(logChannel chan do.Option) {
	node.logQueue.Unsubscribe(logChannel)
}
