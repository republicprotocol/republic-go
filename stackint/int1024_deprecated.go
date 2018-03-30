package stackint

////////////////////////////////////////////////////
// Code that may be useful but not currently used //
////////////////////////////////////////////////////

/* ******************************************************





MONTGOMERY





****************************************************** */

// /*
//  * References:
//  *
//  * https://www.nayuki.io/page/montgomery-reduction-algorithm
//  * https://github.com/indutny/bn.js/blob/master/lib/bn.js#L3374
//  *
//  */

// type Montgomery struct {
// 	m            Int1024
// 	M            MontInt
// 	shift        uint
// 	mask         Int1024
// 	r            Int1024
// 	r2           Int1024
// 	rInv         Int1024
// 	factor       Int1024
// 	mInv         Int1024
// 	WordLookup map[Word]MontInt
// }

// func (mont *Montgomery) ToMont(x *Int1024) MontInt {
// 	newX := x.ShiftLeft(mont.shift)
// 	newX = newX.Mod(&mont.m)
// 	return MontInt{
// 		newX, mont,
// 	}
// }

// func PrimeReduction() *Montgomery {
// 	m := FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
// 	shift := uint(m.BitLength())
// 	if shift%WORDSIZE != 0 {
// 		shift += (WORDSIZE - (shift % WORDSIZE))
// 	}
// 	r := one.ShiftLeft(shift)
// 	// r2 := r.MulModulo(&r, &m)

// 	// r^2 overflows 2048, so calculate manually
// 	r2 := FromUint(11025)
// 	rinv := r.ModInverse(&m)
// 	mask := r.Sub(&one)

// 	// fmt.Printf("Reciprocal: %v\n", rinv.String()) // good!

// 	factor := r.Mul(&rinv)
// 	factor.Dec(&one)
// 	factor = factor.Div(&m)
// 	// (self.reducer*self.reciprocal - 1) // mod

// 	// fmt.Printf("Factor: %v\n", factor.String())

// 	minv := rinv.Mul(&r)
// 	minv.Dec(&one)
// 	minv = minv.Div(&m)
// 	minv = minv.Mod(&r)
// 	minv = r.Sub(&minv)

// 	mont := &Montgomery{
// 		m:            m,
// 		shift:        shift,
// 		mask:         mask,
// 		r:            r,
// 		r2:           r2,
// 		rInv:         rinv,
// 		mInv:         minv,
// 		factor:       factor,
// 		WordLookup: make(map[Word]MontInt),
// 	}

// 	mont.M = mont.ToMont(&m)

// 	// var i Word
// 	// for i = 0; i < 1000; i++ {
// 	// 	x := FromUint(i)
// 	// 	mont.WordLookup[i] = mont.ToMont(&x)
// 	// }

// 	return mont
// }

// // var PrimeM Montgomery
// // var OneC MontInt

// var PrimeM = PrimeReduction()
// var OneC = PrimeM.One()

// type MontInt struct {
// 	Int1024
// 	mont *Montgomery
// }

// func (a *MontInt) ToInt1024() Int1024 {
// 	tmp := a.Int1024.MulModulo(&a.mont.rInv, &a.mont.m)
// 	return tmp
// }

// func (a *MontInt) MontMul(b *MontInt) MontInt {
// 	mont := a.mont
// 	if a.IsZero() || b.IsZero() {
// 		return mont.ToMont(&zero)
// 	}

// 	product := a.Int1024.Mul(&b.Int1024)

// 	temp := product.Mask(mont.shift)
// 	temp = temp.Mul(&mont.factor)
// 	temp.MaskInPlace(mont.shift)

// 	temp2 := temp.Mul(&mont.m)

// 	if mont.shift >= SIZE/2 {
// 		productRight := product.Mask(mont.shift)
// 		temp2Right := temp2.Mask(mont.shift)
// 		productRight.Inc(&temp2Right)
// 		reducedRight := productRight.words[mont.shift/WORDSIZE] > 0

