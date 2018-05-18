package orderbook

import (
	"context"
	"errors"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

type Client interface {

	// OpenOrder by sending an order.EncryptedFragment to an
	// identity.MultiAddress. The order.EncryptedFragment will be stored by the
	// Server hosted at the identity.MultiAddress.
	OpenOrder(context.Context, identity.MultiAddress, order.EncryptedFragment) error
}

type Server interface {
	OpenOrder(context.Context, order.EncryptedFragment) error
}

type Listener interface {
	OnConfirmOrderMatch(order.Order, order.Order)
}

type Orderbooker interface {

	// Sync orders and order states from the Ren Ledger to this local
	// Orderbooker. Returns a list of changes that were made to this local
	// Orderbooker during the synchronization.
	Sync() ([]OrderUpdate, error)

	// OrderFragment stored in this local Orderbooker. These are received from
	// other Orderbookers calling Orderbooker.OpenOrder to send an
	// order.EncryptedFragment to this local Orderbooker.
	OrderFragment(order.ID) (order.Fragment, error)

	// Order that has been reconstructed and stored in this local Orderbooker.
	// This only happens for orders that have been matched and confirmed.
	Order(order.ID) (order.Order, error)

	// Confirm an order match.
	ConfirmOrderMatch(order.ID, order.ID) error
}

type orderbook struct {
	storer Storer
	syncer Syncer
}

func NewOrderbook(storer Storer, syncer Syncer) Orderbooker {
	return &orderbook{
		storer: storer,
		syncer: syncer,
	}
}

func (book *orderbook) OpenOrder(ctx context.Context, orderFragment order.EncryptedFragment) error {
	panic("unimplemented")
}

func (book *orderbook) Sync() ([]OrderUpdate, error) {
	return book.syncer.Sync()
}

func (book *orderbook) OrderFragment(id order.ID) (order.Fragment, error) {
	return book.storer.OrderFragment(id)
}

func (book *orderbook) Order(id order.ID) (order.Order, error) {
	return book.storer.Order(id)
}

func (book *orderbook) ConfirmOrderMatch(buy order.ID, sell order.ID) error {
	panic("uimplemented")
}
