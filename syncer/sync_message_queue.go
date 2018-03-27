package syncer

import (
	"fmt"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/network/rpc"
)

const SyncMessageQueueLimit = 100

type SyncMessageQueue struct {
	stream rpc.Dark_SyncServer
	write  chan *rpc.SyncBlock
	quit   chan bool
}

func NewSyncMessageQueue(stream rpc.Dark_SyncServer) SyncMessageQueue {
	return SyncMessageQueue{
		stream : stream,
		write:    make(chan *rpc.SyncBlock, SyncMessageQueueLimit),
		quit:     make(chan bool, 1),
	}
}

func (syncMessageQueue SyncMessageQueue) Run() error {
	return syncMessageQueue.writeAll()
}

func (syncMessageQueue SyncMessageQueue) Shutdown() error {
	syncMessageQueue.quit <- true
	return nil
}

func (syncMessageQueue SyncMessageQueue) Send(message dispatch.Message) error {
	block , ok := message.(*rpc.SyncBlock)
	if !ok {
		return fmt.Errorf("wrong message type, has %T expect *rpc.SyncBlock", message)
	}
	syncMessageQueue.write <- block
	return nil
}

func (syncMessageQueue SyncMessageQueue) Recv() (dispatch.Message, bool) {
	// since it's a server-side stream, so we will not receive any message
	// from the client
	return struct {}{}, false
}

func (syncMessageQueue SyncMessageQueue) writeAll() error {
	for {
		select {
		case <- syncMessageQueue.quit:
			return nil
		case message := <-syncMessageQueue.write:
			if err := syncMessageQueue.stream.Send(message); err != nil {
				return err
			}
		}
	}
}
