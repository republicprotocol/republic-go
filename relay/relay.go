package relay

import (
	"context"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc"
)

type OrderFragmentMapping map[[32]byte][]order.Fragment

type Relayer interface {
	OpenOrder(signature [65]byte, orderID order.ID, orderType order.Type, orderParity order.Parity, orderExpiry int64, orderFragmentMapping OrderFragmentMapping) error
	CancelOrder(signature [65]byte, orderID order.ID) error
}

type Relay struct {
	darkpool    cal.Darkpool
	renLedger   cal.RenLedger
	swarmClient rpc.SwarmClient
	smpcClient  rpc.SmpcClient
}

func NewRelay(darkpool cal.Darkpool, renLedger cal.RenLedger, swarmClient rpc.SwarmClient, smpcClient rpc.SmpcClient) Relay {
	return Relay{
		darkpool:    darkpool,
		renLedger:   renLedger,
		swarmClient: swarmClient,
		smpcClient:  smpcClient,
	}
}

func (relay *Relay) OpenOrder(signature [65]byte, orderID order.ID, orderType order.Type, orderParity order.Parity, orderExpiry int64, orderFragmentMapping OrderFragmentMapping) error {
	// TODO: Verify that the signature is valid before sending it to the
	// RenLedger. This is not strictly necessary but it can save the Relay some
	// gas.
	if err := relay.renLedger.OpenOrder(orderID, orderType, orderParity, orderExpiry, signature); err != nil {
		return err
	}
	if err := relay.renLedger.WaitForOpenOrder(orderID); err != nil {
		return err
	}
	return relay.openOrderFragments(orderFragmentMapping)
}

func (relay *Relay) CancelOrder(signature [65]byte, orderID order.ID) error {
	// TODO: Verify that the signature is valid before sending it to the
	// RenLedger. This is not strictly necessary but it can save the Relay some
	// gas.
	if err := relay.renLedger.CancelOrder(orderID, signature); err != nil {
		return err
	}
	if err := relay.renLedger.WaitForCancelOrder(orderID); err != nil {
		return err
	}
	return nil
}

func (relay *Relay) openOrderFragments(orderFragmentMapping OrderFragmentMapping) error {
	pools, err := relay.darkpool.Pools()
	if err != nil {
		return fmt.Errorf("cannot get pools from darkpool: %v", err)
	}

	errs := make([]error, 0, len(pools))
	poolDidReceiveFragments := false
	for _, pool := range pools {
		orderFragments := orderFragmentMapping[pool.Hash]
		if orderFragments != nil && len(orderFragments) > 0 {
			if err := relay.sendOrderFragmentsToPool(pool, orderFragments); err == nil {
				errs = append(errs, err)
				poolDidReceiveFragments = true
			}
		}
	}
	if !poolDidReceiveFragments {
		if len(errs) == 0 {
			return fmt.Errorf("cannot open order fragments: no pool received an order fragments")
		}
		return fmt.Errorf("cannot open order fragments: no pool received an order fragments: %v", errs[0])
	}
	return nil
}

func (relay *Relay) sendOrderFragmentsToPool(pool cal.Pool, orderFragments []order.Fragment) error {

	// Map order fragments to their respective Darknodes
	orderFragmentIndexMapping := map[int64]order.Fragment{}
	for _, orderFragment := range orderFragments {
		orderFragmentIndexMapping[orderFragment.PriceShare.Key] = orderFragment
	}

	errs := make(chan error, len(pool.Darknodes))
	go func() {
		defer close(errs)

		dispatch.CoForAll(pool.Darknodes, func(i int) {
			orderFragment, ok := orderFragmentIndexMapping[int64(i)]
			if !ok {
				return
			}
			darknode := pool.Darknodes[i]

			// Send the order fragment to the Darknode
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			darknodeMultiAddr, err := relay.swarmClient.Query(ctx, darknode, -1)
			if err != nil {
				errs <- fmt.Errorf("cannot send query to %v: %v", darknode, err)
				return
			}
			if err := relay.smpcClient.OpenOrder(ctx, darknodeMultiAddr, orderFragment); err != nil {
				errs <- fmt.Errorf("cannot send order fragment to %v: %v", darknode, err)
				return
			}
		})
	}()

	// Capture all errors and keep the first error that occurred.
	var errNum int
	var err error
	for errLocal := range errs {
		if errLocal != nil {
			errNum++
			if err == nil {
				err = errLocal
			}
		}
	}

	// Check if at least 2/3 of the nodes in the specified pool have recieved
	// the order fragments.
	errNumMax := len(orderFragments) - (2 * (len(pool.Darknodes) + 1) / 3)
	if len(pool.Darknodes) > 0 && errNum > errNumMax {
		return fmt.Errorf("cannot send order fragments to %v nodes (out of %v nodes) in pool %v", errNum, len(pool.Darknodes), pool.Hash)
	}

	return nil
}

// IsOrderFragmentsComplete returns true if there are enough order.Fragments
// for a cal.Pool to have a chance at matching the order.Order. Returns false
// otherwise.
func (relay *Relay) IsOrderFragmentsComplete(fragmentCount, poolSize int) bool {
	return fragmentCount >= 2/3*poolSize && fragmentCount <= poolSize
}
