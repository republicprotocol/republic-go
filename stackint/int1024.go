package stackint

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big" // Used for converting to/from big.Ints
	"strconv"
	"strings"

	"github.com/republicprotocol/republic-go/stackint/asm"
)

// Stackint is a big number library offering 1024-bit numbers
// whose operators do not not allocate any heap memory

// For debugging to check for heap allocations:
// go build -a -gcflags='-m -m'
// grep with: `2>&1 | grep -E "(leaking|escapes)"`

// SIZE is the number of bits stored by Int1024
const SIZE = 1024

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
	words  [INT1024WORDS]asm.Word
	length uint16
}

// FromUint returns a new Int1024 from a Word
func FromUint(n uint) Int1024 {
	x := Zero()
	x.SetUint(n)
	return x
}

// SetUint sets x's value to n
func (x *Int1024) SetUint(n uint) {
	var i uint16
	for i = 1; i < x.length; i++ {
		x.words[i] = 0
	}
	x.words[0] = asm.Word(n)
	if n == 0 {
		x.length = 0
	} else {
		x.length = 1
	}
}

// ToUint converts an Int1024 to a Word if it is small enough
func (x *Int1024) ToUint() (uint, error) {
	// Check that all other words are zero
	var i uint16
	for i = 1; i < x.length; i++ {
		if x.words[i] != 0 {
			return 0, errors.New("too big for uint")
		}
	}
	return uint(x.words[0]), nil
}

// FromString returns a new Int1024 from a string
func FromString(number string) (Int1024, error) {
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
		return Int1024{}, errors.New("string must not be empty")
	}

	if number == "0" {
		return self, nil
	}

	// Break up into blocks of size 19 (log10(2 ** 64))
	blockSize := 1
	limit := asm.Word(1<<63-1) / asm.Word(base)
	for basePower := asm.Word(1); basePower < limit; basePower *= asm.Word(base) {
		blockSize++
	}

	// Number of blocks
	count := (length / blockSize)
	if length%blockSize != 0 {
		count++
	}

	shift := FromUint(uint(base))
	blockSizeInt := FromUint(uint(blockSize))
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
			return Int1024{}, err
		}

		wordI := FromUint(uint(word))
		wordI = wordI.Mul(&shiftAcc)
		self = self.Add(&wordI)

		shiftAcc = shiftAcc.Mul(&shift)
	}

	return self, nil
}

