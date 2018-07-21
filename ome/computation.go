package ome

import (
	"bytes"
	"encoding/base64"
	"time"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
)

// ComputationID is used to distinguish between different combinations of
// orders that are being matched against each other.
type ComputationID [32]byte

// NewComputationID returns the crypto.Keccak256 of a buy order.ID concatenated
// with a sell order.ID.
func NewComputationID(buy, sell order.ID, depth order.FragmentEpochDepth) ComputationID {
	depthBytes := [4]byte{
		byte(depth >> 24),
		byte(depth >> 16),
		byte(depth >> 8),
		byte(depth),
	}
	comID := ComputationID{}
	copy(comID[:], crypto.Keccak256(buy[:], sell[:], depthBytes[:]))
	return comID
}

// String returns a human-readable representation of the ComputationID.
func (id ComputationID) String() string {
	runes := []rune(base64.StdEncoding.EncodeToString(id[:]))
	return string(runes[:4])
}

// ComputationState is used to track the state of a Computation as it changes
// over its lifetime. This prevents duplicated work in the system.
type ComputationState int

// Values for a ComputationState
const (
	ComputationStateNil ComputationState = iota
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
	default:
		return "unsupported state"
	}
}

// Computations is an alias type.
type Computations []Computation

// A Computation is a combination of a buy order.Order and a sell order.Order.
type Computation struct {
	Timestamp  time.Time                `json:"timestamp"`
	ID         ComputationID            `json:"id"`
	Buy        order.Fragment           `json:"buy"`
	Sell       order.Fragment           `json:"sell"`
	Epoch      [32]byte                 `json:"epoch"`
	EpochDepth order.FragmentEpochDepth `json:"epochDepth"`

	State ComputationState `json:"state"`
	Match bool             `json:"match"`
}

// NewComputation returns a pending Computation between a buy order.Order and a
// sell order.Order. It initialized the ComputationID to the Keccak256 hash of
// the buy order.ID and the sell order.ID.
func NewComputation(epoch [32]byte, buy, sell order.Fragment, state ComputationState, match bool) Computation {
	com := Computation{
		Buy:        buy,
		Sell:       sell,
		Epoch:      epoch,
		EpochDepth: buy.EpochDepth,
		State:      state,
		Match:      match,
	}
	com.Timestamp = time.Now()
	com.ID = NewComputationID(buy.OrderID, sell.OrderID, com.EpochDepth)
	return com
}

// Equal returns true when Computations are equal in value and state, and
// returns false otherwise.
func (com *Computation) Equal(arg *Computation) bool {
	return com.Timestamp.Equal(arg.Timestamp) &&
		bytes.Equal(com.ID[:], arg.ID[:]) &&
		com.Buy.Equal(&arg.Buy) &&
		com.Sell.Equal(&arg.Sell) &&
		com.State == arg.State &&
		com.Match == arg.Match
}
