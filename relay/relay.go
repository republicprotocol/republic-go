package relay

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/swarm"
)

var ErrInvalidNumberOfOrderFragments = errors.New("invalid number of order fragments")

type OrderFragmentMapping map[[32]byte][]order.Fragment

type Relayer interface {
	OpenOrder(signature [65]byte, orderID order.ID, orderFragmentMapping OrderFragmentMapping) error
	CancelOrder(signature [65]byte, orderID order.ID) error
}

type Relay struct {
	darkpool  cal.Darkpool
	renLedger cal.RenLedger
	swarmer   swarm.Swarmer
	smpcer    smpc.Smpcer
}

func NewRelay(darkpool cal.Darkpool, renLedger cal.RenLedger, swarmer swarm.Swarmer, smpcer smpc.Smpcer) Relay {
	return Relay{
		darkpool:  darkpool,
		renLedger: renLedger,
		swarmer:   swarmer,
		smpcer:    smpcer,
	}
}

func (relay *Relay) OpenOrder(signature [65]byte, orderID order.ID, orderFragmentMapping OrderFragmentMapping) error {
	// TODO: Verify that the signature is valid before sending it to the
	// RenLedger. This is not strictly necessary but it can save the Relay some
	// gas.
	// FIXME: Re-enable this interaction once signatures have been figured out.
	if err := relay.renLedger.OpenOrder(signature, orderID); err != nil {
		return err
	}
	return relay.openOrderFragments(orderFragmentMapping)
}

func (relay *Relay) CancelOrder(signature [65]byte, orderID order.ID) error {
	// TODO: Verify that the signature is valid before sending it to the
	// RenLedger. This is not strictly necessary but it can save the Relay some
	// gas.
	if err := relay.renLedger.CancelOrder(signature, orderID); err != nil {
		return err
	}
	return nil
}

func (relay *Relay) openOrderFragments(orderFragmentMapping OrderFragmentMapping) error {
	pods, err := relay.darkpool.Pods()
	if err != nil {
		return fmt.Errorf("cannot get pods from darkpool: %v", err)
	}

	errs := make([]error, 0, len(pods))
	podDidReceiveFragments := false
	for _, pod := range pods {
		log.Printf("pod %v", base64.StdEncoding.EncodeToString(pod.Hash[:]))
		orderFragments := orderFragmentMapping[pod.Hash]
		if orderFragments != nil && len(orderFragments) > 0 {
			if err := relay.sendOrderFragmentsToPod(pod, orderFragments); err != nil {
				errs = append(errs, err)
				continue
			}
			podDidReceiveFragments = true
		}
	}
	if !podDidReceiveFragments {
		if len(errs) == 0 {
			return fmt.Errorf("cannot open order fragments: no pod received an order fragments")
		}
		return fmt.Errorf("cannot open order fragments: no pod received an order fragments: %v", errs[0])
	}
	return nil
}

func (relay *Relay) sendOrderFragmentsToPod(pod cal.Pod, orderFragments []order.Fragment) error {
	threshold := 2 * (len(pod.Darknodes) + 1) / 3
	if len(orderFragments) < threshold || len(orderFragments) > len(pod.Darknodes) {
		return ErrInvalidNumberOfOrderFragments
	}

	// Map order fragments to their respective Darknodes
	orderFragmentIndexMapping := map[int64]order.Fragment{}
	for _, orderFragment := range orderFragments {
		orderFragmentIndexMapping[orderFragment.PriceShare.Key] = orderFragment
	}

	errs := make(chan error, len(pod.Darknodes))
	go func() {
		defer close(errs)

		dispatch.CoForAll(pod.Darknodes, func(i int) {
			orderFragment, ok := orderFragmentIndexMapping[int64(i)]
			if !ok {
				return
			}
			darknode := pod.Darknodes[i]

			// Send the order fragment to the Darknode
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			darknodeMultiAddr, err := relay.swarmer.Query(ctx, darknode, -1)
			if err != nil {
				errs <- fmt.Errorf("cannot send query to %v: %v", darknode, err)
				return
			}
			if err := relay.smpcer.OpenOrder(ctx, darknodeMultiAddr, orderFragment); err != nil {
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

	// Check if at least 2/3 of the nodes in the specified pod have recieved
	// the order fragments.
	errNumMax := len(orderFragments) - threshold
	if len(pod.Darknodes) > 0 && errNum > errNumMax {
		return fmt.Errorf("cannot send order fragments to %v nodes (out of %v nodes) in pod %v: %v", errNum, len(pod.Darknodes), base64.StdEncoding.EncodeToString(pod.Hash[:]), err)
	}
	return nil
}
