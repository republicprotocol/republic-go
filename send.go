package x

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x/rpc"
	"golang.org/x/net/context"
)

// SendOrderFragment is used to forward an rpc.OrderFragment through the X
// Network to its destination Node. This forwarding is done using a distributed
// Dijkstra search, using the XOR distance between identity.Addresses as the
// distance heuristic.
func (node *Node) SendOrderFragment(ctx context.Context, orderFragment *rpc.OrderFragment) (*rpc.MultiAddress, error) {
	// Check for errors in the context.
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, err
	}

	// Spawn a goroutine to evaluate the return value.
	type retType struct {
		*rpc.MultiAddress
		error
	}
	wait := make(chan retType)
	go func() {
		defer close(wait)
		multi, err := node.handleSendOrderFragment(orderFragment)
		wait <- retType{multi, err}
	}()

	select {
	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, ctx.Err()

	// Select the value passed by the goroutine.
	case ret := <-wait:
		return ret.MultiAddress, ret.error
	}
}

func (node *Node) handleSendOrderFragment(orderFragment *rpc.OrderFragment) (*rpc.MultiAddress, error) {

	target := identity.Address(orderFragment.To)
	if target == node.DHT.Address {
		node.Delegate.OnOrderFragmentReceived()
		return &rpc.MultiAddress{Multi: node.MultiAddress.String()}, nil
	}

	targetMultiMu := new(sync.Mutex)
	var targetMulti *identity.MultiAddress

	// Initialize the closed list
	closed := make(map[string]bool)
	self, err := node.MultiAddress.Address()
	if err != nil {
		return nil, err

	}
	closed[self.String()] = true
	multiFrom, err := identity.NewMultiAddressFromString(orderFragment.From)
	if err != nil {
		return nil, err
	}
	from, err := multiFrom.Address()
	if err != nil {
		return nil, err
	}
	closed[from.String()] = true

	// Initialize the open list
	openMu := new(sync.Mutex)
	open := identity.MultiAddresses{}
	for _, peer := range node.DHT.MultiAddresses() {
		address, err := peer.Address()
		if err != nil {
			return nil, err
		}
		if _, ok := closed[address.String()]; !ok {
			open = append(open, peer)
		}
	}
	if len(open) == 0 {
		return nil, errors.New("empty dht")
	}
	sort.Slice(open, func(i, j int) bool {
		left, _ := open[i].Address()
		right, _ := open[j].Address()
		closer, _ := identity.Closer(left, right, target)
		return closer
	})

	// If we know the target,send the order fragment to the target directly
	closestNode, err := open[0].Address()
	if err != nil {
		return nil, err
	}
	if closestNode == target {
		_, err := node.RPCSendOrderFragment(open[0], orderFragment)
		return &rpc.MultiAddress{Multi: open[0].String()}, err
	}

	// TODO: We are only using one dht.Bucket to search the network. If this
	//       dht.Bucket is not sufficient, we should also search the
	//       neighborhood of the dht.Bucket. The neighborhood should expand by
	//       a distance of one, until the entire DHT has been searched, or the
	//       targetMulti has been found.
	for len(open) > 0 {
		asyncRoutines := len(open)
		if len(open) > 3 {
			asyncRoutines = α
		}
		var wg sync.WaitGroup
		wg.Add(asyncRoutines)

		// Take the first α multi-addresses from the open list and expand them
		// concurrently. This moves them from the open list to the closed list,
		// preventing the same multi-address from being expanded more than
		// once.
		for i := 0; i < asyncRoutines; i++ {
			if len(open) == 0 {
				break
			}
			multi := open[0]
			open = open[1:]
			address, err := multi.Address()
			if err != nil {
				return nil, err
			}
			closed[address.String()] = true

			go func() {
				defer wg.Done()

				// Get all peers of this multi-address. This is the expansion
				// step of the search.
				peers, err := node.RPCPeers(multi)
				if err != nil {
					return
				}

				// Traverse all peers and collect them into the openNext, a
				// list of peers that we want to add the open list.
				openNext := make([]identity.MultiAddress, 0, len(peers))
				for _, peer := range peers {
					address, err := peer.Address()
					if err != nil {
						return
					}
					if target == address {
						// If we have found the target, set the targetMulti and
						// exit the loop. There is no point acquiring more
						// peers for the open list.
						targetMultiMu.Lock()
						if targetMulti == nil {
							targetMulti = &peer
						}
						targetMultiMu.Unlock()
						break
					}
					// Otherwise, store this peer's multi-address in the
					// openNext list. It will be added to the open list if it
					// has not already been closed.
					openNext = append(openNext, peer)
				}

				targetMultiMu.Lock()
				if targetMulti == nil {
					// We have not found the targetMulti yet.
					targetMultiMu.Unlock()
					openMu.Lock()
					// Add new peers to the open list if they have not been
					// closed.
					for _, next := range openNext {
						address, err := next.Address()
						if err != nil {
							return
						}
						if _, ok := closed[address.String()]; !ok {
							open = append(open, next)
						}
					}
					openMu.Unlock()
				} else {
					targetMultiMu.Unlock()
				}
			}()
		}
		wg.Wait()
		if targetMulti != nil {
			// If we have found the targetMulti then we can exit the search
			// loop.
			break
		}

		// Otherwise, sort the open list by distance to the target.
		sort.Slice(open, func(i, j int) bool {
			left, _ := open[i].Address()
			right, _ := open[j].Address()
			closer, _ := identity.Closer(left, right, target)
			return closer
		})
	}

	if targetMulti == nil {
		return nil, fmt.Errorf("cannot find target")
	}
	response, err := node.RPCSendOrderFragment(*targetMulti, orderFragment)
	if err != nil {
		return nil, err
	}

	return &rpc.MultiAddress{Multi: response.String()}, nil
}
