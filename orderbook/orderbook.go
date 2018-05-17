package orderbook

import (
	"context"

	"github.com/pkg/errors"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

type Client interface {
	OpenOrder(context.Context, identity.MultiAddress, order.EncryptedFragment) error
}

type Server interface {
	OpenOrder(context.Context, order.Fragment) error
}

type Orderbooker interface {

	// Sync orders and order states from the Ren Ledger to this local
	// Orderbooker. Returns a list of changes that were made to this local
	// Orderbooker during the synchronization.
	Sync() ([]OrderUpdate, error)

	// OpenOrder by sending an order.EncryptedFragment to an
	// identity.MultiAddress. The order.EncryptedFragment should not be stored
	// in the local Orderbooker, it will be stored by the Orderbook hosted at
	// the identity.MultiAddress.
	OpenOrder(context.Context, identity.MultiAddress, order.EncryptedFragment) error

	// OrderFragment stored in this local Orderbooker. These are received from
	// other Orderbookers calling Orderbooker.OpenOrder to send an
	// order.EncryptedFragment to this local Orderbooker.
	OrderFragment(order.ID) (order.Fragment, error)

	// Order that has been reconstructed and stored in this local Orderbooker.
	// This only happens for orders that have been matched and confirmed.
	Order(order.ID) (order.Order, error)
}

type orderbook struct {
	client Client
	storer Storer
	syncer Syncer
}

func NewOrderbook(client Client, storer Storer, syncer Syncer) Orderbooker {
	return &orderbook{
		client: client,
		storer: storer,
		syncer: syncer,
	}
}

func (book *orderbook) Sync() ([]OrderUpdate, error) {
	return book.syncer.Sync()
}

func (book *orderbook) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.EncryptedFragment) error {
	return book.client.OpenOrder(ctx, multiAddr, orderFragment)
}

func (book *orderbook) OrderFragment(id order.ID) (order.Fragment, error) {
	return book.storer.OrderFragment(id)
}

func (book *orderbook) Order(id order.ID) (order.Order, error) {
	return book.storer.Order(id)
}
