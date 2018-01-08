package dht

import (
	"sort"
	"time"

	"github.com/republicprotocol/go-identity"
)

// Constants for use in the DHT.
const (
	MaxBucketSize  = 20
	IDLengthInBits = identity.IDLength * 8
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

// Update an identity.MultiAddress by adding it to its respective Bucket.
// Returns an error if the Bucket is full, or any error that happens while
// finding the required Bucket.
func (dht *DHT) Update(multi identity.MultiAddress) error {
	target, err := multi.Address()
	if err != nil {
		return err
	}
	bucket, err := dht.Bucket(target)
	if err != nil {
		return err
	}
	if bucket.IsFull() {
		return ErrFullBucket
	}
	*bucket = append(*bucket, Entry{multi, time.Now()})
	return nil
}

// Find the identity.MultiAddress associated with the target identity.Address.
// Returns nil if the target is not in the DHT, or an error.
func (dht *DHT) Find(target identity.Address) (*identity.MultiAddress, error) {
	bucket, err := dht.Bucket(target)
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
	same, err := dht.Address.SamePrefixLength(target)
	if err != nil {
		return nil, err
	}
	index := len(dht.Buckets) - same - 1
	if index < 0 || index > len(dht.Buckets)-1 {
		panic("runtime error: index out of range")
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
	i := 0
	multis := make(identity.MultiAddresses, numMultis)
	for _, bucket := range dht.Buckets[start:end] {
		for _, entry := range bucket {
			multis[i] = entry.MultiAddress
			i++
		}
	}
	return multis, nil
}

// FindBucket uses the target identity.Address and returns the respective
// Bucket. The target does not have to be in the DHT. Returns the Bucket, or an
// error.
func (dht *DHT) FindBucket(target identity.Address) (*Bucket, error) {
	return dht.Bucket(target)
}

// FindNeighborhoodBuckets uses the target identity.Address to find Buckets
// within a given neighborhood of the target Bucket. It does not include the
// actual target Bucket, which can be found using FindBucket. The target does
// not have to be in the DHT. Returns the Buckets, or an error.
func (dht *DHT) FindNeighborhoodBuckets(target identity.Address) (*Bucket, error) {
	return dht.Bucket(target)
}

// MultiAddresses returns all MultiAddresses from all Buckets in the DHT.
func (dht *DHT) MultiAddresses() identity.MultiAddresses {
	numMultis := 0
	for _, bucket := range dht.Buckets {
		numMultis += len(bucket)
	}
	i := 0
	multis := make(identity.MultiAddresses, numMultis)
	for _, bucket := range dht.Buckets {
		for _, entry := range bucket {
			multis[i] = entry.MultiAddress
			i++
		}
	}
	return multis
}

// Bucket returns the respective Bucket for the target identity.Address, or an
// error.
func (dht *DHT) Bucket(target identity.Address) (*Bucket, error) {
	same, err := dht.Address.SamePrefixLength(target)
	if err != nil {
		return nil, err
	}
	index := len(dht.Buckets) - same - 1
	if index < 0 || index > len(dht.Buckets)-1 {
		panic("runtime error: index out of range")
	}
	return &dht.Buckets[index], nil
}

// Bucket is a mapping of Addresses to Entries. In standard Kademlia, a list is
// used because Buckets need to be sorted.
type Bucket []Entry

// Find a target Address in the Bucket. Returns nil if the target Address
// cannot be found.
func (bucket Bucket) Find(target identity.Address) *identity.MultiAddress {
	for _, entry := range bucket {
		address, err := entry.MultiAddress.Address()
		if err == nil && address == target {
			return &entry.MultiAddress
		}
	}
	return nil
}

// Sort the Bucket by the time at which Entries were added.
func (bucket Bucket) Sort() {
	sort.Slice(bucket, func(i, j int) bool {
		return bucket[i].Time.Before(bucket[j].Time)
	})
}

// IsFull returns true if, and only if, the number of Entries in the Bucket is
// equal to the maximum number of Entries allowed.
func (bucket Bucket) IsFull() bool {
	return len(bucket) == MaxBucketSize
}

// An Entry in a Bucket. It holds a MultiAddress, and a timestamp for when it
// was added to the Bucket.
type Entry struct {
	identity.MultiAddress
	time.Time
}
