package stackint

// Int1024 provides a 1024 bit number optimised to never use the heap
type Int1024 struct {
}

// NewInt1024 returns a new Int1024 from a uint64
func NewInt1024(x uint64) Int1024 {
	panic("Not implemented!")
}

// Equals returns true of x and y represent the same Int1024
func (x *Int1024) Equals(y *Int1024) bool {
	panic("Not implemented!")
}

// Add returns x+y
func (x *Int1024) Add(y *Int1024) Int1024 {
	panic("Not implemented!")
}

// Sub returns x-y
func (x *Int1024) Sub(y *Int1024) Int1024 {
	panic("Not implemented!")
}

// Mul returns x*y
func (x *Int1024) Mul(y *Int1024) Int1024 {
	panic("Not implemented!")
}

// Div returns the quotient of x/y. If y is 0, a run-time panic occurs.
func (x *Int1024) Div(y *Int1024) Int1024 {
	panic("Not implemented!")
}

// Mod returns the modulus x%y. If y is 0, a run-time panic occurs.
func (x *Int1024) Mod(y *Int1024) Int1024 {
	panic("Not implemented!")
}

// ModInverse returns the multiplicative inverse of x in the ring ℤ/nℤ.
// If x and n are not relatively prime, the result in undefined.
func (x *Int1024) ModInverse(n *Int1024) Int1024 {
	panic("Not implemented!")
}
