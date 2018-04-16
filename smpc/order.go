package smpc

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/order"
)

// ProcessOrderFragments by reading order fragments from an input channel, and
// writing OrderTuples to an output channel.
func ProcessOrderFragments(ctx context.Context, orderFragmentChIn <-chan order.Fragment, sharedOrderTable *SharedOrderTable, bufferLimit int) (<-chan OrderTuple, <-chan error) {
	orderTupleCh := make(chan OrderTuple, bufferLimit)
	errCh := make(chan error)

	log.Println("ProcessOrderFragments")

	go func() {
		defer close(orderTupleCh)
		defer close(errCh)

		buffer := make([]OrderTuple, bufferLimit)
		tick := time.NewTicker(time.Millisecond)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-tick.C:
				for i, n := 0, sharedOrderTable.OrderTuples(buffer[:]); i < n; i++ {
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case orderTupleCh <- buffer[i]:
					}
				}
			}
		}
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case orderFragment, ok := <-orderFragmentChIn:
				log.Println("Inserting order fragment")
				if !ok {
					return
				}
				if orderFragment.OrderParity == order.ParityBuy {
					sharedOrderTable.InsertBuyOrder(orderFragment)
				} else {
					sharedOrderTable.InsertSellOrder(orderFragment)
				}
			}
		}
	}()

	return orderTupleCh, errCh
}

// OrderTuple involving a buy order, and a sell order. Orders are stores as
// pointers to fragments. The pointers must not be used to modify the
// underlying fragment.
type OrderTuple struct {
	BuyOrderFragment  *order.Fragment
	SellOrderFragment *order.Fragment
}

// OrderState of an order in the SharedOrderTable determines whether or not it
// is active.
type OrderState int

const (
	// OrderStateOn allows an order to be considered in the construction of
	// OrderTuples.
	OrderStateOn = 1

	// OrderStateOff disallows an order to be considered in the construction of
	// OrderTuples.
	OrderStateOff = 2
)

// A SharedOrderTable stores order fragments and generates OrderTuples that
// will be used to perform order matching computations. It is safe for
// concurrent use.
type SharedOrderTable struct {
	mu *sync.Mutex

	// State maps store a fragment as either being seen, or as not being seen.
	buyOrderStates  map[string]OrderState
	sellOrderStates map[string]OrderState

	buyOrderFragments  []order.Fragment
	sellOrderFragments []order.Fragment
	orderTuples        []OrderTuple
}

// NewSharedOrderTable returns an empty SharedOrderTable.
func NewSharedOrderTable() SharedOrderTable {
	return SharedOrderTable{
		mu:                 new(sync.Mutex),
		buyOrderStates:     map[string]OrderState{},
		sellOrderStates:    map[string]OrderState{},
		buyOrderFragments:  []order.Fragment{},
		sellOrderFragments: []order.Fragment{},
		orderTuples:        []OrderTuple{},
	}
}

// InsertBuyOrder into the SharedOrderTable. This will generate a list of
// OrderTuples between the inserted buy order fragment and all sell order
// fragments currently stored in the SharedOrderTable.
func (table *SharedOrderTable) InsertBuyOrder(buyOrderFragment order.Fragment) {
	table.mu.Lock()
	defer table.mu.Unlock()

	buyOrderID := string(buyOrderFragment.OrderID)
	if _, ok := table.buyOrderStates[buyOrderID]; ok {
		return
	}
	table.buyOrderStates[buyOrderID] = OrderStateOn
	table.buyOrderFragments = append(table.buyOrderFragments, buyOrderFragment)

	for i := range table.sellOrderFragments {
		if table.sellOrderStates[string(table.sellOrderFragments[i].OrderID)] != OrderStateOn {
			continue
		}
		table.orderTuples = append(table.orderTuples, OrderTuple{
			BuyOrderFragment:  &table.buyOrderFragments[len(table.buyOrderFragments)-1],
			SellOrderFragment: &table.sellOrderFragments[i],
		})
	}
}

// InsertSellOrder into the SharedOrderTable. This will generate a list of
// OrderTuples between the inserted sell order fragment and all buy order
// fragments currently stored in the SharedOrderTable.
func (table *SharedOrderTable) InsertSellOrder(sellOrderFragment order.Fragment) {
	table.mu.Lock()
	defer table.mu.Unlock()

	sellOrderID := string(sellOrderFragment.OrderID)
	if _, ok := table.sellOrderStates[sellOrderID]; ok {
		return
	}
	table.sellOrderStates[sellOrderID] = OrderStateOn
	table.sellOrderFragments = append(table.sellOrderFragments, sellOrderFragment)

	for i := range table.buyOrderFragments {
		if table.buyOrderStates[string(table.buyOrderFragments[i].OrderID)] != OrderStateOn {
			continue
		}
		table.orderTuples = append(table.orderTuples, OrderTuple{
			BuyOrderFragment:  &table.buyOrderFragments[i],
			SellOrderFragment: &table.sellOrderFragments[len(table.sellOrderFragments)-1],
		})
	}
}

// SetBuyOrderState for an order in the SharedOrderTable. Setting the state to
// OrderStateOff will remove all OrderTuples involving the order.
func (table *SharedOrderTable) SetBuyOrderState(orderID order.ID, state OrderState) {
	table.mu.Lock()
	defer table.mu.Unlock()

	table.buyOrderStates[string(orderID)] = state

	// If the order has been turned off, then remove it from the existing set
	// of OrderTuples.
	if state != OrderStateOn {
		n := len(table.orderTuples)
		d := 0
		for i := 0; i < n; i++ {
			if table.orderTuples[i].BuyOrderFragment.OrderID.Equal(orderID) {
				table.orderTuples[i] = table.orderTuples[n-d-1]
				d++
			}
		}
		table.orderTuples = table.orderTuples[:n-d]
	}
}

// SetSellOrderState for an order in the SharedOrderTable. Setting the state to
// OrderStateOff will remove all OrderTuples involving the order.
func (table *SharedOrderTable) SetSellOrderState(orderID order.ID, state OrderState) {
	table.mu.Lock()
	defer table.mu.Unlock()

	table.sellOrderStates[string(orderID)] = state

	// If the order has been turned off, then remove it from the existing set
	// of OrderTuples.
	if state != OrderStateOn {
		n := len(table.orderTuples)
		d := 0
		for i := 0; i < n; i++ {
			if table.orderTuples[i].SellOrderFragment.OrderID.Equal(orderID) {
				table.orderTuples[i] = table.orderTuples[n-d-1]
				d++
			}
		}
		table.orderTuples = table.orderTuples[:n-d]
	}
}

// OrderTuples writes a list of OrderTuples to the buffer and returns the
// number of OrderTuples written.
func (table *SharedOrderTable) OrderTuples(buffer []OrderTuple) int {
	table.mu.Lock()
	defer table.mu.Unlock()

	m := len(table.orderTuples)
	n := len(buffer)
	if m == 0 || n == 0 {
		return 0
	}

	i := 0
	for ; i < m && i < n; i++ {
		buffer[i] = table.orderTuples[m-i-1]
	}
	if i >= m {
		table.orderTuples = table.orderTuples[0:0]
	} else {
		table.orderTuples = table.orderTuples[:m-i]
	}

	return i
}
