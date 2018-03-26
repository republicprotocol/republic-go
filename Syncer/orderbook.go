package Syncer

import "github.com/republicprotocol/republic-go/order"

type OrderBook struct {
	orderBookCache *OrderBookCache
	orderBookDB *OrderBookDB
	orderBookStreamer *OrderBookStreamer
}

func NewOrderBook(maxConnections int) *OrderBook {
	return &OrderBook{
		orderBookCache: NewOrderBookCache(),
		orderBookDB: NewOrderBookDB(),
		orderBookStreamer: NewOrderBookStreamer(maxConnections),
	}
}

func (orderBook OrderBook) Subscribe(id string , listener chan OrderStatusEvent){
	orderBook.orderBookStreamer.Subscribe(id, listener)
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



