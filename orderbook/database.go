package orderbook

import (
	"github.com/republicprotocol/republic-go/order"
	"github.com/syndtr/goleveldb/leveldb"
)

// Database stores the order history in a local db file.
type Database struct {
	*leveldb.DB
}

// NewDatabase creates a new Database struct and stores the db file
// in the given path.
func NewDatabase(path string) (Database, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return Database{}, err
	}
	return Database{
		DB: db,
	}, nil
}

func (database Database) Close() {
	database.Close()
}

func (database Database) Open(ord order.Order) error {
	panic("unimplemented")
}

func (database Database) Match(ord order.Order) error {
	panic("unimplemented")
}

func (database Database) Confirm(ord order.Order) error {
	panic("unimplemented")
}

func (database Database) Release(ord order.Order) error {
	panic("unimplemented")
}

func (database Database) Settle(ord order.Order) error {
	panic("unimplemented")
}

func (database Database) Cancel(ord order.Order) error {
	panic("unimplemented")
}

func (database Database) Blocks() []Entry {
	panic("unimplemented")
}

func (database Database) Order(id order.ID) Entry {
	panic("unimplemented")
}
