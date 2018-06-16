package ingress

import (
	"bytes"
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
	"github.com/republicprotocol/republic-go/logger"
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

// NumBackgroundWorkers is the number of background workers that the Ingress
// will use.
var NumBackgroundWorkers = runtime.NumCPU() * 4

// An OrderFragmentMapping maps pods to encrypted order fragments.
type OrderFragmentMapping map[[32]byte][]OrderFragment

// OrderFragment has an order.EncryptedFragment, encrypted by the trader before
// being sent to the Ingress, and the required index that identifies which set
// shares are held by the order.EncryptedFragment.
type OrderFragment struct {
	order.EncryptedFragment
	Index int64
}

// Ingress interface can open and cancel orders on behalf of a user.
type Ingress interface {

	// Sync the epoch.
	Sync(<-chan struct{}) <-chan error

	// OpenOrder on the Ren Ledger and on the Darkpool. A signature from the
	// trader identifies them as the owner, the order ID is submitted to the
	// Ren Ledger along with the necessary fee, and the order fragment mapping
	// is used to send order fragments to pods in the Darkpool.
	OpenOrder(signature [65]byte, orderID order.ID, orderFragmentMapping OrderFragmentMapping) error

	// CancelOrder on the Ren Ledger. A signature from the trader is needed to
	// verify the cancelation.
	CancelOrder(signature [65]byte, orderID order.ID) error

	// ProcessRequests in the background. Closing the done channel will stop
	// all processing. Running this background worker is required to open and
	// cancel orders.
	ProcessRequests(done <-chan struct{}) <-chan error
}

type ingress struct {
	darkpool        cal.Darkpool
	renLedger       cal.RenLedger
	swarmer         swarm.Swarmer
	orderbookClient orderbook.Client

	podsMu *sync.RWMutex
	pods   map[[32]byte]cal.Pod

	queueRequests              chan Request
	queueOrderFragmentMappings chan OrderFragmentMapping
}

// NewIngress returns an Ingress. The background services of the Ingress must
// be started separately by calling Ingress.OpenOrderProcess and
// Ingress.OpenOrderFragmentsProcess.
func NewIngress(darkpool cal.Darkpool, renLedger cal.RenLedger, swarmer swarm.Swarmer, orderbookClient orderbook.Client) Ingress {
	ingress := &ingress{
		darkpool:        darkpool,
		renLedger:       renLedger,
		swarmer:         swarmer,
		orderbookClient: orderbookClient,

		podsMu: new(sync.RWMutex),
		pods:   map[[32]byte]cal.Pod{},

		queueRequests:              make(chan Request, 1024),
		queueOrderFragmentMappings: make(chan OrderFragmentMapping, 1024),
	}
	return ingress
}

// Sync implements the Ingress interface.
func (ingress *ingress) Sync(done <-chan struct{}) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		epochIntervalBig, err := ingress.darkpool.MinimumEpochInterval()
		if err != nil {
			errs <- err
			return
		}
		epochInterval := epochIntervalBig.Int64()
		if epochInterval < 100 {
			// An Ingress will not trigger epochs faster than once every 100
			// blocks
			epochInterval = 100
		}
		epoch := cal.Epoch{}

		ticks := int64(0)
		ticker := time.NewTicker(time.Second * 14)
		defer ticker.Stop()

		for {
			ticks++
			if ticks%epochInterval == 0 {
				logger.Info(fmt.Sprintf("queueing syncing of epoch"))
				select {
				case <-done:
				case ingress.queueRequests <- EpochRequest{}:
				}
			}

			func() {
				nextEpoch, err := ingress.darkpool.Epoch()
				if err != nil {
					select {
					case <-done:
					case errs <- err:
					}
					return
				}
				if bytes.Equal(epoch.Hash[:], nextEpoch.Hash[:]) {
					return
				}
				epoch = nextEpoch
				if err := ingress.syncFromEpoch(epoch); err != nil {
					select {
					case <-done:
					case errs <- err:
					}
					return
				}
			}()

			select {
			case <-done:
				return
			case <-ticker.C:
			}
		}
	}()

	return errs
}

func (ingress *ingress) OpenOrder(signature [65]byte, orderID order.ID, orderFragmentMapping OrderFragmentMapping) error {
	// TODO: Verify that the signature is valid before sending it to the
	// RenLedger. This is not strictly necessary but it can save the Ingress
	// some gas.
	if err := ingress.verifyOrderFragmentMapping(orderFragmentMapping); err != nil {
		return err
	}

	go func() {
		logger.Info(fmt.Sprintf("queueing opening of order %v", orderID))
		ingress.queueRequests <- OpenOrderRequest{
			signature:            signature,
			orderID:              orderID,
			orderFragmentMapping: orderFragmentMapping,
		}
	}()
	return nil
}

