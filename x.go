package x

import (
	"reflect"
	"runtime"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

// AssignXHash assigns an X Hash to all Miners. For each Miner, the epoch Hash
// and the commitment Hash are combined and hashed to produce the X Hash. All
// hashing is done using the Keccak256 hashing function. The list of miners
// will be sorted by their X Hashes.
func AssignXHash(epoch Hash, miners []Miner) {

	// Calculate workload size per CPU.
	numCPUs := runtime.NumCPU()
	numHashesPerCPU := (len(miners) / numCPUs) + 1

	// Calculate list of output hashes in parallel.
	var wg sync.WaitGroup
	wg.Add(numCPUs)
	for i := 0; i < len(miners); i += numHashesPerCPU {
		go func(i int) {
			defer wg.Done()
			for j := i; j < i+numHashesPerCPU && j < len(miners); j++ {
				miners[j].X = crypto.Keccak256([]byte(epoch), []byte(miners[j].Commitment))
			}
		}(i)
	}
	wg.Wait()

	// Sort the list of output hashes.
	sort.Slice(miners, func(i, j int) bool {
		return miners[i].X.LessThan(miners[j].X)
	})
}

// AssignClass will assign a class to each miner. This function can only be
// called after the assignXHash function. If any miner in the list does not
// have an X Hash then this function will do nothing. If the miners are not
// sorted then this function will do nothing.
func AssignClass(numberOfMNetworks int, miners []Miner) {
	if !RequireXHashes(miners) {
		return
	}

	// Calculate workload size per CPU.
	numCPUs := runtime.NumCPU()
	numHashesPerCPU := (len(miners) / numCPUs) + 1

	// Calculate list of output hashes in parallel.
	var wg sync.WaitGroup
	wg.Add(numCPUs)
	for i := 0; i < len(miners); i += numHashesPerCPU {
		go func(i int) {
			defer wg.Done()
			for j := i; j < i+numHashesPerCPU && j < len(miners); j++ {
				miners[j].Class = j/numberOfMNetworks + 1
			}
		}(i)
	}
	wg.Wait()
}

// AssignMNetwork will assigned an M Network to each miner. This function must
// not be called before the assignXHash function. If any miner in the list does
// not have an X Hash, this function will do nothing. If the miners are not
// sorted then this function will do nothing.
func AssignMNetwork(numberOfMNetworks int, miners []Miner) {
	if !RequireXHashes(miners) {
		return
	}
	forAll(miners, func(k int) {
		miners[k].MNetwork = k % numberOfMNetworks
	})
}

// RequireXHashes checks that every Miner in the list has a valid X Hash. It
// also guarantees that the list is sorted by these X Hashes. Returns true if
// all Miners have a valid X Hash, otherwise false.
func RequireXHashes(miners []Miner) bool {
	// Require that all miners have an X Hash.
	for _, miner := range miners {
		if miner.X == nil {
			return false
		}
	}
	// Require that the miners are sorted.
	isSorted := sort.SliceIsSorted(miners, func(i, j int) bool {
		return miners[i].X.LessThan(miners[j].X)
	})
	if !isSorted {
		sort.Slice(miners, func(i, j int) bool {
			return miners[i].X.LessThan(miners[j].X)
		})
	}
	return true
}

func forAll(data interface{}, f func(i int)) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		// Calculate workload size per CPU.
		length := reflect.ValueOf(data).Len()
		numCPUs := runtime.NumCPU()
		numIterationsPerCPU := (length / numCPUs) + 1

		// Calculate list of output hashes in parallel.
		var wg sync.WaitGroup
		wg.Add(numCPUs)
		for offset := 0; offset < length; offset += numIterationsPerCPU {
			go func(offset int) {
				defer wg.Done()
				for i := offset; i < offset+numIterationsPerCPU && i < length; i++ {
					f(i)
				}
			}(offset)
		}
		wg.Wait()
	}
}
