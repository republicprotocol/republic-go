package dark

import (
	"errors"
	"math/big"
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/identity"
)

// Ocean of Pools.
type Ocean struct {
	Pools             Pools
	DarkNodeRegistrar dnr.DarkNodeRegistrar
}

// NewOcean uses a DarkNodeRegistrar to read all registered nodes and sort them
// into Pools.
func NewOcean(darkNodeRegistrar dnr.DarkNodeRegistrar) (*Ocean, error) {
	allNodes, err := darkNodeRegistrar.GetAllNodes()
	if err != nil {
		return &Ocean{}, err
	}

	blockhash := /* TODO: Get from contract */ big.NewInt(1234567)
	poolsize := /* TODO: Get from contract? */ 72

	// Find the prime smaller or equal to the number of registered nodes
	// Start at +2 because it has to greater than the maximum (x+1)
	previousPrime := big.NewInt(int64(len(allNodes) + 2))
	// https://golang.org/src/math/big/prime.go
	// ProbablyPrime is 100% accurate for inputs less than 2^64.
	for !previousPrime.ProbablyPrime(0) {
		previousPrime = previousPrime.Sub(previousPrime, big.NewInt(1))
	}

	inverse := blockhash.ModInverse(blockhash, previousPrime)

	// Integer division
	numberOfPools := big.NewInt(0).Div(previousPrime, big.NewInt(int64(poolsize)))
	if numberOfPools.Int64() == 0 {
		numberOfPools = big.NewInt(1)
	}

	pools := make(Pools, numberOfPools.Int64())
	for x := range allNodes {
		// Add one so that
		xPlusOne := big.NewInt(int64(x + 1))
		i := big.NewInt(0).Mod(big.NewInt(0).Mul(xPlusOne, inverse), previousPrime)

		pool := big.NewInt(0).Mod(i, numberOfPools).Int64()
		pools[pool].Append(Node{
			ID:           allNodes[x],
			MultiAddress: nil,
		})
	}

	return &Ocean{
		Pools: pools,
	}, err
}

// FindPool with the given node ID. Returns the Pool, or nil if no Pool can be
// found.
func (ocean *Ocean) FindPool(id identity.ID) *Pool {
	for _, pool := range ocean.Pools {
		if pool.Has(id) != nil {
			return pool
		}
	}
	return nil
}

// WatchForDarkOceanChanges returns a channel through which it will send an update every epoch
// Will check if a new epoch has been triggered and then sleep for 5 minutes
// Blocking function
func WatchForDarkOceanChanges(registrar dnr.DarkNodeRegistrarInterface, channel chan do.Option) {

	// This function runs until the channel is closed
	defer func() { recover() }()

	var currentBlockhash [32]byte

	// TODO loop until an epoch, call calculateDarkOcean()
	for {
		epoch, err := registrar.CurrentEpoch()
		if err != nil {
			channel <- do.Err(errors.New("Couldn't retrieve current epoch"))
		}

		if epoch.Blockhash != currentBlockhash {
			currentBlockhash = epoch.Blockhash
			darkOceanOverlay, err := GetDarkPools(registrar)
			if err != nil {
				channel <- do.Err(errors.New("Couldn't retrieve dark ocean overlay"))
			} else {
				channel <- do.Ok(darkOceanOverlay)
			}
		}
		time.Sleep(5 * 60 * time.Second)
	}
}
