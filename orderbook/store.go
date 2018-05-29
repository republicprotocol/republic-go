package orderbook

import (
	"errors"

	"github.com/republicprotocol/republic-go/order"
)

// ErrOrderFragmentNotFound is return when attempting to load an order that
// cannot be found.
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

// ErrOrderNotFound is return when attempting to load an order that cannot be
// found.
var ErrOrderNotFound = errors.New("order not found")

type Storer interface {
	InsertOrderFragment(orderFragment order.Fragment) error
	InsertOrder(order order.Order) error
	OrderFragment(id order.ID) (order.Fragment, error)
	Order(id order.ID) (order.Order, error)
	RemoveOrderFragment(id order.ID) error
	RemoveOrder(id order.ID) error
}
