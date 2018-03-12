package stackint

// go build -a -gcflags='-m -m' stackint.go

// SIZE is the number of bits stored by Int1024
const SIZE = 1024

// WORDSIZE is 64 for Word
const WORDSIZE = 64

// Word is the internal type
type Word uint64

// WORDMAX represents the largest word value
const WORDMAX = 1<<WORDSIZE - 1

// INT1024WORDS is 1024 / 64
const INT1024WORDS = SIZE / WORDSIZE

// Int1024 provides a 1024 bit number optimised to never use the heap
type Int1024 struct {
	words [INT1024WORDS]Word
}

// Zero returns a new Int1024 representing 0
func Zero() Int1024 {
	var words [INT1024WORDS]Word
	for i := 0; i < INT1024WORDS; i++ {
		words[i] = 0
	}
	return Int1024{
		words: words,
	}
}

var zero = Zero()

// Int1024FromUint64 returns a new Int1024 from a Word
func Int1024FromUint64(n uint64) Int1024 {
	z := Zero()
	z.words[INT1024WORDS-1] = Word(n)
	return z
}

// Int1024FromString returns a new Int1024 from a string
func Int1024FromString(x string) Int1024 {
	panic("Not implemented!")
}

// Clone returns a new Int1024 representing the same value as x
func (x *Int1024) Clone() Int1024 {
	var words [16]Word
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

// LessThan returns x<y
func (x *Int1024) LessThan(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] < y.words[i] {
			return true
		}
	}
	return false
}

// LessThanOrEqual returns x<=y
func (x *Int1024) LessThanOrEqual(y *Int1024) bool {
	for i := 0; i < INT1024WORDS; i++ {
		if x.words[i] > y.words[i] {
			return false
		}
	}
	return true
}

// shiftLeftByOne shifts to the left x by one
func (x *Int1024) shiftLeftByOne() Int1024 {
	z := Zero()
	overflow := Word(0)
	for i := INT1024WORDS - 1; i >= 0; i-- {
		// Shift word to the right
		// If previous word overflowed, add 1
		z.words[i] = (x.words[i] << 1) | overflow
		// Calculate if word overflows
		overflow = (x.words[i] >> (WORDSIZE - 1)) & 1
	}

	if overflow == 1 {
		// WARNING: Overflow occured (not important for Shift)
	}

	return z
}

// shiftRightByOne shifts to the right x by one
func (x *Int1024) shiftRightByOne() Int1024 {
	z := Zero()
	overflow := Word(0)
	for i := 0; i < INT1024WORDS; i++ {
		// Shift word to the right
		// If previous word overflowed, add 1
		z.words[i] = (x.words[i] >> 1) | overflow
		// Calculate if word overflows
		overflow = (x.words[i] >> (WORDSIZE - 1)) & (1 << (WORDSIZE - 1))
	}

	if overflow == 1 {
		// WARNING: Overflow occured (not important for Shift)
	}

	return z
}

// Or returns x|y
func (x *Int1024) Or(y *Int1024) Int1024 {
	z := Zero()
	for i := 0; i < INT1024WORDS; i++ {
		z.words[i] = x.words[i] | y.words[i]
	}
	return z
}

// Not returns ~x
func (x *Int1024) Not() Int1024 {
	z := Zero()
	for i := 0; i < INT1024WORDS; i++ {
		z.words[i] = ^x.words[i]
	}
	return z
}

// Add returns x+y
func (x *Int1024) Add(y *Int1024) Int1024 {
	z := Zero()

	// Loop through each of the 16 pair of words and add them
	// Storing the overflow if necessary
	var overflow Word
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		z.words[i] = x.words[i] + y.words[i] + overflow
		if x.words[i] > WORDMAX-y.words[i]-overflow {
			overflow = 1
		} else {
			overflow = 0
		}
	}

	if overflow == 1 {
		// WARNING: Overflow occured
	}

	return z
}

// overwritingAdd sets x to x+y (used for multiplication)
func (x *Int1024) overwritingAdd(y *Int1024) {
	var overflow Word
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		previousOverflow := overflow
		if x.words[i] > WORDMAX-y.words[i]-previousOverflow {
			overflow = 1
		} else {
			overflow = 0
		}
		x.words[i] = x.words[i] + y.words[i] + previousOverflow
	}

	if overflow == 1 {
		// WARNING: Overflow occured
	}
}

// Sub returns x-y
func (x *Int1024) Sub(y *Int1024) Int1024 {
	z := Zero()

	var overflow Word
	overflow = 0
	for i := INT1024WORDS - 1; i >= 0; i-- {
		z.words[i] = x.words[i] - y.words[i] - overflow
		if x.words[i] < y.words[i]+overflow {
			overflow = 1
		} else {
			overflow = 0
		}
	}

	if overflow == 1 {
		// WARNING: Overflow occured
	}

	return z
}

// Mul returns x*y
func (x *Int1024) Mul(y *Int1024) Int1024 {
	// Naïve inplementation!
	// Uses up to 16384 uint64 additions (worst case)
	// TODO: Rewrite using more efficient algorithm
	z := Zero()
	shifted := *x

	for i := INT1024WORDS - 1; i >= 0; i-- {
		word := y.words[i]
		for j := uint(0); j < WORDSIZE; j++ {
			bit := (word >> j) & 1
			if bit == 1 {
				z.overwritingAdd(&shifted)
			}
			shifted = shifted.shiftLeftByOne()
		}
	}

	return z
}

// DivMod returns (x/y, x%y). If y is 0, a run-time panic occurs.
func (x *Int1024) DivMod(y *Int1024) (Int1024, Int1024) {
	dividend := x.Clone()
	denom := y.Clone()
	current := Int1024FromUint64(1)
	answer := Zero()

	if denom.Equals(&zero) {
		// TODO: Panic division by zero
		return zero, zero
	}

	limit := answer.Not()
	limit = limit.shiftRightByOne()
	overflowed := false
	for denom.LessThanOrEqual(&dividend) {
		if !denom.LessThan(&limit) {
			overflowed = true
			break
		}
		denom = denom.shiftLeftByOne()
		current = current.shiftLeftByOne()
	}

	if !overflowed {
		denom = denom.shiftRightByOne()
		current = current.shiftRightByOne()
	}

	for !current.Equals(&zero) {
		if !dividend.LessThan(&denom) {
			dividend = dividend.Sub(&denom)
			answer = answer.Or(&current)
		}
		current = current.shiftRightByOne()
		denom = denom.shiftRightByOne()
	}

	return answer, dividend
}

// Div returns the quotient of x/y. If y is 0, a run-time panic occurs.
func (x *Int1024) Div(y *Int1024) Int1024 {
	div, _ := x.DivMod(y)
	return div
}

// Mod returns the modulus x%y. If y is 0, a run-time panic occurs.
func (x *Int1024) Mod(y *Int1024) Int1024 {
	_, mod := x.DivMod(y)
	return mod
}

// ModInverse returns the multiplicative inverse of x in the ring ℤ/nℤ.
// If x and n are not relatively prime, the result in undefined.
func (x *Int1024) ModInverse(n *Int1024) Int1024 {
	panic("Not implemented!")
}

// Words returns the internal [16]Word that stores the value of x
func (x *Int1024) Words() [INT1024WORDS]Word {
	return x.words
}
