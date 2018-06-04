package ingress

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/swarm"
)

// ErrUnknownPod is returned when an unknown pod is mapped.
var ErrUnknownPod = errors.New("unknown pod id")

// ErrInvalidNumberOfPods is returned when an insufficient number of pods are
// mapped.
var ErrInvalidNumberOfPods = errors.New("invalid number of pods")

// ErrInvalidNumberOfOrderFragments is returned when a pod is mapped to an
// insufficient number of order fragments, or too many order fragments.
var ErrInvalidNumberOfOrderFragments = errors.New("invalid number of order fragments")

// ErrCannotOpenOrderFragments is returned when none of the pods were available
// to receive order fragments
var ErrCannotOpenOrderFragments = errors.New("cannot open order fragments: no pod received an order fragment")

// An OrderFragmentMapping maps pods to encrypted order fragments.
type OrderFragmentMapping map[[32]byte][]OrderFragment

type OrderFragment struct {
	order.EncryptedFragment
	Index int64
}

type OpenOrderRequest struct {
	order.ID
	signature            [65]byte
	orderFragmentMapping OrderFragmentMapping
}

// Ingress interface can open and cancel orders on behalf of a user.
type Ingress interface {

	// OpenOrder on the Ren Ledger and on the Darkpool. A signature from the
	// trader identifies them as the owner, the order ID is submitted to the
	// Ren Ledger along with the necessary fee, and the order fragment mapping
	// is used to send order fragments to pods in the Darkpool.
	OpenOrder(signature [65]byte, orderID order.ID, orderFragmentMapping OrderFragmentMapping) error

	// OpenOrderProcess starts reading from the openOrderQueue to process
	// new open order requests. A done channel must be passed and when this
	// done channel is closed by the user, the openOrderQueue will be closed.
	OpenOrderProcess(done chan struct{}) <-chan error

	// OpenOrderFragmentsProcess starts reading from the openOrderFragmentsQueue
	// to process order fragments. A done channel must be passed and when this
	// done channel is closed by the user, the openOrderFragmentsQueue will be
	// closed.
	OpenOrderFragmentsProcess(done chan struct{}) <-chan error

	// CancelOrder on the Ren Ledger. A signature from the trader is needed to
	// verify the cancelation.
	CancelOrder(signature [65]byte, orderID order.ID) error

	// Sync Darkpool to ensure an up-to-date state.
	Sync(<-chan struct{}) <-chan error
}

type ingress struct {
	darkpool                cal.Darkpool
	renLedger               cal.RenLedger
	swarmer                 swarm.Swarmer
	orderbookClient         orderbook.Client
	podsMu                  sync.RWMutex
	pods                    map[[32]byte]cal.Pod
	openOrderQueue          chan OpenOrderRequest
	openOrderFragmentsQueue chan OpenOrderRequest
}

// NewIngress creates an ingress and starts reading from openOrderQueue
// and openOrderFragmentsQueue. A call to NewIngress must be completed
// with a call to the StopIngress function which will close all the channels.
func NewIngress(darkpool cal.Darkpool, renLedger cal.RenLedger, swarmer swarm.Swarmer, orderbookClient orderbook.Client) Ingress {
	ingress := &ingress{
		darkpool:        darkpool,
		renLedger:       renLedger,
		swarmer:         swarmer,
		orderbookClient: orderbookClient,
		pods:            map[[32]byte]cal.Pod{},

		openOrderQueue:          make(chan OpenOrderRequest),
		openOrderFragmentsQueue: make(chan OpenOrderRequest),
	}
	return ingress
}

func (ingress *ingress) OpenOrder(signature [65]byte, orderID order.ID, orderFragmentMapping OrderFragmentMapping) error {
	// TODO: Verify that the signature is valid before sending it to the
	// RenLedger. This is not strictly necessary but it can save the Ingress
	// some gas.
	if err := ingress.verifyOrderFragments(orderFragmentMapping); err != nil {
		return err
	}

	ingress.openOrderQueue <- OpenOrderRequest{
		ID:                   orderID,
		signature:            signature,
		orderFragmentMapping: orderFragmentMapping,
	}

	return nil
}

func (ingress *ingress) OpenOrderProcess(done chan struct{}) <-chan error {
	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		for {
			select {
			case <-done:
				close(ingress.openOrderQueue)
				return
			case request, ok := <-ingress.openOrderQueue:
				if !ok {
					return
				}
				if err := ingress.processOpenOrderRequest(request); err != nil {
					select {
					case <-done:
						return
					case errs <- err:
					}
				}
			}
		}
	}()
	return errs
}

func (ingress *ingress) OpenOrderFragmentsProcess(done chan struct{}) <-chan error {
	errs := make(chan error, runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		go func(i int) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case request, ok := <-ingress.openOrderFragmentsQueue:
					if !ok {
						return
					}
					ingress.podsMu.RLock()
					if err := ingress.processOpenOrderFragmentsRequest(request, ingress.pods); err != nil {
						select {
						case <-done:
						case errs <- err:
						}
					}
					ingress.podsMu.RUnlock()
				}
			}
		}(i)
	}

	go func() {
		defer close(errs)
		wg.Wait()
	}()

	return errs
}

func (ingress *ingress) CancelOrder(signature [65]byte, orderID order.ID) error {
	// TODO: Verify that the signature is valid before sending it to the
	// RenLedger. This is not strictly necessary but it can save the Ingress
	// some gas.
	if err := ingress.renLedger.CancelOrder(signature, orderID); err != nil {
		return err
	}
	return nil
}

