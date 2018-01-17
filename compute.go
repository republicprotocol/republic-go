package compute

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

type Computation struct {
	ID                OrderFragmentID
	BuyOrderFragment  *OrderFragment
	SellOrderFragment *OrderFragment
}

func NewComputation(left *OrderFragment, right *OrderFragment) (*Computation, error) {
	if err := left.IsCompatible(right); err != nil {
		return nil, err
	}
	com := &Computation{}
	if left.OrderBuySell != 0 {
		com.BuyOrderFragment = left
		com.SellOrderFragment = right
	} else {
		com.BuyOrderFragment = right
		com.SellOrderFragment = left
	}
	com.ID = OrderFragmentID(crypto.Keccak256(com.BuyOrderFragment.ID[:], com.SellOrderFragment.ID[:]))
	return com, nil
}

func (com *Computation) Add(prime *big.Int) (*ComputedOrderFragment, error) {
	return com.BuyOrderFragment.Add(com.SellOrderFragment, prime)
}

func (com *Computation) Sub(prime *big.Int) (*ComputedOrderFragment, error) {
	return com.BuyOrderFragment.Sub(com.SellOrderFragment, prime)
}

type ComputationMatrix struct {
	orderFragments []*OrderFragment

	computationsGuard  *sync.Cond
	computationsLeft   int
	computations       []*Computation
	computationMarkers map[string]struct{}
}

func NewComputationMatrix() *ComputationMatrix {
	return &ComputationMatrix{
		orderFragments: []*OrderFragment{},

		computationsGuard:  sync.NewCond(new(sync.Mutex)),
		computationsLeft:   0,
		computations:       []*Computation{},
		computationMarkers: map[string]struct{}{},
	}
}

func (matrix *ComputationMatrix) AddOrderFragment(orderFragment *OrderFragment) {
	matrix.computationsGuard.L.Lock()
	defer matrix.computationsGuard.L.Unlock()

	for _, rhs := range matrix.orderFragments {
		if orderFragment.ID.Equals(rhs.ID) {
			return
		}
	}

	for _, rhs := range matrix.orderFragments {
		com, err := NewComputation(orderFragment, rhs)
		if err != nil {
			continue
		}
		matrix.computations = append(matrix.computations, com)
		matrix.computationsLeft++
	}

	matrix.orderFragments = append(matrix.orderFragments, orderFragment)
	if matrix.computationsLeft > 0 {
		matrix.computationsGuard.Signal()
	}
}

func (matrix *ComputationMatrix) WaitForComputations(max int) []*Computation {
	matrix.computationsGuard.L.Lock()
	defer matrix.computationsGuard.L.Unlock()
	for matrix.computationsLeft == 0 {
		matrix.computationsGuard.Wait()
	}

	computations := make([]*Computation, 0, max)
	for _, com := range matrix.computations {
		if _, ok := matrix.computationMarkers[string(com.ID)]; !ok {
			matrix.computationMarkers[string(com.ID)] = struct{}{}
			computations = append(computations, com)
		}
	}

	matrix.computationsLeft -= len(computations)
	if matrix.computationsLeft > 0 {
		matrix.computationsGuard.Signal()
	}
	return computations
}
