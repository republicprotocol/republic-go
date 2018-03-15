package stackint

// Add returns x+y
func (x *Int1024) Add(y *Int1024) Int1024 {

	// xB, _ := big.NewInt(0).SetString(x.ToBinary(), 2)
	// yB, _ := big.NewInt(0).SetString(y.ToBinary(), 2)
	// expected := big.NewInt(0).Add(xB, yB)

	z := x.Clone()
	z.overwritingAdd(y)

	// actual, _ := big.NewInt(0).SetString(z.ToBinary(), 2)
	// if expected.Cmp(actual) != 0 && expected.BitLen() <= 1024 {
	// 	panic(fmt.Sprintf("Addition failed! for %s and %s.\n\nExpected %b\n\nGot %b", x.ToBinary(), y.ToBinary(), expected, actual))
	// }

	return z
}

// overwritingAdd sets x to x+y (used for multiplication)
func (x *Int1024) overwritingAdd(y *Int1024) {
	var overflow Word
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		previousOverflow := overflow
		if x.words[i] > WORDMAX-y.words[i] || x.words[i] > WORDMAX-y.words[i]-previousOverflow {
			overflow = 1
		} else {
			overflow = 0
		}
		x.words[i] = x.words[i] + y.words[i] + previousOverflow
	}

	if overflow == 1 {
		// WARNING: Overflow occured
	}
}

// Sub returns x-y
func (x *Int1024) Sub(y *Int1024) Int1024 {

	// xB, _ := big.NewInt(0).SetString(x.ToBinary(), 2)
	// yB, _ := big.NewInt(0).SetString(y.ToBinary(), 2)
	// expected := big.NewInt(0).Sub(xB, yB)

	z := zero()

	var overflow Word
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		z.words[i] = x.words[i] - y.words[i] - overflow
		if x.words[i] < y.words[i]+overflow || y.words[i] == WORDMAX && overflow == 1 {
			overflow = 1
		} else {
			overflow = 0
		}
	}

	if overflow == 1 {
		// WARNING: Overflow occured
	}

	// actual, _ := big.NewInt(0).SetString(z.ToBinary(), 2)
	// if expected.Cmp(actual) != 0 && expected.BitLen() <= 1024 {
	// 	panic(fmt.Sprintf("Subtraction failed! for %s and %s.\n\nExpected %b\n\nGot %b", x.ToBinary(), y.ToBinary(), expected, actual))
	// }

	return z
}

// Mul returns x*y
func (x *Int1024) Mul(y *Int1024) Int1024 {

	// xB, _ := big.NewInt(0).SetString(x.ToBinary(), 2)
	// yB, _ := big.NewInt(0).SetString(y.ToBinary(), 2)
	// expected := big.NewInt(0).Mul(xB, yB)

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
				z = z.Add(&shifted)
			}
			shifted.overwritingShiftLeftByOne()
		}
	}

	// actual, _ := big.NewInt(0).SetString(z.ToBinary(), 2)
	// if expected.Cmp(actual) != 0 && expected.BitLen() <= 1024 {
	// 	panic(fmt.Sprintf("Multiplication failed for %s and %s.\n\nExpected %b\n\nGot %b", x.ToBinary(), y.ToBinary(), expected, actual))
	// }

	return z
}

// DivMod returns (x/y, x%y). If y is 0, a run-time panic occurs.
func (x *Int1024) DivMod(y *Int1024) (Int1024, Int1024) {
	dividend := x.Clone()
	denom := y.Clone()
	current := FromUint64(1)
	answer := zero()

	if denom.Equals(&ZERO) {
		panic("division by zero")
	}

	limit := MAXINT1024.Clone()
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

	for !current.Equals(&ZERO) {
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

// // ModInverse returns the multiplicative inverse of x in the ring ℤ/nℤ.
// // If x and n are not relatively prime, the result in undefined.
// func (x *Int1024) ModInverse(n *Int1024) Int1024 {
// 	panic("Not implemented!")
// }

// ModInverse sets z to the multiplicative inverse of g in the ring ℤ/nℤ
// and returns z. If g and n are not relatively prime, the result is undefined.
// Code adapted from https://www.di-mgt.com.au/euclidean.html
func (x *Int1024) ModInverse(n *Int1024) Int1024 {
	u := x.Clone()
	v := n.Clone()
	// unsigned int inv, u1, u3, v1, v3, t1, t3, q;
	// int iter;
	/* Step X1. Initialise */
	u1 := FromUint64(1)
	u3 := u.Clone()
	v1 := zero()
	v3 := v.Clone()
	/* Remember odd/even iterations */
	iter := 1
	/* Step X2. Loop while v3 != 0 */
	for !v3.IsZero() {
		/* Step X3. Divide and "Subtract" */
		q, t3 := u3.DivMod(&v3)
		tmp := q.Mul(&v1)
		t1 := u1.Add(&tmp)
		/* Swap */
		u1, v1, u3, v3 = v1, t1, v3, t3
		iter = -iter
	}

	/* Make sure u3 = gcd(u,v) == 1 */
	if !u3.Equals(&ONE) {
		// return zero() /* Error: No inverse exists */
		panic("not relatively prime")
	}
	/* Ensure a positive result */

	inv := u1
	if iter < 0 {
		inv = v.Sub(&inv)
	}
	return inv
}

// Exp returns x**y
func (x *Int1024) Exp(y *Int1024) Int1024 {
	if y.IsZero() {
		return ONE
	} else if y.Equals(&ONE) {
		return *x
	} else if y.IsEven() {
		square := x.Mul(x)
		half := y.Div(&TWO)
		return square.Exp(&half)
	}
	square := x.Mul(x)
	ySubOne := y.Sub(&ONE)
	half := ySubOne.Div(&TWO)
	power := square.Exp(&half)
	return x.Mul(&power)
}
