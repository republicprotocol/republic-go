package stackint

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// go build -a -gcflags='-m -m' int1024.go int1024_arithmetic.go int1024_bitwise.go int1024_comparison.go int1024_internal.go

// SIZE is the number of bits stored by Int1024
const SIZE = 2048

// WORDSIZE is 64 for Word
const WORDSIZE = 64

// Word is the internal type
// type Word uint64

// WORDMAX represents the largest word value
const WORDMAX = 1<<WORDSIZE - 1

// BYTECOUNT represents the number of bytes that can be stores
const BYTECOUNT = SIZE / 8

// INT1024WORDS is 1024 / 64 = 16
const INT1024WORDS = SIZE / WORDSIZE

// Int1024 provides a 1024 bit number optimised to never use the heap
type Int1024 struct {
	words  [INT1024WORDS]uint64
	length uint16
}

// FromUint64 returns a new Int1024 from a Word
func FromUint64(n uint64) Int1024 {
	return Int1024{
		words:  [INT1024WORDS]uint64{n},
		length: 1,
	}
}

// SetUint64 sets x's value to n
func (x *Int1024) SetUint64(n uint64) {
	var i uint16
	for i = 1; i < x.length; i++ {
		x.words[i] = 0
	}
	x.words[0] = n
	x.length = 1
}

// ToUint64 converts an Int1024 to a uint64 if it is small enough
func (x *Int1024) ToUint64() uint64 {
	// Check that all other words are zero
	var i uint16
	for i = 1; i < x.length; i++ {
		if x.words[i] != 0 {
			panic("Int1024 is too large to be converted to uint64")
		}
	}
	return uint64(x.words[0])
}

// FromString returns a new Int1024 from a string
func FromString(number string) Int1024 {
	self := Zero()

	// Base to convert from
	base := 10

	if len(number) > 2 {
		if number[0:2] == "0x" {
			number = number[2:]
			base = 16
		} else if number[0:2] == "0b" {
			number = number[2:]
			base = 2
		}
	}

	// Length of string
	length := len(number)

	if length == 0 {
		panic("invalid number")
	}

	// Break up into blocks of size 19 (log10(2 ** 64))
	blockSize := 1
	limit := uint64(1<<63-1) / uint64(base)
	for basePower := uint64(1); basePower < limit; basePower *= uint64(base) {
		blockSize++
	}

	// Number of blocks
	count := (length / blockSize)
	if length%blockSize != 0 {
		count++
	}

	// TODO: Replace with 10.Pow(blockSize)
	shift := FromUint64(uint64(base))
	blockSizeInt := FromUint64(uint64(blockSize))
	shift = shift.Exp(&blockSizeInt)
	shiftAcc := One()

	// Loop through each block. Multiply block by (10**19)**i and add to number.
	for i := 0; i < count; i++ {
		end := length - i*blockSize
		start := length - (i+1)*blockSize
		if start < 0 {
			start = 0
		}
		word, err := strconv.ParseUint(number[start:end], base, 64)
		if err != nil {
			panic(err)
		}

		wordI := FromUint64(word)
		wordI = wordI.Mul(&shiftAcc)
		self = self.Add(&wordI)

		shiftAcc = shiftAcc.Mul(&shift)
	}

	return self
}

func (x *Int1024) String() string {
	blockSize := 19
	blockSize1024 := FromUint64(uint64(blockSize))
	base := FromUint64(10)
	base = base.Exp(&blockSize1024)
	q := *x
	var r Int1024
	ret := ""

	for !q.IsZero() {
		q, r = q.DivMod(&base)
		chars := strconv.FormatUint(uint64(r.words[0]), 10)
		zeroes := strings.Repeat("0", blockSize-len(chars))
		ret = zeroes + chars + ret
	}

	ret = strings.TrimLeft(ret, "0")

	if ret == "" {
		return "0"
	}

	return ret
}

// ToBytes returns an array of BYTECOUNT (128) bytes (Big Endian)
func (x *Int1024) ToBytes() []byte {

	bytesAll := make([]byte, BYTECOUNT)
	b8 := make([]byte, 8)

	var i uint16
	for i = 0; i < x.length; i++ {
		word := x.words[i]
		binary.BigEndian.PutUint64(b8, word)
		var j uint16
		for j = 0; j < 8; j++ {
			bytesAll[(INT1024WORDS-1-i)*8+j] = b8[j]
		}
	}

	return bytesAll
}

