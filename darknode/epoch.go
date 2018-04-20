package darknode

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc"
	"github.com/republicprotocol/republic-go/smpc"
)

// RunEpochProcess until the done channel is closed. Epochs define the Pools in
// which Darknodes cooperate to match orders by receiving order.Fragments from
// traders, performing secure multi-party computations. The EpochProcess uses a
// DarkOcean to determine the epoch, and a Router to receive messages.
func RunEpochProcess(done <-chan struct{}, id identity.ID, darkOcean DarkOcean, router *Router) (<-chan smpc.Delta, <-chan error) {
	deltas := make(chan smpc.Delta)
	errs := make(chan error, 1)

	go func() {
		defer close(deltas)

		pool, err := darkOcean.Pool(id)
		if err != nil {
			errs <- err
			return
		}

		// Open connections with all Darknodes in the Pool
		mu := new(sync.Mutex)
		senders := map[identity.Address]chan<- *rpc.Computation{}
		defer func() {
			for key := range senders {
				close(senders[key])
			}
		}()
		receivers := map[identity.Address]<-chan *rpc.Computation{}
		errors := map[identity.Address]<-chan error{}

		addresses := pool.Addresses()
		dispatch.CoForAll(addresses, func(i int) {
			if bytes.Compare(addresses[i].ID()[:], id[:]) == 0 {
				return
			}

			sender := make(chan *rpc.Computation)
			receiver, errs := router.Compute(done, addresses[i], sender)

			mu.Lock()
			defer mu.Unlock()
			senders[addresses[i]] = sender
			receivers[addresses[i]] = receiver
			errors[addresses[i]] = errs
		})

		n := int64(pool.Size())
		k := (n + 1) * 2 / 3
		smpcer := smpc.NewComputer(id, n, k)

		orderFragments, orderFragmentErrs := router.OrderFragments(done)
		go dispatch.Pipe(done, orderFragmentErrs, errs)

		deltaFragments := make(chan smpc.DeltaFragment)
		defer close(deltaFragments)

		deltaFragmentsComputed, deltasComputed := smpcer.ComputeOrderMatches(done, orderFragments, deltaFragments)

		dispatch.CoBegin(func() {

			// Receive smpc.DeltaFragments from other Darknodes in the Pool
			dispatch.CoForAll(receivers, func(receiver identity.Address) {
				for {
					select {
					case <-done:
						return
					case computation, ok := <-receivers[receiver]:
						if !ok {
							return
						}
						if computation.DeltaFragment != nil {
							deltaFragment, err := rpc.UnmarshalDeltaFragment(computation.DeltaFragment)
							if err != nil {
								errs <- err
								continue
							}
							select {
							case <-done:
								return
							case deltaFragments <- deltaFragment:
							}
						}
					}
				}
			})
		}, func() {

			// Broadcast computed smpc.DeltaFragments to other Darknodes in the
			// Pool
			for {
				select {
				case <-done:
					return
				case deltaFragment, ok := <-deltaFragmentsComputed:
					if !ok {
						return
					}
					computation := &rpc.Computation{DeltaFragment: rpc.MarshalDeltaFragment(&deltaFragment)}
					for _, sender := range senders {
						ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
						select {
						case <-done:
							cancel()
							return
						case <-ctx.Done():
							cancel()
							continue
						case sender <- computation:
						}
						cancel()
					}
				}
			}
		}, func() {

			// Output computed smpc.Deltas
			dispatch.Pipe(done, deltasComputed, deltas)
		})
	}()

	return deltas, errs
}

// RunEpochWatcher until the done channel is closed. An EpochWatcher will watch
// for changes to the DarknodeRegistry epoch. Returns a read-only channel that
// can be used to read epochs as they change.
func RunEpochWatcher(done <-chan struct{}, darknodeRegistry contracts.DarkNodeRegistry) (<-chan contracts.Epoch, <-chan error) {
	changes := make(chan contracts.Epoch)
	errs := make(chan error, 1)

	go func() {
		defer close(changes)
		defer close(errs)

		minimumEpochInterval, err := darknodeRegistry.MinimumEpochInterval()
		if err != nil {
			errs <- fmt.Errorf("cannot get minimum epoch interval: %v", err)
			return
		}

		currentEpoch, err := darknodeRegistry.CurrentEpoch()
		if err != nil {
			errs <- fmt.Errorf("cannot get current epoch: %v", err)
			return
		}

		for {
			// Signal that the epoch has changed
			select {
			case <-done:
				return
			case changes <- currentEpoch:
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

				nextEpoch, err := darknodeRegistry.CurrentEpoch()
				if err != nil {
					errs <- fmt.Errorf("cannot get next epoch: %v", err)
					return
				}
				if !bytes.Equal(currentEpoch.Blockhash[:], nextEpoch.Blockhash[:]) {
					currentEpoch = nextEpoch
					break
				}

				time.Sleep(time.Minute)
			}
		}
	}()

	return changes, errs
}
