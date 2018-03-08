package darkocean

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/identity"
)

// DarkPool is a list of node multiaddresses
type DarkPool struct {
	do.GuardedObject

	Nodes identity.MultiAddresses
}

// Add will append a multiaddress to the darkpool
func (darkpool DarkPool) Add(targets ...identity.MultiAddress) {
	darkpool.Enter(nil)
	defer darkpool.Exit()
	darkpool.Nodes = append(darkpool.Nodes, targets...)
}

// IDDarkPool is a list of node ids
type IDDarkPool []identity.ID

// Overlay contains a list of dark pools
type Overlay struct {
	Pools []IDDarkPool
}

// FindDarkPool returns the pool containing a prticular ID
func (ocean Overlay) FindDarkPool(id identity.ID) (IDDarkPool, error) {

	for _, pool := range ocean.Pools {
		for _, node := range pool {
			if bytes.Compare(node, id) == 0 {
				return pool, nil
			}
		}
	}

	return nil, errors.New("Node is not a part of a pool")
}

// GetDarkPools gets the full list of nodes and sorts them into pools
func GetDarkPools(darkNodeRegistrar *dnr.DarkNodeRegistrar) (*Overlay, error) {
	allNodes, err := darkNodeRegistrar.GetAllNodes()
	if err != nil {
		return &Overlay{}, err
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

	pools := make([]IDDarkPool, numberOfPools.Int64())

	for x := range allNodes {
		// Add one so that
		xPlusOne := big.NewInt(int64(x + 1))
		i := big.NewInt(0).Mod(big.NewInt(0).Mul(xPlusOne, inverse), previousPrime)

		assignedPool := big.NewInt(0).Mod(i, numberOfPools).Int64()

		pools[assignedPool] = append(pools[assignedPool], allNodes[x][:])
	}

	return &Overlay{
		Pools: pools,
	}, err
}
