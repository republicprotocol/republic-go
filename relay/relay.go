package relay

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/swarm"
)

// ErrUnknownPod is returned when an unknown pod is mapped.
var ErrUnknownPod = errors.New("invalid number of pods")

// ErrInvalidNumberOfPods is returned when an insufficient number of pods are
// mapped.
var ErrInvalidNumberOfPods = errors.New("invalid number of pods")

// ErrInvalidNumberOfOrderFragments is returned when a pod is mapped to an
// insufficient number of order fragments, or too many order fragments.
var ErrInvalidNumberOfOrderFragments = errors.New("invalid number of order fragments")

// An OrderFragmentMapping maps pods to order fragments. The order fragments
// are expected to be encrypted.
type OrderFragmentMapping map[[32]byte][]order.Fragment

// Relayer interface can open and cancel orders on behalf of a trader.
type Relayer interface {

	// OpenOrder on the Ren Ledger and on the Darkpool. A signature from the
	// trader identifies them as the owner, the order ID is submitted to the
	// Ren Ledger along with the necessary fee, and the order fragment mapping
	// is used to send order fragments to pods in the Darkpool.
	OpenOrder(signature [65]byte, orderID order.ID, orderFragmentMapping OrderFragmentMapping) error

	// CancelOrder on the Ren Ledger. A signature from the trader is needed to
	// verify the cancelation.
	CancelOrder(signature [65]byte, orderID order.ID) error

	// SyncDarkpool to ensure an up-to-date state.
	SyncDarkpool() error
}

type Relay struct {
	darkpool  cal.Darkpool
	renLedger cal.RenLedger
	swarmer   swarm.Swarmer
	smpcer    smpc.Smpcer
	pods      map[[32]byte]cal.Pod
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
	if err := relay.verifyOrderFragments(orderFragmentMapping); err != nil {
		return err
	}
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

func (relay *Relay) SyncDarkpool() error {
	pods, err := relay.darkpool.Pods()
	if err != nil {
		return fmt.Errorf("cannot get pods from darkpool: %v", err)
	}
	relay.pods = map[[32]byte]cal.Pod{}
	for _, pod := range pods {
		relay.pods[pod.Hash] = pod
	}
	return nil
}

func (relay *Relay) verifyOrderFragments(orderFragmentMapping OrderFragmentMapping) error {
	if len(orderFragmentMapping) == 0 || len(orderFragmentMapping) > len(relay.pods) {
		return ErrInvalidNumberOfPods
	}
	for hash, orderFragments := range orderFragmentMapping {
		pod, ok := relay.pods[hash]
		if !ok {
			return ErrUnknownPod
		}
		if len(orderFragments) > len(pod.Darknodes) || len(orderFragments) < pod.Threshold() {
			return ErrInvalidNumberOfOrderFragments
		}
	}
	return nil
}

func (relay *Relay) openOrderFragments(orderFragmentMapping OrderFragmentMapping) error {
	errs := make([]error, 0, len(relay.pods))
	podDidReceiveFragments := false
	for hash, pod := range relay.pods {
		orderFragments := orderFragmentMapping[hash]
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
	if len(orderFragments) < pod.Threshold() || len(orderFragments) > len(pod.Darknodes) {
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
	errNumMax := len(orderFragments) - pod.Threshold()
	if len(pod.Darknodes) > 0 && errNum > errNumMax {
		return fmt.Errorf("cannot send order fragments to %v nodes (out of %v nodes) in pod %v: %v", errNum, len(pod.Darknodes), base64.StdEncoding.EncodeToString(pod.Hash[:]), err)
	}
	return nil
}
