package x

import (
	"errors"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x/rpc"
	"golang.org/x/net/context"
)

// Peers is used to return the rpc.MultiAddresses to which a Node is connected.
// The rpc.MultiAddresses returned are not guaranteed to provide healthy
// connections and should be pinged.
func (node *Node) Peers(ctx context.Context, sender *rpc.Nothing) (*rpc.MultiAddresses, error) {
	if err := ctx.Err(); err != nil {
		return &rpc.MultiAddresses{}, err
	}

	// Spawn a goroutine to evaluate the return value.
	wait := make(chan *rpc.MultiAddresses)
	go func() {
		defer close(wait)
		wait <- node.handlePeers()
	}()

	select {
	case ret := <-wait:
		return ret, nil

	// Select the timeout from the context.
	case <-ctx.Done():
		return &rpc.MultiAddresses{}, ctx.Err()
	}
}

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

// ForwardOrderFragment forwards the order fragment to the miners so that they
// can transmit the order fragment to the target. Return nil if forward
// successfully, or an error.
func (node *Node) ForwardOrderFragment(orderFragment *rpc.OrderFragment) error {

	target := identity.Address(orderFragment.To)
	open := node.DHT.MultiAddresses()
	if len(open) == 0 {
		return errors.New("empty dht")
	}

	// Sort the nodes we already know
	sort.SliceStable(open, func(i, j int) bool {
		left, _ := open[i].Address()
		right, _ := open[j].Address()
		closer, _ := identity.Closer(left, right, target)
		return closer
	})

	// If we know the target,send the order fragment to the target directly
	closestNode, err := open[0].Address()
	if err != nil {
		return err
	}
	if string(closestNode) == string(target) {
		_, err := node.RPCSendOrderFragment(open[0], orderFragment)
		return err
	}

	// Otherwise forward the fragment to the closest α nodes simultaneously
	for len(open) > 0 {
		asyncRoutines := len(open)
		if len(open) > 3 {
			asyncRoutines = α
		}

		var wg sync.WaitGroup
		wg.Add(asyncRoutines)
		targetFound := int32(0)

		// Forward order fragment
		for i := 0; i < asyncRoutines; i++ {
			multi := open[0]
			open = open[1:]
			go func() {
				defer wg.Done()
				response, _ := node.RPCSendOrderFragment(multi, orderFragment)
				if response != nil {
					atomic.StoreInt32(&targetFound, 1)
				}
			}()
			if len(open) == 0 {
				break
			}
		}
		wg.Wait()
		if targetFound != 0 {
			return nil
		}
	}
	return errors.New("we can't find the target")
}
