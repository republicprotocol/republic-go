package compute

import (
	"github.com/republicprotocol/republic-go/stackint"

	do "github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/order"
)

type DeltaBuilder struct {
	do.GuardedObject

	k                      int64
	prime                  *stackint.Int1024
	deltas                 map[string]*Delta
	deltaFragments         map[string]*DeltaFragment
	deltasToDeltaFragments map[string][]*DeltaFragment
}

func NewDeltaBuilder(k int64, prime *stackint.Int1024) *DeltaBuilder {
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
	defer builder.ExitReadOnly()
	return builder.hasDelta(deltaID)
}

func (builder *DeltaBuilder) hasDelta(deltaID DeltaID) bool {
	_, ok := builder.deltas[string(deltaID)]
	return ok
}

func (builder *DeltaBuilder) HasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	builder.EnterReadOnly(nil)
	defer builder.ExitReadOnly()
	return builder.hasDeltaFragment(deltaFragmentID)
}

func (builder *DeltaBuilder) hasDeltaFragment(deltaFragmentID DeltaFragmentID) bool {
	_, ok := builder.deltaFragments[string(deltaFragmentID)]
	return ok
}

func (builder *DeltaBuilder) SetK(k int64) {
	builder.Enter(nil)
	defer builder.Exit()
	builder.setK(k)
}

func (builder *DeltaBuilder) setK(k int64) {
	builder.k = k
}

type DeltaFragmentMatrix struct {
	do.GuardedObject

	prime                  *stackint.Int1024
	buyOrderFragments      map[string]*order.Fragment
	sellOrderFragments     map[string]*order.Fragment
	buySellDeltaFragments  map[string]map[string]*DeltaFragment
	completeOrderFragments map[string]bool
}

func NewDeltaFragmentMatrix(prime *stackint.Int1024) *DeltaFragmentMatrix {
	return &DeltaFragmentMatrix{
		GuardedObject:          do.NewGuardedObject(),
		prime:                  prime,
		buyOrderFragments:      map[string]*order.Fragment{},
		sellOrderFragments:     map[string]*order.Fragment{},
		buySellDeltaFragments:  map[string]map[string]*DeltaFragment{},
		completeOrderFragments: map[string]bool{},
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
	if _, ok := matrix.buyOrderFragments[string(buyOrderFragment.OrderID)]; ok {
		return []*DeltaFragment{}, nil
	}
	if _, ok := matrix.completeOrderFragments[string(buyOrderFragment.OrderID)]; ok {
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
		deltaFragmentsMap[string(matrix.sellOrderFragments[i].OrderID)] = deltaFragment
	}

	matrix.buyOrderFragments[string(buyOrderFragment.OrderID)] = buyOrderFragment
	matrix.buySellDeltaFragments[string(buyOrderFragment.OrderID)] = deltaFragmentsMap
	return deltaFragments, nil
}

func (matrix *DeltaFragmentMatrix) insertSellOrderFragment(sellOrderFragment *order.Fragment) ([]*DeltaFragment, error) {
	if _, ok := matrix.sellOrderFragments[string(sellOrderFragment.OrderID)]; ok {
		return []*DeltaFragment{}, nil
	}
	if _, ok := matrix.completeOrderFragments[string(sellOrderFragment.OrderID)]; ok {
		return []*DeltaFragment{}, nil
	}

	deltaFragments := make([]*DeltaFragment, 0, len(matrix.buyOrderFragments))
	for i := range matrix.buyOrderFragments {
		deltaFragment := NewDeltaFragment(matrix.buyOrderFragments[i], sellOrderFragment, matrix.prime)
		if deltaFragment == nil {
			continue
		}
		if _, ok := matrix.buySellDeltaFragments[string(matrix.buyOrderFragments[i].OrderID)]; ok {
			deltaFragments = append(deltaFragments, deltaFragment)
			matrix.buySellDeltaFragments[string(matrix.buyOrderFragments[i].OrderID)][string(sellOrderFragment.OrderID)] = deltaFragment
		}
	}

	matrix.sellOrderFragments[string(sellOrderFragment.OrderID)] = sellOrderFragment
	return deltaFragments, nil
}

func (matrix *DeltaFragmentMatrix) RemoveOrderFragment(orderID order.ID) error {
	matrix.Enter(nil)
	defer matrix.Exit()
	if err := matrix.removeBuyOrderFragment(orderID); err != nil {
		return err
	}
	return matrix.removeSellOrderFragment(orderID)
}

func (matrix *DeltaFragmentMatrix) removeBuyOrderFragment(buyOrderID order.ID) error {
	if _, ok := matrix.buyOrderFragments[string(buyOrderID)]; !ok {
		return nil
	}

	delete(matrix.buyOrderFragments, string(buyOrderID))
	delete(matrix.buySellDeltaFragments, string(buyOrderID))

	matrix.completeOrderFragments[string(buyOrderID)] = true
	return nil
}

func (matrix *DeltaFragmentMatrix) removeSellOrderFragment(sellOrderID order.ID) error {
	if _, ok := matrix.sellOrderFragments[string(sellOrderID)]; !ok {
		return nil
	}

	for i := range matrix.buySellDeltaFragments {
		delete(matrix.buySellDeltaFragments[i], string(sellOrderID))
	}

	matrix.completeOrderFragments[string(sellOrderID)] = true
	return nil
}
