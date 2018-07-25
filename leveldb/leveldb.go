package leveldb

import (
	"path"
	"time"

	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/syndtr/goleveldb/leveldb"
)

// Constants for use in the OrderbookOrderTable. Keys in the
// OrderbookOrderTable have a length of 32 bytes, and so 32 bytes of padding is
// needed to ensure that keys are 64 bytes.
var (
	OrderbookOrderTableBegin   = []byte{0x01, 0x00}
	OrderbookOrderTablePadding = paddingBytes(0x00, 32)
	OrderbookOrderIterBegin    = paddingBytes(0x00, 32)
	OrderbookOrderIterEnd      = paddingBytes(0xFF, 32)
)

// Constants for use in the OrderbookOrderFragmentTable. Keys in the
// OrderbookOrderFragmentTable have a length of 64 bytes, 32 bytes for the
// epoch and 32 bytes for the order ID, and so no padding is needed to ensure
// that keys are 64 bytes.
var (
	OrderbookOrderFragmentTableBegin   = []byte{0x02, 0x00}
	OrderbookOrderFragmentTablePadding = paddingBytes(0x00, 0)
	OrderbookOrderFragmentIterBegin    = paddingBytes(0x00, 32)
	OrderbookOrderFragmentIterEnd      = paddingBytes(0xFF, 32)
)

// Constants for use in the OrderbookPointerTable. Keys in the
// OrderbookPointerTable have a length of 0 bytes, and so 64 bytes of padding
// is needed to ensure that keys are 64 bytes.
var (
	OrderbookPointerTableBegin   = []byte{0x03, 0x00}
	OrderbookPointerTablePadding = paddingBytes(0x00, 64)
)

// Constants for use in the SomerComputationTable. Keys in the
// SomerComputationTable have a length of 32 bytes, and so 32 bytes of padding
// is needed to ensure that keys are 64 bytes.
var (
	SomerComputationTableBegin   = []byte{0x10, 0x00}
	SomerComputationTablePadding = paddingBytes(0x00, 32)
	SomerComputationIterBegin    = paddingBytes(0x00, 32)
	SomerComputationIterEnd      = paddingBytes(0xFF, 32)
)

// Constants for use in the SomerBuyOrderFragmentTable. Keys in the
// SomerBuyOrderFragmentTable have a length of 64 bytes, 32 bytes for the
// epoch and 32 bytes for the order ID, and so no padding is needed to ensure
// that keys are 64 bytes.
var (
	SomerBuyOrderFragmentTableBegin   = []byte{0x11, 0x00}
	SomerBuyOrderFragmentTablePadding = paddingBytes(0x00, 0)
	SomerBuyOrderFragmentIterBegin    = paddingBytes(0x00, 32)
	SomerBuyOrderFragmentIterEnd      = paddingBytes(0xFF, 32)
)

// Constants for use in the SomerSellOrderFragmentTable. Keys in the
// SomerSellOrderFragmentTable have a length of 64 bytes, 32 bytes for the
// epoch and 32 bytes for the order ID, and so no padding is needed to ensure
// that keys are 64 bytes.
var (
	SomerSellOrderFragmentTableBegin   = []byte{0x12, 0x00}
	SomerSellOrderFragmentTablePadding = paddingBytes(0x00, 0)
	SomerSellOrderFragmentIterBegin    = paddingBytes(0x00, 32)
	SomerSellOrderFragmentIterEnd      = paddingBytes(0xFF, 32)
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

	somerComputationTable   *SomerComputationTable
	somerOrderFragmentTable *SomerOrderFragmentTable
}

// NewStore returns a new Store with a new LevelDB instances that use the
// the given directory as the root for all LevelDB instances. A call to
// Store.Release is needed to ensure that no resources are leaked when
// the Store is no longer needed. Each Store must have a unique directory.
func NewStore(dir string, expiry time.Duration) (*Store, error) {
	db, err := leveldb.OpenFile(path.Join(dir, "db"), nil)
	if err != nil {
		return nil, err
	}
	return &Store{
		db: db,

		orderbookOrderTable:         NewOrderbookOrderTable(db, expiry),
		orderbookOrderFragmentTable: NewOrderbookOrderFragmentTable(db, expiry),
		orderbookPointerTable:       NewOrderbookPointerTable(db),

		somerComputationTable:   NewSomerComputationTable(db),
		somerOrderFragmentTable: NewSomerOrderFragmentTable(db, expiry),
	}, nil
}

// Release the resources required by the Store.
func (store *Store) Release() error {
	return store.db.Close()
}

// Prune the Store by deleting expired data.
func (store *Store) Prune() (err error) {
	if localErr := store.orderbookOrderTable.Prune(); localErr != nil {
		err = localErr
	}
	if localErr := store.orderbookOrderFragmentTable.Prune(); localErr != nil {
		err = localErr
	}
	if localErr := store.somerComputationTable.Prune(); localErr != nil {
		err = localErr
	}
	if localErr := store.somerOrderFragmentTable.Prune(); localErr != nil {
		err = localErr
	}
	return err
}

// OrderbookOrderStore returns the OrderbookOrderTable used by the Store. It
// implements the orderbook.OrderStorer interface.
func (store *Store) OrderbookOrderStore() orderbook.OrderStorer {
	return store.orderbookOrderTable
}

// OrderbookOrderFragmentStore returns the OrderbookOrderFragmentTable used by
// the Store. It implements the orderbook.OrderFragmentStorer interface.
func (store *Store) OrderbookOrderFragmentStore() orderbook.OrderFragmentStorer {
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
	return store.somerComputationTable
}

// SomerOrderFragmentStore returns the SomerOrderFragmentTable used by the Store.
// It implements the ome.OrderFragmentStorer interface.
func (store *Store) SomerOrderFragmentStore() ome.OrderFragmentStorer {
	return store.somerOrderFragmentTable
}

func paddingBytes(value byte, num int) []byte {
	padding := make([]byte, num)
	for i := range padding {
		padding[i] = value
	}
	return padding
}
