package syncer

import (
	"fmt"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/network/rpc"
)

const (
	// SyncMessageQueueLimit defines the number of messages that can be
	// buffered in the message queue.
	SyncMessageQueueLimit = 100

	// QuitChannelSize is the size of the quit channel in SyncMessageQueue
	QuitChannelSize = 1
)

// A SyncMessageQueue is a implementation of the dispatch.MessageQueue
// specifically for synchronization. It will only be able to write to the
// stream and ignore any messages it receives.
type SyncMessageQueue struct {
	stream rpc.Dark_SyncServer
	write  chan *rpc.SyncBlock
	quit   chan bool
}

// NewSyncMessageQueue returns a SyncMessageQueue that owns a gRPC stream.
func NewSyncMessageQueue(stream rpc.Dark_SyncServer) SyncMessageQueue {
	return SyncMessageQueue{
		stream: stream,
		write:  make(chan *rpc.SyncBlock, SyncMessageQueueLimit),
		quit:   make(chan bool, QuitChannelSize),
	}
}

// Run starts running the SyncMessageQueue and returns the first error happens.
func (syncMessageQueue SyncMessageQueue) Run() error {
	return syncMessageQueue.writeAll()
}

// Shutdown the SyncMessageQueue.
func (syncMessageQueue SyncMessageQueue) Shutdown() error {
	syncMessageQueue.quit <- true
	return nil
}

// Send a message to the SyncMessageQueue. It will return an error if the message
// is not of type rpc.SyncBlock. It will also block if the queue is full.
func (syncMessageQueue SyncMessageQueue) Send(message dispatch.Message) error {
	block, ok := message.(*rpc.SyncBlock)
	if !ok {
		return fmt.Errorf("wrong message type, has %T expect *rpc.SyncBlock", message)
	}
	syncMessageQueue.write <- block
	return nil
}

// Recv will always return an empty struct and false, it should not be called.
func (syncMessageQueue SyncMessageQueue) Recv() (dispatch.Message, bool) {
	// since it's a server-side stream, so we will not receive any message
	// from the client
	return struct{}{}, false
}

// Write all messages from the SyncMessageQueue to the stream.
func (syncMessageQueue SyncMessageQueue) writeAll() error {
	for {
		select {
		case <-syncMessageQueue.quit:
			return nil
		case message := <-syncMessageQueue.write:
			if err := syncMessageQueue.stream.Send(message); err != nil {
				return err
			}
		}
	}
}