// FromBytes deserializes an array of BYTECOUNT (128) bytes to an Int1024 (Big Endian)
func FromBytes(bytesAll []byte) Int1024 {

	x := Zero()

	numWords := len(bytesAll) / 8
	if len(bytesAll)%8 != 0 {
		numWords++
	}

	if numWords > INT1024WORDS {
		numWords = INT1024WORDS
	}

	// mod := 8 - len(bytesAll)%8

	var firstPositive uint16
	for i := 0; i < numWords; i++ {
		b8 := make([]byte, 8)
		start := len(bytesAll) - i*8
		end := start - 8
		if end < 0 {
			end = 0
		}
		for j := 0; j < start-end; j++ {
			b8[7-j] = bytesAll[start-j-1]
		}
		x.words[i] = binary.BigEndian.Uint64(b8)
		if x.words[i] != 0 {
			firstPositive = uint16(i)
		}
	}

	x.length = firstPositive + 1

	return x
}

// // ToLittleEndianBytes returns an array of BYTECOUNT (128) bytes (Little Endian)
// func (x *Int1024) ToLittleEndianBytes() []byte {

// 	bytesAll := make([]byte, BYTECOUNT)
// 	b8 := make([]byte, 8)

// 	for i := range x.words {
// 		binary.LittleEndian.PutUint64(b8, uint64(x.words[i]))
// 		for j := 0; j < 8; j++ {
// 			bytesAll[i*8+j] = b8[j]
// 		}
// 	}

// 	return bytesAll
// }

// // FromLittleEndianBytes deserializes an array of 128 bytes to an Int1024 (LittleBig Endian)
// func FromLittleEndianBytes(bytesAll []byte) Int1024 {

// 	x := Zero()

// 	numWords := len(bytesAll) / 8

// 	for i := 0; i < numWords; i++ {
// 		b8 := bytesAll[i*8 : (i+1)*8]
// 		x.words[i] = binary.LittleEndian.Uint64(b8)
// 	}

// 	x.length = max(1, uint16(numWords))

// 	return x
// }

// Clone returns a new Int1024 representing the same value as x
func (x *Int1024) Clone() Int1024 {
	var words [INT1024WORDS]uint64
	var i uint16
	for i = 0; i < x.length; i++ {
		words[i] = x.words[i]
	}
	return Int1024{
		words:  words,
		length: x.length,
	}
}

// Words returns a clone of the [16]Word used by x as its internal representation
func (x *Int1024) Words() [INT1024WORDS]uint64 {
	var words [INT1024WORDS]uint64
	var i uint16
	for i = 0; i < x.length; i++ {
		words[i] = x.words[i]
	}
	return words
}

// ToBinary returns the binary representation of x as a string
func (x *Int1024) ToBinary() string {

	str := fmt.Sprintf("%b", x.words[x.length-1])
	var i int16
	for i = int16(x.length) - 2; i >= 0; i-- {
		str = str + fmt.Sprintf("%064b", x.words[i])
	}
	if str == "" {
		return "0"
	}

	return str
}

// ToBigInt converts x to a big.Int
func (x *Int1024) ToBigInt() *big.Int {
	return big.NewInt(0).SetBytes(x.ToBytes())
}

// FromBigInt converts a big.Int to an Int1024
func FromBigInt(bg *big.Int) Int1024 {
	return FromBytes(bg.Bytes())
}

/* CONSTANTS */

// Zero returns a new Int1024 representing 0
func Zero() Int1024 {
	var words [INT1024WORDS]uint64
	// words := make([]uint64, INT1024WORDS)
	// Not needed?
	for i := 0; i < len(words); i++ {
		words[i] = 0
	}
	return Int1024{
		words:  words,
		length: 1,
	}
}

// One is the Int1024 that represents 1
var One = func() Int1024 { return FromUint64(1) }

// Two is the Int1024 that represents 2
var Two = func() Int1024 { return FromUint64(2) }

// HalfMax represents max / 2
var HalfMax = func() Int1024 {
	max := MAXINT1024()
	max.ShiftRightInPlace(1)
	return max
}

// Do not call overwriting functions on these!
var zero = Zero()
var one = One()

var two = Two()
var halfMax = HalfMax()

// maxInt returns a new Int1024 representing 2**1024 - 1
func maxInt() Int1024 {
	var words [INT1024WORDS]uint64
	// words := make([]uint64, INT1024WORDS)
	for i := 0; i < len(words); i++ {
		words[i] = WORDMAX
	}
	return Int1024{
		words:  words,
		length: INT1024WORDS,
	}
}

// MAXINT1024 is the Int1024 that represents 2**1024 - 1
var MAXINT1024 = func() Int1024 { return maxInt() }

func max(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

// func (x *Int1024) Verify() {
// 	var i uint16
// 	for i = x.length; i < INT1024WORDS; i++ {
// 		if x.words[i] != 0 {
// 			fmt.Println(x)
// 			panic("Length too small")
// 		}
// 	}
// 	if x.words[x.length-1] == 0 && x.length != 1 {
// 		fmt.Println(x)
// 		panic("Length too big")
// 	}
// 	if x.length == 0 {
// 		panic("length is zero!")
// 	}
// }