// 		product.ShiftRightInPlace(mont.shift)
// 		temp2.ShiftRightInPlace(mont.shift)

// 		product.Inc(&temp2)

// 		if reducedRight {
// 			product.Inc(&one)
// 		}
// 	} else {
// 		product.Inc(&temp2)
// 	}

// 	if product.GreaterThan(&mont.m) {
// 		product.Dec(&mont.m)
// 	}

// 	return MontInt{
// 		product, mont,
// 	}
// }

// func (a *MontInt) MontInv() MontInt {
// 	res := a.ModInverse(&a.mont.m)
// 	res = res.MulModulo(&a.mont.r2, &a.mont.m)
// 	return MontInt{
// 		res, a.mont,
// 	}
// }

// func (x *Int1024) Overflows(y *Int1024) bool {

// 	var a, b *Int1024
// 	// var aLen, bLen int
// 	if x.length > y.length {
// 		a, b = x, y
// 		// aLen, bLen = xLen, yLen
// 	} else {
// 		a, b = y, x
// 		// aLen, bLen = yLen, xLen
// 	}

// 	var overflow Word
// 	overflow = 0
// 	var i uint16
// 	for i = 0; i < b.length; i++ {
// 		previousOverflow := overflow
// 		if a.words[i] > WORDMAX-b.words[i] || a.words[i] > WORDMAX-b.words[i]-previousOverflow || b.words[i] > WORDMAX-previousOverflow {
// 			overflow = 1
// 		} else {
// 			overflow = 0
// 		}
// 		x.words[i] = a.words[i] + b.words[i] + previousOverflow
// 	}
// 	if overflow > 0 {
// 		for i = b.length; i < a.length; i++ {
// 			previousOverflow := overflow
// 			if a.words[i] > WORDMAX || a.words[i] > WORDMAX-previousOverflow {
// 				overflow = 1
// 			} else {
// 				overflow = 0
// 			}
// 			x.words[i] = a.words[i] + previousOverflow
// 		}
// 	}

// 	x.length = a.length

// 	ret := (overflow == 1) && a.length == (INT1024WORDS/2)
// 	return ret

// }

// func (x *MontInt) MontClone() MontInt {
// 	return MontInt{
// 		x.Int1024.Clone(), x.mont,
// 	}
// }

// func (x *MontInt) MontAdd(y *MontInt) MontInt {
// 	i := x.Int1024.AddModulo(&y.Int1024, &x.mont.m)
// 	return MontInt{
// 		i, x.mont,
// 	}
// }

// func (x *MontInt) MontSub(y *MontInt) MontInt {
// 	i := x.Int1024.SubModulo(&y.Int1024, &x.mont.m)
// 	return MontInt{
// 		i, x.mont,
// 	}
// }

// func (m *Montgomery) One() MontInt {
// 	one := One()
// 	return m.ToMont(&one)
// }

// func (m *Montgomery) FromUint(x Word) MontInt {
// 	if val, ok := m.WordLookup[x]; ok {
// 		return val.MontClone()
// 	}
// 	tmp := FromUint(x)
// 	ret := m.ToMont(&tmp)
// 	m.WordLookup[x] = ret
// 	return ret
// }

/* ******************************************************





KARATSUBA





****************************************************** */

// const karatsubaThreshold = 1024

// func karatsubaLen(n int) int {
// 	i := uint(0)
// 	for n > karatsubaThreshold {
// 		n >>= 1
// 		i++
// 	}
// 	return n << i
// }

// func (x *Int1024) split(n uint) (Int1024, Int1024) {
// 	b := x.ShiftRight(n)
// 	a := Zero()
// 	var i uint16
// 	len := min(uint16(n), x.length*WORDSIZE)
// 	var firstPositive uint16
// 	for i = 0; i < len/WORDSIZE; i++ {
// 		a.words[i] = x.words[i]
// 		if a.words[i] != 0 {
// 			firstPositive = i
// 		}
// 	}
// 	mod := len % WORDSIZE
// 	if mod != 0 {
// 		a.words[i] = x.words[i] & (1<<mod - 1)
// 		if a.words[i] != 0 {
// 			firstPositive = i
// 		}
// 	}
// 	a.length = firstPositive + 1
// 	return a, b
// }

