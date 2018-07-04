package orderbook

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
)

// ErrServerIsRunning is returned when a Server that has already been started
// is started again.
var ErrServerIsRunning = errors.New("server is running")

// ErrServerIsNotRunning is returned when a Server receives the shutdown signal
// while handling an RPC, or when a Server handles an RPC before being started.
var ErrServerIsNotRunning = errors.New("server is not running")

// Client for invoking the Server.OpenOrder ROC on a remote Server.
type Client interface {

	// OpenOrder by sending an order.EncryptedFragment to an
	// identity.MultiAddress. The order.EncryptedFragment will be stored by the
	// Server.
	OpenOrder(context.Context, identity.MultiAddress, order.EncryptedFragment) error
}

// A Server expose RPCs for opening orders with a Darknode by sending it an
// order.EncryptedFragment.
type Server interface {

	// OpenOrder is the RPC invoked by a Client. It accepts an
	// order.EncryptedFragment and decrypts it. Until this RPC is invoked, the
	// Server cannot participate in the respective secure multi-party
	// computation.
	OpenOrder(context.Context, order.EncryptedFragment) error
}

// An Orderbook combines the Server interface with synchronising order.IDs from
// Ethereum. Once a synchronised order.ID reaches the order.Open status, and
// the respective order.EncryptedFragment is received, the decrypted
// order.Fragment stored inside a Notification and produced. For all other
// changes, a Notification is produced directly from the change.
type Orderbook interface {
	Server

	// Sync status changes from Ethereum, receive order.EncryptedFragments from
	// traders, and produce Notifications. Stop once the done channel is
	// closed. An error is returned when a call to Orderbook.Sync happens
	// before the previous call has stopped its done channel.
	Sync(done <-chan struct{}) (<-chan Notification, <-chan error)

	// OnChangeEpoch should be called whenever a change to the registry.Epoch
	// is detected.
	OnChangeEpoch(epoch registry.Epoch)
}

type orderbook struct {
	doneMu *sync.RWMutex
	done   chan struct{}

	rsaKey             crypto.RsaKey
	orderStore         OrderStorer
	orderFragmentStore OrderFragmentStorer
	contractBinder     ContractBinder

	syncerMu                 *sync.Mutex
	syncerCurrDone           chan struct{}
	syncerCurrOrderFragments chan order.Fragment
	syncerPrevDone           chan struct{}
	syncerPrevOrderFragments chan order.Fragment

	notificationMerger chan (<-chan Notification)
	errMerger          chan (<-chan error)
}

// NewOrderbook returns an Orderbok that uses a crypto.RsaKey to decrypt the
// order.EncryptedFragments that it receives, and stores them in a Storer.
func NewOrderbook(rsaKey crypto.RsaKey, orderStore OrderStorer, orderFragmentStore OrderFragmentStorer, contractBinder ContractBinder) Orderbook {
	return &orderbook{
		doneMu: new(sync.RWMutex),
		done:   nil,

		rsaKey:             rsaKey,
		orderStore:         orderStore,
		orderFragmentStore: orderFragmentStore,
		contractBinder:     contractBinder,

		syncerMu:                 new(sync.Mutex),
		syncerCurrDone:           nil,
		syncerCurrOrderFragments: nil,
		syncerPrevDone:           nil,
		syncerPrevOrderFragments: nil,

		notificationMerger: make(chan (<-chan Notification)),
		errMerger:          make(chan (<-chan error)),
	}
}

// OpenOrder implements the Server interface.
func (orderbook *orderbook) OpenOrder(ctx context.Context, encryptedOrderFragment order.EncryptedFragment) error {
	orderbook.doneMu.RLock()
	defer orderbook.doneMu.RUnlock()

	if orderbook.done == nil {
		return ErrServerIsNotRunning
	}

	orderFragment, err := encryptedOrderFragment.Decrypt(*orderbook.rsaKey.PrivateKey)
	if err != nil {
		return err
	}
	if orderFragment.OrderParity == order.ParityBuy {
		logger.BuyOrderReceived(logger.LevelDebugLow, orderFragment.OrderID.String(), orderFragment.ID.String())
	} else {
		logger.SellOrderReceived(logger.LevelDebugLow, orderFragment.OrderID.String(), orderFragment.ID.String())
	}

	switch orderFragment.Depth {
	case 0:
		select {
		case <-ctx.Done():
			return nil
		case <-orderbook.done:
			return ErrServerIsNotRunning
		case orderbook.syncerCurrOrderFragments <- orderFragment:
			return nil
		}
	case 1:
		select {
		case <-ctx.Done():
			return nil
		case <-orderbook.done:
			return ErrServerIsNotRunning
		case orderbook.syncerPrevOrderFragments <- orderFragment:
			return nil
		}
	default:
		return nil
	}
}

