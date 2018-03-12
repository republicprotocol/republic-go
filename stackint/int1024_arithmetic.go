package stackint

// Add returns x+y
func (x *Int1024) Add(y *Int1024) Int1024 {
	z := zero()

	// Loop through each of the 16 pair of words and add them
	// Storing the overflow if necessary
	var overflow Word
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		z.words[i] = x.words[i] + y.words[i] + overflow
		if x.words[i] > WORDMAX-y.words[i]-overflow {
			overflow = 1
		} else {
			overflow = 0
		}
	}

	if overflow == 1 {
		// WARNING: Overflow occured
	}

	return z
}

// Sub returns x-y
func (x *Int1024) Sub(y *Int1024) Int1024 {
	z := zero()

	var overflow Word
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		z.words[i] = x.words[i] - y.words[i] - overflow
		if x.words[i] < y.words[i]+overflow {
			overflow = 1
		} else {
			overflow = 0
		}
	}

	if overflow == 1 {
		// WARNING: Overflow occured
	}

	return z
}

// Mul returns x*y
func (x *Int1024) Mul(y *Int1024) Int1024 {
	// Naïve inplementation!
	// Uses up to 16384 uint64 additions (worst case)
	// TODO: Rewrite using more efficient algorithm
	z := zero()
	shifted := x.Clone()

	for i := INT1024WORDS - 1; i >= 0; i-- {
		word := y.words[i]
		for j := uint(0); j < WORDSIZE; j++ {
			bit := (word >> j) & 1
			if bit == 1 {
				z.overwritingAdd(&shifted)
			}
			shifted.overwritingShiftLeftByOne()
		}
	}

	return z
}

// DivMod returns (x/y, x%y). If y is 0, a run-time panic occurs.
func (x *Int1024) DivMod(y *Int1024) (Int1024, Int1024) {
	dividend := x.Clone()
	denom := y.Clone()
	current := Int1024FromUint64(1)
	answer := zero()

	if denom.Equals(&Zero) {
		panic("division by zero")
	}

	limit := answer.NOT()
	limit.overwritingShiftRightByOne()
	overflowed := false
	for denom.LessThanOrEqual(&dividend) {
		if !denom.LessThan(&limit) {
			overflowed = true
			break
		}
		denom.overwritingShiftLeftByOne()
		current.overwritingShiftLeftByOne()
	}

	if !overflowed {
		denom.overwritingShiftRightByOne()
		current.overwritingShiftRightByOne()
	}

	for !current.Equals(&Zero) {
		if !dividend.LessThan(&denom) {
			dividend = dividend.Sub(&denom)
			answer = answer.OR(&current)
		}
		current.overwritingShiftRightByOne()
		denom.overwritingShiftRightByOne()
	}

	return answer, dividend
}

// Div returns the quotient of x/y. If y is 0, a run-time panic occurs.
func (x *Int1024) Div(y *Int1024) Int1024 {
	div, _ := x.DivMod(y)
	return div
}

// Mod returns the modulus x%y. If y is 0, a run-time panic occurs.
func (x *Int1024) Mod(y *Int1024) Int1024 {
	_, mod := x.DivMod(y)
	return mod
}

// ModInverse returns the multiplicative inverse of x in the ring ℤ/nℤ.
// If x and n are not relatively prime, the result in undefined.
func (x *Int1024) ModInverse(n *Int1024) Int1024 {
	panic("Not implemented!")
}
