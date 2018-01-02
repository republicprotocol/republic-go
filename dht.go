package dht

import (
	"github.com/republicprotocol/go-identity"
)

const (
	IDLengthInBits  = identity.IDLength * 8
	MaxBucketLength = 20
)

type DHT struct {
	Address identity.Address
	Buckets [IDLengthInBits]Bucket
}

func NewDHT(address identity.Address) *DHT {
	return &DHT{
		Address: address,
		Buckets: [IDLengthInBits]Bucket{},
	}
}

func (dht *DHT) Find(target identity.Address) (*identity.MultiAddress, error) {
	same, err := dht.Address.SamePrefixLen(target)
	if err != nil {
		return nil, err
	}
	index := IDLengthInBits - 1 - same
	if index < 0 || index > IDLengthInBits-1 {
		return nil, NewErrIndexOutOfRange(index)
	}
	return dht.Buckets[index].Find(target), nil
}

func (dht *DHT) MultiAddresses() identity.MultiAddresses {
	numPeers := 0
	for _, bucket := range dht.Buckets {
		numPeers += len(bucket)
	}
	i := 0
	peers := make(identity.MultiAddresses, numPeers)
	for _, bucket := range dht.Buckets {
		for j := range bucket {
			peers[i] = bucket[j]
			i++
		}
	}
	return peers
}

type Bucket map[identity.Address]identity.MultiAddress

func (bucket Bucket) Find(target identity.Address) *identity.MultiAddress {
	multi, ok := bucket[target]
	if !ok {
		return nil
	}
	return &multi
}

func (bucket Bucket) IsFull() bool {
	return len(bucket) == MaxBucketLength
}
