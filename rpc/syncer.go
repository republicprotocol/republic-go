package rpc

import (
	"fmt"
	"io"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"google.golang.org/grpc"
)

// SyncerOptions defines the options specifically for syncer service
type SyncerOptions struct {
	MaxConnections int `json:"maxConnections"`
}

// SyncerService implements the syncer gRPC service. It creates
// MessageQueues for each gRPC stream.
type SyncerService struct {
	Options

	Logger    *logger.Logger
	OrderBook *orderbook.OrderBook
}

// NewSyncerService creates a SyncerService with the given options.
func NewSyncerService(options Options, logger *logger.Logger, orderbook *orderbook.OrderBook) *SyncerService {
	return &SyncerService{
		Options:   options,
		Logger:    logger,
		OrderBook: orderbook,
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

	// Create a MessageQueue that owns this gRPC stream
	messageQueue := NewSyncerServerStreamQueue(stream, service.MessageQueueLimit)

	// Shutdown the MessageQueue when the quit signal is received
	go func() {
		<-quit
		messageQueue.Shutdown()
	}()

	// Subscribe to the orderbook to received updates
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

	msg, ok := message.(*orderbook.Message)
	if !ok {
		return fmt.Errorf("wrong message type, has %T expect *orderbook.Message", message)
	}
	queue.write <- ToBlock(msg)
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
	var epochHash [32]byte
	copy(epochHash[:], message.EpochHash)
	return orderbook.NewMessage(ord, status, epochHash), true
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

func ToBlock(message *orderbook.Message) *SyncBlock {
	syncBlock := &SyncBlock{
		Signature: message.Ord.Signature,
		Timestamp: time.Now().Unix(),
		EpochHash: message.EpochHash[:],
	}
	switch message.Status {
	case order.Open:
		syncBlock.OrderBlock = &SyncBlock_Open{
			Open: MarshalOrder(&message.Ord),
		}
	case order.Canceled:
		syncBlock.OrderBlock = &SyncBlock_Canceled{
			Canceled: MarshalOrder(&message.Ord),
		}
	case order.Unconfirmed:
		syncBlock.OrderBlock = &SyncBlock_Unconfirmed{
			Unconfirmed: MarshalOrder(&message.Ord),
		}
	case order.Confirmed:
		syncBlock.OrderBlock = &SyncBlock_Confirmed{
			Confirmed: MarshalOrder(&message.Ord),
		}
	case order.Settled:
		syncBlock.OrderBlock = &SyncBlock_Settled{
			Settled: MarshalOrder(&message.Ord),
		}
	}
	return syncBlock
}
