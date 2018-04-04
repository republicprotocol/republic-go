package stackint_test

import (
	"math/big"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/stackint"
)

var zero = Zero()
var one = FromUint(1)
var two = FromUint(2)

var three = FromUint(3)
var four = FromUint(4)
var five = FromUint(5)
var six = FromUint(6)
var seven = FromUint(7)
var eleven = FromUint(11)
var twelve = FromUint(12)
var oneWord = FromUint(WORDMAX)
var two64 = oneWord.Add(&one)
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

func MaxSquaredStr() string {
	one := big.NewInt(1)
	lim := big.NewInt(2)
	pow := big.NewInt(SIZE)
	lim = lim.Exp(lim, pow, nil)
	max := big.NewInt(0).Sub(lim, one)
	sqr := max.Mul(max, max)
	mod := sqr.Mod(sqr, lim)
	return mod.String()
}

var maxSquaredStr = MaxSquaredStr()

func TC(in ...interface{}) []interface{} {
	return in
}

var _ = Describe("Int1024", func() {

	Context("when converting from and to Words", func() {
		It("should return the right result for 1024 bit numbers", func() {
			cases := []uint{
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
				fromInt := FromUint(uint(n))
				toInt, err := fromInt.ToUint()
				Ω(err).Should(BeNil())
				Ω(toInt).Should(Equal(n))
				tmp := Zero()
				tmp.SetUint(n)
				Ω(tmp).Should(Equal(fromInt))
			}
		})

		It("should panic when converting a number bigger than MAX to Word", func() {
			_, err := max.ToUint()
			Ω(err).Should(Not(BeNil()))
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
				TC("0xFF", FromUint(255)),
				TC("0xff", FromUint(255)),
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

		It("should return error for invalid string", func() {
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
				_, err := FromString(tc[0].(string))
				Ω(err).Should(Not(BeNil()))
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

	Context("when converting to string", func() {
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

			tmp := FromUint(111)
			actual = tmp.String()
			expected = "111"
			Ω(actual).Should(Equal(expected))
		})
	})

	Context("when serializing to bytes", func() {
		It("should return the right result for 1024 bit numbers", func() {
			array := []Int1024{zero, one, two, three, four, five, six, seven, eleven, twelve, oneWord, max}
			for _, num := range array {
				actual, err := FromBytes(num.Bytes())
				Ω(err).Should(BeNil())
				Ω(actual).Should(Equal(num))

				actual, err = FromLittleEndianBytes(num.LittleEndianBytes())
				Ω(err).Should(BeNil())
				Ω(actual).Should(Equal(num))
			}

			str := "156110199609722120002645975834934187153674084697980344259599400078744195864483123168001725978362465713804593874868304438459220080111195600585730100927755271978903140799951022170241026510196255297991522400685742295892348482226518075857613157769551309646160118720740138838217231149054483993553648924213524999209"
			stackint, err := FromString(str)
			Ω(err).Should(BeNil())
			bigint, _ := big.NewInt(0).SetString(str, 10)
			Ω(stackint.ToBigInt().Cmp(bigint)).Should(Equal(0))
		})
	})
})
