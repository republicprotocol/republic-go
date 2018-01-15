package x

import identity "github.com/republicprotocol/go-identity"

// Hash is an alias for a slice of bytes. It represents a Keccak256 hash.
type Hash []byte

// LessThan returns true if the left hand hash is less than the right hand
// hash, otherwise false.
func (lhs Hash) LessThan(rhs Hash) bool {
	for i, _ := range lhs {
		if lhs[i] < rhs[i] {
			return true
		} else if lhs[i] > rhs[i] {
			return false
		}
	}
	return false
}

// A Miner represents a miner in the X Network. They are identified by an
// identity.ID and organized by combining their commitment hash with the Epoch
// hash to produce an X hash.
type Miner struct {
	identity.ID
	Commitment Hash

	X        Hash
	Class    int
	MNetwork int
}

// NewMiner returns a Miner with the given ID and commitment Hash. The miner
// will have no X Hash, no class assignment, and no M Network assignment. These
// values must be generated using the AssignX, AssignClass, and AssignMNetwork
// functions respectively.
func NewMiner(id identity.ID, commitment Hash) Miner {
	return Miner{
		ID:         id,
		Commitment: commitment,
	}
}
