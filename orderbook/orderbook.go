package orderbook

import (
	"context"
	"errors"
	"sync"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
)

// ErrOrderbookIsAlreadySyncing is returned when Orderbook.Sync is called while
// the Orderbook is already syncing. A call to Orderbook.Sync should only ever
// happen once.
var ErrOrderbookIsAlreadySyncing = errors.New("orderbook is already syncing")

// Client for invoking RPCs on a remote Server.
type Client interface {

	// OpenOrder by sending an order.EncryptedFragment to an
	// identity.MultiAddress. The order.EncryptedFragment will be stored by the
	// Server.
	OpenOrder(context.Context, identity.MultiAddress, order.EncryptedFragment) error
}

// Server for opening order.EncryptedFragments. This RPC should only be called
// after the respective order.Order has been opened on the Ethereum blockchain
// otherwise it will be ignored by the Server.
type Server interface {
	OpenOrder(context.Context, order.EncryptedFragment) error
}

// An Orderbook is responsible for receiving orders. It should read order.Order
// statuses from a Syncer, and order.EncryptedFragemnts from a Server. It
// outputs Notifications that should be consumed by the user to understand the
// changing state of the Orderbook. During boot, the Orderbook should output
// initialising Notifications for the current state of the Orderbook.
type Orderbook interface {

	// Sync the Orderbook with the Ethereum blockchain until the done channel
	// is closed.
	Sync(dont <-chan struct{}) (<-chan Notification, <-chan error)
}

type orderbook struct {
	doneMu *sync.Mutex
	done   <-chan struct{}

	rsaKey crypto.RsaKey

	syncer              Syncer
	filter              Filter
	filterNotifications chan Notification
}

// NewOrderbook returns an Orderbok that uses a crypto.RsaKey to decrypt
// sensitive data. It implements the Server interface to receive
// order.EncryptedFragments. It uses a Syncer to create Notifications about
// updates to the Orderbook, and a Filter to filter these Notifications.
func NewOrderbook(rsaKey crypto.RsaKey, syncer Syncer, filter Filter) Orderbook {
	return &orderbook{
		doneMu: new(sync.Mutex),
		done:   nil,

		rsaKey: rsaKey,

		syncer:              syncer,
		filter:              filter,
		filterNotifications: make(chan Notification),
	}
}

// Sync implements the Syncer interface.
func (orderbook *orderbook) Sync(done <-chan struct{}) (<-chan Notification, <-chan error) {
	notifications := make(chan Notification)
	errs := make(chan error)

	// Set the done channel
	orderbook.doneMu.Lock()
	defer orderbook.doneMu.Unlock()
	if orderbook.done != nil {
		errs <- ErrOrderbookIsAlreadySyncing
		close(notifications)
		close(errs)
		return notifications, errs
	}
	orderbook.done = done

	var wg sync.WaitGroup
	wg.Add(2)

	syncNotifications, syncErrs := orderbook.syncer.Sync(orderbook.done)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-orderbook.done:
			case notification, ok := <-syncNotifications:
				if !ok {
					return
				}
				select {
				case <-orderbook.done:
				case orderbook.filterNotifications <- notification:
				}
			case err, ok := <-syncErrs:
				if !ok {
					return
				}
				select {
				case <-orderbook.done:
				case errs <- err:
				}
			}
		}
	}()

	filterNotifications, filterErrs := orderbook.filter.Filter(orderbook.done, orderbook.filterNotifications)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-orderbook.done:
			case notification, ok := <-filterNotifications:
				if !ok {
					return
				}
				select {
				case <-orderbook.done:
				case notifications <- notification:
				}
			case err, ok := <-filterErrs:
				if !ok {
					return
				}
				select {
				case <-orderbook.done:
				case errs <- err:
				}
			}
		}
	}()

	go func() {
		defer close(notifications)
		defer close(errs)
		wg.Wait()
	}()

	return notifications, errs
}

// OpenOrder implements the Server interface.
func (orderbook *orderbook) OpenOrder(ctx context.Context, encryptedOrderFragment order.EncryptedFragment) error {
	orderFragment, err := encryptedOrderFragment.Decrypt(*orderbook.rsaKey.PrivateKey)
	if err != nil {
		return err
	}
	if orderFragment.OrderParity == order.ParityBuy {
		logger.BuyOrderReceived(logger.LevelDebugLow, orderFragment.OrderID.String(), orderFragment.ID.String())
	} else {
		logger.SellOrderReceived(logger.LevelDebugLow, orderFragment.OrderID.String(), orderFragment.ID.String())
	}

	notification := NotificationSyncOrderFragment{
		OrderFragment: orderFragment,
	}
	select {
	case <-orderbook.done:
	case orderbook.filterNotifications <- notification:
	}
	return nil
}
