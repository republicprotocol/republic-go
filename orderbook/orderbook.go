package orderbook

import (
	"context"
	"errors"
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
	rsaKey             crypto.RsaKey
	orderFragmentStore OrderFragmentStorer
	orderFragments     chan order.Fragment

	syncerMu        *sync.Mutex
	syncerCurrEpoch *syncer
	syncerPrevEpoch *syncer

	doneMu *sync.RWMutex
	done   chan struct{}
}

// NewOrderbook returns an Orderbok that uses a crypto.RsaKey to decrypt the
// order.EncryptedFragments that it receives, and stores them in a Storer.
func NewOrderbook(rsaKey crypto.RsaKey, orderFragmentStore OrderFragmentStorer) Orderbook {
	return &orderbook{
		rsaKey:             rsaKey,
		orderFragmentStore: orderFragmentStore,
		orderFragments:     make(chan order.Fragment),

		syncerMu:        new(sync.Mutex),
		syncerCurrEpoch: nil,
		syncerPrevEpoch: nil,

		doneMu: new(sync.RWMutex),
		done:   nil,
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

	select {
	case <-ctx.Done():
		return nil
	case <-orderbook.done:
		return ErrServerIsNotRunning
	case orderbook.orderFragments <- orderFragment:
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

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				orderbook.sync(done, notifications, errs)
			}
		}
	}()

	return notifications, errs
}

func (orderbook *orderbook) sync(done <-chan struct{}, notifications chan<- Notification, errs chan<- error) {
	orderbook.syncerMu.Lock()
	defer orderbook.syncerMu.Unlock()

	if orderbook.syncerPrevEpoch != nil {
		orderbook.syncerPrevEpoch.sync(done, notifications, errs)
	}
	if orderbook.syncerCurrEpoch != nil {
		orderbook.syncerCurrEpoch.sync(done, notifications, errs)
	}
}

// OnChangeEpoch implements the Orderbook interface.
func (orderbook *orderbook) OnChangeEpoch(epoch registry.Epoch) {
	orderbook.syncerMu.Lock()
	defer orderbook.syncerMu.Unlock()

	// FIXME: Set the syncers
}
