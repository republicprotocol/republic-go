package smpc

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/order"
)

// ProcessOrderFragments by reading order fragments from an input channel, and
// writing OrderFragmentComputations to an output channel.
func ProcessOrderFragments(ctx context.Context, orderFragmentChIn <-chan order.Fragment, bufferLimit int) (<-chan OrderFragmentComputation, <-chan error) {
	// orderFragmentCh := make(chan order.Fragment)
	errCh := make(chan error)

	orderFragmentTable := NewOrderFragmentTable()
	consumerErrCh := consumeOrderFragments(ctx, orderFragmentChIn, &orderFragmentTable)

	producerErrCtx, producerCtxCancel := context.WithCancel(context.Background())
	orderFragmentComputationCh, producerErrCh := produceOrderFragmentComputations(producerErrCtx, &orderFragmentTable, bufferLimit)

	go func() {
		defer producerCtxCancel()
		defer close(errCh)
		for {
			select {
			case err, ok := <-consumerErrCh:
				if !ok {
					continue
				}
				if err == context.Canceled {
					continue
				}
				errCh <- err
			case err, ok := <-producerErrCh:
				if !ok {
					return
				}
				errCh <- err
				if err == context.Canceled {
					return
				}
			}
		}
	}()

	return orderFragmentComputationCh, errCh
}

// consumeOrderFragments from an input channel and store
// OrderFragmentComputations in an OrderFragmentTable.
func consumeOrderFragments(ctx context.Context, orderFagmentChIn <-chan order.Fragment, orderFragmentTable *OrderFragmentTable) <-chan error {
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
func produceOrderFragmentComputations(ctx context.Context, orderFragmentTable *OrderFragmentTable, bufferLimit int) (<-chan OrderFragmentComputation, <-chan error) {
	orderFragmentComputationCh := make(chan OrderFragmentComputation, bufferLimit)
	errCh := make(chan error)

	go func() {
		defer close(orderFragmentComputationCh)
		defer close(errCh)

		buffer := make([]OrderFragmentComputation, bufferLimit)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("shutdown")
				errCh <- ctx.Err()
				return
			case <-ticker.C:
				log.Println("tick")
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
	buyOrderFragments         []order.Fragment
	sellOrderFragments        []order.Fragment
	buyOrderFragmentStatus    map[string]bool
	sellOrderFragmentStatus   map[string]bool
	orderFragmentComputations []OrderFragmentComputation
}

// NewOrderFragmentTable returns an empty OrderFragmentTable.
func NewOrderFragmentTable() OrderFragmentTable {
	return OrderFragmentTable{
		mu:                        new(sync.Mutex),
		buyOrderFragments:         []order.Fragment{},
		sellOrderFragments:        []order.Fragment{},
		buyOrderFragmentStatus:    map[string]bool{},
		sellOrderFragmentStatus:   map[string]bool{},
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
	if _, ok := table.buyOrderFragmentStatus[buyOrderID]; ok {
		return
	}
	table.buyOrderFragmentStatus[buyOrderID] = false
	table.buyOrderFragments = append(table.buyOrderFragments, buyOrderFragment)

	for i := range table.sellOrderFragments {
		if table.sellOrderFragmentStatus[string(table.sellOrderFragments[i].OrderID)] {
			continue
		}
		table.orderFragmentComputations = append(table.orderFragmentComputations, OrderFragmentComputation{
			BuyOrderFragment:  &table.buyOrderFragments[len(table.buyOrderFragments)-1],
			SellOrderFragment: &table.sellOrderFragments[i],
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
	if _, ok := table.sellOrderFragmentStatus[sellOrderID]; ok {
		return
	}
	table.sellOrderFragmentStatus[sellOrderID] = false
	table.sellOrderFragments = append(table.sellOrderFragments, sellOrderFragment)

	for i := range table.buyOrderFragments {
		if table.buyOrderFragmentStatus[string(table.buyOrderFragments[i].OrderID)] {
			continue
		}
		table.orderFragmentComputations = append(table.orderFragmentComputations, OrderFragmentComputation{
			BuyOrderFragment:  &table.buyOrderFragments[i],
			SellOrderFragment: &table.sellOrderFragments[len(table.sellOrderFragments)-1],
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

	table.buyOrderFragmentStatus[string(orderID)] = true

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

	table.sellOrderFragmentStatus[string(orderID)] = true

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
