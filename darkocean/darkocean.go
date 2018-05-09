package darkocean

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/identity"
)

// ErrAccessDenied will be returned when number of dark nodes are lesser than
// minimum number of dark nodes required to run a dark pool
var ErrAccessDenied = errors.New("error while trying to connect to the dark ocean")

// A DarkOcean of Darknodes non-deterministically shuffled into Pools by the
// DarknodeRegistry.
type DarkOcean struct {
	epoch [32]byte
	pools Pools
}

// NewDarkOcean returns a new DarkOcean that uses the Darknodes, and
// DarknodeRegistry to create Pools.
func NewDarkOcean(registry *dnr.DarknodeRegistry, epoch [32]byte, darknodes [][]byte) (DarkOcean, error) {
	numberOfNodesInPool, err := registry.MinimumDarkPoolSize()
	if err != nil {
		return DarkOcean{}, err
	}
	if len(darknodes) < int(numberOfNodesInPool.ToBigInt().Int64()) {
		return DarkOcean{}, ErrAccessDenied
	}
	epochVal := big.NewInt(0).SetBytes(epoch[:])
	numberOfDarkNodes := big.NewInt(int64(len(darknodes)))
	x := big.NewInt(0).Mod(epochVal, numberOfDarkNodes)
	positionInOcean := make([]int, len(darknodes))
	for i := 0; i < len(darknodes); i++ {
		positionInOcean[i] = -1
	}
	pools := make([]Pool, (len(darknodes) / int(numberOfNodesInPool.ToBigInt().Int64())))
	for i := 0; i < len(darknodes); i++ {
		isRegistered, err := registry.IsRegistered(darknodes[x.Int64()])
		if err != nil {
			return DarkOcean{}, err
		}
		for !isRegistered || positionInOcean[x.Int64()] != -1 {
			x.Add(x, big.NewInt(1))
			x.Mod(x, numberOfDarkNodes)
			isRegistered, err = registry.IsRegistered(darknodes[x.Int64()])
			if err != nil {
				return DarkOcean{}, err
			}
		}
		positionInOcean[x.Int64()] = i
		poolID := i % len(darknodes) / int(numberOfNodesInPool.ToBigInt().Int64())
		pools[poolID].addresses = append(pools[poolID].addresses, identity.ID(darknodes[x.Int64()]).Address())
		x.Mod(x.Add(x, epochVal), numberOfDarkNodes)
	}

	return DarkOcean{
		epoch: epoch,
		pools: pools,
	}, nil
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
			if bytes.Equal(id[:], darkOcean.pools[i].addresses[j].ID()[:]) {
				return darkOcean.pools[i], nil
			}
		}
	}
	return Pool{}, ErrPoolNotFound
}

// PoolIndex returns the index of the Pool that contains the a Darknode with
// the given ID. Returns -1 if no such Pool can be found.
func (darkOcean *DarkOcean) PoolIndex(id identity.ID) int {
	for i := range darkOcean.pools {
		for j := range darkOcean.pools[i].addresses {
			if bytes.Equal(id[:], darkOcean.pools[i].addresses[j].ID()[:]) {
				return i
			}
		}
	}
	return -1
}

// Pools returns all Pools in the DarkOcean.
func (darkOcean *DarkOcean) Pools() Pools {
	return darkOcean.pools
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

// Size of the Pool.
func (pool *Pool) Size() int {
	return len(pool.addresses)
}

// ErrPoolNotFound is returned when no Pool can be found for a given ID.
var ErrPoolNotFound = errors.New("pool not found")
