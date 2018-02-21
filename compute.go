package compute

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-do"
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

func (computation *Computation) Sub(prime *big.Int) (*ResultFragment, error) {
	return computation.BuyOrderFragment.Sub(computation.SellOrderFragment, prime)
}

type ComputationShardID []byte

type ComputationShard struct {
	ID           ComputationShardID
	Computations []*Computation
}

func NewComputationShard(computations []*Computation) ComputationShard {
	computationIDs := make([]byte, 0, len(computations)*32)
	for _, computation := range computations {
		computationIDs = append(computationIDs, []byte(computation.ID)...)
	}
	return ComputationShard{
		ID:           ComputationShardID(crypto.Keccak256(computationIDs)),
		Computations: computations,
	}
}

func (shard ComputationShard) Compute(prime *big.Int) []*ResultFragment {
	resultFragments := make([]*ResultFragment, len(shard.Computations))
	for i := range resultFragments {
		// FIXME: We are processing computations in bulk with the expectation
		// that some of them will fail (hopefully 2/3rds of participiants will
		// succeed). Errors are dropped here.
		resultFragments[i], _ = shard.Computations[i].Sub(prime)
	}
	return resultFragments
}

type ComputationBid int64

const (
	ComputationBidYes = 1
	ComputationBidNo  = 2
)

type ComputationShardBid struct {
	ID   ComputationShardID
	Bids map[string]ComputationBid
}

type HiddenOrderBook struct {
	do.GuardedObject

	orderFragments []*OrderFragment
	shardSize      int

	pendingComputations              []*Computation
	pendingComputationsReadyForShard *do.Guard
}

// NewHiddenOrderBook returns a new HiddenOrderBook with no OrderFragments, or
// Computations.
func NewHiddenOrderBook(shardSize int) *HiddenOrderBook {
	orderBook := new(HiddenOrderBook)
	orderBook.GuardedObject = do.NewGuardedObject()
	orderBook.orderFragments = make([]*OrderFragment, 0)
	orderBook.shardSize = shardSize
	orderBook.pendingComputations = make([]*Computation, 0)
	orderBook.pendingComputationsReadyForShard = orderBook.Guard(func() bool { return len(orderBook.pendingComputations) >= shardSize })
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

func (orderBook *HiddenOrderBook) WaitForComputationShard() ComputationShard {
	orderBook.Enter(orderBook.pendingComputationsReadyForShard)
	defer orderBook.Exit()
	return orderBook.preemptComputationShard()
}

func (orderBook *HiddenOrderBook) PreemptComputationShard() ComputationShard {
	orderBook.Enter(nil)
	defer orderBook.Exit()
	return orderBook.preemptComputationShard()
}

func (orderBook *HiddenOrderBook) preemptComputationShard() ComputationShard {
	shardSize := orderBook.shardSize
	if shardSize > len(orderBook.pendingComputations) {
		shardSize = len(orderBook.pendingComputations)
	}
	pendingComputations := make([]*Computation, 0, shardSize)
	pendingComputations = append(pendingComputations, orderBook.pendingComputations[len(orderBook.pendingComputations)-shardSize:]...)
	orderBook.pendingComputations = orderBook.pendingComputations[0 : len(orderBook.pendingComputations)-shardSize]
	return NewComputationShard(pendingComputations)
}
