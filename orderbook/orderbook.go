package orderbook

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/registry"
)

// Client for invoking the Server.OpenOrder ROC on a remote Server.
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
	pointerStore       PointerStorer
	orderStore         OrderStorer
	orderFragmentStore OrderFragmentStorer
	contractBinder     ContractBinder
	interval           time.Duration

	aggMu   *sync.RWMutex
	aggCurr Aggregator
	aggPrev Aggregator

	syncer        Syncer
	notifications chan Notification
	errs          chan error
}

// NewOrderbook returns an Orderbok that uses a crypto.RsaKey to decrypt the
// order.EncryptedFragments that it receives, and stores them in a Storer.
func NewOrderbook(rsaKey crypto.RsaKey, pointerStore PointerStorer, orderStore OrderStorer, orderFragmentStore OrderFragmentStorer, contractBinder ContractBinder, interval time.Duration, limit int) Orderbook {
	return &orderbook{
		rsaKey:             rsaKey,
		pointerStore:       pointerStore,
		orderStore:         orderStore,
		orderFragmentStore: orderFragmentStore,
		contractBinder:     contractBinder,
		interval:           interval,

		aggMu:   new(sync.RWMutex),
		aggCurr: nil,
		aggPrev: nil,

		syncer:        NewSyncer(pointerStore, orderStore, contractBinder, limit),
		notifications: make(chan Notification),
		errs:          make(chan error),
	}
}

// OpenOrder implements the Server interface.
func (orderbook *orderbook) OpenOrder(ctx context.Context, encryptedOrderFragment order.EncryptedFragment) error {
	orderFragment, err := encryptedOrderFragment.Decrypt(orderbook.rsaKey.PrivateKey)
	if err != nil {
		return err
	}
	if encryptedOrderFragment.OrderParity == order.ParityBuy {
		logger.BuyOrderReceived(logger.LevelDebugLow, encryptedOrderFragment.OrderID.String(), encryptedOrderFragment.ID.String())
	} else {
		logger.SellOrderReceived(logger.LevelDebugLow, encryptedOrderFragment.OrderID.String(), encryptedOrderFragment.ID.String())
	}
	return orderbook.routeOrderFragment(ctx.Done(), orderFragment)
}

// Sync implements the Orderbook interface.
func (orderbook *orderbook) Sync(done <-chan struct{}) (<-chan Notification, <-chan error) {
	notifications := make(chan Notification)
	errs := make(chan error)

	go func() {
		defer close(notifications)
		defer close(errs)

		dispatch.CoBegin(
			func() {
				dispatch.Forward(done, orderbook.notifications, notifications)
			},
			func() {
				dispatch.Forward(done, orderbook.errs, errs)
			},
			func() {
				orderbook.sync(done)
			})
	}()

	return notifications, errs
}

// OnChangeEpoch implements the Orderbook interface.
func (orderbook *orderbook) OnChangeEpoch(epoch registry.Epoch) {
	orderbook.aggMu.Lock()
	defer orderbook.aggMu.Unlock()

	orderbook.aggPrev = orderbook.aggCurr
	orderbook.aggCurr = NewAggregator(epoch, orderbook.orderStore, orderbook.orderFragmentStore)
}

func (orderbook *orderbook) sync(done <-chan struct{}) {
	ticker := time.NewTicker(orderbook.interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			notifications, err := orderbook.syncer.Sync()
			if err != nil {
				select {
				case <-done:
					return
				case orderbook.errs <- err:
				}
			}
			if notifications != nil {
				for _, notification := range notifications {
					if err := orderbook.routeNotification(done, notification); err != nil {
						select {
						case <-done:
							return
						case orderbook.errs <- err:
						}
					}
				}
			}
		}
	}
}

func (orderbook *orderbook) routeNotification(done <-chan struct{}, notification Notification) error {
	switch n := notification.(type) {
	case NotificationOpenOrder:
		return orderbook.routeOrder(done, n.OrderID, order.Open, n.Trader)
	default:
		select {
		case <-done:
		case orderbook.notifications <- notification:
		}
	}
	return nil
}

func (orderbook *orderbook) routeOrder(done <-chan struct{}, orderID order.ID, orderStatus order.Status, trader string) error {

	ns, err := func() (ns [2]Notification, err error) {
		orderbook.aggMu.RLock()
		defer orderbook.aggMu.RUnlock()

		if orderbook.aggCurr == nil {
			return
		}
		ns[0], err = orderbook.aggCurr.InsertOrder(orderID, orderStatus, trader)
		if orderbook.aggPrev == nil {
			return
		}
		ns[1], err = orderbook.aggPrev.InsertOrder(orderID, orderStatus, trader)
		return
	}()

	for _, n := range ns {
		if n == nil {
			continue
		}
		select {
		case <-done:
		case orderbook.notifications <- n:
		}
	}
	return err
}

// TODO: Using the block number of the order, the orderbook should infer
// which epoch the order fragment is destined for. If the epoch is unknown
// then the orderbook should sleep here (this is safe given that this
// function is generally called in a background goroutine) and try again.
// Failing a second time should see the order fragment dropped. This helps
// with robust acceptance of order fragments at the turn of an epoch where
// the Darknode and the trader might briefly have different ideas about the
// "current" epoch.
func (orderbook *orderbook) routeOrderFragment(done <-chan struct{}, orderFragment order.Fragment) error {

	n, err := func() (Notification, error) {
		orderbook.aggMu.RLock()
		defer orderbook.aggMu.RUnlock()

		switch orderFragment.EpochDepth {
		case 0:
			if orderbook.aggCurr == nil {
				return nil, nil
			}
			return orderbook.aggCurr.InsertOrderFragment(orderFragment)
		case 1:
			if orderbook.aggPrev == nil {
				return nil, nil
			}
			return orderbook.aggPrev.InsertOrderFragment(orderFragment)
		default:
			log.Printf("[error] (sync) unexpected depth = %v", orderFragment.EpochDepth)
			return nil, nil
		}
	}()

	if n != nil {
		select {
		case <-done:
		case orderbook.notifications <- n:
		}
	}

	return err
}
