package stackint

import (
	"github.com/republicprotocol/republic-go/stackint/asm"
)

// divW sets x to x divided by the single-word y and returns x%y
// preconditions:
//   x > y
//   y != 0
func (x *Int1024) divW(y asm.Word) asm.Word {
	if y == 1 {
		return 0
	}
	// m > 0
	// r := DivWVW_g(x.words[:], 0, x.words[:], y)
	r := asm.Word(0)
	var first uint16
	for i := int(x.length - 1); i >= 0; i-- {
		x.words[i], r = asm.DivWW_g(r, x.words[i], y)
		if first == 0 && x.words[i] != 0 {
			first = uint16(i) + 1
		}
	}
	x.length = first
	return r
}

// DivMod returns (x/y, x%y)
func (x *Int1024) DivMod(y *Int1024) (Int1024, Int1024) {
	if y.IsZero() {
		panic("division by zero")
	}
	if x.Cmp(y) < 0 {
		return Zero(), x.Clone()
	}
	if y.length == 1 {
		q := x.Clone()
		rr := q.divW(y.words[0])
		r := FromUint(uint(rr))
		return q, r
	}
	q, r := x.divLarge(y)
	return q, r
}

// greaterThan reports whether (x1<<_W + x2) > (y1<<_W + y2)
func greaterThan(x1, x2, y1, y2 asm.Word) bool {
	return x1 > y1 || x1 == y1 && x2 > y2
}

// divLarge returns (x/y, x%y)
func (x *Int1024) divLarge(y *Int1024) (Int1024, Int1024) {

	v := y.words
	uIn := x.words

	// z is nil
	// u is nil
	// uIn is x
	// v is y

	n := y.length
	m := x.length - n
	// determine if z can be reused
	// TODO(gri) should find a better solution - this if statement
	//           is very costly (see e.g. time pidigits -s -n 10000)

	var q [INT1024WORDS]asm.Word
	var highestQ uint16

	var qhatv [INT1024WORDS + 1]asm.Word
	var u [INT1024WORDS + 1]asm.Word
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

	var rWords [INT1024WORDS]asm.Word
	copy(rWords[:], u[:INT1024WORDS])
	var highestR uint16
	for i := 0; i < int(min(INT1024WORDS, uint16(len(u)))); i++ {
		if rWords[i] != 0 {
			highestR = uint16(i) + 1
		}
	}

	return Int1024{q, highestQ}, Int1024{rWords, highestR}
}
