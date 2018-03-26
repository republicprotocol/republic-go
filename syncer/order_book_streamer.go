package syncer

import (
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"time"
	"fmt"
	"sync"
)

type OrderBookSyncer interface {
	Open(ord *order.Order)
	Match(ord *order.Order)
	Confirm(ord *order.Order)
	Release(ord *order.Order)
	Settle(ord *order.Order)
}

type OrderBookStreamer struct {
	subscribersMu *sync.Mutex
	subscribers    map[string]chan *rpc.SyncBlock

	maxConnections int
}

func NewOrderBookStreamer(maxConnection int) OrderBookStreamer {
	return OrderBookStreamer{
		subscribersMu:             new(sync.Mutex),
		subscribers:    map[string]chan *rpc.SyncBlock{},

		maxConnections: maxConnection,
	}
}

func (orderBookStreamer *OrderBookStreamer) Subscribe(id string, listener chan *rpc.SyncBlock) error {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	if len(orderBookStreamer.subscribers) >= orderBookStreamer.maxConnections {
		return fmt.Errorf("cannot subscribe %s: connection limit reached", id)
	}
	orderBookStreamer.subscribers[id] = listener

	return nil
}

func (orderBookStreamer *OrderBookStreamer) Unsubscribe(id string) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	delete(orderBookStreamer.subscribers, id)
}

func (orderBookStreamer *OrderBookStreamer) Open(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		// Allow back-pressure to cause blocking (this is meant to be mitigated
		// by dropping dead clients, or reducing the maximum connections)
		subscriber <- orderToSyncBlock(ord, order.Open)
	}
}

func (orderBookStreamer *OrderBookStreamer) Match(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- orderToSyncBlock(ord, order.Unconfirmed)
	}
}

func (orderBookStreamer *OrderBookStreamer) Confirm(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- orderToSyncBlock(ord, order.Confirmed)
	}
}

func (orderBookStreamer *OrderBookStreamer) Release(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- orderToSyncBlock(ord, order.Open)
	}
}

func (orderBookStreamer *OrderBookStreamer) Settle(ord *order.Order) {
	orderBookStreamer.subscribersMu.Lock()
	defer orderBookStreamer.subscribersMu.Unlock()

	for _, subscriber := range orderBookStreamer.subscribers {
		subscriber <- orderToSyncBlock(ord, order.Settled)
	}
}

func orderToSyncBlock(ord *order.Order, status order.Status) *rpc.SyncBlock{
	block := new(rpc.SyncBlock)
	block.Timestamp = time.Now().Unix()
	block.Signature = []byte{} // todo : will be finished later
	switch status{
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