// Sync implements the Orderbook interface.
func (orderbook *orderbook) Sync(done <-chan struct{}) (<-chan Notification, <-chan error) {
	notifications := make(chan Notification)
	errs := make(chan error, 1)

	// Check whether the Server has already been started
	orderbook.doneMu.RLock()
	serverIsRunning := orderbook.done != nil
	orderbook.doneMu.RUnlock()
	if serverIsRunning {
		errs <- ErrServerIsRunning
		close(notifications)
		close(errs)
		return notifications, errs
	}

	// Open a new done channel that signals the Server has started
	orderbook.doneMu.Lock()
	orderbook.done = make(chan struct{})
	orderbook.doneMu.Unlock()

	// Wait for the shutdown signal and then stop the Server
	go func() {
		<-done

		close(orderbook.done)
		orderbook.doneMu.Lock()
		orderbook.done = nil
		orderbook.doneMu.Unlock()
	}()

	// Merge all outputs from the syncers
	go func() {
		defer close(notifications)
		orderbook.mergeNotifications(done, notifications)
	}()
	go func() {
		defer close(errs)
		orderbook.mergeErrors(done, errs)
	}()

	return notifications, errs
}

// OnChangeEpoch implements the Orderbook interface.
func (orderbook *orderbook) OnChangeEpoch(epoch registry.Epoch) {
	orderbook.syncerMu.Lock()
	defer orderbook.syncerMu.Unlock()

	// Close the previous epoch
	if orderbook.syncerPrevDone == nil {
		close(orderbook.syncerPrevDone)
		close(orderbook.syncerPrevOrderFragments)
	}
	// Transition the current epoch to be the previous epoch and create a new
	// epoch setup
	orderbook.syncerPrevDone = orderbook.syncerCurrDone
	orderbook.syncerPrevOrderFragments = orderbook.syncerCurrOrderFragments
	orderbook.syncerCurrDone = make(chan struct{})
	orderbook.syncerCurrOrderFragments = make(chan order.Fragment)

	// Get the minimum epoch interval and retry several times on failure
	var minimumEpochInterval *big.Int
	var err error
	for i := 0; i < 3; i++ {
		minimumEpochInterval, err = orderbook.contractBinder.MinimumEpochInterval()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
	}
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get minimum epoch interval: %v, defaulting to: 10_000", err))
		minimumEpochInterval = big.NewInt(10000)
	}

	// Start the syncer for this epoch
	syncer := newSyncer(
		orderbook.orderStore,
		orderbook.orderFragmentStore,
		orderbook.contractBinder,
		big.NewInt(0).SetUint64(uint64(epoch.BlockNumber)), // Block offset
		minimumEpochInterval,                               // Block limit
		0,                                                  // Order synchronising offset
		1024,                                               // Order synchronising limit
		epoch,                                              // Epoch
	)
	notifications, errs := syncer.sync(orderbook.syncerCurrDone, orderbook.syncerCurrOrderFragments)

	// Signal that the outputs of this syncer should be accepted by the merger
	select {
	case <-orderbook.done:
	case orderbook.notificationMerger <- notifications:
	}
	select {
	case <-orderbook.done:
	case orderbook.errMerger <- errs:
	}
}

func (orderbook *orderbook) mergeNotifications(done <-chan struct{}, notifications chan<- Notification) {
	for {
		select {
		case <-done:
			return
		case ch, ok := <-orderbook.notificationMerger:
			if !ok {
				return
			}
			go func() {
				for {
					select {
					case <-done:
						return
					case notification, ok := <-ch:
						if !ok {
							return
						}
						select {
						case <-done:
						case notifications <- notification:
						}
					}
				}
			}()
		}
	}
}

func (orderbook *orderbook) mergeErrors(done <-chan struct{}, errs chan<- error) {
	for {
		select {
		case <-done:
			return
		case ch, ok := <-orderbook.errMerger:
			if !ok {
				return
			}
			go func() {
				for {
					select {
					case <-done:
						return
					case err, ok := <-ch:
						if !ok {
							return
						}
						select {
						case <-done:
						case errs <- err:
						}
					}
				}
			}()
		}
	}
}
