package dark

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/republic-go/stackint"

	"github.com/republicprotocol/go-do"

	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// Ocean of Pools.
type Ocean struct {
	do.GuardedObject

	logger           *logger.Logger
	pools            Pools
	darkNodeRegistry dnr.DarkNodeRegistry
}

// NewOcean uses a DarkNodeRegistry to read all registered nodes and sort them
// into Pools.
func NewOcean(logger *logger.Logger, darkNodeRegistry dnr.DarkNodeRegistry) (*Ocean, error) {
	ocean := &Ocean{
		GuardedObject:    do.NewGuardedObject(),
		logger:           logger,
		pools:            Pools{},
		darkNodeRegistry: darkNodeRegistry,
	}
	return ocean, ocean.update()
}

// FindPool with the given node ID. Returns the Pool, or nil if no Pool can be
// found.
func (ocean *Ocean) FindPool(id identity.ID) *Pool {
	ocean.EnterReadOnly(nil)
	defer ocean.ExitReadOnly()
	for _, pool := range ocean.pools {
		if pool.Has(id) != nil {
			return pool
		}
	}
	return nil
}

// Update updates the dark ocean from the registrar contract
func (ocean *Ocean) Update() error {
	ocean.Enter(nil)
	defer ocean.Exit()
	return ocean.update()
}

func (ocean *Ocean) update() error {
	epoch, err := ocean.darkNodeRegistry.CurrentEpoch()
	if err != nil {
		return err
	}
	blockhash, err := stackint.FromBytes(epoch.Blockhash[:])
	if err != nil {
		return err
	}

	poolsize, err := ocean.darkNodeRegistry.MinimumDarkPoolSize()
	if err != nil {
		return err
	}

	nodeIDs, err := ocean.darkNodeRegistry.GetAllNodes()
	if err != nil {
		return err
	}

	// Find the prime smaller or equal to the number of registered nodes
	// Start at +2 because it has to greater than the maximum (x+1)
	previousPrimeBig := big.NewInt(int64(len(nodeIDs) + 2))

	// ProbablyPrime is 100% accurate for inputs less than 2^64.
	// https://golang.org/src/math/big/prime.go
	for !previousPrimeBig.ProbablyPrime(0) {
		previousPrimeBig = previousPrimeBig.Sub(previousPrimeBig, big.NewInt(1))
	}

	previousPrime, err := stackint.FromBigInt(previousPrimeBig)
	if err != nil {
		return err
	}

	// TODO: This has a bias
	blockhash = blockhash.Mod(&previousPrime)
	if blockhash.IsZero() {
		blockhash = stackint.FromUint(1)
	}

	// Integer division
	numberOfPools := previousPrime.Div(&poolsize)

	if numberOfPools.IsZero() {
		numberOfPools = stackint.FromUint(1)
	}
	poolCount, err := numberOfPools.ToUint()
	if err != nil {
		return err
	}
	pools := make(Pools, poolCount)
	for i := range pools {
		pools[i] = NewPool()
	}

	// Calcualte the pool assignment for each node
	inverse := blockhash.ModInverse(&previousPrime)
	for n := range nodeIDs {
		nPlusOne := stackint.FromUint(uint(n + 1))
		i := nPlusOne.MulModulo(&inverse, &previousPrime)
		poolInt := i.Mod(&numberOfPools)
		pool, err := poolInt.ToUint()
		if err != nil {
			return err
		}

		pools[pool].Append(NewNode(nodeIDs[n]))
	}

	ocean.pools = pools
	return nil
}

// Watch for changes to the Ocean. This function is a blocking function that
// sleeps and wakes once per period to check for a change in epoch. It accepts
// a channel that is pinged whenever the Ocean changes.
func (ocean *Ocean) Watch(changes chan struct{}) {
	// Recover from writing to a closed channel
	defer func() { recover() }()

	minInterval, err := ocean.darkNodeRegistry.MinimumEpochInterval()
	if err != nil {
		ocean.logger.Error(fmt.Sprintf("cannot retrieve minimum epoch interval: %s", err.Error()))
		return
	}

	var currentBlockhash [32]byte
	if err := ocean.Update(); err != nil {
		ocean.logger.Error(fmt.Sprintf("cannot update dark ocean: %s", err.Error()))
		return
	}
	changes <- struct{}{}
	for {
		epoch, err := ocean.darkNodeRegistry.CurrentEpoch()
		if err != nil {
			ocean.logger.Error(fmt.Sprintf("cannot update epoch: %s", err.Error()))
			return
		}
		if !bytes.Equal(currentBlockhash[:], epoch.Blockhash[:]) {
			currentBlockhash = epoch.Blockhash
			if err := ocean.Update(); err != nil {
				ocean.logger.Error(fmt.Sprintf("cannot update dark ocean: %s", err.Error()))
				return
			}
			changes <- struct{}{}
		}

		nextTime := epoch.Timestamp.Add(&minInterval)
		unix, err := nextTime.ToUint()
		if err != nil {
			// Either minInterval is really big, or unix epoch time has overflowed uint64s.
			ocean.logger.Error(fmt.Sprintf("epoch timestamp has overflowed: %s", err.Error()))
			return
		}

		toWait := time.Second * time.Duration(int64(unix)-time.Now().Unix())

		// Wait at least one second
		if toWait < 1*time.Second {
			toWait = 1 * time.Second
		}

		// Try again within a minute
		if toWait > time.Minute {
			toWait = time.Minute
		}

		time.Sleep(toWait)
	}
}
