package leveldb

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/swarm"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// ErrNonceTooLow is returned if the nonce of the multiaddress is lower than the
// one present in the local store.
var ErrNonceTooLow = errors.New("nonce too low")

// MultiAddressValue is the storage format for multiaddresses being stored in
// LevelDB. It contains additional timestamping information so that LevelDB can
// provide pruning.
type MultiAddressValue struct {
	Nonce        uint64
	MultiAddress identity.MultiAddress
	Timestamp    time.Time
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
	db     *leveldb.DB
	expiry time.Duration
}

// NewMultiAddressTable returns a new MultiAddressTable that uses the
// given LevelDB instance to store and load values from the disk.
func NewMultiAddressTable(db *leveldb.DB) *MultiAddressTable {
	return &MultiAddressTable{db: db}
}

// PutMultiAddress implements the swarm.MultiAddressStorer interface.
func (table *MultiAddressTable) PutMultiAddress(multiaddress identity.MultiAddress, nonce uint64) (bool, error) {
	isNew := false
	value := MultiAddressValue{
		Nonce:        nonce,
		MultiAddress: multiaddress,
		Timestamp:    time.Now(),
	}

	oldMultiAddr, oldNonce, err := table.MultiAddress(multiaddress.Address())
	if err != nil {
		isNew = true
	}
	// Return err if nonce is too low
	if oldNonce > nonce {
		return isNew, ErrNonceTooLow
	}
	// If there is a change in the multiaddress stored, then return true
	if err == nil && oldMultiAddr.String() != multiaddress.String() {
		isNew = true
	}

	data, err := json.Marshal(value)
	if err != nil {
		return isNew, err
	}
	return isNew, table.db.Put(table.key(multiaddress.Address().Hash()), data, nil)
}

// MultiAddress implements the swarm.MultiAddressStorer interface.
func (table *MultiAddressTable) MultiAddress(address identity.Address) (identity.MultiAddress, uint64, error) {
	data, err := table.db.Get(table.key(address.Hash()), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = swarm.ErrMultiAddressNotFound
		}
		return identity.MultiAddress{}, 0, err
	}
	value := MultiAddressValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return identity.MultiAddress{}, 0, err
	}
	return value.MultiAddress, value.Nonce, nil
}

// MultiAddresses implements the swarm.MultiAddressStorer interface.
func (table *MultiAddressTable) MultiAddresses() (swarm.MultiAddressIterator, error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(MultiAddressIterBegin), Limit: table.key(MultiAddressIterEnd)}, nil)
	return newMultiAddressIterator(iter), nil
}

func (table *MultiAddressTable) key(k []byte) []byte {
	return append(append(MultiAddressTableBegin, k...), MultiAddressTablePadding...)
}

// Prune iterates over all multiaddresses and deletes those that have expired.
func (table *MultiAddressTable) Prune() (err error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(MultiAddressIterBegin), Limit: table.key(MultiAddressIterEnd)}, nil)
	defer iter.Release()

	now := time.Now()
	for iter.Next() {
		key := iter.Key()
		value := MultiAddressValue{}
		if localErr := json.Unmarshal(iter.Value(), &value); localErr != nil {
			err = localErr
			continue
		}
		if value.Timestamp.Add(table.expiry).Before(now) {
			if localErr := table.db.Delete(key, nil); localErr != nil {
				err = localErr
			}
		}
	}
	return err
}
