package orderbook

import (
	"context"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

var ErrOrderFragmentNotFound

type Orderbook interface {
	OpenOrder(context.Context, identity.MultiAddress, order.Fragment) error
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
}

func (book *orderbook) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.Fragment) error {
	return book.client.OpenOrder(ctx, multiAddr, orderFragment)
}
