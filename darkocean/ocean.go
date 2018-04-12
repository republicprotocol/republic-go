package darkocean

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/identity"
)

// Ocean of Pools.
type Ocean struct {
	do.GuardedObject

	pools             Pools
	darkNodeRegistrar contracts.DarkNodeRegistry
}

// Oceans is a list of dark Oceans .
type Oceans []Ocean

// NewOcean uses a DarkNodeRegistry to read all registered nodes and sort them
// into Pools.
func NewOcean(darkNodeRegistrar contracts.DarkNodeRegistry) (*Ocean, error) {
	ocean := &Ocean{
		GuardedObject:     do.NewGuardedObject(),
		pools:             Pools{},
		darkNodeRegistrar: darkNodeRegistrar,
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
	epoch, err := ocean.darkNodeRegistrar.CurrentEpoch()
	if err != nil {
		return fmt.Errorf("cannot get current epoch :%v", err)
	}

	nodeIDs, err := ocean.darkNodeRegistrar.GetAllNodes()
	if err != nil {
		return fmt.Errorf("cannot get all nodes: %v", err)
	}

	nodePositionHashesToIDs := map[string][]byte{}
	nodePositionHashes := make([][]byte, len(nodeIDs))
	for i := range nodeIDs {
		nodePositionHashes[i] = crypto.Keccak256(epoch.Blockhash[:], nodeIDs[i])
		nodePositionHashesToIDs[string(nodePositionHashes[i])] = nodeIDs[i]
	}

	sort.Slice(nodePositionHashes, func(i, j int) bool {
		return bytes.Compare(nodePositionHashes[i], nodePositionHashes[j]) < 0
	})

	minimumDarkPoolSizeBig, err := ocean.darkNodeRegistrar.MinimumDarkPoolSize()
	if err != nil {
		return fmt.Errorf("cannot get minimum dark pool size: %v", err)
	}
	minimumDarkPoolSize, err := minimumDarkPoolSizeBig.ToUint()
	if err != nil {
		return fmt.Errorf("cannot parse minimum dark pool size: %v", err)
	}

	numberOfPools := len(nodeIDs) / int(minimumDarkPoolSize)

	pools := make(Pools, numberOfPools)
	for i := range pools {
		pools[i] = NewPool()
	}
	for i := range nodePositionHashes {
		id := identity.ID(nodePositionHashesToIDs[string(nodePositionHashes[i])])
		pools[i%numberOfPools].Append(NewNode(id))
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

	minInterval, err := ocean.darkNodeRegistrar.MinimumEpochInterval()
	if err != nil {
		log.Printf(fmt.Sprintf("cannot retrieve minimum epoch interval: %s", err.Error()))
		return
	}

	var currentBlockhash [32]byte
	if err := ocean.Update(); err != nil {
		log.Printf(fmt.Sprintf("cannot update dark ocean: %s", err.Error()))
		return
	}
	changes <- struct{}{}
	for {
		epoch, err := ocean.darkNodeRegistrar.CurrentEpoch()
		if err != nil {
			log.Printf(fmt.Sprintf("cannot update epoch: %s", err.Error()))
			return
		}
		if !bytes.Equal(currentBlockhash[:], epoch.Blockhash[:]) {
			currentBlockhash = epoch.Blockhash
			if err := ocean.Update(); err != nil {
				log.Printf(fmt.Sprintf("cannot update dark ocean: %s", err.Error()))
				return
			}
			changes <- struct{}{}
		}
		// TODO: Retrieve sleep time from epoch.Timestamp and minimumEpochInterval
		nextTime := epoch.Timestamp.Add(&minInterval)
		unix, err := nextTime.ToUint()
		if err != nil {
			// Either minInterval is really big, or unix epoch time has overflowed uint64s.
			log.Printf(fmt.Sprintf("epoch timestamp has overflowed: %s", err.Error()))
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

// GetPools returns dark pools in the dark ocean
func (ocean *Ocean) GetPools() Pools {
	return ocean.pools
}

// Equal checks if two dark oceans have the same view.
func (ocean *Ocean) Equal(other *Ocean) bool {
	if len(ocean.pools) != len(other.pools) {
		return false
	}

	for _, pool := range ocean.pools {
		has := false
		for _, otherPool := range other.pools {
			if pool.Equal(otherPool) {
				has = true
				break
			}
		}
		if !has {
			return false
		}
	}

	return true
}
