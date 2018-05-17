package orderbook

import (
	"context"

	"github.com/pkg/errors"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

var ErrOrderFragmentNotFound = errors.New("order fragment not found")

type Orderbooker interface {
	OpenOrder(context.Context, identity.MultiAddress, order.Fragment) error
	Sync() ([]OrderUpdate, error)
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
	syncer Syncer
}

func NewOrderbook(client OrderbookClient, storer Storer, syncer Syncer) orderbook{
	return orderbook{
		client: client,
		storer: storer,
		syncer : syncer,
	}
}

func (book *orderbook) OpenOrder(ctx context.Context, multiAddr identity.MultiAddress, orderFragment order.Fragment) error {
	return book.client.OpenOrder(ctx, multiAddr, orderFragment)
}

func (book *orderbook) 	Sync() ([]OrderUpdate, error) {
	return book.syncer.Sync()
}

func (book *orderbook) OrderFragment(id order.ID) (order.Fragment, error){
	return book.storer.OrderFragment(id)
}

func (book *orderbook) Order(id order.ID) (order.Order, error){
	return book.storer.Order(id)
}
