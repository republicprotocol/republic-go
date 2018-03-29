package orderbook

import "github.com/republicprotocol/republic-go/order"

// todo : finish this
type OrderBookDB struct {
}

func NewOrderBookDB() OrderBookDB {
	return OrderBookDB{}
}

func (orderBookDB *OrderBookDB) Open(ord order.Order) {
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Match(ord order.Order) {
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Confirm(ord order.Order) {
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Release(ord order.Order) {
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Settle(ord order.Order) {
	// TODO: Implement key/value file store
}
