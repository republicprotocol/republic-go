package stackint

// Add returns x+y
func (x *Int1024) Add(y *Int1024) Int1024 {

	z := x.Clone()
	z.Inc(y)

	return z
}

// Inc sets x to x+y (used for multiplication)
func (x *Int1024) Inc(y *Int1024) {

	// expected := big.NewInt(0).Add(x.ToBigInt(), y.ToBigInt())

	// a, b := x, y

	var a, b *Int1024
	// var aLen, bLen int
	if x.length > y.length {
		a, b = x, y
		// aLen, bLen = xLen, yLen
	} else {
		a, b = y, x
		// aLen, bLen = yLen, xLen
	}

	var firstPositive uint16
	var overflow uint64
	overflow = 0
	var i uint16
	for i = 0; i < b.length; i++ {
		previousOverflow := overflow
		if a.words[i] > WORDMAX-b.words[i] || a.words[i] > WORDMAX-b.words[i]-previousOverflow || b.words[i] > WORDMAX-previousOverflow {
			overflow = 1
		} else {
			overflow = 0
		}
		x.words[i] = a.words[i] + b.words[i] + previousOverflow
		if x.words[i] != 0 {
			firstPositive = i
		}
	}
	for i = b.length; i < a.length; i++ {
		previousOverflow := overflow
		if a.words[i] > WORDMAX || a.words[i] > WORDMAX-previousOverflow {
			overflow = 1
		} else {
			overflow = 0
		}
		x.words[i] = a.words[i] + previousOverflow
		if x.words[i] != 0 {
			firstPositive = i
		}
	}

	x.length = a.length

	if overflow == 1 && x.length < INT1024WORDS {
		x.length++
		x.words[x.length-1] = 1
		firstPositive = x.length - 1
		// WARNING: Overflow occured
	}

	x.length = firstPositive + 1

	// if expected.Cmp(x.ToBigInt()) != 0 && expected.BitLen() <= SIZE {
	// 	panic("Addition failed")
	// }
}

// Sub returns x-y
func (x *Int1024) Sub(y *Int1024) Int1024 {

	z := x.Clone()
	z.Dec(y)

	return z
}

// Dec sets x to x-y
func (x *Int1024) Dec(y *Int1024) {

	// expected := big.NewInt(0).Sub(x.ToBigInt(), y.ToBigInt())

	var overflow uint64
	var i uint16
	var firstPositive uint16
	if x.length >= y.length {
		for i = 0; i < y.length; i++ {
			newOverflow := uint64(0)
			if x.words[i] < y.words[i]+overflow || y.words[i] == WORDMAX && overflow == 1 {
				newOverflow = 1
			}
			x.words[i] = x.words[i] - y.words[i] - overflow
			overflow = newOverflow
			if x.words[i] != 0 {
				firstPositive = i
			}
		}
		for i = y.length; i < x.length; i++ {
			newOverflow := uint64(0)
			if x.words[i] < overflow {
				newOverflow = 1
			}
			x.words[i] = x.words[i] - overflow
			overflow = newOverflow
			if x.words[i] != 0 {
				firstPositive = i
			}
		}
	} else {
		for i = 0; i < x.length; i++ {
			newOverflow := uint64(0)
			if x.words[i] < y.words[i]+overflow || y.words[i] == WORDMAX && overflow == 1 {
				newOverflow = 1
			}
			x.words[i] = x.words[i] - y.words[i] - overflow
			overflow = newOverflow
			if x.words[i] != 0 {
				firstPositive = i
			}
		}
		for i = x.length; i < y.length; i++ {
			newOverflow := uint64(0)
			if 0 < y.words[i]+overflow || y.words[i] == WORDMAX && overflow == 1 {
				newOverflow = 1
			}
			x.words[i] = 0 - y.words[i] - overflow
			overflow = newOverflow
			if x.words[i] != 0 {
				firstPositive = i
			}
		}
	}

	if overflow > 0 {
		for ; i < INT1024WORDS; i++ {
			x.words[i] = (1 << 64) - 1
		}
		firstPositive = INT1024WORDS - 1
	}
	x.length = firstPositive + 1

	// if expected.Cmp(x.ToBigInt()) != 0 && expected.Cmp(big.NewInt(0)) >= 0 {
	// 	panic("Subtraction failed")
	// }
}

func tmp(zz, xx []uint64, yy uint64) uint64 {
	ll := len(zz)
	var c uint64
	for i := 0; i < ll; i++ {
		z1, z0 := mulAddWWW_g(xx[i], yy, zz[i])
		c, zz[i] = addWW_g(z0, c, 0)
		c += z1
	}
	return c
}

