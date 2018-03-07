package compute

import (
	"math/big"

	do "github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/order"
)

type DeltaBuilder struct {
	do.GuardedObject

	k                      int64
	prime                  *big.Int
	deltas                 map[string]*Delta
	deltaFragments         map[string]*DeltaFragment
	deltasToDeltaFragments map[string][]*DeltaFragment
}

func NewDeltaBuilder(k int64, prime *big.Int) *DeltaBuilder {
	return &DeltaBuilder{
		GuardedObject:          do.NewGuardedObject(),
		k:                      k,
		prime:                  prime,
		deltas:                 map[string]*Delta{},
		deltaFragments:         map[string]*DeltaFragment{},
		deltasToDeltaFragments: map[string][]*DeltaFragment{},
	}
}

func (builder *DeltaBuilder) InsertDeltaFragment(deltaFragment *DeltaFragment) *Delta {
	builder.Enter(nil)
	defer builder.Exit()
	return builder.insertDeltaFragment(deltaFragment)
}

func (builder *DeltaBuilder) insertDeltaFragment(deltaFragment *DeltaFragment) *Delta {
	// Is the delta already built, or are we adding a delta fragment that we
	// have already seen
	if builder.hasDelta(deltaFragment.DeltaID) {
		return nil // Only return new deltas
	}
	if builder.hasDeltaFragment(deltaFragment.ID) {
		return nil // Only return new deltas
	}

	// Add the delta fragment to the builder and attach it to the appropriate
	// delta
	builder.deltaFragments[string(deltaFragment.ID)] = deltaFragment
	if deltaFragments, ok := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]; ok {
		deltaFragments = append(deltaFragments, deltaFragment)
		builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] = deltaFragments
	} else {
		builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)] = []*DeltaFragment{deltaFragment}
	}

	// Build the delta if possible and return it
	deltaFragments := builder.deltasToDeltaFragments[string(deltaFragment.DeltaID)]
	if int64(len(deltaFragments)) >= builder.k {
		delta := NewDelta(deltaFragments, builder.prime)
		if delta == nil {
			return nil
		}
		builder.deltas[string(delta.ID)] = delta
		return delta
	}

	return nil
}

func (builder *DeltaBuilder) HasDelta(deltaID DeltaID) bool {
	builder.EnterReadOnly(nil)
	defer builder.Exit()
	return builder.hasDelta(deltaID)
}

func (builder *DeltaBuilder) hasDelta(deltaID DeltaID) bool {
	_, ok := builder.deltas[string(deltaID)]
	return ok
}

func (builder *DeltaBuilder) HasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	builder.EnterReadOnly(nil)
	defer builder.Exit()
	return builder.hasDeltaFragment(deltaFragmentID)
}

func (builder *DeltaBuilder) hasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	_, ok := builder.deltaFragments[string(deltaFragmentID)]
	return ok
}

type DeltaFragmentMatrix struct {
	do.GuardedObject

	prime                  *big.Int
	buyOrderFragments      map[string]*order.Fragment
	sellOrderFragments     map[string]*order.Fragment
	buySellDeltaFragments  map[string]map[string]*DeltaFragment
	completeOrderFragments map[string]*order.Fragment
}

func NewDeltaFragmentMatrix(prime *big.Int) *DeltaFragmentMatrix {
	return &DeltaFragmentMatrix{
		GuardedObject:          do.NewGuardedObject(),
		prime:                  prime,
		buyOrderFragments:      map[string]*order.Fragment{},
		sellOrderFragments:     map[string]*order.Fragment{},
		buySellDeltaFragments:  map[string]map[string]*DeltaFragment{},
		completeOrderFragments: map[string]*order.Fragment{},
	}
}

func (matrix *DeltaFragmentMatrix) InsertOrderFragment(orderFragment *order.Fragment) ([]*DeltaFragment, error) {
	matrix.Enter(nil)
	defer matrix.Exit()
	if orderFragment.OrderParity == order.ParityBuy {
		return matrix.insertBuyOrderFragment(orderFragment)
	}
	return matrix.insertSellOrderFragment(orderFragment)
}

