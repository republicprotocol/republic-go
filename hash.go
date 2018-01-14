package x

// Hash is an alias for a slice of bytes. It represents a Keccak256 hash.
type Hash []byte

// LessThen returns true if the left hand hash is less than the right hand
// hash, otherwise false.
func (lhs Hash) LessThen(rhs Hash) bool {
}