func (x *Int1024) String() string {
	blockSize := 19
	blockSize1024 := FromUint(uint(blockSize))
	base := FromUint(10)
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

// Bytes returns an array of BYTECOUNT (128) bytes (Big Endian)
func (x *Int1024) Bytes() []byte {

	bitlen := x.BitLength()
	i := uint16(bitlen / 8)
	if bitlen%8 > 0 {
		i++
	}
	buf := make([]byte, i)

	var k uint16
	for k = 0; k < x.length; k++ {
		d := x.words[k]
		for j := 0; j < asm.S; j++ {
			if d > 0 {
				buf[i-1] = byte(d)
			}
			d >>= 8
			i--
		}
	}

	return buf
}

// FromBytes deserializes an array of BYTECOUNT (128) bytes to a new Int1024
// (Big Endian)
func FromBytes(bytesAll []byte) (Int1024, error) {
	x := Zero()
	err := x.SetBytes(bytesAll)
	return x, err
}

// SetBytes deserializes an array of BYTECOUNT (128) bytes in place
// (Big Endian)
func (x *Int1024) SetBytes(bytesAll []byte) error {

	leadingZeros := 0
	for i, b := range bytesAll {
		if b != 0 {
			break
		}
		leadingZeros = i
	}
	bytesAll = bytesAll[leadingZeros:]

	numWords := len(bytesAll) / 8
	if len(bytesAll) > 0 && len(bytesAll)%8 != 0 {
		numWords++
	}

	if numWords > INT1024WORDS {
		return errors.New("bytes array too long")
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
		x.words[i] = asm.Word(binary.BigEndian.Uint64(b8))
		if x.words[i] != 0 {
			firstPositive = uint16(i) + 1
		}
	}
	for i := numWords; i < INT1024WORDS; i++ {
		x.words[i] = 0
	}

	x.length = firstPositive

	return nil
}

// // LittleEndianBytes returns an array of BYTECOUNT (128) bytes (Little Endian)
// func (x *Int1024) LittleEndianBytes() []byte {

// 	bitlen := x.BitLength()
// 	i := uint16(bitlen / 8)
// 	if bitlen%8 > 0 {
// 		i++
// 	}
// 	buf := make([]byte, i)

// 	var index uint16
// 	var k uint16
// 	for k = 0; k < x.length; k++ {
// 		d := x.words[k]
// 		for j := 0; j < asm.S && d > 0; j++ {
// 			if d > 0 {
// 				buf[index] = byte(d)
// 			}
// 			d >>= 8
// 			index++
// 		}
// 	}
// 	return buf
// }

// // SetLittleEndianBytes deserializes an array of 128 bytes to an Int1024
// // (LittleBig Endian)
// func FromLittleEndianBytes(bytesAll []byte) (Int1024, error) {
// 	x := Zero()
// 	err := x.SetLittleEndianBytes(bytesAll)
// 	return x, err
// }

// // SetLittleEndianBytes deserializes an array of 128 bytes to an Int1024
// // (LittleBig Endian)
// func (x *Int1024) SetLittleEndianBytes(bytesAll []byte) error {

// 	fmt.Println(bytesAll)
// 	len := len(bytesAll)

// 	if len == 0 {
// 		return nil
// 	}

// 	numWords := len / 8
// 	if len%8 != 0 {
// 		numWords++
// 	}

// 	if numWords > INT1024WORDS {
// 		return errors.New("bytes array too long")
// 	}

// 	for i := 0; i < numWords; i++ {
// 		b8 := make([]byte, 8)
// 		last := (i + 1) * 8
// 		if last > len {
// 			last = len
// 		}
// 		copy(b8, bytesAll[i*8:last])
// 		x.words[i] = asm.Word(binary.LittleEndian.Uint64(b8))
// 	}
// 	for i := numWords; i < INT1024WORDS; i++ {
// 		x.words[i] = 0
// 	}

// 	x.length = uint16(numWords)
// 	fmt.Println(x)

// 	return nil
// }

// MarshalJSON implements the json.Marshaler interface.
func (x Int1024) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Bytes())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (x *Int1024) UnmarshalJSON(data []byte) error {
	var bytes []byte
	if err := json.Unmarshal(data, &bytes); err != nil {
		return err
	}
	return x.SetBytes(bytes)
}

// Clone returns a new Int1024 representing the same value as x
func (x *Int1024) Clone() Int1024 {
	var words [INT1024WORDS]asm.Word
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
func (x *Int1024) Words() [INT1024WORDS]uint {
	var words [INT1024WORDS]uint
	var i uint16
	for i = 0; i < x.length; i++ {
		words[i] = uint(x.words[i])
	}
	return words
}

// ToBinary returns the binary representation of x as a string
func (x *Int1024) ToBinary() string {
	if x.length == 0 {
		return "0"
	}

	str := fmt.Sprintf("%b", x.words[x.length-1])
	var i int16
	for i = int16(x.length) - 2; i >= 0; i-- {
		str = str + fmt.Sprintf("%064b", x.words[i])
	}

	return str
}

// ToBigInt converts x to a big.Int
func (x *Int1024) ToBigInt() *big.Int {
	return big.NewInt(0).SetBytes(x.Bytes())
}

// FromBigInt converts a big.Int to an Int1024
func FromBigInt(bg *big.Int) (Int1024, error) {
	return FromBytes(bg.Bytes())
}

/* CONSTANTS */

// Zero returns a new Int1024 representing 0
func Zero() Int1024 {
	return Int1024{
		words:  [INT1024WORDS]asm.Word{},
		length: 0,
	}
}

// One is the Int1024 that represents 1
var One = func() Int1024 { return FromUint(1) }

// Two is the Int1024 that represents 2
var Two = func() Int1024 { return FromUint(2) }

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
	var words [INT1024WORDS]asm.Word
	// words := make([]Word, INT1024WORDS)
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

// max returns max(a,b) for uint16s
func max(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

// min retrusn min(a,b) for uint16s
func min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

// setLength recalculates x's length
func (x *Int1024) setLength() {
	var firstPositive uint16
	for i := int(INT1024WORDS) - 1; i >= 0; i-- {
		if x.words[i] != 0 {
			firstPositive = uint16(i) + 1
			break
		}
	}
	x.length = firstPositive
}
