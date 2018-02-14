package compute

import (
	"math/big"

	"github.com/republicprotocol/go-do"

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
	if left.OrderParity == OrderParityBuy {
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

type Computer struct {
	do.GuardedObject

	orderFragments []*OrderFragment

	computations           []*Computation
	computationsIsNotEmpty *do.Guard
}

// NewComputer returns a new Computer with no OrderFragments, or Computations.
func NewComputer() *Computer {
	computer := new(Computer)
	computer.GuardedObject = do.NewGuardedObject()
	computer.orderFragments = make([]*OrderFragment, 0)
	computer.computations = make([]*Computation, 0)
	computer.computationsIsNotEmpty = computer.Guard(func() bool { return len(computer.computations) > 0 })
	return computer
}

func (computer *Computer) AddOrderFragment(orderFragment *OrderFragment) {
	computer.Enter(nil)
	defer computer.Exit()
	computer.addOrderFragment(orderFragment)
}

func (computer *Computer) addOrderFragment(orderFragment *OrderFragment) {
	// Check that the OrderFragment has not been added.
	for _, rhs := range computer.orderFragments {
		if orderFragment.ID.Equals(rhs.ID) {
			return
		}
	}

	// For all other OrderFragment in the Computer, create a new Computation
	// that needs to be processed.
	for _, other := range computer.orderFragments {
		if orderFragment.OrderID.Equals(other.OrderID) {
			continue
		}
		if err := orderFragment.IsCompatible(other); err != nil {
			continue
		}
		computation, err := NewComputation(orderFragment, other)
		if err != nil {
			continue
		}
		computer.computations = append(computer.computations, computation)
	}
	computer.orderFragments = append(computer.orderFragments, orderFragment)
}

func (computer *Computer) Compute() []*Result {
	computer.Enter(computer.computationsIsNotEmpty)
	defer computer.Exit()
	return computer.compute()
}

func (computer *Computer) compute() []*Result {
	panic("unimplemented")
}
