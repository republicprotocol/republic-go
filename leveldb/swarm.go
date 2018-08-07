package leveldb

import (
	"encoding/json"
	"time"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/swarm"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// SwarmMultiAddressValue is the storage format for multiAddresses being stored in
// LevelDB. It contains additional timestamp information so that LevelDB can
// provide pruning.
type SwarmMultiAddressValue struct {
	MultiAddress identity.MultiAddress `json:"multiAddress"`
	Timestamp    time.Time             `json:"timestamp"`
}

// SwarmMultiAddressesIterator implements the swarm.MultiAddressStorer using a
// LevelDB iterator.
type SwarmMultiAddressesIterator struct {
	inner iterator.Iterator
}

func newSwarmMultiAddressIterator(iter iterator.Iterator) *SwarmMultiAddressesIterator {
	return &SwarmMultiAddressesIterator{
		inner: iter,
	}
}

// Next implements the swarm.MultiAddressIterator interface.
func (iter *SwarmMultiAddressesIterator) Next() bool {
	return iter.inner.Next()
}

// Cursor implements the swarm.MultiAddressIterator interface.
func (iter *SwarmMultiAddressesIterator) Cursor() (identity.MultiAddress, error) {
	if !iter.inner.Valid() {
		return identity.MultiAddress{}, swarm.ErrCursorOutOfRange
	}

	value := SwarmMultiAddressValue{}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &value); err != nil {
		return identity.MultiAddress{}, swarm.ErrCursorOutOfRange
	}

	return value.MultiAddress, iter.inner.Error()
}

// Collect implements the swarm.MultiAddressIterator interface.
func (iter *SwarmMultiAddressesIterator) Collect() (identity.MultiAddresses, error) {
	multiaddresses := identity.MultiAddresses{}
	for iter.Next() {
		multiaddress, err := iter.Cursor()
		if err != nil {
			return multiaddresses, err
		}

		multiaddresses = append(multiaddresses, multiaddress)
	}
	return multiaddresses, iter.inner.Error()
}

// Release implements the swarm.MultiAddressIterator interface.
func (iter *SwarmMultiAddressesIterator) Release() {
	iter.inner.Release()
}

// SwarmMultiAddressTable implements the swarm.MultiAddressStorer interface using
// LevelDB.
type SwarmMultiAddressTable struct {
	db     *leveldb.DB
	expiry time.Duration
}

// NewSwarmMultiAddressTable returns a new SwarmMultiAddressTable that uses the
// given LevelDB instance to store and load values from the disk.
func NewSwarmMultiAddressTable(db *leveldb.DB, expiry time.Duration) *SwarmMultiAddressTable {
	return &SwarmMultiAddressTable{db: db, expiry: expiry}
}

// InsertMultiAddress implements the swarm.MultiAddressStorer interface.
func (table *SwarmMultiAddressTable) InsertMultiAddress(multiAddress identity.MultiAddress) error {
	value := SwarmMultiAddressValue{
		MultiAddress: multiAddress,
		Timestamp:    time.Now(),
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return table.db.Put(table.key(multiAddress.Address().Hash()), data, nil)
}

// MultiAddress implements the swarm.MultiAddressStorer interface.
func (table *SwarmMultiAddressTable) MultiAddress(address identity.Address) (identity.MultiAddress, error) {
	data, err := table.db.Get(table.key(address.Hash()), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = swarm.ErrMultiAddressNotFound
		}
		return identity.MultiAddress{}, err
	}
	value := SwarmMultiAddressValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return identity.MultiAddress{}, err
	}
	return value.MultiAddress, nil
}

// MultiAddresses implements the swarm.MultiAddressStorer interface.
func (table *SwarmMultiAddressTable) MultiAddresses() (swarm.MultiAddressIterator, error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(SwarmMultiAddressIterBegin), Limit: table.key(SwarmMultiAddressIterEnd)}, nil)
	return newSwarmMultiAddressIterator(iter), nil
}

// Prune iterates over all multiAddresses and deletes those that have expired.
func (table *SwarmMultiAddressTable) Prune() (err error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(SwarmMultiAddressIterBegin), Limit: table.key(SwarmMultiAddressIterEnd)}, nil)
	defer iter.Release()

	now := time.Now()
	for iter.Next() {
		key := iter.Key()
		value := SwarmMultiAddressValue{}
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

func (table *SwarmMultiAddressTable) key(k []byte) []byte {
	return append(append(SwarmMultiAddressTableBegin, k...), SwarmMultiAddressTablePadding...)
}
