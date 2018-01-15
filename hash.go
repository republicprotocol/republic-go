package x

import identity "github.com/republicprotocol/go-identity"

// Hash is an alias for a slice of bytes. It represents a Keccak256 hash.
type Hash []byte

// LessThan returns true if the left hand hash is less than the right hand
// hash, otherwise false.
func (lhs Hash) LessThan(rhs Hash) bool {
	return false
}

// A Miner represents a miner in the X Network. They are identified by an
// identity.ID and organized by combining their commitment hash with the Epoch
// hash to produce an X hash.
type Miner struct {
	identity.ID
	Commitment Hash
	X          Hash
}
