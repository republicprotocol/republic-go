package stackint

// go build -a -gcflags='-m -m' int1024.go int1024_arithmetic.go int1024_bitwise.go int1024_comparison.go int1024_internal.go

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

// Int1024FromUint64 returns a new Int1024 from a Word
func Int1024FromUint64(n uint64) Int1024 {
	z := zero()
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

// Words returns the internal [16]Word that stores the value of x
func (x *Int1024) Words() [INT1024WORDS]Word {
	return x.words
}

/* CONSTANTS */

// zero returns a new Int1024 representing 0
func zero() Int1024 {
	var words [INT1024WORDS]Word
	for i := 0; i < INT1024WORDS; i++ {
		words[i] = 0
	}
	return Int1024{
		words: words,
	}
}

// ZERO is the Int1024 that represents 0
var ZERO = zero()

// maxInt returns a new Int1024 representing 2**1024 - 1
func maxInt() Int1024 {
	var words [INT1024WORDS]Word
	for i := 0; i < INT1024WORDS; i++ {
		words[i] = WORDMAX
	}
	return Int1024{
		words: words,
	}
}

// MAXINT1024 is the Int1024 that represents 2**1024 - 1
var MAXINT1024 = maxInt()

// maxInt returns a new Int1024 representing 2**1024 - 1
func twoPow1023() Int1024 {
	var words [INT1024WORDS]Word
	words[0] = 1 << 63
	return Int1024{
		words: words,
	}
}

// TWOPOW1023 is the Int1024 that represents 2**1023
var TWOPOW1023 = twoPow1023()
