package orderbook

import (
	"context"
	"errors"
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

	notificationMerger chan (<-chan Notification)
	errMerger          chan (<-chan error)
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

	orderbook.syncerMu.RLock()
	defer orderbook.syncerMu.RUnlock()

	switch orderFragment.Depth {
	case 0:
		if orderbook.syncerCurrOrderFragments == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return nil
		case <-orderbook.done:
			return ErrServerIsNotRunning
		case orderbook.syncerCurrOrderFragments <- orderFragment:
			return nil
		}
	case 1:
		if orderbook.syncerPrevOrderFragments == nil {
			return nil
		}
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
		dispatch.Merge(done, notifications, orderbook.notificationMerger)
	}()
	go func() {
		defer close(errs)
		dispatch.Merge(done, errs, orderbook.errMerger)
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

	// Clone a new PointerStorer for the new syncer
	pointerStore, err := orderbook.pointerStore.Clone()
	if err != nil {

	}

	// Start the syncer for this epoch
	syncer := newSyncer(
		epoch,
		pointerStore,
		orderbook.orderStore,
		orderbook.orderFragmentStore,
		orderbook.contractBinder,
		orderbook.interval,
		orderbook.limit,
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
