package compute

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

type PendingResultFragment struct {
	ID                ResultFragmentID
	BuyOrderFragment  *OrderFragment
	SellOrderFragment *OrderFragment
}

func NewPendingResultFragment(left *OrderFragment, right *OrderFragment) (*PendingResultFragment, error) {
	if err := left.IsCompatible(right); err != nil {
		return nil, err
	}
	resultFragment := &PendingResultFragment{}
	if left.OrderBuySell != 0 {
		resultFragment.BuyOrderFragment = left
		resultFragment.SellOrderFragment = right
	} else {
		resultFragment.BuyOrderFragment = right
		resultFragment.SellOrderFragment = left
	}
	resultFragment.ID = ResultFragmentID(crypto.Keccak256(resultFragment.BuyOrderFragment.ID[:], resultFragment.SellOrderFragment.ID[:]))
	return resultFragment, nil
}

func (resultFragment *PendingResultFragment) Add(prime *big.Int) (*ResultFragment, error) {
	return resultFragment.BuyOrderFragment.Add(resultFragment.SellOrderFragment, prime)
}

func (resultFragment *PendingResultFragment) Sub(prime *big.Int) (*ResultFragment, error) {
	return resultFragment.BuyOrderFragment.Sub(resultFragment.SellOrderFragment, prime)
}

type PendingResultsMatrix struct {
	orderFragments []*OrderFragment

	pendingResultFragmentsGuard   *sync.Cond
	pendingResultFragmentsLeft    int
	pendingResultFragments        []*PendingResultFragment
	pendingResultFragmentsMarkers map[string]struct{}
}

func NewPendingResultsMatrix() *PendingResultsMatrix {
	return &PendingResultsMatrix{
		orderFragments: []*OrderFragment{},

		pendingResultFragmentsGuard:   sync.NewCond(new(sync.Mutex)),
		pendingResultFragmentsLeft:    0,
		pendingResultFragments:        []*PendingResultFragment{},
		pendingResultFragmentsMarkers: map[string]struct{}{},
	}
}

func (matrix *PendingResultsMatrix) AddOrderFragment(orderFragment *OrderFragment) {
	matrix.pendingResultFragmentsGuard.L.Lock()
	defer matrix.pendingResultFragmentsGuard.L.Unlock()

	for _, rhs := range matrix.orderFragments {
		if orderFragment.ID.Equals(rhs.ID) {
			return
		}
	}

	for _, rhs := range matrix.orderFragments {
		resultFragment, err := NewPendingResultFragment(orderFragment, rhs)
		if err != nil {
			continue
		}
		matrix.pendingResultFragments = append(matrix.pendingResultFragments, resultFragment)
		matrix.pendingResultFragmentsLeft++
	}

	matrix.orderFragments = append(matrix.orderFragments, orderFragment)
	if matrix.pendingResultFragmentsLeft > 0 {
		matrix.pendingResultFragmentsGuard.Signal()
	}
}

func (matrix *PendingResultsMatrix) WaitForPendingResultFragments(max int) []*PendingResultFragment {
	matrix.pendingResultFragmentsGuard.L.Lock()
	defer matrix.pendingResultFragmentsGuard.L.Unlock()
	for matrix.pendingResultFragmentsLeft == 0 {
		matrix.pendingResultFragmentsGuard.Wait()
	}

	pendingResultFragments := make([]*PendingResultFragment, 0, max)
	for _, com := range matrix.pendingResultFragments {
		if _, ok := matrix.pendingResultFragmentsMarkers[string(com.ID)]; !ok {
			matrix.pendingResultFragmentsMarkers[string(com.ID)] = struct{}{}
			pendingResultFragments = append(pendingResultFragments, com)
		}
	}

	matrix.pendingResultFragmentsLeft -= len(pendingResultFragments)
	if matrix.pendingResultFragmentsLeft > 0 {
		matrix.pendingResultFragmentsGuard.Signal()
	}
	return pendingResultFragments
}
