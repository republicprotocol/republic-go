package testutils

import (
	"context"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
)

type Orderbook struct {
}

// OpenOrder implements the orderbook.Orderbook interface.
func (mock *Orderbook) OpenOrder(context.Context, order.EncryptedFragment) error {
	return nil
}

// Sync implements the orderbook.Orderbook interface.
func (mock *Orderbook) Sync(done <-chan struct{}) (<-chan orderbook.Notification, <-chan error) {
	notifications := make(chan orderbook.Notification)
	errs := make(chan error)

	go func() {
		for {
			select {
			case <-done:
			}
		}
	}()

	return notifications, errs
}

// OnChangeEpoch implements the orderbook.Orderbook interface.
func (mock *Orderbook) OnChangeEpoch(epoch registry.Epoch) {
}
