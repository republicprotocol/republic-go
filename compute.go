package compute

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

type ComputationID []byte

type Computation struct {
	ID                ComputationID
	BuyOrderFragment  *OrderFragment
	SellOrderFragment *OrderFragment
}

func NewComputation(left *OrderFragment, right *OrderFragment) (*Computation, error) {
	if err := left.IsCompatible(right); err != nil {
		return nil, err
	}
	computation := &Computation{}
	if left.OrderBuySell != 0 {
		computation.BuyOrderFragment = left
		computation.SellOrderFragment = right
	} else {
		computation.BuyOrderFragment = right
		computation.SellOrderFragment = left
	}
	computation.ID = ComputationID(crypto.Keccak256(computation.BuyOrderFragment.ID[:], computation.SellOrderFragment.ID[:]))
	return computation, nil
}

func (computation *Computation) Add(prime *big.Int) (*ResultFragment, error) {
	return computation.BuyOrderFragment.Add(computation.SellOrderFragment, prime)
}

func (computation *Computation) Sub(prime *big.Int) (*ResultFragment, error) {
	return computation.BuyOrderFragment.Sub(computation.SellOrderFragment, prime)
}

type ComputationMatrix struct {
	orderFragments []*OrderFragment

	computationsMu       *sync.Mutex
	computationsLeftCond *sync.Cond
	computations         []*Computation
	computationsLeft     int
	computationsMarker   map[string]struct{}

	resultsMu       *sync.Mutex
	results         map[string]*Result
	resultFragments map[string][]*ResultFragment
}

func NewComputationMatrix() *ComputationMatrix {
	return &ComputationMatrix{
		orderFragments: []*OrderFragment{},

		computationsMu:       new(sync.Mutex),
		computationsLeftCond: sync.NewCond(new(sync.Mutex)),
		computations:         []*Computation{},
		computationsLeft:     0,
		computationsMarker:   map[string]struct{}{},

		resultsMu:       new(sync.Mutex),
		results:         map[string]*Result{},
		resultFragments: map[string][]*ResultFragment{},
	}
}

func (matrix *ComputationMatrix) AddOrderFragment(orderFragment *OrderFragment) {
	matrix.computationsMu.Lock()
	defer matrix.computationsMu.Unlock()

	for _, rhs := range matrix.orderFragments {
		if orderFragment.ID.Equals(rhs.ID) {
			return
		}
	}

	for _, other := range matrix.orderFragments {
		if orderFragment.OrderID.Equals(other.OrderID) {
			continue
		}
		computation, err := NewComputation(orderFragment, other)
		if err != nil {
			continue
		}
		matrix.computations = append(matrix.computations, computation)
		matrix.computationsLeft++
	}

	matrix.orderFragments = append(matrix.orderFragments, orderFragment)
	if matrix.computationsLeft > 0 {
		matrix.computationsLeftCond.Signal()
	}
}

func (matrix *ComputationMatrix) WaitForComputations(max int) []*Computation {
	matrix.computationsMu.Lock()
	defer matrix.computationsMu.Unlock()

	matrix.computationsLeftCond.L.Lock()
	defer matrix.computationsLeftCond.L.Unlock()
	for matrix.computationsLeft == 0 {
		matrix.computationsLeftCond.Wait()
	}

	computations := make([]*Computation, 0, max)
	for _, computation := range matrix.computations {
		if _, ok := matrix.computationsMarker[string(computation.ID)]; !ok {
			matrix.computationsMarker[string(computation.ID)] = struct{}{}
			computations = append(computations, computation)
		}
	}
	matrix.computationsLeft -= len(computations)
	return computations
}

func (matrix *ComputationMatrix) AddResultFragment(resultFragment *ResultFragment) (*Result, error) {
	matrix.resultsMu.Lock()
	defer matrix.resultsMu.Unlock()

	resultID := ResultID(crypto.Keccak256(resultFragment.BuyOrderID[:], resultFragment.SellOrderID[:]))
	matrix.resultFragments[string(resultID)] = append(matrix.resultFragments[string(resultID)], resultFragment)
}
