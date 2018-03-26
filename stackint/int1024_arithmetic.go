package stackint

import (
	"github.com/republicprotocol/republic-go/stackint/asm"
)

// Add returns x+y
func (x *Int1024) Add(y *Int1024) Int1024 {

	z := x.Clone()
	z.Inc(y)

	return z
}

// Inc sets x to x+y (used for multiplication)
func (x *Int1024) Inc(y *Int1024) {

	if y.IsZero() {
		return
	} else if x.IsZero() {
		var i uint16
		for i = 0; i < y.length; i++ {
			x.words[i] = y.words[i]
		}
		x.length = y.length
		return
	}

	a, b := x, y
	if x.length < y.length {
		a, b = y, x
	}

	m := a.length
	n := b.length

	c := asm.AddVV(x.words[0:n], a.words[:], b.words[:])
	if m > n {
		c = asm.AddVW(x.words[n:m], a.words[n:], c)
	}
	x.length = m
	if c > 0 {
		if m < INT1024WORDS {
			x.words[x.length] = c
			x.length++
		} else {
			x.setLength()
		}
	}
}

// Sub returns x-y
func (x *Int1024) Sub(y *Int1024) Int1024 {

	z := x.Clone()
	z.Dec(y)

	return z
}

// Dec sets x to x-y
func (x *Int1024) Dec(y *Int1024) {

	if y.IsZero() {
		return
	}

	m := x.length
	n := y.length

	c := asm.SubVV(x.words[0:n], x.words[:], y.words[:])
	if m > n {
		c = asm.SubVW(x.words[n:], x.words[n:], c)
		if c != 0 {
			panic("!!!")
		}
	}

	if c != 0 {
		for i := n; i < INT1024WORDS; i++ {
			x.words[i] = WORDMAX
		}
		x.length = INT1024WORDS
	} else {
		x.setLength()
	}
}

func (x *Int1024) BasicMulBig(y *Int1024) [INT1024WORDS * 2]asm.Word {

	var words [INT1024WORDS * 2]asm.Word
	var i uint16
	var j uint16
	l := uint16(x.length)
	for i = 0; i < y.length; i++ {
		d := y.words[i]
		if d != 0 {
			var c asm.Word
			for j = i; j < i+l; j++ {
				var z0, z1 asm.Word
				z1, zz0 := asm.MulWW(x.words[j-i], d)
				if z0 = zz0 + words[j]; z0 < zz0 {
					z1++
				}
				c, words[j] = asm.AddWW(z0, c, 0)
				c += z1
			}
			words[l+i] = c
		}
	}

	return words
}

// BasicMul returns x*y using the shift and add method
func (x *Int1024) BasicMul(y *Int1024) Int1024 {

	var words [INT1024WORDS]asm.Word
	var i uint16
	var j uint16
	l := uint16(x.length)
	for i = 0; i < y.length; i++ {
		d := y.words[i]
		if d != 0 {
			var c asm.Word
			for j = i; j < i+l; j++ {
				var z0, z1 asm.Word
				z1, zz0 := asm.MulWW(x.words[j-i], d)
				if z0 = zz0 + words[j]; z0 < zz0 {
					z1++
				}
				c, words[j] = asm.AddWW(z0, c, 0)
				if words[j] != 0 {
				}
				c += z1
			}
			words[l+i] = c
			if words[l+i] != 0 {
			}
		}
	}
	var highest uint16
	for i := x.length + y.length - 1; i > 0; i-- {
		if words[i] > 0 {
			highest = i
			break
		}
	}
	return Int1024{
		words, highest + 1,
	}
}

func mulAddWW(x *Int1024, y asm.Word) Int1024 {

	m := x.length
	z := Zero()
	nxt := asm.MulAddVWW(z.words[0:m], x.words[0:m], y, 0)
	z.length = m
	if m < INT1024WORDS && nxt > 0 {
		z.words[m] = nxt
		z.length++
	}

	return z
}

