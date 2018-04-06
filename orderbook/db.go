package orderbook

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// OrderBookDB store the order history in a local db file.
type OrderBookDB struct {
	*leveldb.DB
}

// NewOrderBookDB creates a new OrderBookDB and store the db file
// in the given path. It use the default
func NewOrderBookDB(path string) (OrderBookDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return OrderBookDB{}, err
	}
	return OrderBookDB{
		DB: db,
	}, nil
}

func (orderBookDB OrderBookDB) Close() {
	orderBookDB.Close()
}

func (orderBookDB OrderBookDB) Open(message *Message) {
	// TODO: Implement key/value file store
}

func (orderBookDB OrderBookDB) Match(message *Message) {
	// TODO: Implement key/value file store
}

func (orderBookDB OrderBookDB) Confirm(message *Message) {
	// TODO: Implement key/value file store
}

func (orderBookDB OrderBookDB) Release(message *Message) {
	// TODO: Implement key/value file store
}

func (orderBookDB OrderBookDB) Settle(message *Message) {
	// TODO: Implement key/value file store
}
