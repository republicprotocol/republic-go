package orderbook

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

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
	contractBinder     ContractBinder
	orderFragmentStore OrderFragmentStorer

	syncerMu       *sync.Mutex
	syncerCurrDone chan struct{}
	syncerPrevDone chan struct{}

	orderFragments  chan order.Fragment
	notificationChs chan (<-chan Notification)
	errChs          chan (<-chan error)
}

// NewOrderbook returns an Orderbok that uses a crypto.RsaKey to decrypt the
// order.EncryptedFragments that it receives, and stores them in a Storer.
func NewOrderbook(rsaKey crypto.RsaKey, contractBinder ContractBinder, orderFragmentStore OrderFragmentStorer) Orderbook {
	return &orderbook{
		doneMu: new(sync.RWMutex),
		done:   nil,

		rsaKey:             rsaKey,
		contractBinder:     contractBinder,
		orderFragmentStore: orderFragmentStore,

		syncerMu:       new(sync.Mutex),
		syncerCurrDone: nil,
		syncerPrevDone: nil,

		orderFragments:  make(chan order.Fragment),
		notificationChs: make(chan (<-chan Notification)),
		errChs:          make(chan (<-chan error)),
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

	if orderbook.syncerPrevDone == nil {
		close(orderbook.syncerPrevDone)
	}
	orderbook.syncerPrevDone = orderbook.syncerCurrDone
	orderbook.syncerCurrDone = make(chan struct{})

	minimumEpochInterval, err := orderbook.contractBinder.MinimumEpochInterval()
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get minimum epoch interval: %v", err))
		return
	}
	notifications, errs := newSyncer(
		orderbook.contractBinder,
		big.NewInt(0).SetUint64(uint64(epoch.BlockNumber)), // Block offset
		minimumEpochInterval,                               // Block limit
		0,                                                  // Order synchronising offset
		1024,                                               // Order synchronising limit
	).sync(orderbook.syncerCurrDone)

	select {
	case <-orderbook.done:
	case orderbook.notificationChs <- notifications:
	}

	select {
	case <-orderbook.done:
	case orderbook.errChs <- errs:
	}
}

func (orderbook *orderbook) mergeNotifications(done <-chan struct{}, notifications chan<- Notification) {
	for {
		select {
		case <-done:
			return
		case notificationCh, ok := <-orderbook.notificationChs:
			if !ok {
				return
			}
			go func() {
				for {
					select {
					case <-done:
						return
					case notification, ok := <-notificationCh:
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
		case errCh, ok := <-orderbook.errChs:
			if !ok {
				return
			}
			go func() {
				for {
					select {
					case <-done:
						return
					case err, ok := <-errCh:
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
