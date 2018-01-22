package dht

import (
	"sort"
	"sync"
	"time"

	"github.com/republicprotocol/go-identity"
)

// Constants for use in the DHT.
const (
	IDLengthInBits = identity.IDLength * 8
)

// A DHT is a Distributed Hash Table. Each instance has an identity.Address and
// several Buckets of identity.MultiAddresses that are directly connected to
// that identity.Address. It uses a modified Kademlia approach to storing
// identity.MultiAddresses in each Bucket, favoring old connections over new
// connections. It is safe to use concurrently.
type DHT struct {
	μ       *sync.RWMutex
	Address identity.Address
	Buckets [IDLengthInBits]Bucket
}

// NewDHT returns a new DHT with the given Address, and empty Buckets.
func NewDHT(address identity.Address, maxBucketLength int) *DHT {
	dht := &DHT{
		μ:       new(sync.RWMutex),
		Address: address,
		Buckets: [IDLengthInBits]Bucket{},
	}
	for i := range dht.Buckets {
		dht.Buckets[i] = NewBucket(maxBucketLength)
	}
	return dht
}

// Update an identity.MultiAddress by adding it to its respective Bucket.
// Returns an error if the Bucket is full, or any error that happens while
// finding the required Bucket. If the identity.MultiAddress is already in the
// Bucket then it is put at the top.
func (dht *DHT) Update(multiAddress identity.MultiAddress) error {
	dht.μ.Lock()
	defer dht.μ.Unlock()
	return dht.update(multiAddress)
}

// Remove an identity.MultiAddress by removing it from its respective Bucket.
// Nothing happens if the identity.MultiAddress is not in the DHT. Returns any
// error that happens while finding the required Bucket.
func (dht *DHT) Remove(multi identity.MultiAddress) error {
	dht.μ.Lock()
	defer dht.μ.Unlock()
	return dht.remove(multi)
}

// FindMultiAddress finds the identity.MultiAddress associated with the target
// identity.Address. Returns nil if the target is not in the DHT, or an error.
func (dht *DHT) FindMultiAddress(target identity.Address) (*identity.MultiAddress, error) {
	dht.μ.RLock()
	defer dht.μ.RUnlock()
	return dht.findMultiAddress(target)
}

// FindMultiAddressNeighbors finds the closest identity.MultiAddresses to the
// target identity.Address. Returns up to α identity.MultiAddresses, or an
// error.
func (dht *DHT) FindMultiAddressNeighbors(target identity.Address, α int) (identity.MultiAddresses, error) {
	dht.μ.RLock()
	defer dht.μ.RUnlock()
	return dht.findMultiAddressNeighbors(target, α)
}

// FindBucket uses the target identity.Address and returns the respective
// Bucket. The target does not have to be in the DHT. Returns the Bucket, or an
// error.
func (dht *DHT) FindBucket(target identity.Address) (*Bucket, error) {
	dht.μ.RLock()
	defer dht.μ.RUnlock()
	return dht.findBucket(target)
}

// FindBucketNeighbors uses the target identity.Address to find Buckets
// that are close to the target Bucket. The target does not have to be in the
// DHT. Returns up to α Bucket, or an error.
func (dht *DHT) FindBucketNeighbors(target identity.Address, α int) (Buckets, error) {
	dht.μ.RLock()
	defer dht.μ.RUnlock()
	return dht.findBucketNeighbors(target, α)
}

// Neighborhood returns the start and end indices of a α-sized neighborhood
// around the Bucket associated with the target identity.Address.
func (dht *DHT) Neighborhood(target identity.Address, α int) (int, int, error) {
	dht.μ.RLock()
	defer dht.μ.RUnlock()
	return dht.neighborhood(target, α)
}

// MultiAddresses returns all identity.MultiAddresses in all Buckets.
func (dht *DHT) MultiAddresses() identity.MultiAddresses {
	dht.μ.RLock()
	defer dht.μ.RUnlock()
	return dht.multiAddresses()
}

func (dht *DHT) update(multiAddress identity.MultiAddress) error {
	target, err := multiAddress.Address()
	if err != nil {
		return err
	}
	bucket, err := dht.findBucket(target)
	if err != nil {
		return err
	}
	bucket.UpdateEntry(target)

	// Remove the target if it is already in the Bucket.
	
	prevMultiAddress := bucket.FindMultiAddress(target)
	if prevMultiAddress != nil {
		pivot := -1
		for i, entry := range bucket.Entries {
			address, err := entry.MultiAddress.Address()
			if err != nil {
				return err
			}
			if address == target {
				pivot = i
				break
			}
		}
		if pivot >= 0 {
			for i := pivot + 1; i < 
		}
	}

	if bucket.IsFull() {
		return ErrFullBucket
	}
	bucket.Entries = append(bucket.Entries, Entry{multi, time.Now()})
	return nil
}

func (dht *DHT) remove(multi identity.MultiAddress) error {
	target, err := multi.Address()
	if err != nil {
		return err
	}
	bucket, err := dht.findBucket(target)
	if err != nil {
		return err
	}
	removeIndex := -1
	for i, entry := range bucket.Entries {
		address, err := entry.MultiAddress.Address()
		if err != nil {
			return err
		}
		if address == target {
			removeIndex = i
			break
		}
	}
	if removeIndex >= 0 {
		if removeIndex == bucket.Length()-1 {
			bucket.Entries = bucket.Entries[:removeIndex]
		} else {
			bucket.Entries = append(bucket.Entries[:removeIndex], bucket.Entries[removeIndex+1:]...)
		}
	}
	return nil
}

