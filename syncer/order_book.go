package syncer

import (
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
)

type Broadcaster interface {
	Subscribe(id string, listener chan *rpc.SyncBlock) error
	Unsubscribe(id string)
}

type OrderBook struct {
	orderBookCache OrderBookCache
	orderBookDB OrderBookDB
	orderBookStreamer OrderBookStreamer
}

func NewOrderBook(maxConnections int) *OrderBook {
	return &OrderBook{
		orderBookCache: NewOrderBookCache(),
		orderBookDB: NewOrderBookDB(),
		orderBookStreamer: NewOrderBookStreamer(maxConnections),
	}
}

func (orderBook OrderBook) Subscribe(id string, listener chan *rpc.SyncBlock) error {
	// todo : implement this
	//err := orderBook.orderBookStreamer.Subscribe(id, listener)
	//if err != nil {
	//	return err
	//}
	//blocks := orderBook.orderBookCache.Orders()
	//go func() {
	//	for _, block := range blocks {
	//		listener <- block
	//	}
	//}()
	return nil
}

func (orderBook OrderBook) Unsubscribe(id string){
	orderBook.orderBookStreamer.Unsubscribe(id)
}

func (orderBook OrderBook) Open(ord *order.Order){
	orderBook.orderBookCache.Open(ord)
	orderBook.orderBookDB.Open(ord)
	orderBook.orderBookStreamer.Open(ord)
}

func (orderBook OrderBook) Match(ord *order.Order){
	orderBook.orderBookCache.Match(ord)
	orderBook.orderBookDB.Match(ord)
	orderBook.orderBookStreamer.Match(ord)
}

func (orderBook OrderBook) Confirm(ord *order.Order){
	orderBook.orderBookCache.Confirm(ord)
	orderBook.orderBookDB.Confirm(ord)
	orderBook.orderBookStreamer.Confirm(ord)
}

func (orderBook OrderBook) Release(ord *order.Order){
	orderBook.orderBookCache.Release(ord)
	orderBook.orderBookDB.Release(ord)
	orderBook.orderBookStreamer.Release(ord)
}

func (orderBook OrderBook) Settle(ord *order.Order){
	orderBook.orderBookCache.Settle(ord)
	orderBook.orderBookDB.Settle(ord)
	orderBook.orderBookStreamer.Settle(ord)
}




