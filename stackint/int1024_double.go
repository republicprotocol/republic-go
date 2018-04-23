package stackint

import "github.com/republicprotocol/republic-go/stackint/asm"

/**
 * DoubleInt is used to calculate MulModulo without overflowing
 *
 *
 */

// DoubleInt is used when an operation will overflow an Int1024
type DoubleInt struct {
	words  [INT1024WORDS * 2]asm.Word
	length uint16
}

// MulModuloBig will return (x*y)%y, even if (x*y) > Int1024Max
func (x *Int1024) MulModuloBig(y, n *Int1024) Int1024 {

	words := x.BasicMulBig(y)
	var highest uint16
	for i := int(x.length+y.length) - 1; i >= 0; i-- {
		if words[i] > 0 {
			highest = uint16(i) + 1
			break
		}
	}

	xyDouble := DoubleInt{
		words, highest,
	}

	if (highest) <= INT1024WORDS {
		var words2 [INT1024WORDS]asm.Word
		copy(words2[:], words[:INT1024WORDS])
		xy := Int1024{
			words2, highest,
		}
		return xy.Mod(n)
	}
	_, r := xyDouble.divDouble(n)
	var words2 [INT1024WORDS]asm.Word
	copy(words2[:], r.words[:INT1024WORDS])
	return Int1024{
		words2, r.length,
	}
}

// divDouble returns (x/y, x%y) for when x is a DoubleInt
// Preconditions:
//    len(v) >= 2,
func (x *DoubleInt) divDouble(y *Int1024) (DoubleInt, DoubleInt) {

	v := y.words
	uIn := x.words

	n := y.length
	if n < 2 {
		panic("y too small for divDouble")
	}
	m := x.length - n

	var q [INT1024WORDS * 2]asm.Word
	var highestQ uint16

	var qhatv [(INT1024WORDS * 2) + 1]asm.Word
	var u [(INT1024WORDS * 2) + 1]asm.Word
	// D1.
	shift := asm.Nlz(v[n-1])
	if shift > 0 {
		// do not modify v, it may be used by another goroutine simultaneously

		var v1 [INT1024WORDS]asm.Word
		asm.ShlVU_g(v1[:], v[:], shift)
		v = v1
	}
	u[int(x.length)] = asm.ShlVU_g(u[:x.length], uIn[:], shift)
	// D2.
	vn1 := v[n-1]
	for jj := int(m); jj >= 0; jj-- {
		j := uint16(jj)
		// D3.
		qhat := asm.Word(asm.M)
		if ujn := u[j+n]; ujn != vn1 {
			var rhat asm.Word
			qhat, rhat = asm.DivWW_g(ujn, u[j+n-1], vn1)
			// x1 | x2 = q̂v_{n-2}
			vn2 := v[n-2]
			x1, x2 := asm.MulWW(qhat, vn2)
			// test if q̂v_{n-2} > br̂ + u_{j+n-2}
			ujn2 := u[j+n-2]
			for greaterThan(x1, x2, rhat, ujn2) {
				qhat--
				prevRhat := rhat
				rhat += vn1
				// v[n-1] >= 0, so this tests for overflow.
				if rhat < prevRhat {
					break
				}
				x1, x2 = asm.MulWW(qhat, vn2)
			}
		}
		// D4.

		// Inlined
		// qhatv[n] = MulAddVWW(qhatv[0:n], v[:], qhat, 0)
		c := asm.Word(0)
		var i uint16
		for i = 0; i < n; i++ {
			c, qhatv[i] = asm.MulAddWWW_g(v[i], qhat, c)
		}
		qhatv[n] = c

		// Inlined
		c = asm.SubVV_g(u[j:j+n+1], u[j:], qhatv[:])

		if c != 0 {
			// Inlined
			c := asm.AddVV_g(u[j:j+n], u[j:], v[:])

			u[j+n] += c
			qhat--
		}
		q[j] = qhat
		if (j+1) > highestQ && q[j] > 0 {
			highestQ = j + 1
		}
	}
	asm.ShrVU_g(u[:x.length+1], u[:uint(x.length)+1], shift)

	var rWords [INT1024WORDS * 2]asm.Word

	copy(rWords[:], u[:INT1024WORDS*2])

	var highestR uint16
	for i := 0; i < int(min(INT1024WORDS*2, uint16(len(u)))); i++ {
		if rWords[i] != 0 {
			highestR = uint16(i) + 1
		}
	}

	return DoubleInt{q, highestQ}, DoubleInt{rWords, highestR}
}