func (matrix *DeltaFragmentMatrix) insertBuyOrderFragment(buyOrderFragment *order.Fragment) ([]*DeltaFragment, error) {
	if _, ok := matrix.buyOrderFragments[string(buyOrderFragment.ID)]; ok {
		return []*DeltaFragment{}, nil
	}
	if _, ok := matrix.completeOrderFragments[string(buyOrderFragment.ID)]; ok {
		return []*DeltaFragment{}, nil
	}

	deltaFragments := make([]*DeltaFragment, 0, len(matrix.sellOrderFragments))
	deltaFragmentsMap := map[string]*DeltaFragment{}
	for i := range matrix.sellOrderFragments {
		deltaFragment := NewDeltaFragment(buyOrderFragment, matrix.sellOrderFragments[i], matrix.prime)
		if deltaFragment == nil {
			continue
		}
		deltaFragments = append(deltaFragments, deltaFragment)
		deltaFragmentsMap[string(matrix.sellOrderFragments[i].ID)] = deltaFragment
	}

	matrix.buyOrderFragments[string(buyOrderFragment.ID)] = buyOrderFragment
	matrix.buySellDeltaFragments[string(buyOrderFragment.ID)] = deltaFragmentsMap
	return deltaFragments, nil
}

func (matrix *DeltaFragmentMatrix) insertSellOrderFragment(sellOrderFragment *order.Fragment) ([]*DeltaFragment, error) {
	if _, ok := matrix.sellOrderFragments[string(sellOrderFragment.ID)]; ok {
		return []*DeltaFragment{}, nil
	}
	if _, ok := matrix.completeOrderFragments[string(sellOrderFragment.ID)]; ok {
		return []*DeltaFragment{}, nil
	}

	deltaFragments := make([]*DeltaFragment, 0, len(matrix.buyOrderFragments))
	for i := range matrix.buyOrderFragments {
		deltaFragment := NewDeltaFragment(matrix.buyOrderFragments[i], sellOrderFragment, matrix.prime)
		if deltaFragment == nil {
			continue
		}
		if _, ok := matrix.buySellDeltaFragments[string(matrix.buyOrderFragments[i].ID)]; ok {
			deltaFragments = append(deltaFragments, deltaFragment)
			matrix.buySellDeltaFragments[string(matrix.buyOrderFragments[i].ID)][string(sellOrderFragment.ID)] = deltaFragment
		}
	}

	matrix.sellOrderFragments[string(sellOrderFragment.ID)] = sellOrderFragment
	return deltaFragments, nil
}

func (matrix *DeltaFragmentMatrix) RemoveOrderFragment(orderFragment *order.Fragment) error {
	matrix.Enter(nil)
	defer matrix.Exit()
	if orderFragment.OrderParity == order.ParityBuy {
		return matrix.removeBuyOrderFragment(orderFragment)
	}
	return matrix.removeSellOrderFragment(orderFragment)
}

func (matrix *DeltaFragmentMatrix) removeBuyOrderFragment(buyOrderFragment *order.Fragment) error {
	if _, ok := matrix.buyOrderFragments[string(buyOrderFragment.ID)]; !ok {
		return nil
	}

	delete(matrix.buyOrderFragments, string(buyOrderFragment.ID))
	delete(matrix.buySellDeltaFragments, string(buyOrderFragment.ID))

	matrix.completeOrderFragments[string(buyOrderFragment.ID)] = buyOrderFragment
	return nil
}

func (matrix *DeltaFragmentMatrix) removeSellOrderFragment(sellOrderFragment *order.Fragment) error {
	if _, ok := matrix.sellOrderFragments[string(sellOrderFragment.ID)]; !ok {
		return nil
	}

	for i := range matrix.buySellDeltaFragments {
		delete(matrix.buySellDeltaFragments[i], string(sellOrderFragment.ID))
	}

	matrix.completeOrderFragments[string(sellOrderFragment.ID)] = sellOrderFragment
	return nil
}

func (matrix *DeltaFragmentMatrix) DeltaFragment(buyOrderFragmentID, sellOrderFragmentID order.FragmentID) *DeltaFragment {
	matrix.EnterReadOnly(nil)
	defer matrix.ExitReadOnly()
	return matrix.deltaFragment(buyOrderFragmentID, sellOrderFragmentID)
}

func (matrix *DeltaFragmentMatrix) deltaFragment(buyOrderFragmentID, sellOrderFragmentID order.FragmentID) *DeltaFragment {
	if deltaFragments, ok := matrix.buySellDeltaFragments[string(buyOrderFragmentID)]; ok {
		if deltaFragment, ok := deltaFragments[string(sellOrderFragmentID)]; ok {
			return deltaFragment
		}
	}
	return nil
}
