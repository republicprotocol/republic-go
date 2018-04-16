package relay

import (
	"fmt"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/orderbook"
)

type WriteOnlyChannelQueue struct {
	out   chan<- orderbook.Entry
	write chan orderbook.Entry
	quit  chan struct{}
}

func NewWriteOnlyChannelQueue(out chan orderbook.Entry, messageQueueLimit int) WriteOnlyChannelQueue {
	return WriteOnlyChannelQueue{
		out:   out,
		write: make(chan orderbook.Entry, messageQueueLimit),
		quit:  make(chan struct{}),
	}
}

func (queue WriteOnlyChannelQueue) Run() error {
	return queue.writeAll()
}

func (queue WriteOnlyChannelQueue) Shutdown() error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on shutdown: %v", r)
		}
	}()

	close(queue.quit)
	return err
}

func (queue WriteOnlyChannelQueue) Send(message dispatch.Message) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on send: %v", r)
		}
	}()

	msg, ok := message.(orderbook.Entry)
	if !ok {
		return fmt.Errorf("wrong message type, has %T expect orderbook.Entry", message)
	}
	queue.write <- msg

	return err
}

func (queue WriteOnlyChannelQueue) Recv() (dispatch.Message, bool) {
	panic("read from a read only message queue")
}

func (queue WriteOnlyChannelQueue) writeAll() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on send all: %v", r)
		}
	}()

	for {
		select {
		case <-queue.quit:
			return nil
		case message := <-queue.write:
			queue.out <- message
		}
	}
}
