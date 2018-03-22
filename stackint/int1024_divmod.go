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
		x.words[i], r = divWW_g(r, x.words[i], y)
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
	return x.divmodLarge(y)
}

// // q = (x-r)/v, with 0 <= r < y
// // Uses z as storage for q, and u as storage for r if possible.
// // See Knuth, Volume 2, section 4.3.1, Algorithm D.
// // Preconditions:
// //    len(v) >= 2
// //    len(x) >= len(v)
// func (x Int1024) divLarge(y Int1024) (Int1024, Int1024) {
// 	n := y.length
// 	m := x.length - n
// 	// u = u.make(len(x) + 1)
// 	// u.clear() // TODO(gri) no need to clear if we allocated a new u
// 	// D1.
// 	// var v1p *nat
// 	shift := nlz(y.words[n-1])
// 	if shift > 0 {
// 		// do not modify v, it may be used by another goroutine simultaneously
// 		v := y.ShiftLeft(shift)
// 	}
// 	u := Zero()
// 	u[x.length] = shlVU(u[0:len(x)], x, shift)
// 	// D2.
// 	vn1 := v[n-1]
// 	for j := m; j >= 0; j-- {
// 		// D3.
// 		qhat := uint64(_M)
// 		if ujn := u[j+n]; ujn != vn1 {
// 			var rhat Word
// 			qhat, rhat = divWW(ujn, u[j+n-1], vn1)
// 			// x1 | x2 = q̂v_{n-2}
// 			vn2 := v[n-2]
// 			x1, x2 := mulWW(qhat, vn2)
// 			// test if q̂v_{n-2} > br̂ + u_{j+n-2}
// 			ujn2 := u[j+n-2]
// 			for greaterThan(x1, x2, rhat, ujn2) {
// 				qhat--
// 				prevRhat := rhat
// 				rhat += vn1
// 				// v[n-1] >= 0, so this tests for overflow.
// 				if rhat < prevRhat {
// 					break
// 				}
// 				x1, x2 = mulWW(qhat, vn2)
// 			}
// 		}
// 		// D4.
// 		qhatv[n] = mulAddVWW(qhatv[0:n], v, qhat, 0)
// 		c := subVV(u[j:j+len(qhatv)], u[j:], qhatv)
// 		if c != 0 {
// 			c := addVV(u[j:j+n], u[j:], v)
// 			u[j+n] += c
// 			qhat--
// 		}
// 		q[j] = qhat
// 	}
// 	if v1p != nil {
// 		putNat(v1p)
// 	}
// 	putNat(qhatvp)
// 	q = q.norm()
// 	shrVU(u, u, shift)
// 	r = u.norm()
// 	return q, r
// }