// BasicMul returns x*y using the shift and add method
func (x *Int1024) BasicMul(y *Int1024) Int1024 {

	// words := make([]uint64, x.length+y.length)
	var words [INT1024WORDS]uint64
	var i uint16
	var j uint16
	var highest uint16
	l := uint16(x.length)
	for i = 0; i < y.length; i++ {
		d := y.words[i]
		if d != 0 {
			var c uint64
			for j = i; j < i+l; j++ {
				var z0, z1 uint64
				z1, zz0 := mulWW(x.words[j-i], d)
				if z0 = zz0 + words[j]; z0 < zz0 {
					z1++
				}
				c, words[j] = addWW_g(z0, c, 0)
				if words[j] != 0 {
					highest = max(highest, j)
				}
				c += z1
			}
			words[l+i] = c
			if words[l+i] != 0 {
				highest = max(highest, l+i)
			}
		}
	}
	return Int1024{
		words, highest + 1,
	}
}

// Mul returns x*y
func (x *Int1024) Mul(y *Int1024) Int1024 {

	z := x.karatsuba(y)

	return z
}

const karatsubaThreshold = 1024

func karatsubaLen(n int) int {
	i := uint(0)
	for n > karatsubaThreshold {
		n >>= 1
		i++
	}
	return n << i
}

func (x *Int1024) split(n uint) (Int1024, Int1024) {
	b := x.ShiftRight(n)
	a := Zero()
	var i uint16
	len := min(uint16(n), x.length*WORDSIZE)
	var firstPositive uint16
	for i = 0; i < len/WORDSIZE; i++ {
		a.words[i] = x.words[i]
		if a.words[i] != 0 {
			firstPositive = i
		}
	}
	mod := len % WORDSIZE
	if mod != 0 {
		a.words[i] = x.words[i] & (1<<mod - 1)
		if a.words[i] != 0 {
			firstPositive = i
		}
	}
	a.length = firstPositive + 1
	return a, b
}

// Mul returns x*y
func (x *Int1024) karatsuba(y *Int1024) Int1024 {

	lenX := x.BitLength()
	lenY := y.BitLength()
	if lenX < karatsubaThreshold || lenY < karatsubaThreshold {
		return x.BasicMul(y)
	}

	var n uint
	if lenX > lenY {
		n = uint(lenX)
	} else {
		n = uint(lenY)
	}

	n = (n + 1) / 2

	a, b := x.split(n)
	c, d := y.split(n)

	ac := a.Mul(&c)
	bd := b.Mul(&d)
	aPlusB := a
	aPlusB.Inc(&b)
	cPlusD := c
	cPlusD.Inc(&d)
	abcd := aPlusB.Mul(&cPlusD)
	abcd.Dec(&ac)
	abcd.Dec(&bd)

	res := ac
	abcd.ShiftLeftInPlace(n)
	res.Inc(&abcd)
	bd.ShiftLeftInPlace(2 * n)
	res.Inc(&bd)

	return res
}

// DivMod returns (x/y, x%y). If y is 0, a run-time panic occurs.
func (x *Int1024) divmodLarge(y *Int1024) (Int1024, Int1024) {

	// expected1, expected2 := big.NewInt(0).DivMod(x.ToBigInt(), y.ToBigInt(), big.NewInt(1))

	dividend := x.Clone()
	denom := y.Clone()
	current := FromUint64(1)
	answer := Zero()

	if denom.IsZero() {
		panic("division by zero")
	}

	limit := MAXINT1024()
	limit.ShiftRightInPlace(1)
	overflowed := false
	for denom.LessThanOrEqual(&dividend) {
		if !denom.LessThan(&limit) {
			overflowed = true
			break
		}
		denom.ShiftLeftInPlace(1)
		current.ShiftLeftInPlace(1)
	}

	if !overflowed {
		denom.ShiftRightInPlace(1)
		current.ShiftRightInPlace(1)
	}

	for !current.IsZero() {
		if dividend.GreaterThanOrEqual(&denom) {
			dividend.Dec(&denom)
			answer.ORInPlace(&current)
		}
		current.ShiftRightInPlace(1)
		denom.ShiftRightInPlace(1)
	}

	// actual1 := answer.ToBigInt()
	// if expected1.Cmp(actual1) != 0 {
	// 	panic(fmt.Sprintf("AddModulo failed!\nFor (%v / %v)\n.\n\nExp: %b\n\nGot: %b", x, y, expected1, actual1))
	// }

	// actual2 := dividend.ToBigInt()
	// if expected2.Cmp(actual2) != 0 {
	// 	panic(fmt.Sprintf("AddModulo failed!\nFor (%v mod %v)\n.\n\nExp: %b\n\nGot: %b", x, y, expected2, actual2))
	// }

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
	switch x.Cmp(n) {
	case -1:
		return x.Clone()
	case 0:
		return Zero()
	case 1:
		mod := x.mod(n)
		return mod
	default:
		panic("unexpected cmp result (expecting -1, 0 or 1)")
	}
}

