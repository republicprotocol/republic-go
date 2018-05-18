package orderbook

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
)

var ErrNotFoundInStore = errors.New("not found in store")

type Storer interface {
	InsertOrderFragment(orderFragment order.Fragment) error
	InsertOrder(order order.Order) error
	OrderFragment(id order.ID) (order.Fragment, error)
	Order(id order.ID) (order.Order, error)
	RemoveOrderFragment(id order.ID) error
	RemoveOrder(id order.ID) error
}
