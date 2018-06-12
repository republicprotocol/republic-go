package ome

import (
	"bytes"
	"encoding/base64"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

// ComputationID is used to distinguish between different combinations of
// orders that are being matched against each other.
type ComputationID [32]byte

// String returns a human-readable representation of the ComputationID.
func (id ComputationID) String() string {
	return base64.StdEncoding.EncodeToString(id[:8])
}

// ComputationState is used to track the state of a Computation as it changes
// over its lifetime. This prevents duplicated work in the system.
type ComputationState int

// Values for a ComputationState
const (
	ComputationStateNil = iota
	ComputationStateMatched
	ComputationStateMismatched
	ComputationStateAccepted
	ComputationStateRejected
	ComputationStateSettled
)

// String returns a human-readable representation of the ComputationState.
func (state ComputationState) String() string {
	switch state {
	case ComputationStateNil:
		return "nil"
	case ComputationStateMatched:
		return "matched"
	case ComputationStateMismatched:
		return "mismatched"
	case ComputationStateAccepted:
		return "accepted"
	case ComputationStateRejected:
		return "rejected"
	case ComputationStateSettled:
		return "settled"
	}
	panic("unexpected computation state")
}

// Computations is an alias type.
type Computations []Computation

// A Computation is a combination of a buy order.Order and a sell order.Order.
type Computation struct {
	ID        ComputationID      `json:"id"`
	State     ComputationState   `json:"state"`
	Priority  orderbook.Priority `json:"priority"`
	Match     bool               `json:"match"`
	Timestamp time.Time          `json:"timestamp"`

	Buy       order.ID `json:"buy"`
	Sell      order.ID `json:"sell"`
	EpochHash [32]byte `json:"epochHash"`
}

// NewComputation returns a pending Computation between a buy order.Order and a
// sell order.Order. It initialized the ComputationID to the Keccak256 hash of
// the buy order.ID and the sell order.ID.
func NewComputation(buy, sell order.ID, epochHash [32]byte) Computation {
	com := Computation{
		Buy:       buy,
		Sell:      sell,
		EpochHash: epochHash,
	}
	com.ID = GenerateComputationID(buy, sell)
	return com
}

// Equal returns true when Computations are equal in value and state, and
// returns false otherwise.
func (com *Computation) Equal(arg *Computation) bool {
	return bytes.Equal(com.ID[:], arg.ID[:]) &&
		com.State == arg.State &&
		com.Priority == arg.Priority &&
		com.Match == arg.Match &&
		com.Timestamp.Equal(arg.Timestamp) &&
		com.Buy.Equal(arg.Buy) &&
		com.Sell.Equal(arg.Sell)
}

// GenerateComputationID from the buy order.ID and sell order.ID by
// concatenating them and applying a crypto.Keccak256 hash.
func GenerateComputationID(buy, sell order.ID) ComputationID {
	comID := ComputationID{}
	copy(comID[:], crypto.Keccak256(buy[:], sell[:]))
	return comID
}
