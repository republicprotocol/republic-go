package x

import (
	"runtime"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

// Sort the list of miners. The epoch hash, and each miner hash, is combined
// and hashed to produce an X hash for each miner. The miner are sorted by
// comparing each their X hashes against each other.
func Sort(epochHash Hash, minerHashes []Hash) []Hash {

	// Calculate workload size per CPU.
	xHashes := make([]Hash, len(minerHashes))
	numHashesPerCPU := (len(minerHashes) / runtime.NumCPU()) + 1

	// Calculate list of output hashes in parallel.
	var wg sync.WaitGroup
	wg.Add(len(minerHashes))
	for i := 0; i < len(minerHashes); i += numHashesPerCPU {
		go func(i int) {
			defer wg.Done()
			for j := i; j < i+numHashesPerCPU && j < len(minerHashes); j++ {
				hash := crypto.Keccak256([]byte(epochHash), []byte(minerHashes[j]))
				xHashes[j] = hash
			}
		}(i)
	}
	wg.Wait()

	// Sort the list of output hashes.
	sort.Slice(minerHashes, func(i, j int) bool {
		return xHashes[i].LessThan(xHashes[j])
	})

	return minerHashes
}
