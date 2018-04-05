package orderbook

// todo : finish this
type OrderBookDB struct {
}

func NewOrderBookDB() OrderBookDB {
	return OrderBookDB{}
}

func (orderBookDB *OrderBookDB) Open(message Message) {
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Match(message Message) {
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Confirm(message Message) {
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Release(message Message) {
	// TODO: Implement key/value file store
}

func (orderBookDB *OrderBookDB) Settle(message Message) {
	// TODO: Implement key/value file store
}
