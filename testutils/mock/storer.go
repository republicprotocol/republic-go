package mock

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

// ErrOrderFragmentNotFound is return when attempting to load an order that
// cannot be found.
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

// ErrOrderNotFound is return when attempting to load an order that cannot be
// found.
var ErrOrderNotFound = errors.New("order not found")

// Storer is a mock implementation of the orderbook.Storer interface.
type Storer struct {
	mu             *sync.Mutex
	orderFragments map[order.ID]order.Fragment
	orders         map[order.ID]order.Order
}

// NewStorer creates a new mock Storer.
func NewStorer() orderbook.Storer {
	return &Storer{
		mu:             new(sync.Mutex),
		orderFragments: map[order.ID]order.Fragment{},
		orders:         map[order.ID]order.Order{},
	}
}

// InsertOrderFragment to the Storer.
func (storer *Storer) InsertOrderFragment(orderFragment order.Fragment) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.orderFragments[orderFragment.OrderID] = orderFragment
	return nil
}

// InsertOrder into the Storer.
func (storer *Storer) InsertOrder(order order.Order) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.orders[order.ID] = order
	return nil
}

// OrderFragment returns the order fragment of the given order id.
func (storer *Storer) OrderFragment(id order.ID) (order.Fragment, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	fragment, ok := storer.orderFragments[id]
	if !ok {
		return order.Fragment{}, ErrOrderFragmentNotFound
	}
	return fragment, nil
}

// Order returns the order of the given order id.
func (storer *Storer) Order(id order.ID) (order.Order, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	ord, ok := storer.orders[id]
	if !ok {
		return order.Order{}, ErrOrderNotFound
	}
	return ord, nil
}

// RemoveOrderFragment from the Storer.
func (storer *Storer) RemoveOrderFragment(id order.ID) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	delete(storer.orderFragments, id)
	return nil
}

// RemoveOrder from the Storer.
func (storer *Storer) RemoveOrder(id order.ID) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	delete(storer.orders, id)
	return nil
}
