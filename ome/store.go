package ome

import (
	"errors"

	"github.com/republicprotocol/republic-go/orderbook"
)

// ErrComputationNotFound is returned when the Storer cannot find a Computation
// associated with a ComputationID.
var ErrComputationNotFound = errors.New("computation not found")

// A Storer for the Ome extends the orderbook.Storer by exposing the ability
// to store and load a Computation.
type Storer interface {
	orderbook.Storer

	// InsertComputation into the Storer. The primary use case for this is
	// storing the state of the Computation so that the Ome does not attempt to
	// redo a Computation when it reboots.
	InsertComputation(Computation) error

	// Computation returns a Computation associated with the ComputationID.
	// Returns ErrComputationNotFound if the Computation does not exist in the
	// Storer.
	Computation(ComputationID) (Computation, error)
}
