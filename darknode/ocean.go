package darknode

import (
	"bytes"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/ethereum/contracts"
)

// An Ocean of Darknodes non-deterministically shuffled into Pools.
type Ocean struct {
	pools            Pools
	darknodeRegistry contracts.DarkNodeRegistry
}

// NewOcean returns a new Ocean that uses the DarknodeRegistry to watch for
// updates to the Ocean, and calculates the configuration of the Darknodes and
// Pools.
func NewOcean(darknodeRegistry contracts.DarkNodeRegistry) Ocean {
	ocean := Ocean{
		pools:            Pools{},
		darknodeRegistry: darknodeRegistry,
	}
	return ocean
}

// Watch for changes to the Ocean. It will stop watching for changes when the
// done channel is closed. Returns a signaling channel that is written to
// whenever a change is detected.
func (ocean *Ocean) Watch(done <-chan struct{}) (<-chan struct{}, <-chan error) {
	changes := make(chan struct{})
	errs := make(chan error, 1)

	go func() {
		defer close(changes)
		defer close(errs)

		minimumEpochInterval, err := ocean.darknodeRegistry.MinimumEpochInterval()
		if err != nil {
			errs <- fmt.Errorf("cannot get minimum epoch interval: %v", err)
			return
		}

		currentEpoch, err := ocean.darknodeRegistry.CurrentEpoch()
		if err != nil {
			errs <- fmt.Errorf("cannot get current epoch: %v", err)
			return
		}

		for {
			// Signal that the epoch has changed
			select {
			case <-done:
				return
			case changes <- struct{}{}:
			}

			// Sleep until the next epoch
			nextEpochTime := currentEpoch.Timestamp.Add(&minimumEpochInterval)
			nextEpochTimeUnix, err := nextEpochTime.ToUint()
			if err != nil {
				errs <- fmt.Errorf("cannot convert epoch timestamp to unix timestamp: %v", err)
				return
			}
			delay := time.Duration(int64(nextEpochTimeUnix)-time.Now().Unix()) * time.Second
			time.Sleep(delay)

			// Spin-lock until the new epoch is detected or until the done
			// channel is closed
			for {

				select {
				case <-done:
					return
				default:
				}

				nextEpoch, err := ocean.darknodeRegistry.CurrentEpoch()
				if err != nil {
					errs <- fmt.Errorf("cannot get next epoch: %v", err)
					return
				}
				if !bytes.Equal(currentEpoch.Blockhash[:], nextEpoch.Blockhash[:]) {
					currentEpoch = nextEpoch
					break
				}
			}
		}
	}()

	return changes, errs
}
