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

type HiddenOrderBook struct {
	do.GuardedObject

	orderFragments []*OrderFragment
	blockSize      int

	pendingComputations              []*Computation
	pendingComputationsReadyForBlock *do.Guard
}

// NewHiddenOrderBook returns a new HiddenOrderBook with no OrderFragments, or
// Computations.
func NewHiddenOrderBook(blockSize int) *HiddenOrderBook {
	orderBook := new(HiddenOrderBook)
	orderBook.GuardedObject = do.NewGuardedObject()
	orderBook.orderFragments = make([]*OrderFragment, 0)
	orderBook.blockSize = blockSize
	orderBook.pendingComputations = make([]*Computation, 0)
	orderBook.pendingComputationsReadyForBlock = orderBook.Guard(func() bool { return len(orderBook.pendingComputations) >= blockSize })
	return orderBook
}

func (orderBook *HiddenOrderBook) AddOrderFragment(orderFragment *OrderFragment) {
	orderBook.Enter(nil)
	defer orderBook.Exit()
	orderBook.addOrderFragment(orderFragment)
}

func (orderBook *HiddenOrderBook) addOrderFragment(orderFragment *OrderFragment) {
	// Check that the OrderFragment has not been added.
	for _, rhs := range orderBook.orderFragments {
		if orderFragment.ID.Equals(rhs.ID) {
			return
		}
	}

	// For all other OrderFragment in the Computer, create a new Computation
	// that needs to be processed.
	for _, other := range orderBook.orderFragments {
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
		orderBook.pendingComputations = append(orderBook.pendingComputations, computation)
	}
	orderBook.orderFragments = append(orderBook.orderFragments, orderFragment)
}

func (orderBook *HiddenOrderBook) PendingComputationsForBlock() []*Computation {
	orderBook.Enter(orderBook.pendingComputationsReadyForBlock)
	defer orderBook.Exit()
	return orderBook.pendingComputationsForBlock()
}

func (orderBook *HiddenOrderBook) pendingComputationsForBlock() []*Computation {
	blockSize := orderBook.blockSize
	if blockSize > len(orderBook.pendingComputations) {
		blockSize = len(orderBook.pendingComputations)
	}
	pendingComputations := make([]*Computation, 0, blockSize)
	pendingComputations = append(pendingComputations, orderBook.pendingComputations[len(orderBook.pendingComputations)-blockSize:]...)
	orderBook.pendingComputations = orderBook.pendingComputations[0 : len(orderBook.pendingComputations)-blockSize]
	return pendingComputations
}
