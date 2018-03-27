package orderbook

import (
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

// An OrderBookStreamer will send all the new orders received to the stream.
// It only allow certain amount of connections to it.
type OrderBookStreamer struct {
	dispatch.Splitter

	maxConnections int
}

// NewOrderBookStreamer creates a new OrderBookStreamer by the given limit
// of connections
func NewOrderBookStreamer(maxConnection int) OrderBookStreamer {
	return OrderBookStreamer{
		Splitter:       dispatch.NewSplitter(),
		maxConnections: maxConnection,
	}
}

// Subscribe will add the provided stream as a listener listening for updates.
// It returns an error if max connections have been reached.
func (orderBookStreamer *OrderBookStreamer) Subscribe(id string, stream rpc.Dark_SyncServer) error {

	if orderBookStreamer.Splitter.CurrentConnections() >= orderBookStreamer.maxConnections {
		return fmt.Errorf("cannot subscribe %s: connection limit reached", id)
	}

	messageQueue := NewSyncMessageQueue(stream)

	return orderBookStreamer.Splitter.RunMessageQueue(id, messageQueue)
}

// Unsubscribe will stop listening for updates.
func (orderBookStreamer *OrderBookStreamer) Unsubscribe(id string) {
	orderBookStreamer.Splitter.ShutdownMessageQueue(id)
}

// Open notifies its subscribers that status of an order has been changed
// to 'open'
func (orderBookStreamer *OrderBookStreamer) Open(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Open))
}

// Match notifies its subscribers that status of an order has been changed
// to 'unconfirmed'
func (orderBookStreamer *OrderBookStreamer) Match(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Unconfirmed))
}

// Confirm notifies its subscribers that status of an order has been changed
// to 'confirmed'
func (orderBookStreamer *OrderBookStreamer) Confirm(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Confirmed))
}

// Release notifies its subscribers that status of an order has been changed
// to 'open'
func (orderBookStreamer *OrderBookStreamer) Release(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Open))
}

// Settle notifies its subscribers that status of an order has been changed
// to 'settled'
func (orderBookStreamer *OrderBookStreamer) Settle(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Settled))
}

// orderToSyncBlock convert an order and its new status to a SyncBlock in
// gRPC representation.
func orderToSyncBlock(ord *order.Order, status order.Status) *rpc.SyncBlock {
	block := new(rpc.SyncBlock)
	block.Timestamp = time.Now().Unix()
	block.Signature = []byte{} // todo : will be finished later
	switch status {
	case order.Open:
		block.OrderBlock = &rpc.SyncBlock_Open{
			Open: rpc.SerializeOrder(ord),
		}
	case order.Unconfirmed:
		block.OrderBlock = &rpc.SyncBlock_Unconfirmed{
			Unconfirmed: rpc.SerializeOrder(ord),
		}
	case order.Confirmed:
		block.OrderBlock = &rpc.SyncBlock_Confirmed{
			Confirmed: rpc.SerializeOrder(ord),
		}
	case order.Settled:
		block.OrderBlock = &rpc.SyncBlock_Settled{
			Settled: rpc.SerializeOrder(ord),
		}
	default:
		return nil
	}

	return block
}
