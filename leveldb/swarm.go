package leveldb

import (
	"encoding/json"
	"errors"
	"log"
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

// SwarmMultiAddressValue is the storage format for multiaddresses being stored in
// LevelDB. It contains additional timestamping information so that LevelDB can
// provide pruning.
type SwarmMultiAddressValue struct {
	Nonce        uint64                `json:"multiAddressNonce"`
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
func (iter *SwarmMultiAddressesIterator) Cursor() (identity.MultiAddress, uint64, error) {
	if !iter.inner.Valid() {
		return identity.MultiAddress{}, 0, swarm.ErrCursorOutOfRange
	}

	value := SwarmMultiAddressValue{}
	data := iter.inner.Value()
	if err := json.Unmarshal(data, &value); err != nil {
		return identity.MultiAddress{}, 0, swarm.ErrCursorOutOfRange
	}
	log.Printf("cursor returns ..... %v", value)

	return value.MultiAddress, value.Nonce, iter.inner.Error()
}

// Collect implements the swarm.MultiAddressIterator interface.
func (iter *SwarmMultiAddressesIterator) Collect() ([]identity.MultiAddress, []uint64, error) {
	multiaddresses := []identity.MultiAddress{}
	nonces := []uint64{}
	for iter.Next() {
		multiaddress, nonce, err := iter.Cursor()
		if err != nil {
			return multiaddresses, nonces, err
		}

		log.Printf("iterating..... got %v", multiaddress)

		multiaddresses = append(multiaddresses, multiaddress)
		nonces = append(nonces, nonce)
	}
	return multiaddresses, nonces, iter.inner.Error()
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

// PutMultiAddress implements the swarm.MultiAddressStorer interface.
func (table *SwarmMultiAddressTable) PutMultiAddress(multiAddress identity.MultiAddress, nonce uint64) (bool, error) {
	isNew := false

	value := SwarmMultiAddressValue{
		Nonce:        nonce,
		MultiAddress: multiAddress,
		Timestamp:    time.Now(),
	}
	oldMultiAddr, oldNonce, err := table.MultiAddress(multiAddress.Address())
	if err != nil && err == swarm.ErrMultiAddressNotFound {
		isNew = true
	}
	// Return err if nonce is too low
	if oldNonce > nonce {
		return isNew, ErrNonceTooLow
	}
	// If there is a change in the multiaddress stored, then return true
	if err == nil && oldMultiAddr.String() != multiAddress.String() {
		isNew = true
	}

	data, err := json.Marshal(value)
	if err != nil {
		return isNew, err
	}
	log.Printf("storing........ %v", value.MultiAddress)
	return isNew, table.db.Put(table.key(multiAddress.Address().Hash()), data, nil)
}

// MultiAddress implements the swarm.MultiAddressStorer interface.
func (table *SwarmMultiAddressTable) MultiAddress(address identity.Address) (identity.MultiAddress, uint64, error) {
	data, err := table.db.Get(table.key(address.Hash()), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = swarm.ErrMultiAddressNotFound
		}
		return identity.MultiAddress{}, 0, err
	}
	value := SwarmMultiAddressValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return identity.MultiAddress{}, 0, err
	}
	return value.MultiAddress, value.Nonce, nil
}

// MultiAddresses implements the swarm.MultiAddressStorer interface.
func (table *SwarmMultiAddressTable) MultiAddresses() (swarm.MultiAddressIterator, error) {
	iter := table.db.NewIterator(&util.Range{Start: table.key(SwarmMultiAddressIterBegin), Limit: table.key(SwarmMultiAddressIterEnd)}, nil)
	return newSwarmMultiAddressIterator(iter), nil
}

// Prune iterates over all multiaddresses and deletes those that have expired.
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