func (ingress *ingress) CancelOrder(signature [65]byte, orderID order.ID) error {
	// TODO: Verify that the signature is valid beforNumBackgroundWorkerse sending it to the
	// RenLedger. This is not strictly necessary but it can save the Ingress
	// some gas.
	go func() {
		logger.Info(fmt.Sprintf("queueing cancelation of order %v", orderID))
		ingress.queueRequests <- CancelOrderRequest{
			signature: signature,
			orderID:   orderID,
		}
	}()
	return nil
}

func (ingress *ingress) ProcessRequests(done <-chan struct{}) <-chan error {
	errs := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ingress.processRequests(done, errs)
	}()

	go func() {
		defer wg.Done()
		ingress.processOrderFragmentMappingQueue(done, errs)
	}()

	go func() {
		defer close(errs)
		wg.Wait()
	}()

	return errs
}

func (ingress *ingress) syncFromEpoch(epoch cal.Epoch) error {
	logger.Epoch(epoch.Hash)
	pods, err := ingress.darkpool.Pods()
	if err != nil {
		return err
	}
	ingress.podsMu.Lock()
	ingress.pods = map[[32]byte]cal.Pod{}
	for _, pod := range pods {
		ingress.pods[pod.Hash] = pod
	}
	ingress.podsMu.Unlock()
	return nil
}

func (ingress *ingress) processRequests(done <-chan struct{}, errs chan<- error) {

	// openOrderRequests := make([]OpenOrderRequest, NumBackgroundWorkers)
	// cancelOrderRequests := make([]CancelOrderRequest, NumBackgroundWorkers)

	for {
		select {
		case <-done:
			return
		case request, ok := <-ingress.queueRequests:
			if !ok {
				return
			}

			logger.Info(fmt.Sprintf("received request of type %T", request))

			switch req := request.(type) {
			case EpochRequest:
				ingress.processEpochRequest(req, done, errs)
			case OpenOrderRequest:
				ingress.processOpenOrderRequest(req, done, errs)
			case CancelOrderRequest:
				ingress.processCancelOrderRequest(req, done, errs)
			default:
				logger.Error(fmt.Sprintf("unexpected request type %T", request))
			}

			// openOrderRequestsN := 0
			// cancelOrderRequestsN := 0

			// switch req := request.(type) {
			// case OpenOrderRequest:
			// 	openOrderRequests[openOrderRequestsN] = req
			// 	openOrderRequestsN++
			// case CancelOrderRequest:
			// 	cancelOrderRequests[cancelOrderRequestsN] = req
			// 	cancelOrderRequestsN++
			// default:
			// 	logger.Error(fmt.Sprintf("unexpected request type %T", request))
			// }

			// for i := 1; i < NumBackgroundWorkers; i++ {
			// 	select {
			// 	case request, ok := <-ingress.queueRequests:
			// 		if !ok {
			// 			break
			// 		}
			// 		switch req := request.(type) {
			// 		case OpenOrderRequest:
			// 			openOrderRequests[openOrderRequestsN] = req
			// 			openOrderRequestsN++
			// 		case CancelOrderRequest:
			// 			cancelOrderRequests[cancelOrderRequestsN] = req
			// 			cancelOrderRequestsN++
			// 		default:
			// 			logger.Error(fmt.Sprintf("unexpected request type %T", request))
			// 		}
			// 	default:
			// 	}
			// }

			// if openOrderRequestsN > 0 {
			// 	ingress.processOpenOrderRequests(openOrderRequests[:openOrderRequestsN], done, errs)
			// }
			// if cancelOrderRequestsN > 0 {
			// 	ingress.processCancelOrderRequests(cancelOrderRequests[:cancelOrderRequestsN], done, errs)
			// }
		}
	}
}

func (ingress *ingress) processEpochRequest(req EpochRequest, done <-chan struct{}, errs chan<- error) {
	if _, err := ingress.darkpool.NextEpoch(); err != nil {
		select {
		case <-done:
		case errs <- err:
		}
	}
}

