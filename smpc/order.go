package smpc

import (
	"context"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/order"
)

// ProcessOrderFragments by reading order fragments from an input channel, and
// writing OrderFragmentComputations to an output channel.
func ProcessOrderFragments(ctx context.Context, orderFragmentChIn chan order.Fragment) (chan OrderFragmentComputation, chan error) {
	orderFragmentTable := NewOrderFragmentTable()
	consumerErrCh := consumeOrderFragments(ctx, orderFragmentChIn, &orderFragmentTable)
	orderFragmentComputationCh, producerErrCh := produceOrderFragmentComputations(ctx, &orderFragmentTable)

	errCh := make(chan error)
	go func() {
		defer close(errCh)

		for {
			var err error
			var ok bool
			select {
			case err, ok = <-consumerErrCh:
			case err, ok = <-producerErrCh:
			}
			if !ok {
				return
			}
			errCh <- err
			if err == context.Canceled {
				return
			}
		}
	}()

	return orderFragmentComputationCh, errCh
}

// consumeOrderFragments from an input channel and store
// OrderFragmentComputations in an OrderFragmentTable.
func consumeOrderFragments(ctx context.Context, orderFagmentChIn chan order.Fragment, orderFragmentTable *OrderFragmentTable) chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case orderFragment, ok := <-orderFagmentChIn:
				if !ok {
					return
				}
				if orderFragment.OrderParity == order.ParityBuy {
					orderFragmentTable.InsertBuyOrder(orderFragment)
				} else {
					orderFragmentTable.InsertSellOrder(orderFragment)
				}
			}
		}
	}()

	return errCh
}

// produceOrderFragmentComputations by periodically reading
// OrderFragmentComputations from an OrderFragmentTable and writing them to an
// output channel
func produceOrderFragmentComputations(ctx context.Context, orderFragmentTable *OrderFragmentTable) (chan OrderFragmentComputation, chan error) {
	orderFragmentComputationCh := make(chan OrderFragmentComputation)
	errCh := make(chan error)

	go func() {
		defer close(orderFragmentComputationCh)
		defer close(errCh)

		buffer := make([]OrderFragmentComputation, 128)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-ticker.C:
				for i, n := 0, orderFragmentTable.Computations(buffer[:]); i < n; i++ {
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					case orderFragmentComputationCh <- buffer[i]:
					}
				}
			}
		}
	}()

	return orderFragmentComputationCh, errCh
}

// OrderFragmentComputation involving a buy order fragemnt, and a sell order
// fragment.
type OrderFragmentComputation struct {
	BuyOrderFragment  *order.Fragment
	SellOrderFragment *order.Fragment
}

// An OrderFragmentTable stores order fragments and generates pairwise
// OrderFragmentComputations that should be considered for computation. It is
// safe for concurrent use.
type OrderFragmentTable struct {
	mu                        *sync.Mutex
	buyOrderFragments         map[string]order.Fragment
	sellOrderFragments        map[string]order.Fragment
	orderFragmentComputations []OrderFragmentComputation
}

// NewOrderFragmentTable returns an empty OrderFragmentTable.
func NewOrderFragmentTable() OrderFragmentTable {
	return OrderFragmentTable{
		mu:                        new(sync.Mutex),
		buyOrderFragments:         map[string]order.Fragment{},
		sellOrderFragments:        map[string]order.Fragment{},
		orderFragmentComputations: []OrderFragmentComputation{},
	}
}

// InsertBuyOrder into the OrderFragmentTable. This will generate a list of
// OrderFragmentComputations between the inserted buy order fragment and all
// sell order fragments currently stored in the OrderFragmentTable.
func (table *OrderFragmentTable) InsertBuyOrder(buyOrderFragment order.Fragment) {
	table.mu.Lock()
	defer table.mu.Unlock()

	buyOrderID := string(buyOrderFragment.OrderID)
	if _, ok := table.buyOrderFragments[buyOrderID]; ok {
		return
	}
	table.buyOrderFragments[buyOrderID] = buyOrderFragment

	for _, sellOrderFragment := range table.sellOrderFragments {
		table.orderFragmentComputations = append(table.orderFragmentComputations, OrderFragmentComputation{
			BuyOrderFragment:  &buyOrderFragment,
			SellOrderFragment: &sellOrderFragment,
		})
	}
}

// InsertSellOrder into the OrderFragmentTable. This will generate a list of
// OrderFragmentComputations between the inserted sell order fragment and all
// buy order fragments currently stored in the OrderFragmentTable.
func (table *OrderFragmentTable) InsertSellOrder(sellOrderFragment order.Fragment) {
	table.mu.Lock()
	defer table.mu.Unlock()

	sellOrderID := string(sellOrderFragment.OrderID)
	if _, ok := table.sellOrderFragments[sellOrderID]; ok {
		return
	}
	table.sellOrderFragments[sellOrderID] = sellOrderFragment

	for _, buyOrderFragment := range table.buyOrderFragments {
		table.orderFragmentComputations = append(table.orderFragmentComputations, OrderFragmentComputation{
			BuyOrderFragment:  &buyOrderFragment,
			SellOrderFragment: &sellOrderFragment,
		})
	}
}

// Computations writes a list of OrderFragmentComputations to the buffer and
// returns the number of OrderFragmentComputations written.
func (table *OrderFragmentTable) Computations(buffer []OrderFragmentComputation) int {
	table.mu.Lock()
	defer table.mu.Unlock()

	m := len(table.orderFragmentComputations)
	n := len(buffer)
	if m == 0 || n == 0 {
		return 0
	}

	i := 0
	for ; i < m && i < n; i++ {
		buffer[i] = table.orderFragmentComputations[m-i-1]
	}
	if i >= m {
		table.orderFragmentComputations = table.orderFragmentComputations[0:0]
	} else {
		table.orderFragmentComputations = table.orderFragmentComputations[:m-i]
	}

	return i
}

// RemoveBuyOrder from the OrderFragmentTable. All OrderFragmentComputations
// that involve this order will also be removed.
func (table *OrderFragmentTable) RemoveBuyOrder(orderID order.ID) {
	table.mu.Lock()
	defer table.mu.Unlock()

	delete(table.buyOrderFragments, string(orderID))

	n := len(table.orderFragmentComputations)
	d := 0
	for i := 0; i < n; i++ {
		if table.orderFragmentComputations[i].BuyOrderFragment.OrderID.Equal(orderID) {
			table.orderFragmentComputations[i] = table.orderFragmentComputations[n-d-1]
			d++
		}
	}
	table.orderFragmentComputations = table.orderFragmentComputations[:n-d]
}

// RemoveSellOrder from the OrderFragmentTable. All OrderFragmentComputations
// that involve this order will also be removed.
func (table *OrderFragmentTable) RemoveSellOrder(orderID order.ID) {
	table.mu.Lock()
	defer table.mu.Unlock()

	delete(table.sellOrderFragments, string(orderID))

	n := len(table.orderFragmentComputations)
	d := 0
	for i := 0; i < n; i++ {
		if table.orderFragmentComputations[i].SellOrderFragment.OrderID.Equal(orderID) {
			table.orderFragmentComputations[i] = table.orderFragmentComputations[n-d-1]
			d++
		}
	}
	table.orderFragmentComputations = table.orderFragmentComputations[:n-d]
}