func (dht *DHT) findMultiAddress(target identity.Address) (*identity.MultiAddress, error) {
	bucket, err := dht.findBucket(target)
	if err != nil {
		return nil, err
	}
	return bucket.FindMultiAddress(target), nil
}

func (dht *DHT) findMultiAddressNeighbors(target identity.Address, α int) (identity.MultiAddresses, error) {
	return identity.MultiAddresses{}, nil
}

func (dht *DHT) findBucket(target identity.Address) (*Bucket, error) {
	same, err := dht.Address.SamePrefixLength(target)
	if err != nil {
		return nil, err
	}
	if same == IDLengthInBits {
		return nil, ErrDHTAddress
	}
	index := len(dht.Buckets) - same - 1
	if index < 0 || index > len(dht.Buckets)-1 {
		panic("runtime error: index out of range")
	}
	return &dht.Buckets[index], nil
}

func (dht *DHT) findBucketNeighbors(target identity.Address, α int) (Buckets, error) {
	// Find the index range of the neighborhood.
	start, end, err := dht.neighborhood(target, α)
	if err != nil {
		return nil, err
	}
	return dht.Buckets[start:end], nil
}

func (dht *DHT) neighborhood(target identity.Address, α int) (int, int, error) {
	// Find the index range of the neighborhood.
	same, err := dht.Address.SamePrefixLength(target)
	if err != nil {
		return -1, -1, err
	}
	if same == IDLengthInBits {
		return -1, -1, ErrDHTAddress
	}
	index := len(dht.Buckets) - same - 1
	if index < 0 || index > len(dht.Buckets)-1 {
		panic("runtime error: index out of range")
	}
	start := index - α
	if start < 0 {
		start = 0
	}
	end := index + α
	if end > len(dht.Buckets) {
		end = len(dht.Buckets)
	}
	return start, end, nil
}

func (dht *DHT) multiAddresses() identity.MultiAddresses {
	numMultis := 0
	for _, bucket := range dht.Buckets {
		numMultis += bucket.Length()
	}
	i := 0
	multis := make(identity.MultiAddresses, numMultis)
	for _, bucket := range dht.Buckets {
		for _, entry := range bucket.Entries {
			multis[i] = entry.MultiAddress
			i++
		}
	}
	return multis
}

// Bucket is a mapping of Addresses to Entries. In standard Kademlia, a list is
// used because Buckets need to be sorted.
type Bucket struct {
	identity.MultiAddresses
	MaxLength int
}

// NewBucket returns a new Bucket with an empty set of Entries that can be, at
// most, the given maximum length.
func NewBucket(maxLength int) Bucket {
	return Bucket{
		MaxLength: maxLength,
	}
}

// UpdateMultiAddress adds an identity.MultiAddress to the Bucket. If the
// identity.MultiAddress is already in the Bucket then it is pushed to the end
// of the Bucket.
func (bucket Bucket) UpdateMultiAddress(multiAddress identity.MultiAddress) error {

	// If the identity.MultiAddress is not already in the Bucket then add it to
	// the Bucket.
	cursor, position := bucket.FindMultiAddress(multiAddress)
	if cursor == nil {
		if bucket.IsFull() {
			return ErrFullBucket
		}
		bucket.MultiAddresses = append(bucket.MultiAddresses, multiAddress)
		return nil
	}

	// Otherwise, move the identity.MultiAddress to the end of the Bucket.
	for i := position + 1; i < bucket.Length(); i++ {
		bucket.MultiAddresses[i-1] = bucket.MultiAddresses[i]
	}
	bucket.MultiAddresses[bucket.Length()-1] = *cursor
	return nil
}

// FindMultiAddress finds the identity.MultiAddress associated with a target
// identity.Address in the Bucket. Returns the associated identity.MultiAddress
// and its position in the Bucket. If the target is not in the Bucket then this
// function returns a nil identity.MultiAddress and an invalid position.
func (bucket Bucket) FindMultiAddress(target identity.Address) (*identity.MultiAddress, int) {
	for i, multiAddress := range bucket.MultiAddresses {
		address, err := multiAddress.Address()
		if err == nil && address == target {
			return &multiAddress, i
		}
	}
	return nil, -1
}

// MultiAddresses returns all MultiAddresses in the Bucket.
func (bucket Bucket) MultiAddresses() identity.MultiAddresses {
	return bucket.MultiAddresses
}

// Length returns the number of Entries in the Bucket.
func (bucket Bucket) Length() int {
	return len(bucket.MultiAddresses)
}

// IsFull returns true if, and only if, the number of Entries in the Bucket is
// equal to the maximum number of Entries allowed.
func (bucket Bucket) IsFull() bool {
	return bucket.Length() == bucket.MaxLength
}

// Buckets is an alias.
type Buckets []Bucket

// MultiAddresses returns all MultiAddresses from all Buckets.
func (buckets Buckets) MultiAddresses() identity.MultiAddresses {
	numberOfMultiAddresses := 0
	for _, bucket := range buckets {
		numberOfMultiAddresses += bucket.Length()
	}
	i := 0
	multiAddresses := make(identity.MultiAddresses, numberOfMultiAddresses)
	for _, bucket := range buckets {
		for _, multiAddress := range bucket.MultiAddresses {
			multiAddresses[i] = multiAddress
			i++
		}
	}
	return multiAddresses
}