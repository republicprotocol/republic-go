package syncer

import (
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

type OrderBookSyncer interface {
	Open(ord *order.Order)
	Match(ord *order.Order)
	Confirm(ord *order.Order)
	Release(ord *order.Order)
	Settle(ord *order.Order)
}

type OrderBookStreamer struct {
	dispatch.Splitter

	maxConnections int
}

func NewOrderBookStreamer(maxConnection int) OrderBookStreamer {
	return OrderBookStreamer{
		Splitter:       dispatch.NewSplitter(),
		maxConnections: maxConnection,
	}
}

func (orderBookStreamer *OrderBookStreamer) Subscribe(id string, stream rpc.Dark_SyncServer) error {

	if orderBookStreamer.Splitter.CurrentConnections() >= orderBookStreamer.maxConnections {
		return fmt.Errorf("cannot subscribe %s: connection limit reached", id)
	}

	messageQueue := NewSyncMessageQueue(stream)

	return orderBookStreamer.Splitter.RunMessageQueue(id, messageQueue)
}

func (orderBookStreamer *OrderBookStreamer) Unsubscribe(id string) {
	orderBookStreamer.Splitter.ShutdownMessageQueue(id)
}

func (orderBookStreamer *OrderBookStreamer) Open(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Open))
}

func (orderBookStreamer *OrderBookStreamer) Match(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Unconfirmed))
}

func (orderBookStreamer *OrderBookStreamer) Confirm(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Confirmed))
}

func (orderBookStreamer *OrderBookStreamer) Release(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Open))
}

func (orderBookStreamer *OrderBookStreamer) Settle(ord *order.Order) {
	orderBookStreamer.Splitter.Send(orderToSyncBlock(ord, order.Settled))
}

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
