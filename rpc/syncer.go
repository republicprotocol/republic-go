package rpc

import (
	"fmt"
	"io"
	"log"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"google.golang.org/grpc"
)

type SyncerService struct {
	MultiAddress      identity.MultiAddress
	OrderBook         *orderbook.OrderBook
	MessageQueueLimit int
}

func NewSyncerService(multiAddress identity.MultiAddress, orderbook *orderbook.OrderBook, messageQueueLimit int) *SyncerService {
	return &SyncerService{
		MultiAddress:      multiAddress,
		OrderBook:         orderbook,
		MessageQueueLimit: messageQueueLimit,
	}
}

// Register the SyncerService with a gRPC server.
func (service *SyncerService) Register(server *grpc.Server) {
	RegisterSyncerServer(server, service)
}

func (service *SyncerService) Sync(req *SyncRequest, stream Syncer_SyncServer) error {
	// Use a background MessageQueue to handle the connection until an error
	// is returned by the MessageQueue
	ch := make(chan error, 1)
	quit := make(chan struct{}, 1)
	go func() {
		ch <- service.sync(req, stream, quit)
	}()
	defer close(quit)

	// Select between the context finishing and the background worker
	select {
	case <-stream.Context().Done():
		log.Println("err is ", stream.Context().Err())
		return stream.Context().Err()
	case err := <-ch:
		return err
	}
}

func (service *SyncerService) sync(req *SyncRequest, stream Syncer_SyncServer, quit chan struct{}) error {
	// Verify the sync request
	multiAddress, err := UnmarshalMultiAddress(req.From)
	if err != nil {
		return err
	}

	// Create a MessageQueue that owns this gRPC stream and run it on the
	// Splitter
	messageQueue := NewSyncerServerStreamQueue(stream, service.MessageQueueLimit)
	log.Println(messageQueue)

	// Shutdown the MessageQueue when the quit signal is received
	go func() {
		<-quit
		messageQueue.Shutdown()
	}()
	go func() {
		err := service.OrderBook.SyncHistory(messageQueue)
		if err != nil {
			log.Println(err)
		}
	}()
	return service.OrderBook.Subscribe(multiAddress.Address().String(), messageQueue)
}

type SyncerServerStreamQueue struct {
	stream grpc.Stream
	write  chan *SyncBlock
	quit   chan struct{}
}

func NewSyncerServerStreamQueue(stream Syncer_SyncServer, messageQueueLimit int) SyncerServerStreamQueue {
	return SyncerServerStreamQueue{
		stream: stream,
		write:  make(chan *SyncBlock, messageQueueLimit),
		quit:   make(chan struct{}),
	}
}

func (queue SyncerServerStreamQueue) Run() error {
	return queue.writeAll()
}

func (queue SyncerServerStreamQueue) Shutdown() error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on shutdown: %v", r)
		}
	}()

	close(queue.quit)
	return err
}

func (queue SyncerServerStreamQueue) Send(message dispatch.Message) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on send: %v", r)
		}
	}()

	block, ok := message.(*SyncBlock)
	if !ok {
		return fmt.Errorf("wrong message type, has %T expect *rpc.SyncBlock", message)
	}
	queue.write <- block

	return err
}

func (queue SyncerServerStreamQueue) Recv() (dispatch.Message, bool) {
	panic("read from a syncer server stream queue")
}

func (queue SyncerServerStreamQueue) writeAll() error {
	for {
		select {
		case <-queue.quit:
			return nil
		case message := <-queue.write:
			if err := queue.stream.SendMsg(message); err != nil {
				return err
			}
		}
	}
}

type SyncerClientStreamQueue struct {
	stream grpc.Stream
	read   chan *SyncBlock
	quit   chan struct{}
}

func NewSyncerClientStreamQueue(stream Syncer_SyncClient, messageQueueLimit int) SyncerClientStreamQueue {
	return SyncerClientStreamQueue{
		stream: stream,
		read:   make(chan *SyncBlock, messageQueueLimit),
		quit:   make(chan struct{}),
	}
}

func (queue SyncerClientStreamQueue) Run() error {
	return queue.readAll()
}

func (queue SyncerClientStreamQueue) Shutdown() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on shutdown: %v", r)
		}
	}()

	// If the stream is a client stream, then close the sending channel
	if stream, ok := queue.stream.(Syncer_SyncClient); ok {
		// I hate that Go makes me do this :(  --second that
		err = stream.CloseSend()
	}

	close(queue.quit)
	return nil
}

func (queue SyncerClientStreamQueue) Send(message dispatch.Message) error {
	panic("write to a syncer client stream queue")
}

func (queue SyncerClientStreamQueue) Recv() (dispatch.Message, bool) {
	message, ok := <-queue.read
	if !ok {
		return message, ok
	}

	// Handle different status event
	var status order.Status
	var ord order.Order
	orderType := message.GetOrderBlock()
	switch orderType.(type) {
	case *SyncBlock_Open:
		ord = UnmarshalOrder(message.OrderBlock.(*SyncBlock_Open).Open)
		status = order.Open
	case *SyncBlock_Unconfirmed:
		ord = UnmarshalOrder(message.OrderBlock.(*SyncBlock_Unconfirmed).Unconfirmed)
		status = order.Unconfirmed
	case *SyncBlock_Confirmed:
		ord = UnmarshalOrder(message.OrderBlock.(*SyncBlock_Confirmed).Confirmed)
		status = order.Confirmed
	case *SyncBlock_Settled:
		ord = UnmarshalOrder(message.OrderBlock.(*SyncBlock_Settled).Settled)
		status = order.Settled
	}

	return orderbook.NewMessage(ord, status, nil), true
}

func (queue SyncerClientStreamQueue) readAll() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught on reading all queues: %v", r)
		}
	}()

	defer close(queue.read)
	for {
		message := new(SyncBlock)
		if err = queue.stream.RecvMsg(message); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if message != nil {
			select {
			case <-queue.quit:
				return nil
			case queue.read <- message:
			}
		}
	}
	// todo : why does it not complain about returning nothing ?
}