func (ingress *ingress) Sync(done <-chan struct{}) <-chan error {
	errs := make(chan error, 1)

	epoch, err := ingress.darkpool.Epoch()
	if err != nil {
		errs <- err
		return errs
	}

	pods, err := ingress.darkpool.Pods()
	if err != nil {
		errs <- err
		return errs
	} else {
		ingress.podsMu.Lock()
		ingress.pods = map[[32]byte]cal.Pod{}
		for _, pod := range pods {
			ingress.pods[pod.Hash] = pod
		}
		ingress.podsMu.Unlock()
	}

	currentEpochHash := epoch.Hash

	go func() {
		defer close(errs)
		for {
			// Sync with the approximate block time
			time.Sleep(14 * time.Second)

			select {
			case <-done:
				return
			default:
			}

			currentEpoch, err := ingress.darkpool.Epoch()
			if err != nil {
				select {
				case <-done:
					return
				case errs <- err:
				}
			}

			if currentEpochHash != currentEpoch.Hash {
				pods, err := ingress.darkpool.Pods()
				if err != nil {
					select {
					case <-done:
						return
					case errs <- err:
						continue
					}
				}
				ingress.podsMu.Lock()
				ingress.pods = map[[32]byte]cal.Pod{}
				for _, pod := range pods {
					ingress.pods[pod.Hash] = pod
				}
				ingress.podsMu.Unlock()
			}
		}
	}()
	return errs
}

func (ingress *ingress) verifyOrderFragments(orderFragmentMapping OrderFragmentMapping) error {
	ingress.podsMu.RLock()
	defer ingress.podsMu.RUnlock()

	if len(orderFragmentMapping) == 0 || len(orderFragmentMapping) > len(ingress.pods) {
		return ErrInvalidNumberOfPods
	}

	for hash, orderFragments := range orderFragmentMapping {
		pod, ok := ingress.pods[hash]
		if !ok {
			return ErrUnknownPod
		}
		if len(orderFragments) > len(pod.Darknodes) || len(orderFragments) < pod.Threshold() {
			return ErrInvalidNumberOfOrderFragments
		}
	}
	return nil
}

func (ingress *ingress) sendOrderFragmentsToPod(pod cal.Pod, orderFragments []OrderFragment) error {
	if len(orderFragments) < pod.Threshold() || len(orderFragments) > len(pod.Darknodes) {
		return ErrInvalidNumberOfOrderFragments
	}

	// Map order fragments to their respective Darknodes
	orderFragmentIndexMapping := map[int64]OrderFragment{}
	for _, orderFragment := range orderFragments {
		orderFragmentIndexMapping[orderFragment.Index] = orderFragment
	}

	errs := make(chan error, len(pod.Darknodes))
	go func() {
		defer close(errs)

		dispatch.CoForAll(pod.Darknodes, func(i int) {
			orderFragment, ok := orderFragmentIndexMapping[int64(i)]
			if !ok {
				errs <- fmt.Errorf("no fragment found at index %v", i)
				return
			}
			darknode := pod.Darknodes[i]

			// Send the order fragment to the Darknode
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			darknodeMultiAddr, err := ingress.swarmer.Query(ctx, darknode, -1)
			if err != nil {
				errs <- fmt.Errorf("cannot send query to %v: %v", darknode, err)
				return
			}

			log.Printf("sending order fragment to %v", darknode)
			if err := ingress.orderbookClient.OpenOrder(ctx, darknodeMultiAddr, orderFragment.EncryptedFragment); err != nil {
				log.Printf("cannot send order fragment to %v: %v", darknode, err)
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

func (ingress *ingress) processOpenOrderFragmentsRequest(request OpenOrderRequest, pods map[[32]byte]cal.Pod) error {
	errs := make([]error, 0, len(pods))
	podDidReceiveFragments := int64(0)

	dispatch.CoForAll(pods, func(hash [32]byte) {
		orderFragments := request.orderFragmentMapping[hash]
		if orderFragments != nil && len(orderFragments) > 0 {
			if err := ingress.sendOrderFragmentsToPod(pods[hash], orderFragments); err != nil {
				errs = append(errs, err)
				return
			}
			if atomic.LoadInt64(&podDidReceiveFragments) == int64(0) {
				atomic.AddInt64(&podDidReceiveFragments, 1)
			}
		}
	})

	if atomic.LoadInt64(&podDidReceiveFragments) == int64(0) {
		if len(errs) == 0 {
			return ErrCannotOpenOrderFragments
		} else {
			return fmt.Errorf("%v %v", ErrCannotOpenOrderFragments.Error(), errs[0])
		}
	}
	return nil
}

func (ingress *ingress) processOpenOrderRequest(request OpenOrderRequest) error {
	var orderParity order.Parity
	for _, orderFragments := range request.orderFragmentMapping {
		if len(orderFragments) > 1 {
			orderParity = orderFragments[0].OrderParity
			break
		}
	}

	var err error
	if orderParity == order.ParityBuy {
		err = ingress.renLedger.OpenBuyOrder(request.signature, request.ID)
	} else {
		err = ingress.renLedger.OpenSellOrder(request.signature, request.ID)
	}
	if err != nil {
		return err
	}
	ingress.openOrderFragmentsQueue <- request

	return nil
}
