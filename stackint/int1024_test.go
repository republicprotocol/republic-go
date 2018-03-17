package stackint_test

import (
	"math/big"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var zero = Zero()
var one = FromUint64(1)
var two = FromUint64(2)
var three = FromUint64(3)
var four = FromUint64(4)
var five = FromUint64(5)
var six = FromUint64(6)
var seven = FromUint64(7)
var eleven = FromUint64(11)
var twelve = FromUint64(12)
var oneWord = FromUint64(WORDMAX)
var max = MAXINT1024()

func MaxStr() string {
	one := big.NewInt(1)
	max := big.NewInt(2)
	pow := big.NewInt(SIZE)
	max = max.Exp(max, pow, nil)
	max = max.Sub(max, one)
	return max.String()
}

var maxStr = MaxStr()

func TC(in ...interface{}) []interface{} {
	return in
}

var _ = Describe("Int1024", func() {

	Context("when converting from and to uint64s", func() {
		It("should return the right result for 1024 bit numbers", func() {
			cases := []uint64{
				0,
				1,
				4294967295,
				4294967296,
				8589934591,
				8589934592,
				9223372036854775807,
				9223372036854775808,
				18446744073709551615,
			}

			for _, n := range cases {
				fromInt := FromUint64(n)
				Ω(fromInt.ToUint64()).Should(Equal(n))
			}
		})

		It("should panic when converting a number bigger than MAX to uint64", func() {
			Ω(func() { max.ToUint64() }).Should(Panic())
		})
	})

	Context("when converting from string", func() {
		It("should return the right result for 1024 bit numbers", func() {

			cases := [][]interface{}{
				TC("0", zero),
				TC("1", one),
				TC("0x0", zero),
				TC("0x00", zero),
				TC("0x00", zero),
				TC("0xFF", FromUint64(255)),
				TC("0xff", FromUint64(255)),
				TC("0b0", zero),
				TC("0b00", zero),
				TC("0b01", one),
				TC("0b"+strings.Repeat("1", SIZE), max),
				TC(maxStr, max),
			}
			for _, tc := range cases {
				Ω(FromString(tc[0].(string))).Should(Equal(tc[1]))
			}

		})

		It("should return the right result for 1024 bit numbers", func() {
			cases := [][]interface{}{
				TC("ff"),
				TC("NOT A STRING"),
				TC("1234i"),
				TC("0bA"),
				TC("0x"),
				TC("0b"),
				TC(""),
			}
			for _, tc := range cases {
				Ω(func() { FromString(tc[0].(string)) }).Should(Panic())
			}
		})
	})

	Context("when converting to binary string", func() {
		It("should return the right result for 1024 bit numbers", func() {
			actual := max.ToBinary()
			expected := strings.Repeat("1", SIZE)
			Ω(actual).Should(Equal(expected))

			actual = zero.ToBinary()
			expected = "0"
			Ω(actual).Should(Equal(expected))

			actual = one.ToBinary()
			expected = "1"
			Ω(actual).Should(Equal(expected))

			actual = two.ToBinary()
			expected = "10"
			Ω(actual).Should(Equal(expected))
		})
	})

	Context("when converting to binary string", func() {
		It("should return the right result for 1024 bit numbers", func() {
			actual := max.String()
			expected := maxStr
			Ω(actual).Should(Equal(expected))

			actual = zero.String()
			expected = "0"
			Ω(actual).Should(Equal(expected))

			actual = one.String()
			expected = "1"
			Ω(actual).Should(Equal(expected))

			actual = two.String()
			expected = "2"
			Ω(actual).Should(Equal(expected))

			tmp := FromUint64(111)
			actual = tmp.String()
			expected = "111"
			Ω(actual).Should(Equal(expected))

		})
	})

	Context("when serializing to bytes", func() {
		It("should return the right result for 1024 bit numbers", func() {
			array := []Int1024{zero, one, two, three, four, five, six, seven, eleven, twelve, oneWord, max}
			for _, num := range array {
				actual := FromBytes(num.ToBytes())
				Ω(actual.String()).Should(Equal(num.String()))

				actual = FromLittleEndianBytes(num.ToLittleEndianBytes())
				Ω(actual.String()).Should(Equal(num.String()))
			}
		})
	})

	Context("when retrieving words", func() {
		It("should return the right result for 1024 bit numbers", func() {
			array := []Int1024{zero, one, two, three, four, five, six, seven, eleven, twelve, oneWord, max}
			for _, num := range array {
				actual := FromBytes(num.ToBytes())
				Ω(actual.String()).Should(Equal(num.String()))

				actual = FromLittleEndianBytes(num.ToLittleEndianBytes())
				Ω(actual.String()).Should(Equal(num.String()))
			}
		})
	})
})
