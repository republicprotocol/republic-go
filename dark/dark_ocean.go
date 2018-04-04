package dark

import (
	"sort"
	"bytes"
	"fmt"
	"time"

	"github.com/republicprotocol/go-do"

	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/ethereum/go-ethereum/crypto"
)

// Ocean of Pools.
type Ocean struct {
	do.GuardedObject

	logger            *logger.Logger
	poolSize 		  int
	pools             Pools
	darkNodeRegistrar dnr.DarkNodeRegistrar
}

// NewOcean uses a DarkNodeRegistrar to read all registered nodes and sort them
// into Pools.
func NewOcean(logger *logger.Logger, poolSize int, darkNodeRegistrar dnr.DarkNodeRegistrar) (*Ocean, error) {
	ocean := &Ocean{
		GuardedObject:     do.NewGuardedObject(),
		logger:            logger,
		poolSize: 		   poolSize,
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

func (ocean *Ocean) Update() error {
	ocean.Enter(nil)
	defer ocean.Exit()
	return ocean.update()
}

func (ocean *Ocean) update() error {
	epoch, err := ocean.darkNodeRegistrar.CurrentEpoch();
	if err != nil {
		return err
	}

	nodeIDs, err := ocean.darkNodeRegistrar.GetAllNodes()
	if err != nil {
		return err
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

	numberOfPools := len(nodeIDs) / ocean.poolSize

	pools := make(Pools, numberOfPools)
	for i := range pools {
		pools[i] = NewPool()
	}
	for i := range nodePositionHashes {
		id := identity.ID(nodePositionHashesToIDs[string(nodePositionHashes[i])])
		pools[i % numberOfPools].Append(NewNode(id))
	}
	ocean.pools = pools
	return nil
}

// Watch for changes to the Ocean. This function is a blocking function that
// sleeps and wakes once per period to check for a change in epoch. It accepts
// a channel that is pinged whenever the Ocean changes.
func (ocean *Ocean) Watch(period time.Duration, changes chan struct{}) {
	// Recover from writing to a closed channel
	defer func() { recover() }()

	var currentBlockhash [32]byte
	for {
		epoch, err := ocean.darkNodeRegistrar.CurrentEpoch()
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
		time.Sleep(period)
	}
}

func (ocean *Ocean) GetPools ()(Pools){
	return ocean.pools
}
