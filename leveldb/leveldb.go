package leveldb

import (
	"encoding/json"
	"path"

	"github.com/republicprotocol/republic-go/ome"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Key prefixes to partition data into tables.
var (
	TableOmeComputations = []byte{0x04, 0x0, 0x0, 0x0}
)

// Key postfixes used by iterators to define table ranges.
var (
	TableIterStart = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	TableIterEnd   = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

// Store is an implementation of storage interfaces using LevelDB to load and
// store data to persistent storage. It uses a single database instance and
// high-order bytes for separating data into tables. All keys are 40 bytes
// long, 8 table prefix bytes and a 32 byte key. LevelDB provides a basic type
// type of in-memory caching but has no optimisations that are specific to the
// data.
type Store struct {
	db *leveldb.DB
}

// NewStore returns a LevelDB implementation of storage interfaces. A call to
// Store.Close is required to free resources allocated by the Store.
func NewStore(dir string) (*Store, error) {
	db, err := leveldb.OpenFile(path.Join(dir, "db"), nil)
	if err != nil {
		return nil, err
	}
	return &Store{
		db: db,
	}, nil
}

// Close the internal LevelDB database.
func (store *Store) Close() error {
	return store.db.Close()
}

// PutComputation implements the ome.Storer interface.
func (store *Store) PutComputation(computation ome.Computation) error {
	data, err := json.Marshal(computation)
	if err != nil {
		return err
	}
	return store.db.Put(append(TableOmeComputations, computation.ID[:]...), data, nil)
}

// DeleteComputation implements the ome.Storer interface.
func (store *Store) DeleteComputation(id ome.ComputationID) error {
	return store.db.Delete(append(TableOmeComputations, id[:]...), nil)
}

// Computation implements the ome.Storer interface.
func (store *Store) Computation(id ome.ComputationID) (ome.Computation, error) {
	computation := ome.Computation{}
	data, err := store.db.Get(append(TableOmeComputations, id[:]...), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = ome.ErrComputationNotFound
		}
		return computation, err
	}
	if err := json.Unmarshal(data, &computation); err != nil {
		return computation, err
	}
	return computation, nil
}

// Computations implements the ome.Storer interface.
func (store *Store) Computations() (ome.ComputationIterator, error) {
	iter := store.db.NewIterator(&util.Range{Start: append(TableOmeComputations, TableIterStart...), Limit: append(TableOmeComputations, TableIterEnd...)}, nil)
	return newComputationIterator(iter), nil
}
