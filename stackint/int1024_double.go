package stackint

type DoubleInt struct {
	words  [INT1024WORDS * 2]uint64
	length uint16
}

func (x *Int1024) MulModuloBig(y, n *Int1024) Int1024 {

	words := x.BasicMulBig(y)
	var highest uint16
	var i uint16
	for i = x.length + y.length - 1; i > 0; i-- {
		if words[i] > 0 {
			highest = i
			break
		}
	}

	xyDouble := DoubleInt{
		words, highest + 1,
	}

	if (highest + 1) <= INT1024WORDS {
		var words2 [INT1024WORDS]uint64
		copy(words2[:], words[:INT1024WORDS])
		xy := Int1024{
			words2, highest + 1,
		}
		return xy.Mod(n)
	}
	_, r := xyDouble.divDouble(n)
	var words2 [INT1024WORDS]uint64
	copy(words2[:], r.words[:INT1024WORDS])
	return Int1024{
		words2, r.length,
	}
}

func (x *DoubleInt) divDouble(y *Int1024) (DoubleInt, DoubleInt) {

	v := y.words
	uIn := x.words

	n := y.length
	m := x.length - n

	var q [INT1024WORDS * 2]uint64
	var highestQ uint16

	var qhatv [(INT1024WORDS * 2) + 1]uint64
	var u [(INT1024WORDS * 2) + 1]uint64
	// D1.
	shift := nlz(v[n-1])
	if shift > 0 {
		// do not modify v, it may be used by another goroutine simultaneously

		var v1 [INT1024WORDS]uint64
		shlVU_g(v1[:], v[:], shift)
		v = v1
	}
	u[int(x.length)] = shlVU_g(u[:x.length], uIn[:], shift)
	// D2.
	vn1 := v[n-1]
	for jj := int(m); jj >= 0; jj-- {
		j := uint16(jj)
		// D3.
		qhat := uint64(_M)
		if ujn := u[j+n]; ujn != vn1 {
			var rhat uint64
			qhat, rhat = divWW(ujn, u[j+n-1], vn1)
			// x1 | x2 = q̂v_{n-2}
			vn2 := v[n-2]
			x1, x2 := mulWW(qhat, vn2)
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
				x1, x2 = mulWW(qhat, vn2)
			}
		}
		// D4.

		// Inlined
		// qhatv[n] = mulAddVWW(qhatv[0:n], v[:], qhat, 0)
		c := uint64(0)
		var i uint16
		for i = 0; i < n; i++ {
			c, qhatv[i] = mulAddWWW_g(v[i], qhat, c)
		}
		qhatv[n] = c

		// Inlined
		c = subVV_g(u[j:j+n+1], u[j:], qhatv[:])

		if c != 0 {
			// Inlined
			c := addVV_g(u[j:j+n], u[j:], v[:])

			u[j+n] += c
			qhat--
		}
		q[j] = qhat
		if j > highestQ && q[j] > 0 {
			highestQ = j
		}
	}
	shrVU_g(u[:x.length+1], u[:uint(x.length)+1], shift)

	var rWords [INT1024WORDS * 2]uint64

	copy(rWords[:], u[:INT1024WORDS*2])

	var highestR uint16
	for i := 0; i < int(min(INT1024WORDS*2, uint16(len(u)))); i++ {
		if rWords[i] != 0 {
			highestR = uint16(i)
		}
	}

	return DoubleInt{q, highestQ + 1}, DoubleInt{rWords, highestR + 1}
}
