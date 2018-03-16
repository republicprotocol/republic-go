package stackint

import (
	"fmt"
	"math/big"
)

// Add returns x+y
func (x *Int1024) Add(y *Int1024) Int1024 {

	xB, _ := big.NewInt(0).SetString(x.ToBinary(), 2)
	yB, _ := big.NewInt(0).SetString(y.ToBinary(), 2)
	expected := big.NewInt(0).Add(xB, yB)

	z := x.Clone()
	z.Inc(y)

	actual, _ := big.NewInt(0).SetString(z.ToBinary(), 2)
	if expected.Cmp(actual) != 0 && expected.BitLen() <= 1024 {
		panic(fmt.Sprintf("Addition failed! for %s and %s.\n\nExpected %b\n\nGot %b", x.ToBinary(), y.ToBinary(), expected, actual))
	}

	return z
}

// Inc sets x to x+y (used for multiplication)
func (x *Int1024) Inc(y *Int1024) {

	var overflow Word
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		previousOverflow := overflow
		if x.words[i] > WORDMAX-y.words[i] || x.words[i] > WORDMAX-y.words[i]-previousOverflow || y.words[i] > WORDMAX-previousOverflow {
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

	// var words [INT1024WORDS]Word

	// overflow := Word(0)
	// for i := INT1024WORDS - 1; i >= 0; i-- {
	// 	words[i] = x.words[i] - y.words[i] - overflow
	// 	if x.words[i] < y.words[i]+overflow || y.words[i] == WORDMAX && overflow == 1 {
	// 		overflow = 1
	// 	} else {
	// 		overflow = 0
	// 	}
	// }

	// // if overflow == 1 {
	// // 	// WARNING: Overflow occured
	// // }

	// z := Int1024{
	// 	words: words,
	// }

	z := x.Clone()
	z.Dec(y)

	// actual, _ := big.NewInt(0).SetString(z.ToBinary(), 2)
	// if expected.Cmp(actual) != 0 && expected.BitLen() <= 1024 {
	// 	panic(fmt.Sprintf("Subtraction failed! for %s and %s.\n\nExpected %b\n\nGot %b", x.ToBinary(), y.ToBinary(), expected, actual))
	// }

	return z
}

// Dec sets x to x-y
func (x *Int1024) Dec(y *Int1024) {
	overflow := Word(0)
	for i := INT1024WORDS - 1; i >= 0; i-- {
		newOverflow := Word(0)
		if x.words[i] < y.words[i]+overflow || y.words[i] == WORDMAX && overflow == 1 {
			newOverflow = 1
		}
		x.words[i] = x.words[i] - y.words[i] - overflow
		overflow = newOverflow
	}

	if overflow == 1 {
		// WARNING: Overflow occured
	}
}

// Mul returns x*y
func (x *Int1024) Mul(y *Int1024) Int1024 {

	// xB, _ := big.NewInt(0).SetString(x.ToBinary(), 2)
	// yB, _ := big.NewInt(0).SetString(y.ToBinary(), 2)
	// expected := big.NewInt(0).Mul(xB, yB)

	// Naïve inplementation!
	// Uses up to 16384 uint64 additions (worst case)
	// TODO: Rewrite using more efficient algorithm
	z := Zero()
	shifted := x.Clone()

	for i := INT1024WORDS - 1; i >= 0; i-- {
		word := y.words[i]
		for j := uint(0); j < WORDSIZE; j++ {
			bit := (word >> j) & 1
			if bit == 1 {
				z.Inc(&shifted)
			}
			shifted.ShiftLeftInPlace()
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
	answer := Zero()

	if denom.IsZero() {
		panic("division by zero")
	}

	limit := MAXINT1024()
	limit.ShiftRightInPlace()
	overflowed := false
	for denom.LessThanOrEqual(&dividend) {
		if !denom.LessThan(&limit) {
			overflowed = true
			break
		}
		denom.ShiftLeftInPlace()
		current.ShiftLeftInPlace()
	}

	if !overflowed {
		denom.ShiftRightInPlace()
		current.ShiftRightInPlace()
	}

	for !current.IsZero() {
		if !dividend.LessThan(&denom) {
			dividend = dividend.Sub(&denom)
			answer = answer.OR(&current)
		}
		current.ShiftRightInPlace()
		denom.ShiftRightInPlace()
	}

	return answer, dividend
}

// Div returns the quotient of x/y. If y is 0, a run-time panic occurs.
func (x *Int1024) Div(y *Int1024) Int1024 {
	div, _ := x.DivMod(y)
	return div
}

// Mod returns the modulus x%n. If n is 0, a run-time panic occurs.
func (x *Int1024) Mod(n *Int1024) Int1024 {
	// // Switch not needed. Is it a performance improvement?
	// switch x.Cmp(n) {
	// case -1:
	// 	return x.Clone()
	// case 0:
	// 	return Zero()
	// case 1:
	// 	_, mod := x.DivMod(n)
	// 	return mod
	// default:
	// 	panic("unexpected cmp result (expecting -1, 0 or 1)")
	// }

	xB, _ := big.NewInt(0).SetString(x.ToBinary(), 2)
	yB, _ := big.NewInt(0).SetString(n.ToBinary(), 2)
	expected := big.NewInt(0).Mod(xB, yB)

	_, z := x.DivMod(n)

	actual, _ := big.NewInt(0).SetString(z.ToBinary(), 2)
	if expected.Cmp(actual) != 0 && expected.BitLen() <= 1024 {
		panic(fmt.Sprintf("Modulo failed! for %s and %s.\n\nExpected %b\n\nGot %b", x.ToBinary(), n.ToBinary(), expected, actual))
	}

	return z
}

// SubModulo returns (x - y) % n
func (x *Int1024) SubModulo(y, n *Int1024) Int1024 {
	switch x.Cmp(y) {
	case 1:
		// x - y
		sub := x.Sub(y)
		return sub.Mod(n)
	case 0:
		if n.IsZero() {
			panic("division by zero")
		}
		return Zero()
	case -1:
		sub := y.Sub(x)
		mod := sub.Mod(n)
		if mod.IsZero() {
			return mod
		}
		return n.Sub(&mod)
	default:
		panic("unexpected cmp result (expecting -1, 0 or 1)")
	}
}

// AddModulo returns (x + y) % n (x+y can be larger than 2^1024)
func (x *Int1024) AddModulo(y, n *Int1024) Int1024 {
	xB, _ := big.NewInt(0).SetString(x.ToBinary(), 2)
	yB, _ := big.NewInt(0).SetString(y.ToBinary(), 2)
	nB, _ := big.NewInt(0).SetString(n.ToBinary(), 2)
	expected := big.NewInt(0).Mod(big.NewInt(0).Add(xB, yB), nB)

	z := x.addModulo(y, n)

	actual, _ := big.NewInt(0).SetString(z.ToBinary(), 2)
	if expected.Cmp(actual) != 0 && expected.BitLen() <= 1024 {
		panic(fmt.Sprintf("AddModulo failed! for %s + %s mod %s.\n\nExpected %s\n\nGot %s", x.String(), y.String(), n.String(), expected, actual))
	}

	return z
}

func (x *Int1024) addModulo(y, n *Int1024) Int1024 {

	if !x.IsBitSet(1023) && !y.IsBitSet(1023) {
		tmp := x.Add(y)
		return tmp.Mod(n)
	}

	X := *x
	if x.GreaterThanOrEqual(n) {
		X = x.Mod(n)
	}
	Y := *y
	if y.GreaterThanOrEqual(n) {
		Y = y.Mod(n)
	}

	// Check again
	if !X.IsBitSet(1023) && !Y.IsBitSet(1023) {
		tmp := X.Add(&Y)
		return tmp.Mod(n)
	}

	if n.IsZero() {
		// Can be placed at the top but then might be checked again in Mod()
		panic("division by zero")
	}

	diff := n.Sub(&X)
	if diff.LessThanOrEqual(&Y) {
		ret := Y.Sub(&diff)
		return ret
	}

	return X.Add(&Y)
}

// MulModulo returns (x*y) % n
func (x *Int1024) MulModulo(y, n *Int1024) Int1024 {

	// Not optimized

	// TODO!!!
	// Implement Shrage's Method
	// https://github.com/wdavidw/node-gsl/blob/master/deps/gsl-1.14/rng/schrage.c

	z := Zero()
	shifted := x.Mod(n)

	for i := INT1024WORDS - 1; i >= 0; i-- {
		word := y.words[i]
		for j := uint(0); j < WORDSIZE; j++ {
			bit := (word >> j) & 1
			if bit == 1 {
				z = z.AddModulo(&shifted, n)
			}
			shifted = shifted.AddModulo(&shifted, n)
		}
	}

	return z
}

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
	v1 := Zero()
	v3 := v.Clone()
	/* Remember odd/even iterations */
	iter := 1
	/* Step X2. Loop while v3 != 0 */
	for !v3.IsZero() {
		/* Step X3. Divide and "Subtract" */
		q, t3 := u3.DivMod(&v3)
		tmp := q.Mul(&v1)
		u1.Inc(&tmp)
		/* Swap */
		u1, v1, u3, v3 = v1, u1, v3, t3
		iter = -iter
	}

	/* Make sure u3 = gcd(u,v) == 1 */
	if !u3.EqualsUint64(1) {
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
		return One()
	} else if y.EqualsUint64(1) {
		return *x
	} else if y.IsEven() {
		square := x.Mul(x)
		half := y.Div(&two)
		return square.Exp(&half)
	}
	square := x.Mul(x)
	ySubOne := y.Sub(&one)
	half := ySubOne.Div(&two)
	power := square.Exp(&half)
	return x.Mul(&power)
}
