package orderbook

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
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
	pointerStore       PointerStorer
	orderStore         OrderStorer
	orderFragmentStore OrderFragmentStorer
	contractBinder     ContractBinder
	interval           time.Duration
	limit              int

	syncerMu                 *sync.RWMutex
	syncerCurrDone           chan struct{}
	syncerCurrOrderFragments chan order.Fragment
	syncerPrevDone           chan struct{}
	syncerPrevOrderFragments chan order.Fragment

	broadcastNotifications chan (<-chan Notification)
	broadcastErrs          chan (<-chan error)
}

// NewOrderbook returns an Orderbok that uses a crypto.RsaKey to decrypt the
// order.EncryptedFragments that it receives, and stores them in a Storer.
func NewOrderbook(rsaKey crypto.RsaKey, pointerStore PointerStorer, orderStore OrderStorer, orderFragmentStore OrderFragmentStorer, contractBinder ContractBinder, interval time.Duration, limit int) Orderbook {
	return &orderbook{
		doneMu: new(sync.RWMutex),
		done:   nil,

		rsaKey:             rsaKey,
		pointerStore:       pointerStore,
		orderStore:         orderStore,
		orderFragmentStore: orderFragmentStore,
		contractBinder:     contractBinder,
		interval:           interval,
		limit:              limit,

		syncerMu:                 new(sync.RWMutex),
		syncerCurrDone:           nil,
		syncerCurrOrderFragments: nil,
		syncerPrevDone:           nil,
		syncerPrevOrderFragments: nil,

		broadcastNotifications: make(chan (<-chan Notification)),
		broadcastErrs:          make(chan (<-chan error)),
	}
}

// OpenOrder implements the Server interface.
func (orderbook *orderbook) OpenOrder(ctx context.Context, encryptedOrderFragment order.EncryptedFragment) error {
	orderFragment, err := encryptedOrderFragment.Decrypt(orderbook.rsaKey.PrivateKey)
	if err != nil {
		return err
	}
	if orderFragment.OrderParity == order.ParityBuy {
		logger.BuyOrderReceived(logger.LevelDebugLow, orderFragment.OrderID.String(), orderFragment.ID.String())
	} else {
		logger.SellOrderReceived(logger.LevelDebugLow, orderFragment.OrderID.String(), orderFragment.ID.String())
	}

	return orderbook.routeOrderFragment(ctx, orderFragment)
}

// Sync implements the Orderbook interface.
func (orderbook *orderbook) Sync(done <-chan struct{}) (<-chan Notification, <-chan error) {
	notifications := make(chan Notification, 4096)
	errs := make(chan error, 4096)

	// Check whether the server has already been started
	orderbook.doneMu.RLock()
	serverIsRunning := orderbook.done != nil
	orderbook.doneMu.RUnlock()
	if serverIsRunning {
		errs <- ErrServerIsRunning
		close(notifications)
		close(errs)
		return notifications, errs
	}

	// Open a new done channel that will be used from now on to signal the
	// closure of the server
	orderbook.doneMu.Lock()
	orderbook.done = make(chan struct{})
	orderbook.doneMu.Unlock()
	go func() {
		<-done
		close(orderbook.done)
		orderbook.doneMu.Lock()
		orderbook.done = nil
		orderbook.doneMu.Unlock()
	}()

	// Merge all of the channels on the broadcast channel into the output
	// channel
	go func() {
		defer close(notifications)
		dispatch.Merge(done, orderbook.broadcastNotifications, notifications)
	}()
	go func() {
		defer close(errs)
		dispatch.Merge(done, orderbook.broadcastErrs, errs)
	}()

	return notifications, errs
}

// OnChangeEpoch implements the Orderbook interface.
func (orderbook *orderbook) OnChangeEpoch(epoch registry.Epoch) {
	func() {
		orderbook.doneMu.RLock()
		orderbook.syncerMu.Lock()
		defer orderbook.doneMu.RUnlock()
		defer orderbook.syncerMu.Unlock()

		// Transition the current epoch into the previous epoch and setup a new
		// current epoch
		if orderbook.syncerPrevDone != nil {
			close(orderbook.syncerPrevDone)
			close(orderbook.syncerPrevOrderFragments)
		}
		orderbook.syncerPrevDone = orderbook.syncerCurrDone
		orderbook.syncerPrevOrderFragments = orderbook.syncerCurrOrderFragments
		orderbook.syncerCurrDone = make(chan struct{})
		orderbook.syncerCurrOrderFragments = make(chan order.Fragment)
	}()

	// Clone a new PointerStorer for the new syncer
	pointerStore, err := orderbook.pointerStore.Clone()
	if err != nil {
		logger.Error(fmt.Sprintf("cannot clone pointer store: %v", err))
		return
	}

	syncer := newSyncer(epoch, pointerStore, orderbook.orderStore, orderbook.orderFragmentStore, orderbook.contractBinder, orderbook.interval, orderbook.limit)
	notifications, errs := syncer.sync(orderbook.syncerCurrDone, orderbook.syncerCurrOrderFragments)

	select {
	case <-orderbook.done:
	case orderbook.broadcastNotifications <- notifications:
	}

	select {
	case <-orderbook.done:
	case orderbook.broadcastErrs <- errs:
	}
}

func (orderbook *orderbook) routeOrderFragment(ctx context.Context, orderFragment order.Fragment) error {
	orderbook.doneMu.RLock()
	serverIsRunning := orderbook.done != nil
	orderbook.doneMu.RUnlock()
	if !serverIsRunning {
		return ErrServerIsNotRunning
	}

	orderbook.syncerMu.RLock()
	defer orderbook.syncerMu.RUnlock()

	// TODO: Using the block number of the order, the orderbook should infer
	// which epoch the order fragment is destined for. If the epoch is unknown
	// then the orderbook should sleep here (this is safe given that this
	// function is generally called in a background goroutine) and try again.
	// Failing a second time should see the order fragment dropped. This helps
	// with robust acceptance of order fragments at the turn of an epoch where
	// the Darknode and the trader might briefly have different ideas about the
	// "current" epoch.

	switch orderFragment.EpochDepth {
	case 0:
		logger.Network(logger.LevelInfo, fmt.Sprintf("routing order %v to depth = %v", orderFragment.OrderID, orderFragment.EpochDepth))
		if orderbook.syncerCurrOrderFragments == nil {
			logger.Network(logger.LevelWarn, fmt.Sprintf("cannot route order %v to depth = %v", orderFragment.OrderID, orderFragment.EpochDepth))
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-orderbook.done:
			return ErrServerIsNotRunning
		case orderbook.syncerCurrOrderFragments <- orderFragment:
			return nil
		}
	case 1:
		logger.Network(logger.LevelInfo, fmt.Sprintf("routing order %v to depth = %v", orderFragment.OrderID, orderFragment.EpochDepth))
		if orderbook.syncerPrevOrderFragments == nil {
			logger.Network(logger.LevelWarn, fmt.Sprintf("cannot route order %v to depth = %v", orderFragment.OrderID, orderFragment.EpochDepth))
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-orderbook.done:
			return ErrServerIsNotRunning
		case orderbook.syncerPrevOrderFragments <- orderFragment:
			return nil
		}
	default:
		logger.Network(logger.LevelWarn, fmt.Sprintf("cannot route order %v to depth = %v", orderFragment.OrderID, orderFragment.EpochDepth))
		return nil
	}
}
