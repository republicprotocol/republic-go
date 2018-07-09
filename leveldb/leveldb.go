package leveldb

import (
	"path"
	"time"

	"github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/syndtr/goleveldb/leveldb"
)

// Constants for use in the SomerComputationTable.
var (
	SomerComputationTableBegin   = []byte{0x03, 0x00}
	SomerComputationTablePadding = paddingBytes(0x00, 32)
	SomerComputationIterBegin    = paddingBytes(0x00, 32)
	SomerComputationIterEnd      = paddingBytes(0xFF, 32)
	SomerComputationExpiry       = 72 * time.Hour
)

// Store is an aggregate of all tables that implement storage interfaces. It
// provides access to all of these storage interfaces using different
// underlying LevelDB instances, ensuring that data is shared where possible
// and isolated where needed. For this reason, it is recommended to access all
// storage interfaces through the creation of a Store instance.
type Store struct {
	db *leveldb.DB

	orderbookOrderTable         *OrderbookOrderTable
	orderbookOrderFragmentTable *OrderbookOrderFragmentTable
	orderbookPointerTable       *OrderbookPointerTable

	somerComputationTable *SomerComputationTable
}

// NewStore returns a new Store with a new LevelDB instances that use the
// the given directory as the root for all LevelDB instances. A call to
// Store.Release is needed to ensure that no resources are leaked when
// the Store is no longer needed. Each Store must have a unique directory.
func NewStore(dir string) (*Store, error) {
	db, err := leveldb.OpenFile(path.Join(dir, "db"), nil)
	if err != nil {
		return nil, err
	}
	return &Store{
		db: db,

		orderbookOrderTable:         NewOrderbookOrderTable(db),
		orderbookOrderFragmentTable: NewOrderbookOrderFragmentTable(db),
		orderbookPointerTable:       NewOrderbookPointerTable(),
		orderbookOrderTable:         NewOrderbookOrderTable(db),

		somerComputationTable: NewSomerComputationTable(db),
	}, nil
}

// Release the resources required by the Store.
func (store *Store) Release() error {
	return store.db.Close()
}

// OrderbookOrderStore returns the OrderbookOrderTable used by the Store. It
// implements the orderbook.OrderStorer interface.
func (store *Store) OrderbookOrderStore() orderbook.OrderStorer {
	return store.orderbookOrderTable
}

// OrderbookOrderFragmentStorer returns the OrderbookOrderFragmentTable used by
// the Store. It implements the orderbook.OrderFragmentStorer interface.
func (store *Store) OrderbookOrderFragmentStorer() orderbook.OrderFragmentStorer {
	return store.orderbookOrderFragmentTable
}

// OrderbookPointerStore returns the OrderbookPointerTable used by the Store.
// It implements the orderbook.PointerStorer interface.
func (store *Store) OrderbookPointerStore() orderbook.PointerStorer {
	return store.orderbookPointerTable
}

// SomerComputationStore returns the SomerComputationTable used by the Store.
// It implements the ome.ComputationStorer interface.
func (store *Store) SomerComputationStore() ome.ComputationStorer {
	return store.orderbookPointerTable
}
