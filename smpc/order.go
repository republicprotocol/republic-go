package smpc

import (
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/order"
)

// OrderFragmentsToOrderTuples reads order.Fragments from an input channel and
// stores them in a SharedOrderTable. From this table, it produces OrderTuples
// for all combinations of compatible order.Fragments.
func OrderFragmentsToOrderTuples(done <-chan struct{}, orderFragments <-chan order.Fragment, sharedOrderTable *SharedOrderTable, bufferLimit int) <-chan OrderTuple {
	orderTuples := make(chan OrderTuple, bufferLimit)

	// Insert order.Fragments into the SharedOrderTable
	go func() {
		for {
			select {
			case <-done:
				return
			case orderFragment, ok := <-orderFragments:
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

	// Periodically read OrderTuples from the SharedOrderTable into a buffer
	// and write them to the output channel
	go func() {
		defer close(orderTuples)

		buffer := make([]OrderTuple, bufferLimit)
		tick := time.NewTicker(time.Second)
		defer tick.Stop()

		for {
			select {
			case <-done:
				return
			case <-tick.C:
				for i, n := 0, sharedOrderTable.OrderTuples(buffer[:]); i < n; i++ {
					select {
					case <-done:
						return
					case orderTuples <- buffer[i]:
					}
				}
			}
		}
	}()

	return orderTuples
}

// OrderTuple involving a buy order, and a sell order. Orders are stores as
// pointers to fragments. The pointers must not be used to modify the
// underlying fragment.
type OrderTuple struct {
	BuyOrderFragment  *order.Fragment
	SellOrderFragment *order.Fragment
}

// A SharedOrderTable stores order fragments and generates OrderTuples that
// will be used to perform order matching computations. It is safe for
// concurrent use.
type SharedOrderTable struct {
	mu                 *sync.Mutex
	buyOrderFragments  []order.Fragment
	sellOrderFragments []order.Fragment
	orderTuples        []OrderTuple
}

// NewSharedOrderTable returns an empty SharedOrderTable.
func NewSharedOrderTable() SharedOrderTable {
	return SharedOrderTable{
		mu:                 new(sync.Mutex),
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

	table.buyOrderFragments = append(table.buyOrderFragments, buyOrderFragment)
	for i := range table.sellOrderFragments {
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

	table.sellOrderFragments = append(table.sellOrderFragments, sellOrderFragment)
	for i := range table.buyOrderFragments {
		table.orderTuples = append(table.orderTuples, OrderTuple{
			BuyOrderFragment:  &table.buyOrderFragments[i],
			SellOrderFragment: &table.sellOrderFragments[len(table.sellOrderFragments)-1],
		})
	}
}

// RemoveBuyOrder from the SharedOrderTable. This will also remove all
// OrderTuples involving the order.
func (table *SharedOrderTable) RemoveBuyOrder(orderID order.ID) {
	table.mu.Lock()
	defer table.mu.Unlock()

	// Remove buy order fragment
	n := len(table.buyOrderFragments)
	for i := 0; i < n; i++ {
		if table.buyOrderFragments[i].OrderID.Equal(orderID) {
			table.buyOrderFragments[i] = table.buyOrderFragments[n-1]
			table.buyOrderFragments = table.buyOrderFragments[:n-1]
			break
		}
	}

	// Remove related order tuples
	n = len(table.orderTuples)
	d := 0
	for i := 0; i < n; i++ {
		if table.orderTuples[i].BuyOrderFragment.OrderID.Equal(orderID) {
			table.orderTuples[i] = table.orderTuples[n-d-1]
			d++
		}
	}
	table.orderTuples = table.orderTuples[:n-d]
}

// RemoveSellOrder from the SharedOrderTable. This will also remove all
// OrderTuples involving the order.
func (table *SharedOrderTable) RemoveSellOrder(orderID order.ID) {
	table.mu.Lock()
	defer table.mu.Unlock()

	// Remove sell order fragment
	n := len(table.sellOrderFragments)
	for i := 0; i < n; i++ {
		if table.sellOrderFragments[i].OrderID.Equal(orderID) {
			table.sellOrderFragments[i] = table.sellOrderFragments[n-1]
			table.sellOrderFragments = table.sellOrderFragments[:n-1]
			break
		}
	}

	// Remove related order tuples
	n = len(table.orderTuples)
	d := 0
	for i := 0; i < n; i++ {
		if table.orderTuples[i].SellOrderFragment.OrderID.Equal(orderID) {
			table.orderTuples[i] = table.orderTuples[n-d-1]
			d++
		}
	}
	table.orderTuples = table.orderTuples[:n-d]
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
