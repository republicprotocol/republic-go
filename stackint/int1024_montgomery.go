package stackint

/*
 * References:
 *
 * https://www.nayuki.io/page/montgomery-reduction-algorithm
 * https://github.com/indutny/bn.js/blob/master/lib/bn.js#L3374
 *
 */

type Montgomery struct {
	m      Int1024
	shift  uint
	mask   Int1024
	r      Int1024
	r2     Int1024
	rInv   Int1024
	factor Int1024
	mInv   Int1024
}

func (mont *Montgomery) ToMont(x *Int1024) MontInt {
	newX := x.ShiftLeft(mont.shift)
	newX = newX.Mod(&mont.m)
	return MontInt{
		newX, mont,
	}
}

func PrimeReduction() *Montgomery {
	m := FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	shift := uint(m.BitLength())
	if shift%WORDSIZE != 0 {
		shift += (WORDSIZE - (shift % WORDSIZE))
	}
	r := one.ShiftLeft(shift)
	// r2 := r.MulModulo(&r, &m)

	// r^2 overflows 2048, so calculate manually
	r2 := FromUint64(11025)
	rinv := r.ModInverse(&m)
	mask := r.Sub(&one)

	// fmt.Printf("Reciprocal: %v\n", rinv.String()) // good!

	factor := r.Mul(&rinv)
	factor.Dec(&one)
	factor = factor.Div(&m)
	// (self.reducer*self.reciprocal - 1) // mod

	// fmt.Printf("Factor: %v\n", factor.String())

	minv := rinv.Mul(&r)
	minv.Dec(&one)
	minv = minv.Div(&m)
	minv = minv.Mod(&r)
	minv = r.Sub(&minv)

	return &Montgomery{
		m:      m,
		shift:  shift,
		mask:   mask,
		r:      r,
		r2:     r2,
		rInv:   rinv,
		mInv:   minv,
		factor: factor,
	}
}

// var PrimeM Montgomery
// var OneC MontInt

var PrimeM = PrimeReduction()
var OneC = PrimeM.One()

type MontInt struct {
	Int1024
	mont *Montgomery
}

func (a *MontInt) ToInt1024() Int1024 {
	tmp := a.Int1024.MulModulo(&a.mont.rInv, &a.mont.m)
	return tmp
}

func (a *MontInt) MontMul(b *MontInt) MontInt {
	mont := a.mont
	if a.IsZero() || b.IsZero() {
		return mont.ToMont(&zero)
	}

	product := a.Int1024.Mul(&b.Int1024)

	temp := product.Mask(mont.shift)
	temp = temp.Mul(&mont.factor)
	temp.MaskInPlace(mont.shift)

	temp2 := temp.Mul(&mont.m)

	productRight := product.Mask(mont.shift)
	temp2Right := temp2.Mask(mont.shift)
	reducedRight := productRight.Overflows(&temp2Right)

	product.ShiftRightInPlace(mont.shift)
	temp2.ShiftRightInPlace(mont.shift)

	product.Inc(&temp2)

	if reducedRight {
		product.Inc(&one)
	}

	if product.GreaterThan(&mont.m) {
		product.Dec(&mont.m)
	}

	return MontInt{
		product, mont,
	}
}

func (a *MontInt) MontInv() MontInt {
	res := a.ModInverse(&a.mont.m)
	res = res.Mul(&a.mont.r2)
	res = res.Mod(&a.mont.m)
	return MontInt{
		res, a.mont,
	}
}

func (x *Int1024) Overflows(y *Int1024) bool {

	var a, b *Int1024
	// var aLen, bLen int
	if x.length > y.length {
		a, b = x, y
		// aLen, bLen = xLen, yLen
	} else {
		a, b = y, x
		// aLen, bLen = yLen, xLen
	}

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
	}
	for i = b.length; i < a.length; i++ {
		previousOverflow := overflow
		if a.words[i] > WORDMAX || a.words[i] > WORDMAX-previousOverflow {
			overflow = 1
		} else {
			overflow = 0
		}
		x.words[i] = a.words[i] + previousOverflow
	}

	x.length = a.length

	return overflow == 1
}

func (x *MontInt) MontClone() MontInt {
	return MontInt{
		x.Int1024.Clone(), x.mont,
	}
}

func (x *MontInt) MontAdd(y *MontInt) MontInt {
	i := x.Int1024.AddModulo(&y.Int1024, &x.mont.m)
	return MontInt{
		i, x.mont,
	}
}

func (x *MontInt) MontSub(y *MontInt) MontInt {
	i := x.Int1024.SubModulo(&y.Int1024, &x.mont.m)
	return MontInt{
		i, x.mont,
	}
}

func (m *Montgomery) One() MontInt {
	one := One()
	return m.ToMont(&one)
}

func (m *Montgomery) FromUint64(x uint64) MontInt {
	tmp := FromUint64(x)
	return m.ToMont(&tmp)
}