func (x *Int1024) mod(y *Int1024) Int1024 {

	dividend := x.Clone()
	denom := y.Clone()
	current := 1

	if denom.IsZero() {
		panic("division by zero")
	}

	limit := MAXINT1024()
	limit.ShiftRightInPlace(1)
	overflowed := false

	shift := dividend.BitLength() - denom.BitLength()
	if shift < 0 {
		shift = 0
	}
	denom.ShiftLeftInPlace(uint(shift))
	current += shift

	if denom.LessThanOrEqual(&dividend) {
		if denom.GreaterThanOrEqual(&limit) {
			overflowed = true
		} else {
			denom.ShiftLeftInPlace(1)
			current++
		}
	}

	if !overflowed {
		denom.ShiftRightInPlace(1)
		current--
	}

	for current != 0 {
		if dividend.GreaterThanOrEqual(&denom) {
			dividend.Dec(&denom)
		}
		current--
		denom.ShiftRightInPlace(1)
	}

	return dividend
}

// SubModulo returns (x - y) % n
func (x *Int1024) SubModulo(y, n *Int1024) Int1024 {
	// expected := big.NewInt(0).Sub(x.ToBigInt(), y.ToBigInt())
	// expected = expected.Mod(expected, n.ToBigInt())

	z := x.subModulo(y, n)

	// actual := z.ToBigInt()
	// if expected.Cmp(actual) != 0 && expected.BitLen() <= SIZE {
	// 	panic(fmt.Sprintf("SubModulo failed!\nFor (%v + %v) mod %v\n.\n\nExp: %b\n\nGot: %b", x, y, z, expected, actual))
	// }

	return z
}