// // Mul returns x*y
// func (x *Int1024) karatsuba(y *Int1024) Int1024 {

// 	if x.length+y.length <= INT1024WORDS {
// 		return x.BasicMul(y)
// 	}

// 	lenX := x.BitLength()
// 	lenY := y.BitLength()

// 	var n uint
// 	if lenX > lenY {
// 		n = uint(lenX)
// 	} else {
// 		n = uint(lenY)
// 	}

// 	n = (n + 1) / 2

// 	a, b := x.split(n)
// 	c, d := y.split(n)

// 	ac := a.Mul(&c)
// 	bd := b.Mul(&d)
// 	aPlusB := a
// 	aPlusB.Inc(&b)
// 	cPlusD := c
// 	cPlusD.Inc(&d)
// 	abcd := aPlusB.Mul(&cPlusD)
// 	abcd.Dec(&ac)
// 	abcd.Dec(&bd)

// 	res := ac
// 	abcd.ShiftLeftInPlace(n)
// 	res.Inc(&abcd)
// 	bd.ShiftLeftInPlace(2 * n)
// 	res.Inc(&bd)

// 	return res
// }

/* ******************************************************





Russian Peasant's multiplication and Schrage's multiplication





****************************************************** */

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

// // peasant calculates (x*y) % n
// func (x *Int1024) peasant(y, n *Int1024) Int1024 {
// 	// https://stackoverflow.com/questions/12168348/ways-to-do-modulo-multiplication-with-primitive-types

// 	// mul := x.Mul(y)
// 	// return mul.Mod(n)

// 	b := y.Clone()
// 	a := x.Clone()
// 	m := n
// 	z := Zero()
// 	res := z
// 	// Word_t temp_b;

// 	/* Only needed if b may be >= m */
// 	if b.GreaterThanOrEqual(m) {
// 		// halfMax := maxInt()
// 		// halfMax.ShiftRightInPlace()

// 		// Replace with shift right
// 		// two := FromUint(2)
// 		// halfMax := max.Div(&two)

// 		// if m.GreaterThan(&halfMax) {
// 		if m.IsBitSet(SIZE - 1) {
// 			b = b.Sub(m)
// 			// b -= m;
// 		} else {
// 			b = b.Mod(m)
// 			// b %= m;
// 		}
// 	}

// 	for !a.IsZero() {
// 		if !a.IsEven() {
// 			/* Add b to res, modulo m, without overflow */
// 			m.Dec(&res)
// 			if b.GreaterThanOrEqual(m) { /* Equiv to if (res + b >= m), without overflow */
// 				m.Inc(&res)
// 				res.Dec(m)
// 				// res -= m
// 			} else {
// 				m.Inc(&res)
// 			}
// 			res.Inc(&b)
// 			// res += b;
// 		}
// 		a.ShiftRightInPlace(1)
// 		// a >>= 1;

// 		/* Double b, modulo m */
// 		m.Dec(&b)
// 		if b.GreaterThanOrEqual(m) { /* Equiv to if (2 * b >= m), without overflow */
// 			m.Inc(&b)
// 			// temp_b -= m
// 			// tmpB := b.Sub(m)
// 			b.Inc(&b)
// 			b.Dec(m)
// 			// b.Inc(&tmpB)
// 		} else {
// 			m.Inc(&b)
// 			b.Inc(&b)
// 		}
// 		// b += temp_b
// 	}
// 	return res
// }

// // shrage calculates (x*y) % n
// func (x *Int1024) schrage(y, n *Int1024) Int1024 {

// 	// expected := big.NewInt(0).Mul(x.ToBigInt(), y.ToBigInt())
// 	// expected = expected.Mod(expected, n.ToBigInt())

