package hyper

import (
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/order"
)

type Delegate interface {
	OnOrderReleased(order order.ID)
	OnOrderExecuted(buyOrder order.ID, sellOrder order.ID)
}

type Tuple struct {
	lhs order.ID
	rhs order.ID
}

type ConflictSet struct {
	id     uint64
	open   bool
	tuples []Tuple
}

type Drive struct {
	do.GuardedObject

	conflictSetsNextID        uint64
	conflictSets              map[uint64]*ConflictSet
	tupleToConflictSetMapping map[string]*ConflictSet
}

func (hyperdrive *Drive) AddTuple(tuple Tuple) {
	hyperdrive.Enter(nil)
	defer hyperdrive.Exit()
	hyperdrive.addTuple(tuple)
}

func (hyperdrive *Drive) addTuple(tuple Tuple) {
	lhsConflictSet := hyperdrive.tupleToConflictSetMapping[string(tuple.lhs)]
	rhsConflictSet := hyperdrive.tupleToConflictSetMapping[string(tuple.rhs)]
	if lhsConflictSet != nil && rhsConflictSet != nil {
		if !(lhsConflictSet.open && rhsConflictSet.open) {
			return
		}
		lhsConflictSet.tuples = append(lhsConflictSet.tuples, tuple)
		hyperdrive.mergeConflictSets(lhsConflictSet, rhsConflictSet)
		return
	}
	if lhsConflictSet != nil {
		if !lhsConflictSet.open {
			return
		}
		lhsConflictSet.tuples = append(lhsConflictSet.tuples, tuple)
		return
	}
	if rhsConflictSet != nil {
		if !rhsConflictSet.open {
			return
		}
		rhsConflictSet.tuples = append(rhsConflictSet.tuples, tuple)
		return
	}
	hyperdrive.createConflictSet(tuple)
}

func (hyperdrive *Drive) mergeConflictSets(lhsConflictSet, rhsConflictSet *ConflictSet) {
	if lhsConflictSet == nil || rhsConflictSet == nil {
		return
	}
	for _, tuple := range rhsConflictSet.tuples {
		hyperdrive.tupleToConflictSetMapping[string(tuple.lhs)] = lhsConflictSet
		hyperdrive.tupleToConflictSetMapping[string(tuple.rhs)] = lhsConflictSet
	}
	delete(hyperdrive.conflictSets, rhsConflictSet.id)
}

func (hyperdrive *Drive) createConflictSet(tuple Tuple) {
	conflictSet := &ConflictSet{
		id:     hyperdrive.conflictSetsNextID,
		tuples: []Tuple{tuple},
	}
	hyperdrive.tupleToConflictSetMapping[string(tuple.lhs)] = conflictSet
	hyperdrive.tupleToConflictSetMapping[string(tuple.rhs)] = conflictSet
	hyperdrive.conflictSets[hyperdrive.conflictSetsNextID] = conflictSet
	hyperdrive.conflictSetsNextID++
}
