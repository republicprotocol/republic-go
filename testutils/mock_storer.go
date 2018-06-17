package testutils

import (
	"sync"

	"github.com/republicprotocol/republic-go/orderbook"

	"github.com/pkg/errors"
	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/order"
)

// ErrOrderFragmentNotFound is return when attempting to load an order that
// cannot be found.
var ErrOrderFragmentNotFound = errors.New("order fragment not found")

// ErrOrderNotFound is return when attempting to load an order that cannot be
// found.
var ErrOrderNotFound = errors.New("order not found")

// ErrComputationNotFound is returned when the Storer cannot find a Computation
// associated with a ComputationID.
var ErrComputationNotFound = errors.New("computation not found")

// Storer is a mock implementation of the orderbook.Storer interface.
type Storer struct {
	mu             *sync.Mutex
	orderFragments map[order.ID]order.Fragment
	orders         map[order.ID]order.Order
	computations   map[ome.ComputationID]ome.Computation
	buyPointer     orderbook.SyncPointer
	sellPointer    orderbook.SyncPointer
}

// NewStorer creates a new mock Storer.
func NewStorer() *Storer {
	return &Storer{
		mu:             new(sync.Mutex),
		orderFragments: map[order.ID]order.Fragment{},
		orders:         map[order.ID]order.Order{},
		computations:   map[ome.ComputationID]ome.Computation{},
		buyPointer:     orderbook.SyncPointer(0),
		sellPointer:    orderbook.SyncPointer(0),
	}
}

// InsertOrderFragment implements orderbook.Storer.
func (storer *Storer) InsertOrderFragment(orderFragment order.Fragment) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.orderFragments[orderFragment.OrderID] = orderFragment
	return nil
}

// InsertOrder implements orderbook.Storer.
func (storer *Storer) InsertOrder(order order.Order) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.orders[order.ID] = order
	return nil
}

// InsertComputation implements the ome.Storer.
func (storer *Storer) InsertComputation(computation ome.Computation) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.computations[computation.ID] = computation
	return nil
}

// OrderFragment implements orderbook.Storer.
func (storer *Storer) OrderFragment(id order.ID) (order.Fragment, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	fragment, ok := storer.orderFragments[id]
	if !ok {
		return order.Fragment{}, ErrOrderFragmentNotFound
	}
	return fragment, nil
}

// Order implements orderbook.Storer.
func (storer *Storer) Order(id order.ID) (order.Order, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	ord, ok := storer.orders[id]
	if !ok {
		return order.Order{}, ErrOrderNotFound
	}
	return ord, nil
}

// Computation implements the ome.Storer.
func (storer *Storer) Computation(id ome.ComputationID) (ome.Computation, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	computation, ok := storer.computations[id]
	if !ok {
		return ome.Computation{}, ErrComputationNotFound
	}
	return computation, nil
}

// Computations implements the ome.Storer.
func (storer *Storer) Computations() (ome.Computations, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	computations := make([]ome.Computation, 0, len(storer.computations))
	for _, j := range storer.computations {
		computations = append(computations, j)
	}

	return computations, nil
}

// RemoveOrderFragment implements orderbook.Storer.
func (storer *Storer) RemoveOrderFragment(id order.ID) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	delete(storer.orderFragments, id)
	return nil
}

// RemoveOrder implements orderbook.Storer.
func (storer *Storer) RemoveOrder(id order.ID) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	delete(storer.orders, id)
	return nil
}

// RemoveComputation implements the ome.Storer.
func (storer *Storer) RemoveComputation(id ome.ComputationID) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	delete(storer.computations, id)
	return nil
}

// InsertBuyPointer implements orderbook.SyncStorer.
func (storer *Storer) InsertBuyPointer(ptr orderbook.SyncPointer) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.buyPointer = ptr
	return nil
}

// InsertSellPointer implements orderbook.SyncStorer.
func (storer *Storer) InsertSellPointer(ptr orderbook.SyncPointer) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.sellPointer = ptr
	return nil
}

// BuyPointer implements orderbook.SyncStorer.
func (storer *Storer) BuyPointer() (orderbook.SyncPointer, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	return storer.buyPointer, nil
}

// SellPointer implements orderbook.SyncStorer.
func (storer *Storer) SellPointer() (orderbook.SyncPointer, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	return storer.sellPointer, nil
}

// NumOrderFragments in the Storer.
func (storer *Storer) NumOrderFragments() int {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	return len(storer.orderFragments)
}
