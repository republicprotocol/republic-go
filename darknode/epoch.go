package darknode

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/darkocean"
	"github.com/republicprotocol/republic-go/delta"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/smpc"
)

// RunEpochProcess until the done channel is closed. Epochs define the Pools in
// which Darknodes cooperate to match orders by receiving order.Fragments from
// traders, performing secure multi-party computations. The EpochProcess uses a
// DarkOcean to determine the epoch, and a Router to receive messages.
func (node *Darknode) RunEpochProcess(done <-chan struct{}, ocean darkocean.DarkOcean) (<-chan delta.Delta, <-chan error) {
	deltas := make(chan delta.Delta)
	errs := make(chan error, 1)

	go func() {
		defer close(deltas)

		pool, err := ocean.Pool(node.ID())
		if err != nil {
			errs <- err
			return
		}

		// Open connections with all Darknodes in the Pool
		mu := new(sync.Mutex)
		ctxs, cancels := map[identity.Address]context.Context{}, map[identity.Address]context.CancelFunc{}
		senders := map[identity.Address]chan<- *smpcer.ComputeMessage{}
		defer func() {
			for key := range senders {
				close(senders[key])
			}
		}()
		receivers := map[identity.Address]<-chan *smpcer.ComputeMessage{}
		errors := map[identity.Address]<-chan error{}

		addresses := pool.Addresses()
		dispatch.CoForAll(addresses, func(i int) {

			addr := addresses[i]
			if bytes.Compare(addr.ID()[:], node.ID()[:]) == 0 {
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			mu.Lock()
			ctxs[addr], cancels[addr] = ctx, cancel
			mu.Unlock()

			sender := make(chan *smpcer.ComputeMessage)
			multiAddr, err := node.rpc.SwarmerClient().Query(ctx, addr, 3)
			if err != nil {
				log.Printf("cannot query smpc peer %v: %v", addr, err)
				return
			}
			receiver, errs := node.rpc.SmpcerClient().Compute(ctx, multiAddr, sender)

			mu.Lock()
			senders[addr] = sender
			receivers[addr] = receiver
			errors[addr] = errs
			mu.Unlock()
		})

		n := int64(pool.Size())
		k := (n + 1) * 2 / 3
		smpc := smpc.NewSmpc(node.ID(), n, k)

		orderFragments := node.orderFragments                 // TODO: Splitter for multiple epochs
		orderFragmentsCanceled := node.orderFragmentsCanceled // TODO: Splitter for multiple epochs
		deltaFragments := make(chan delta.Fragment)
		defer close(deltaFragments)

		deltaFragmentsComputed, deltasComputed := smpc.ComputeOrderMatches(done, orderFragments, deltaFragments)

		dispatch.CoBegin(func() {

			// Receive cancelations
			for {
				select {
				case <-done:
					return
				case orderID, ok := <-orderFragmentsCanceled:
					if !ok {
						return
					}
					smpc.SharedOrderTable().RemoveBuyOrder(orderID)
					smpc.SharedOrderTable().RemoveSellOrder(orderID)
				}
			}
		}, func() {

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
						if computation.GetDeltaFragment() != nil {
							deltaFragment, err := smpcer.UnmarshalDeltaFragment(computation.GetDeltaFragment())
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
					computation := &smpcer.ComputeMessage{
						Value: &smpcer.ComputeMessage_DeltaFragment{
							DeltaFragment: smpcer.MarshalDeltaFragment(&deltaFragment),
						},
					}
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
			for {
				select {
				case <-done:
					return
				case deltaComputed, ok := <-deltasComputed:
					if !ok {
						return
					}
					if deltaComputed.IsMatch(smpc.Prime()) {
						smpc.SharedOrderTable().RemoveBuyOrder(deltaComputed.BuyOrderID)
						smpc.SharedOrderTable().RemoveSellOrder(deltaComputed.SellOrderID)
					}
					select {
					case <-done:
						return
					case deltas <- deltaComputed:
					}
				}
			}
		})
	}()

	return deltas, errs
}

// RunEpochWatcher until the done channel is closed. An EpochWatcher will watch
// for changes to the DarknodeRegistry epoch. Returns a read-only channel that
// can be used to read epochs as they change.
func RunEpochWatcher(done <-chan struct{}, darknodeRegistry dnr.DarknodeRegistry) (<-chan dnr.Epoch, <-chan error) {
	changes := make(chan dnr.Epoch)
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