func (ingress *ingress) processOpenOrderRequest(req OpenOrderRequest, done <-chan struct{}, errs chan<- error) {
	// signatures := make([][65]byte, len(reqs))
	// orderIDs := make([]order.ID, len(reqs))
	// orderParities := make([]order.Parity, len(reqs))

	// for i := range reqs {
	// 	var orderParity order.Parity
	// 	for _, orderFragments := range reqs[i].orderFragmentMapping {
	// 		if len(orderFragments) > 1 {
	// 		5000000000ments[0].OrderParity
	// 		5000000000
	// 		}5000000000
	// 	}
	// 	signatures[i] = reqs[i].signature
	// 	orderIDs[i] = reqs[i].orderID
	// 	orderParities[i] = orderParity
	// }

	// n, err := ingress.renLedger.OpenOrders(signatures, orderIDs, orderParities)
	// if err != nil {5000000000
	// 	select {
	// 	case <-done:
	// 	case errs <- err:
	// 	}
	// 	return
	// }

	// for i := 0; i < n; i++ {
	// 	select {
	// 	case <-done:
	// 	case ingress.queueOrderFragmentMappings <- reqs[i].orderFragmentMapping:
	// 	}
	// }
	var orderParity order.Parity
	for _, orderFragments := range req.orderFragmentMapping {
		if len(orderFragments) > 1 {
			orderParity = orderFragments[0].OrderParity
			break
		}
	}

	var err error
	if orderParity == order.ParityBuy {
		err = ingress.renLedger.OpenBuyOrder(req.signature, req.orderID)
	} else {
		err = ingress.renLedger.OpenSellOrder(req.signature, req.orderID)
	}
	if err != nil {
		select {
		case <-done:
		case errs <- err:
		}
	}

	select {
	case <-done:
	case ingress.queueOrderFragmentMappings <- req.orderFragmentMapping:
	}
}

func (ingress *ingress) processCancelOrderRequest(req CancelOrderRequest, done <-chan struct{}, errs chan<- error) {
	// for _, req := range reqs {
	// 	if err := ingress.renLedger.CancelOrder(req.signature, req.orderID); err != nil {
	// 		return err
	// 	}
	// }
	// return nil
	if err := ingress.renLedger.CancelOrder(req.signature, req.orderID); err != nil {
		select {
		case <-done:
		case errs <- err:
		}
	}
}

func (ingress *ingress) processOrderFragmentMappingQueue(done <-chan struct{}, errs chan<- error) {
	dispatch.CoForAll(NumBackgroundWorkers, func(i int) {
		for {
			select {
			case <-done:
				return
			case orderFragmentMapping, ok := <-ingress.queueOrderFragmentMappings:
				if !ok {
					return
				}
				if err := ingress.processOrderFragmentMapping(orderFragmentMapping); err != nil {
					select {
					case <-done:
						return
					case errs <- err:
					}
				}
			}
		}
	})
}

func (ingress *ingress) processOrderFragmentMapping(orderFragmentMapping OrderFragmentMapping) error {
	ingress.podsMu.RLock()
	defer ingress.podsMu.RUnlock()

	errs := make([]error, 0, len(ingress.pods))
	podDidReceiveFragments := int64(0)

	dispatch.CoForAll(ingress.pods, func(hash [32]byte) {
		orderFragments := orderFragmentMapping[hash]
		if orderFragments != nil && len(orderFragments) > 0 {
			if err := ingress.sendOrderFragmentsToPod(ingress.pods[hash], orderFragments); err != nil {
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
		}
		return fmt.Errorf("%v %v", ErrCannotOpenOrderFragments.Error(), errs[0])
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

		logger.Network(logger.LevelInfo, fmt.Sprintf("sending %v order = %v to pod = %v", orderFragments[0].OrderParity, orderFragments[0].OrderID, base64.StdEncoding.EncodeToString(pod.Hash[:8])))

		dispatch.CoForAll(pod.Darknodes, func(i int) {
			orderFragment, ok := orderFragmentIndexMapping[int64(i+1)] // Indices for fragments start at 1
			if !ok {
				errs <- fmt.Errorf("no fragment found at index %v", i)
				return
			}
			darknode := pod.Darknodes[i]

			// Send the order fragment to the Darknode
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			darknodeMultiAddr, err := ingress.swarmer.Query(ctx, darknode, -1)
			if err != nil {
				errs <- fmt.Errorf("cannot send query to %v: %v", darknode, err)
				return
			}

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

	// Check if at least 2/3 of the nodes in the specified pod have received
	// the order fragments.
	errNumMax := len(orderFragments) - pod.Threshold()
	if len(pod.Darknodes) > 0 && errNum > errNumMax {
		return fmt.Errorf("cannot send order fragments to %v nodes (out of %v nodes) in pod %v: %v", errNum, len(pod.Darknodes), base64.StdEncoding.EncodeToString(pod.Hash[:]), err)
	}
	return nil
}

func (ingress *ingress) verifyOrderFragmentMapping(orderFragmentMapping OrderFragmentMapping) error {
	ingress.podsMu.RLock()
	defer ingress.podsMu.RUnlock()

	if len(orderFragmentMapping) == 0 || len(orderFragmentMapping) > len(ingress.pods) {
		logger.Error(fmt.Sprintf("invalid number of pods: got %v, expected %v", len(orderFragmentMapping), len(ingress.pods)))
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
