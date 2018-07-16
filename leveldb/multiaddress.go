package leveldb

import (
	"encoding/json"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/swarm"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// MultiAddressValue is the storage format for multiaddresses being stored in
// LevelDB. It contains additional timestamping information so that LevelDB can
// provide pruning.
type MultiAddressValue struct {
	MultiAddress identity.MultiAddress
	Nonce        uint64
}

// MultiAddressesIterator implements the swarm.MultiAddressStorer using a
// LevelDB iterator.
type MultiAddressesIterator struct {
	inner iterator.Iterator
}

func newMultiAddressIterator(iter iterator.Iterator) *MultiAddressesIterator {
	return &MultiAddressesIterator{
		inner: iter,
	}
}

// Next implements the swarm.MultiAddressIterator interface.
func (iter *MultiAddressesIterator) Next() bool {
	return iter.inner.Next()
}

// Cursor implements the swarm.MultiAddressIterator interface.
func (iter *MultiAddressesIterator) Cursor() (identity.MultiAddress, uint64, error) {
	if !iter.inner.Valid() {
		return identity.MultiAddress{}, 0, swarm.ErrCursorOutOfRange
	}
	value := MultiAddressValue{}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &value); err != nil {
		return identity.MultiAddress{}, 0, swarm.ErrCursorOutOfRange
	}
	return value.MultiAddress, value.Nonce, iter.inner.Error()
}

// Collect implements the swarm.MultiAddressIterator interface.
func (iter *MultiAddressesIterator) Collect() ([]identity.MultiAddress, []uint64, error) {
	multiaddresses := []identity.MultiAddress{}
	nonces := []uint64{}
	for iter.Next() {
		multiaddress, nonce, err := iter.Cursor()
		if err != nil {
			return multiaddresses, nonces, err
		}

		multiaddresses = append(multiaddresses, multiaddress)
		nonces = append(nonces, nonce)
	}
	return multiaddresses, nonces, iter.inner.Error()
}

// Release implements the swarm.MultiAddressIterator interface.
func (iter *MultiAddressesIterator) Release() {
	iter.inner.Release()
}

// MultiAddressTable implements the swarm.MultiAddressStorer interface using
// LevelDB.
type MultiAddressTable struct {
	db *leveldb.DB
	// expiry time.Duration
}

// NewMultiAddressTable returns a new MultiAddressTable that uses the
// given LevelDB instance to store and load values from the disk.
func NewMultiAddressTable(db *leveldb.DB) *MultiAddressTable {
	return &MultiAddressTable{db: db}
}

// PutMultiAddress implements the swarm.MultiAddressStorer interface.
func (table *MultiAddressTable) PutMultiAddress(address identity.Address, multiaddress identity.MultiAddress, nonce uint64) (bool, error) {

	// TODO: check if nonce is lower or equal and also if there is any change in the data
	value := MultiAddressValue{
		MultiAddress: multiaddress,
		Nonce:        nonce,
	}
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return true, table.db.Put(table.key([]byte(address.String())), data, nil)
}

// MultiAddress implements the swarm.MultiAddressStorer interface.
func (table *MultiAddressTable) MultiAddress(address identity.Address) (identity.MultiAddress, error) {
	data, err := table.db.Get(table.key([]byte(address.String())), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = swarm.ErrMultiAddressNotFound
		}
		return identity.MultiAddress{}, err
	}
	value := MultiAddressValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return identity.MultiAddress{}, err
	}
	return value.MultiAddress, nil
}

// MultiAddresses implements the swarm.MultiAddressStorer interface.
func (table *MultiAddressTable) MultiAddresses() (swarm.MultiAddressIterator, error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(MultiAddressIterBegin), Limit: table.key(MultiAddressIterEnd)}, nil)
	return newMultiAddressIterator(iter), nil
}

func (table *MultiAddressTable) key(k []byte) []byte {
	return append(append(MultiAddressTableBegin, k...), MultiAddressTablePadding...)
}

// TODO: Pruning and deleting entries must be implemented
