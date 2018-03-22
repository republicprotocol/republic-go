package stackint

func (x *Int1024) divW(y uint64) uint64 {
	m := len(x.words)
	switch {
	case y == 0:
		panic("division by zero")
	case y == 1:
		return 0
	case m == 0:
		x.SetUint64(0)
		return 0
	}
	// m > 0
	// r := divWVW_g(x.words[:], 0, x.words[:], y)
	r := uint64(0)
	for i := len(x.words) - 1; i >= 0; i-- {
		x.words[i], r = divWW(r, x.words[i], y)
	}
	return r
}

func (x *Int1024) DivMod(y *Int1024) (Int1024, Int1024) {
	if len(y.words) == 0 {
		panic("division by zero")
	}
	if x.Cmp(y) < 0 {
		return Zero(), x.Clone()
	}
	if y.length == 1 {
		q := x.Clone()
		rr := q.divW(y.words[0])
		r := FromUint64(rr)
		return q, r
	}
	return x.divLarge(y)
}

// greaterThan reports whether (x1<<_W + x2) > (y1<<_W + y2)
func greaterThan(x1, x2, y1, y2 uint64) bool {
	return x1 > y1 || x1 == y1 && x2 > y2
}

func (x *Int1024) divLarge(y *Int1024) (qq, r Int1024) {

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

	var q [INT1024WORDS]uint64
	var highestQ uint16

	var qhatv [INT1024WORDS + 1]uint64
	var u [INT1024WORDS + 1]uint64
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

	var rWords [INT1024WORDS]uint64

	var highestR uint16
	for i := 0; i < int(min(INT1024WORDS, uint16(len(u)))); i++ {
		rWords[i] = u[i]
		if rWords[i] != 0 {
			highestR = uint16(i)
		}
	}

	return Int1024{q, highestQ + 1}, Int1024{rWords, highestR + 1}
}