func (x *Int1024) subModulo(y, n *Int1024) Int1024 {
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

// AddModulo returns (x + y) % n (x+y can be larger than 2^SIZE)
func (x *Int1024) AddModulo(y, n *Int1024) Int1024 {

	// expected := big.NewInt(0).Add(x.ToBigInt(), y.ToBigInt())
	// expected = expected.Mod(expected, n.ToBigInt())

	z := x.addModulo(y, n)

	// actual := z.ToBigInt()
	// if expected.Cmp(actual) != 0 && expected.BitLen() <= SIZE {
	// 	panic(fmt.Sprintf("AddModulo failed!\nFor (%v + %v) mod %v\n.\n\nExp: %b\n\nGot: %b", x, y, z, expected, actual))
	// }
	return z
}

func (x *Int1024) addModulo(y, n *Int1024) Int1024 {

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

// // MulModuloSlow returns (x*y) % n
func (x *Int1024) MulModulo(y, n *Int1024) Int1024 {
	z := x.Mul(y)
	return z.Mod(n)
	// return x.peasant(y, n)
}

// // MulModuloSlow returns (x*y) % n but it takes its time
// func (x *Int1024) MulModulo(y, n *Int1024) Int1024 {

// 	expected := big.NewInt(0).Mul(x.ToBigInt(), y.ToBigInt())
// 	expected = expected.Mod(expected, n.ToBigInt())

// 	// Not optimized

// 	// TODO!!!
// 	// Implement Shrage's Method
// 	// https://github.com/wdavidw/node-gsl/blob/master/deps/gsl-1.14/rng/schrage.c

// 	z := Zero()

// 	if x.length < y.length {
// 		x, y = y, x
// 	}

// 	shifted := x.Clone()

// 	var i uint16
// 	for i = 0; i < y.length; i++ {
// 		word := y.words[i]
// 		for j := uint(0); j < WORDSIZE; j++ {
// 			bit := (word >> j) & 1
// 			if bit == 1 {
// 				z = z.AddModulo(&shifted, n)
// 			}
// 			shifted = shifted.AddModulo(&shifted, n)
// 		}
// 	}

// 	actual := z.ToBigInt()
// 	if expected.Cmp(actual) != 0 && expected.BitLen() <= SIZE {
// 		panic(fmt.Sprintf("SubModulo failed!\nFor (%v + %v) mod %v\n.\n\nExp: %b\n\nGot: %b", x, y, n, expected, actual))
// 	}

// 	return z
// }

// peasant calculates (x*y) % n
func (x *Int1024) peasant(y, n *Int1024) Int1024 {
	// https://stackoverflow.com/questions/12168348/ways-to-do-modulo-multiplication-with-primitive-types

	// mul := x.Mul(y)
	// return mul.Mod(n)

	b := y.Clone()
	a := x.Clone()
	m := n
	z := Zero()
	res := z
	// uint64_t temp_b;

	/* Only needed if b may be >= m */
	if b.GreaterThanOrEqual(m) {
		// halfMax := maxInt()
		// halfMax.ShiftRightInPlace()

		// Replace with shift right
		// two := FromUint64(2)
		// halfMax := max.Div(&two)

		// if m.GreaterThan(&halfMax) {
		if m.IsBitSet(SIZE - 1) {
			b = b.Sub(m)
			// b -= m;
		} else {
			b = b.Mod(m)
			// b %= m;
		}
	}

	for !a.IsZero() {
		if !a.IsEven() {
			/* Add b to res, modulo m, without overflow */
			m.Dec(&res)
			if b.GreaterThanOrEqual(m) { /* Equiv to if (res + b >= m), without overflow */
				m.Inc(&res)
				res.Dec(m)
				// res -= m
			} else {
				m.Inc(&res)
			}
			res.Inc(&b)
			// res += b;
		}
		a.ShiftRightInPlace(1)
		// a >>= 1;

		/* Double b, modulo m */
		m.Dec(&b)
		if b.GreaterThanOrEqual(m) { /* Equiv to if (2 * b >= m), without overflow */
			m.Inc(&b)
			// temp_b -= m
			// tmpB := b.Sub(m)
			b.Inc(&b)
			b.Dec(m)
			// b.Inc(&tmpB)
		} else {
			m.Inc(&b)
			b.Inc(&b)
		}
		// b += temp_b
	}
	return res
}

// shrage calculates (x*y) % n
func (x *Int1024) schrage(y, n *Int1024) Int1024 {

	// expected := big.NewInt(0).Mul(x.ToBigInt(), y.ToBigInt())
	// expected = expected.Mod(expected, n.ToBigInt())

	// https://www.thecodingforums.com/threads/extending-schrage-multiplication.558220/

	/**

		} else {

		} // end if
		return result;
	} // end schrageFast()

	*/

	b := y.Clone()
	a := x.Clone()
	m := n
	z := Zero()

	// 	if (a < 2) { return (a * b) % m; }
	if a.LessThan(&two) {
		tmp := a.Mul(&b)
		return tmp.Mod(m)
	}

	// int result;
	result := z
	// Check for b >= m-1
	mLess := m.Sub(&one)

	switch b.Cmp(&mLess) {
	// Check for b >= m-1
	// if (b > m - 1) {
	case +1:
		// result = schrageFast(a, b % m, m);
		tmp := b.Mod(m)
		result = a.MulModulo(&tmp, m)
	// } else if (b == m - 1) {
	case 0:
		// result = (m - a) % m;
		result = m.SubModulo(&a, m)
	case -1:
		/*

			} else {
				// Schrage method
				result = a * (b % quot) - rem * (b / quot);
				if (result < 0) { result += m; }
			} // end if
		*/

		// Check for rem >= quot
		// int quot = m / a;
		// int rem = m % a;
		quot, rem := m.DivMod(&a)
		// if (rem >= quot) {
		if rem.GreaterThanOrEqual(&quot) {
			// result = schrageFast(a/2, b, m);
			tmpA := a.ShiftRight(1)
			// fmt.Println(a.String())
			// fmt.Println("Recursive...")
			result = tmpA.MulModulo(&b, m)
			// result = addMod(result, result, m);
			result = result.AddModulo(&result, m)
			// if (a % 2 == 1) { result = addMod(result, b, m); }
			if a.IsBitSet(0) {
				result = result.AddModulo(&b, m)
				// result = addMod(result, b, m)
			}
		} else {
			// Schrage method
			// result = a * (b % quot) - rem * (b / quot);
			quot2, rem2 := b.DivMod(&quot)
			result = a.Mul(&rem2)
			tmp3 := rem.Mul(&quot2)
			result = result.Sub(&tmp3)
			// 	if (result < 0) { result += m; }
			if tmp3.GreaterThan(&result) {
				result = result.Add(m)
			}
		} // end if
	}

	// actual := result.ToBigInt()
	// if expected.Cmp(actual) != 0 && expected.BitLen() <= SIZE {
	// 	panic(fmt.Sprintf("MulModulo failed!\nFor (%v + %v) mod %v\n.\n\nExp: %b\n\nGot: %b", x, y, n, expected, actual))
	// }

	return result
}

/*

uint64_t mulmod(uint64_t a, uint64_t b, uint64_t m) {
    uint64_t res = 0;
    uint64_t temp_b;

    /* Only needed if b may be >= m
    if (b >= m) {
        if (m > UINT64_MAX / 2u)
            b -= m;
        else
            b %= m;
    }

    while (a != 0) {
        if (a & 1) {
            /* Add b to res, modulo m, without overflow
            if (b >= m - res) /* Equiv to if (res + b >= m), without overflow
                res -= m;
            res += b;
        }
        a >>= 1;

        /* Double b, modulo m
        temp_b = b;
        if (b >= m - b)       /* Equiv to if (2 * b >= m), without overflow
            temp_b -= m;
        b += temp_b;
    }
    return res;
}
*/

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
