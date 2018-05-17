package orderbook

import (
	"context"

	"github.com/pkg/errors"
	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

var ErrOrderFragmentNotFound = errors.New("order fragment not found")

type Orderbooker interface {
	OpenOrder(context.Context, identity.MultiAddress, order.EncryptedFragment) error
	Sync() ([]order.ID, []order.ID, error)
	OrderFragment(order.ID) (order.Fragment, error)
	Order(order.ID) (order.Order, error)
}

type OrderbookClient interface {
	OpenOrder(context.Context, identity.MultiAddress, order.Fragment) error
}

type OrderbookServer interface {
	OpenOrder(context.Context, order.Fragment) error
}

type orderbook struct {
	client OrderbookClient
	storer Storer
}

func NewOrderbook(client OrderbookClient, ledger cal.RenLedger, storer Storer) orderbook {
	return orderbook{
		client: client,
		storer: storer,
	}
}

func (book *orderbook) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.Fragment) error {
	return book.client.OpenOrder(ctx, multiAddr, orderFragment)
}

func (book *orderbook) Sync() ([]order.ID, []order.ID, error) {
	panic("unimplemented")
}
