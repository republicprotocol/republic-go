package x

import (
	"runtime"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

// Sort the list of Miners. The epoch Hash, and each Miner commitment Hash, is
// combined and hashed to produce an X Hash for each Miner. The miner are sorted by
// comparing each their X hashes against each other. All hashing is done using
// the Keccak256 hashing function.
func Sort(epoch Hash, miners []Miner) []Miner {

	// Calculate workload size per CPU.
	xs := make([]Miner, len(miners))
	numCPUs := runtime.NumCPU()
	numHashesPerCPU := (len(miners) / numCPUs) + 1

	// Calculate list of output hashes in parallel.
	var wg sync.WaitGroup
	wg.Add(numCPUs)
	for i := 0; i < len(miners); i += numHashesPerCPU {
		go func(i int) {
			defer wg.Done()
			for j := i; j < i+numHashesPerCPU && j < len(miners); j++ {
				xs[j] = Miner{
					ID:         miners[j].ID,
					Commitment: miners[j].Commitment,
					X:          crypto.Keccak256([]byte(epoch), []byte(miners[j].Commitment)),
				}
			}
		}(i)
	}
	wg.Wait()

	// Sort the list of output hashes.
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].X.LessThan(xs[j].X)
	})

	return xs
}
