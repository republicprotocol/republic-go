package testutils

import (
	"context"
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
)

// An EmptyOrderbook implements the orderbook.Orderbook interface but it never
// produces notifications.
type EmptyOrderbook struct {
}

// NewEmptyOrderbook returns a new EmptyOrderbook.
func NewEmptyOrderbook() *EmptyOrderbook {
	return &EmptyOrderbook{}
}

// OpenOrder implements the orderbook.Orderbook interface.
func (mock *EmptyOrderbook) OpenOrder(ctx context.Context, encryptedOrderFragment order.EncryptedFragment) error {
	return nil
}

// Sync implements the orderbook.Orderbook interface.
func (mock *EmptyOrderbook) Sync(done <-chan struct{}) (<-chan orderbook.Notification, <-chan error) {
	notifications := make(chan orderbook.Notification)
	errs := make(chan error)

	go func() {
		defer close(notifications)
		defer close(errs)

		<-done
	}()

	return notifications, errs
}

// OnChangeEpoch implements the orderbook.Orderbook interface.
func (mock *EmptyOrderbook) OnChangeEpoch(epoch registry.Epoch) {
}

// A RandOrderbook implements the orderbook.Orderbook interface and produces a
// notification for every opened order.Fragment and randomly produces
// cancelations and confirmations.
type RandOrderbook struct {
	rsaKey              crypto.RsaKey
	orderFragments      []order.Fragment
	orderFragmentsQueue []order.Fragment
	traders             []string
}

// NewRandOrderbook returns a new RandOrderbook that is initially empty and
// uses a crypto.RsaKey to decrypt order.EncryptedFragments.
func NewRandOrderbook(rsaKey crypto.RsaKey) *RandOrderbook {
	return &RandOrderbook{
		rsaKey:              rsaKey,
		orderFragments:      []order.Fragment{},
		orderFragmentsQueue: []order.Fragment{},
		traders:             []string{"trader1", "trader2", "trader3", "trader4"},
	}
}

// OpenOrder implements the orderbook.Orderbook interface.
func (mock *RandOrderbook) OpenOrder(ctx context.Context, encryptedOrderFragment order.EncryptedFragment) error {
	orderFragment, err := encryptedOrderFragment.Decrypt(mock.rsaKey.PrivateKey)
	if err != nil {
		return err
	}
	mock.orderFragmentsQueue = append(mock.orderFragmentsQueue, orderFragment)
	return nil
}

// Sync implements the orderbook.Orderbook interface.
func (mock *RandOrderbook) Sync(done <-chan struct{}) (<-chan orderbook.Notification, <-chan error) {
	notifications := make(chan orderbook.Notification)
	errs := make(chan error)

	go func() {
		defer close(notifications)
		defer close(errs)

		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()

		blockNumber := uint64(0)
		for {
			select {
			case <-done:
			case <-ticker.C:
				// Simulate the block number increasing
				blockNumber++

				// Generate an open notification for all pending order.Fagments
				for i := range mock.orderFragmentsQueue {
					notification := orderbook.NotificationOpenOrder{
						OrderID:       mock.orderFragmentsQueue[i].OrderID,
						OrderFragment: mock.orderFragmentsQueue[i],
						Trader:        mock.traders[rand.Intn(len(mock.traders))],
						BlockNumber:   blockNumber,
					}
					select {
					case <-done:
					case notifications <- notification:
					}
				}
				mock.orderFragments = append(mock.orderFragments, mock.orderFragmentsQueue...)
				mock.orderFragmentsQueue = []order.Fragment{}

				// Randomly remove an order.Fragment
				n := rand.Intn(len(mock.orderFragments))
				orderFragment := mock.orderFragments[n]
				mock.orderFragments[n] = mock.orderFragments[len(mock.orderFragments)-1]
				mock.orderFragments = mock.orderFragments[:len(mock.orderFragments)-1]

				// Randomly generate a closure notification for it
				var notification orderbook.Notification
				r := rand.Intn(100)
				if r < 50 {
					notification = orderbook.NotificationConfirmOrder{
						OrderID: orderFragment.OrderID,
					}
				} else {
					notification = orderbook.NotificationCancelOrder{
						OrderID: orderFragment.OrderID,
					}
				}
				select {
				case <-done:
				case notifications <- notification:
				}
			}
		}
	}()

	return notifications, errs
}

// OnChangeEpoch implements the orderbook.Orderbook interface.
func (mock *RandOrderbook) OnChangeEpoch(epoch registry.Epoch) {
}
