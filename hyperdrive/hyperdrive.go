package hyper

import (
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/order"
)

// The Delegate is a callback interface for the Hyperdrive
type Delegate interface {
	OnOrderReleased(order order.ID)
	OnOrderExecuted(buyOrder order.ID, sellOrder order.ID)
}

// A Tuple represents two orders that are compatible with one another
type Tuple struct {
	lhs order.ID
	rhs order.ID
}

// A ConflictSet contains a set of order matches that are in a direct or transitive conflict.
type ConflictSet struct {
	id     uint64
	open   bool
	tuples []Tuple
}

// Drive contains a series of conflict sets and a mapping from tuples to their corresponding conflict sets
type Drive struct {
	do.GuardedObject

	conflictSetsNextID        uint64
	conflictSets              map[uint64]*ConflictSet
	tupleToConflictSetMapping map[string]*ConflictSet
}

// AddTuple inserts a tuple into the Hyperdrive.
// If no conflict set is found, a new one is created.
// If one conflict set is found, the tuple is added to it.
// If two conflict sets are found, they are merged and the tuple added to the combined set.
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
