package dht

import (
	"fmt"
	"sort"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
)

// ErrFullBucket is used when a peer is inserted into a Bucket that already has
// the maximum number of Entries.
var (
	ErrFullBucket = fmt.Errorf("cannot add entry to a full bucket")
	ErrDHTAddress = fmt.Errorf("cannot use the DHT address in the DHT")
)

// Constants for use in the DHT.
const (
	IDLengthInBits = identity.IDLength * 8
)

// A DHT is a Distributed Hash Table. Each instance has an identity.Address and
// several Buckets of identity.MultiAddresses. These identity.MultiAddresses
// represent peers that are directly connected to the DHT identity.Address. The
// DHT uses a modified Kademlia to store identity.MultiAddresses in each Bucket
// and favors old connections over new connections. The identity.MultiAddresses
// stored in each Bucket can be looked up using an identity.Address allowing
// the to serve as a lookup table for identity.Addresses to
// identity.MultiAddresses. It is safe to use concurrently.
type DHT struct {
	μ       *sync.RWMutex
	Address identity.Address
	Buckets [IDLengthInBits]Bucket
}

// NewDHT returns a new DHT with the given Address, and empty Buckets.
func NewDHT(address identity.Address, maxBucketLength int) DHT {
	dht := DHT{
		μ:       new(sync.RWMutex),
		Address: address,
		Buckets: [IDLengthInBits]Bucket{},
	}
	for i := range dht.Buckets {
		dht.Buckets[i] = NewBucket(maxBucketLength)
	}
	return dht
}

// UpdateMultiAddress by adding it to its respective Bucket. If the
// identity.MultiAddress is already in the Bucket then it is moved to the end
// Returns an error if the Bucket is full, or any error that happens while
// finding the required Bucket.
func (dht *DHT) UpdateMultiAddress(multiAddress identity.MultiAddress) error {
	dht.μ.Lock()
	defer dht.μ.Unlock()
	return dht.updateMultiAddress(multiAddress)
}

// RemoveMultiAddress from its respective Bucket. Nothing happens if the
// identity.MultiAddress is not in the DHT. Returns any error that happens
// while finding the required Bucket.
func (dht *DHT) RemoveMultiAddress(multiAddress identity.MultiAddress) error {
	dht.μ.Lock()
	defer dht.μ.Unlock()
	return dht.removeMultiAddress(multiAddress)
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
	bucket, _, err := dht.findBucket(target)
	return bucket, err
}

// MultiAddresses returns all identity.MultiAddresses in all Buckets.
func (dht *DHT) MultiAddresses() identity.MultiAddresses {
	dht.μ.RLock()
	defer dht.μ.RUnlock()
	return dht.multiAddresses()
}

func (dht *DHT) updateMultiAddress(multiAddress identity.MultiAddress) error {
	address := multiAddress.Address()
	bucket, _, err := dht.findBucket(address)
	if err != nil {
		return err
	}
	return bucket.UpdateMultiAddress(multiAddress)
}

func (dht *DHT) removeMultiAddress(multiAddress identity.MultiAddress) error {
	target := multiAddress.Address()
	bucket, _, err := dht.findBucket(target)
	if err != nil {
		return err
	}
	removeIndex := -1
	for i, multiAddress := range bucket.MultiAddresses {
		address := multiAddress.Address()
		if address == target {
			removeIndex = i
			break
		}
	}
	if removeIndex >= 0 {
		if removeIndex == bucket.Length()-1 {
			bucket.MultiAddresses = bucket.MultiAddresses[:removeIndex]
		} else {
			bucket.MultiAddresses = append(bucket.MultiAddresses[:removeIndex], bucket.MultiAddresses[removeIndex+1:]...)
		}
	}
	return nil
}

func (dht *DHT) findMultiAddress(target identity.Address) (*identity.MultiAddress, error) {
	bucket, _, err := dht.findBucket(target)
	if err != nil {
		return nil, err
	}
	if bucket == nil {
		return nil, nil
	}
	cursor, _ := bucket.FindMultiAddress(target)
	return cursor, nil
}

func (dht *DHT) findMultiAddressNeighbors(target identity.Address, α int) (identity.MultiAddresses, error) {
	// Get all identity.MultiAddresses.
	multiAddresses := dht.MultiAddresses()
	if len(multiAddresses) < α {
		α = len(multiAddresses)
	}

	// Sort them by their distance to the target identity.Address.
	var outerErr error
	sort.SliceStable(multiAddresses, func(i, j int) bool {
		left := multiAddresses[i].Address()
		right := multiAddresses[j].Address()
		closer, err := identity.Closer(left, right, target)
		if outerErr == nil {
			outerErr = err
		}
		return closer
	})
	return multiAddresses[:α], outerErr
}

func (dht *DHT) findBucket(target identity.Address) (*Bucket, int, error) {
	same, err := dht.Address.SamePrefixLength(target)
	if err != nil {
		return nil, -1, err
	}
	if same == IDLengthInBits {
		return nil, -1, ErrDHTAddress
	}
	index := len(dht.Buckets) - same - 1
	if index < 0 || index > len(dht.Buckets)-1 {
		panic("runtime error: index out of range")
	}
	return &dht.Buckets[index], index, nil
}

func (dht *DHT) multiAddresses() identity.MultiAddresses {
	return Buckets(dht.Buckets[:]).MultiAddresses()
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
		MultiAddresses: make(identity.MultiAddresses, 0, maxLength),
		MaxLength:      maxLength,
	}
}

// UpdateMultiAddress adds an identity.MultiAddress to the Bucket. If the
// identity.MultiAddress is already in the Bucket then it is pushed to the end
// of the Bucket.
func (bucket *Bucket) UpdateMultiAddress(multiAddress identity.MultiAddress) error {

	// If the identity.MultiAddress is not already in the Bucket then add it to
	// the Bucket.
	cursor, position := bucket.FindMultiAddress(multiAddress.Address())
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
	bucket.MultiAddresses[bucket.Length()-1] = multiAddress
	return nil
}

// FindMultiAddress finds the identity.MultiAddress associated with a target
// identity.Address in the Bucket. Returns the associated identity.MultiAddress
// and its position in the Bucket. If the target is not in the Bucket then this
// function returns a nil identity.MultiAddress and an invalid position.
func (bucket *Bucket) FindMultiAddress(target identity.Address) (*identity.MultiAddress, int) {
	for i, multiAddress := range bucket.MultiAddresses {
		address := multiAddress.Address()
		if address == target {
			return &multiAddress, i
		}
	}
	return nil, -1
}

// Length returns the number of Entries in the Bucket.
func (bucket *Bucket) Length() int {
	return len(bucket.MultiAddresses)
}

// IsFull returns true if, and only if, the number of Entries in the Bucket is
// equal to the maximum number of Entries allowed.
func (bucket *Bucket) IsFull() bool {
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
