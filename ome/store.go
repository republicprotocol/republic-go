package ome

import (
	"errors"

	"github.com/republicprotocol/republic-go/orderbook"
)

var ErrComputationNotFound = errors.New("computation not found")

type Storer interface {
	orderbook.Storer

	InsertComputation(Computation) error
	Computation(ComputationID) (Computation, error)
}
