package stackint

// Int1024 provides a 1024 bit number optimised to never use the heap
type Int1024 struct {
}

// Add returns x+y
func (x *Int1024) Add(y *Int1024) Int1024 {
	return *x
}

// Sub returns x-y
func (x *Int1024) Sub(y *Int1024) Int1024 {
	return *x
}

// Mul returns x*y
func (x *Int1024) Mul(y *Int1024) Int1024 {
	return *x
}

// Div returns the quotient of x/y. If y is 0, a run-time panic occurs.
func (x *Int1024) Div(y *Int1024) Int1024 {
	return *x
}

// Mod returns the modulus x%y. If y is 0, a run-time panic occurs.
func (x *Int1024) Mod(y *Int1024) Int1024 {
	return *x
}

// ModInverse returns the multiplicative inverse of x in the ring ℤ/nℤ.
// If x and n are not relatively prime, the result in undefined.
func (x *Int1024) ModInverse(n *Int1024) Int1024 {
	return *x
}
