package stackint

// go build -a -gcflags='-m -m' stackint.go

// INT1024WORDS is 1024 / 64
const INT1024WORDS = 16

// MAXUINT64 represents the largest uint64 value
const MAXUINT64 = 1<<64 - 1

// Int1024 provides a 1024 bit number optimised to never use the heap
type Int1024 struct {
	words [INT1024WORDS]uint64
}

// Zero returns a new Int1024 representing 0
func Zero() Int1024 {
	var words [INT1024WORDS]uint64
	for i := 0; i < INT1024WORDS; i++ {
		words[i] = 0
	}
	return Int1024{
		words: words,
	}
}

// Int1024FromUint64 returns a new Int1024 from a uint64
func Int1024FromUint64(n uint64) Int1024 {
	z := Zero()
	z.words[INT1024WORDS-1] = n
	return z
}

// Int1024FromString returns a new Int1024 from a string
func Int1024FromString(x string) Int1024 {
	panic("Not implemented!")
}

// Clone returns a new Int1024 representing the same value as x
func (x *Int1024) Clone() Int1024 {
	var words [16]uint64
	for i := 0; i < INT1024WORDS; i++ {
		words[i] = x.words[i]
	}
	return Int1024{
		words: words,
	}
}

// Equals returns true of x and y represent the same Int1024
func (x *Int1024) Equals(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] != y.words[i] {
			return false
		}
	}
	return true
}

// Add returns x+y
func (x *Int1024) Add(y *Int1024) Int1024 {
	z := Zero()

	var overflow uint64
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		z.words[i] = x.words[i] + y.words[i] + overflow
		if x.words[i] > MAXUINT64-y.words[i]-overflow {
			overflow = 1
		} else {
			overflow = 0
		}
	}

	return z
}

// Sub returns x-y
func (x *Int1024) Sub(y *Int1024) Int1024 {
	panic("Not implemented!")
}

// Mul returns x*y
func (x *Int1024) Mul(y *Int1024) Int1024 {
	panic("Not implemented!")
}

// Div returns the quotient of x/y. If y is 0, a run-time panic occurs.
func (x *Int1024) Div(y *Int1024) Int1024 {
	panic("Not implemented!")
}

// Mod returns the modulus x%y. If y is 0, a run-time panic occurs.
func (x *Int1024) Mod(y *Int1024) Int1024 {
	panic("Not implemented!")
}

// ModInverse returns the multiplicative inverse of x in the ring ℤ/nℤ.
// If x and n are not relatively prime, the result in undefined.
func (x *Int1024) ModInverse(n *Int1024) Int1024 {
	panic("Not implemented!")
}

// Words returns the internal [16]uint64 that stores the value of x
func (x *Int1024) Words() [INT1024WORDS]uint64 {
	return x.words
}
