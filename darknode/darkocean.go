package darknode

import (
	"bytes"
	"fmt"

	"github.com/republicprotocol/republic-go/identity"
)

// A DarkOcean of Darknodes non-deterministically shuffled into Pools by the
// DarknodeRegistry.
type DarkOcean struct {
	epoch [32]byte
	pools Pools
}

// NewDarkOcean returns a new DarkOcean that uses the Darknodes, and
// DarknodeRegistry to create Pools.
func NewDarkOcean(epoch [32]byte, darknodes [][]byte) DarkOcean {
	return DarkOcean{
		epoch: epoch,
		pools: Pools{},
	}
}

// Epoch returns the [32]byte epoch hash used by the DarkOcean.
func (darkOcean *DarkOcean) Epoch() [32]byte {
	return darkOcean.epoch
}

// Pool returns the Pool that contains the a Darknode with the given ID.
// Returns ErrPoolNotFound if no such Pool can be found.
func (darkOcean *DarkOcean) Pool(id identity.ID) (Pool, error) {
	for i := range darkOcean.pools {
		for j := range darkOcean.pools[i].addresses {
			if bytes.Equal(id[:], darkOcean.pools[i].addresses[i].ID()[:]) {
				return darkOcean.pools[i], nil
			}
		}
	}
	return Pool{}, ErrPoolNotFound
}

// Pools is an alias of a slice.
type Pools []Pool

// A Pool of Darknodes.
type Pool struct {
	id        [32]byte
	addresses identity.Addresses
}

// ID returns the ID of the Pool, as a [32]byte hash of all Darknode IDs.
func (pool *Pool) ID() [32]byte {
	return pool.id
}

// Addresses returns the identity.Addresses of all Darknodes in the Pool.
func (pool *Pool) Addresses() identity.Addresses {
	return pool.addresses
}

// ErrPoolNotFound is returned when no Pool can be found for a given ID.
var ErrPoolNotFound = fmt.Sprintf("pool not found")