// 	// https://www.thecodingforums.com/threads/extending-schrage-multiplication.558220/

// 	/**

// 		} else {

// 		} // end if
// 		return result;
// 	} // end schrageFast()

// 	*/

// 	b := y.Clone()
// 	a := x.Clone()
// 	m := n
// 	z := Zero()

// 	// 	if (a < 2) { return (a * b) % m; }
// 	if a.LessThan(&two) {
// 		tmp := a.Mul(&b)
// 		return tmp.Mod(m)
// 	}

// 	// int result;
// 	result := z
// 	// Check for b >= m-1
// 	mLess := m.Sub(&one)

// 	switch b.Cmp(&mLess) {
// 	// Check for b >= m-1
// 	// if (b > m - 1) {
// 	case +1:
// 		// result = schrageFast(a, b % m, m);
// 		tmp := b.Mod(m)
// 		result = a.MulModulo(&tmp, m)
// 	// } else if (b == m - 1) {
// 	case 0:
// 		// result = (m - a) % m;
// 		result = m.SubModulo(&a, m)
// 	case -1:
// 		/*

// 			} else {
// 				// Schrage method
// 				result = a * (b % quot) - rem * (b / quot);
// 				if (result < 0) { result += m; }
// 			} // end if
// 		*/

// 		// Check for rem >= quot
// 		// int quot = m / a;
// 		// int rem = m % a;
// 		quot, rem := m.DivMod(&a)
// 		// if (rem >= quot) {
// 		if rem.GreaterThanOrEqual(&quot) {
// 			// result = schrageFast(a/2, b, m);
// 			tmpA := a.ShiftRight(1)
// 			// fmt.Println(a.String())
// 			// fmt.Println("Recursive...")
// 			result = tmpA.MulModulo(&b, m)
// 			// result = addMod(result, result, m);
// 			result = result.AddModulo(&result, m)
// 			// if (a % 2 == 1) { result = addMod(result, b, m); }
// 			if a.IsBitSet(0) {
// 				result = result.AddModulo(&b, m)
// 				// result = addMod(result, b, m)
// 			}
// 		} else {
// 			// Schrage method
// 			// result = a * (b % quot) - rem * (b / quot);
// 			quot2, rem2 := b.DivMod(&quot)
// 			result = a.Mul(&rem2)
// 			tmp3 := rem.Mul(&quot2)
// 			result = result.Sub(&tmp3)
// 			// 	if (result < 0) { result += m; }
// 			if tmp3.GreaterThan(&result) {
// 				result = result.Add(m)
// 			}
// 		} // end if
// 	}

// 	// actual := result.ToBigInt()
// 	// if expected.Cmp(actual) != 0 && expected.BitLen() <= SIZE {
// 	// 	panic(fmt.Sprintf("MulModulo failed!\nFor (%v + %v) mod %v\n.\n\nExp: %b\n\nGot: %b", x, y, n, expected, actual))
// 	// }

// 	return result
// }

/*

Word_t mulmod(Word_t a, Word_t b, Word_t m) {
    Word_t res = 0;
    Word_t temp_b;

    /* Only needed if b may be >= m
    if (b >= m) {
        if (m > WordMAX / 2u)
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

// func (x *Int1024) Mask(n uint) Int1024 {
// 	z := x.Clone()
// 	z.MaskInPlace(n)
// 	return z
// }

// func (x *Int1024) MaskInPlace(n uint) {
// 	nn := uint16(n / WORDSIZE)
// 	if n%WORDSIZE != 0 {
// 		panic("not implemented")
// 	} else {
// 		if x.length > nn {
// 			for i := nn; i < x.length; i++ {
// 				x.words[i] = 0
// 			}
// 			firstPositive := uint16(0)
// 			for i := nn - 1; i > 0; i-- {
// 				if x.words[i] != 0 {
// 					firstPositive = uint16(i)
// 					break
// 				}
// 			}
// 			x.length = firstPositive + 1
// 		}
// 	}
// }
