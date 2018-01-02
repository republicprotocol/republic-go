package dht

import (
	"time"

	"github.com/republicprotocol/go-identity"
)

// Constants for use in the DHT.
const (
	IDLengthInBits  = identity.IDLength * 8
	MaxBucketLength = 20
)

// A DHT is a Distributed Hash Table. Each instance has an Address, and several
// Buckets of MultiAddresses that are directly connected to that Address. It
// uses a modified Kademlia approach to storing MultiAddresses in each Bucket
// and favoring old connections over new connections.
type DHT struct {
	Address identity.Address
	Buckets [IDLengthInBits]Bucket
}

// NewDHT returns a new DHT with the given Address, and empty Buckets.
func NewDHT(address identity.Address) DHT {
	return DHT{
		Address: address,
		Buckets: [IDLengthInBits]Bucket{},
	}
}

// Update the target identity.MultiAddress by adding it to its respective
// Bucket. Returns an error if the Bucket is full, or any error that happens
// while finding the respective Bucket.
func (dht *DHT) Update(target identity.MultiAddress) error {
	address, err := target.Address()
	if err != nil {
		return err
	}
	bucket, err := dht.Bucket(address)
	if err != nil {
		return err
	}
	if bucket.IsFull() {
		return NewErrFullBucket()
	}
	bucket[address] = target
	return nil
}

// Find the identity.MultiAddress associated with the target identity.Address.
// Returns nil if the target is not in the DHT, or an error.
func (dht *DHT) Find(target identity.Address) (*identity.MultiAddress, error) {
	bucket, err := dht.Bucket(address)
	if err != nil {
		return nil, err
	}
	return bucket.Find(target), nil
}

// FindNeighborhood returns the identity.MultiAddresses in the same Bucket as
// the target identity.Address. It also returns identity.MultiAddresses in
// Buckets within the neighborhood of the target Bucket.
func (dht *DHT) FindNeighborhood(target identity.Address, neighborhood uint) (identity.MultiAddresses, error) {

	// Find the index range of the neighborhood.
	same, err := dht.Address.SamePrefixLen(target)
	if err != nil {
		return nil, err
	}
	index := len(dht.Buckets) - same - 1
	if index < 0 || index > len(dht.Buckets)-1 {
		return nil, NewErrIndexOutOfRange(index)
	}
	start := index - int(neighborhood)
	if start < 0 {
		start = 0
	}
	end := index + int(neighborhood)
	if end > len(dht.Buckets) {
		end = len(dht.Buckets)
	}

	// Get the total number of identity.MultiAddresses in the neighborhood.
	numMultis := 0
	for i := start; i < end; i++ {
		numMultis += len(dht.Buckets[i])
	}

	// Fill out a perfectly sized slice.
	m := 0
	multis := make(identity.MultiAddresses, numMultis)
	for i := start; i < end; i++ {
		for j := range dht.Buckets[i] {
			multis[m] = dht.Buckets[i][j].MultiAddress
			m++
		}
	}
	return multis, nil
}

// Entries returns all Entries from all Buckets in the DHT.
func (dht *DHT) Entries() Entries {
	numEntries := 0
	for _, bucket := range dht.Buckets {
		numEntries += len(bucket)
	}
	i := 0
	entries := make(Entries, numEntries)
	for _, bucket := range dht.Buckets {
		for j := range bucket {
			entries[i] = bucket[j]
			i++
		}
	}
	return entries
}

// MultiAddresses returns all MultiAddresses from all Buckets in the DHT.
func (dht *DHT) MultiAddresses() identity.MultiAddresses {
	numPeers := 0
	for _, bucket := range dht.Buckets {
		numPeers += len(bucket)
	}
	i := 0
	peers := make(identity.MultiAddresses, numPeers)
	for _, bucket := range dht.Buckets {
		for j := range bucket {
			peers[i] = bucket[j].MultiAddress
			i++
		}
	}
	return peers
}

// Bucket returns the respective Bucket for the target identity.Address, or an
// error.
func (dht *DHT) Bucket(target identity.Address) (*Bucket, error) {
	same, err := dht.Address.SamePrefixLen(target)
	if err != nil {
		return nil, err
	}
	index := len(dht.Buckets) - same - 1
	if index < 0 || index > len(dht.Buckets)-1 {
		return nil, NewErrIndexOutOfRange(index)
	}
	return &dht.Buckets[index], nil
}

// Bucket is a mapping of Addresses to Entries. In standard Kademlia, a list is
// used because Buckets need to be sorted.
type Bucket map[identity.Address]Entry

// Buckets is an alias.
type Buckets []Bucket

// Find a target Address in the Bucket. Returns nil if the target Address
// cannot be found.
func (bucket Bucket) Find(target identity.Address) *identity.MultiAddress {
	multi, ok := bucket[target]
	if !ok {
		return nil
	}
	return &multi.MultiAddress
}

// IsFull returns true if, and only if, the number of Entries in the Bucket is
// equal to the maximum number of Entries allowed.
func (bucket Bucket) IsFull() bool {
	return len(bucket) == MaxBucketLength
}

// An Entry in a Bucket. It holds a MultiAddress, and a timestamp for when it
// was added to the Bucket.
type Entry struct {
	identity.MultiAddress
	time.Time
}

// Entries is an alias.
type Entries []Entry