// Mul returns x*y
func (x *Int1024) Mul(y *Int1024) Int1024 {

	m := x.length
	n := y.length
	switch {
	case m < n:
		x, y = y, x
	case x.IsZero() || y.IsZero():
		return Zero()
	case n == 1:
		return mulAddWW(x, y.words[0])
	}

	if x.length+y.length <= INT1024WORDS {
		return x.BasicMul(y)
	}
	words := x.BasicMulBig(y)
	var words2 [INT1024WORDS]asm.Word
	var highest uint16
	var i uint16
	for i = 0; i < INT1024WORDS; i++ {
		words2[i] = words[i]
		if words2[i] > 0 {
			highest = i
		}
	}
	return Int1024{words2, highest + 1}
}

// Div returns the quotient of x/y. If y is 0, a run-time panic occurs.
func (x *Int1024) Div(y *Int1024) Int1024 {
	div, _ := x.DivMod(y)
	return div
}

// Mod returns the modulus x%n. If n is 0, a run-time panic occurs.
func (x *Int1024) Mod(n *Int1024) Int1024 {
	// Switch not needed. Is it a performance improvement?
	switch x.Cmp(n) {
	case -1:
		return x.Clone()
	case 0:
		return Zero()
	case 1:
		_, mod := x.DivMod(n)
		return mod
	default:
		panic("unexpected cmp result (expecting -1, 0 or 1)")
	}
}

// SubModulo returns (x - y) % n
func (x *Int1024) SubModulo(y, n *Int1024) Int1024 {

	// exp := big.NewInt(0).Sub(x.ToBigInt(), y.ToBigInt())
	// exp = exp.Mod(exp, n.ToBigInt())

	z := x.subModulo(y, n)

	// actual := z.ToBigInt()
	// if exp.Cmp(actual) != 0 {
	// 	panic("Panic in SubModulo")
	// }

	return z
}

func (x *Int1024) subModulo(y, n *Int1024) Int1024 {

	switch x.Cmp(y) {
	case 1: // x > y
		// x - y
		sub := x.Sub(y)
		return sub.Mod(n)
	case 0: // x == y
		if n.IsZero() {
			panic("division by zero")
		}
		return Zero()
	case -1: // x < y
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

// AddModulo returns (x + y) % n (x+y can be larger than 2^SIZE)
func (x *Int1024) AddModulo(y, n *Int1024) Int1024 {

	if !x.IsBitSet(SIZE-1) && !y.IsBitSet(SIZE-1) {
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
	if !X.IsBitSet(SIZE-1) && !Y.IsBitSet(SIZE-1) {
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
	if x.length+y.length > INT1024WORDS {
		return x.MulModuloBig(y, n)
	}
	z := x.Mul(y)
	return z.Mod(n)
}

// ModInverse sets z to the multiplicative inverse of g in the ring ℤ/nℤ
// and returns z. If g and n are not relatively prime, the result is undefined.
// Code adapted from https://www.di-mgt.com.au/euclidean.html
func (x *Int1024) ModInverse(n *Int1024) Int1024 {
	/* Step X1. Initialise */
	lastX := FromUint(1)
	A := *x
	X := Zero()
	N := n.Clone()
	B := N
	/* Remember odd/even iterations */
	iter := 1
	/* Step X2. Loop while B != 0 */

	for !B.IsZero() {
		/* Step X3. Divide and "Subtract" */
		q, r := A.DivMod(&B)
		tmp := q.Mul(&X)
		lastX.Inc(&tmp)
		/* Swap */
		lastX, X, A, B = X, lastX, B, r
		iter ^= 1
	}

	/* Make sure A = gcd(u,v) == 1 */
	if !A.EqualsWord(1) {
		// return zero() /* Error: No inverse exists */
		panic("not relatively prime")
	}
	/* Ensure a positive result */

	inv := lastX
	if iter == 0 {
		N.Dec(&inv)
		return N
		// inv = v.Sub(&inv)
	}
	return inv
}

// Exp returns x**y
func (x *Int1024) Exp(y *Int1024) Int1024 {
	if y.IsZero() {
		return One()
	} else if y.EqualsWord(1) {
		return *(x)
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
