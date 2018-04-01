package smpc

import (
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

type Computation struct {
	BuyOrderFragment  *order.Fragment
	SellOrderFragment *order.Fragment
}

type ComputationMatrix struct {
	mu                 *sync.Mutex
	buyOrderFragments  map[string]order.Fragment
	sellOrderFragments map[string]order.Fragment
	computations       []Computation
}

func NewComputationMatrix() ComputationMatrix {
	return ComputationMatrix{
		mu:                 new(sync.Mutex),
		buyOrderFragments:  map[string]order.Fragment{},
		sellOrderFragments: map[string]order.Fragment{},
		computations:       []Computation{},
	}
}

func (matrix *ComputationMatrix) InsertBuyOrder(buyOrderFragment order.Fragment) {
	matrix.mu.Lock()
	defer matrix.mu.Unlock()

	buyOrderID := string(buyOrderFragment.OrderID)
	if _, ok := matrix.buyOrderFragments[buyOrderID]; ok {
		return
	}
	matrix.buyOrderFragments[buyOrderID] = buyOrderFragment

	for _, sellOrderFragment := range matrix.sellOrderFragments {
		matrix.computations = append(matrix.computations, Computation{
			BuyOrderFragment:  &buyOrderFragment,
			SellOrderFragment: &sellOrderFragment,
		})
	}
}

func (matrix *ComputationMatrix) InsertSellOrder(sellOrderFragment order.Fragment) {
	matrix.mu.Lock()
	defer matrix.mu.Unlock()

	sellOrderID := string(sellOrderFragment.OrderID)
	if _, ok := matrix.sellOrderFragments[sellOrderID]; ok {
		return
	}
	matrix.sellOrderFragments[sellOrderID] = sellOrderFragment

	for _, buyOrderFragment := range matrix.buyOrderFragments {
		matrix.computations = append(matrix.computations, Computation{
			BuyOrderFragment:  &buyOrderFragment,
			SellOrderFragment: &sellOrderFragment,
		})
	}
}

func (matrix *ComputationMatrix) Computations(computations []Computation) int {
	matrix.mu.Lock()
	defer matrix.mu.Unlock()

	m := len(matrix.computations)
	n := len(computations)
	if m == 0 || n == 0 {
		return 0
	}

	i := 0
	for ; i < m && i < n; i++ {
		computations[i] = matrix.computations[m-i-1]
	}
	if i >= m {
		matrix.computations = matrix.computations[0:0]
	} else {
		matrix.computations = matrix.computations[:m-i]
	}

	return i
}

func (matrix *ComputationMatrix) RemoveOrder(orderID order.ID) {
	matrix.mu.Lock()
	defer matrix.mu.Unlock()

	delete(matrix.buyOrderFragments, string(orderID))
	delete(matrix.sellOrderFragments, string(orderID))

	n := len(matrix.computations)
	d := 0
	for i := 0; i < n; i++ {
		if matrix.computations[i].BuyOrderFragment.OrderID.Equal(orderID) || matrix.computations[i].SellOrderFragment.OrderID.Equal(orderID) {
			matrix.computations[i] = matrix.computations[n-d-1]
			d++
		}
	}
	matrix.computations = matrix.computations[:n-d]
}
