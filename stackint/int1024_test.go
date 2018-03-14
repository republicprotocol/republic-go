package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var zero = ZERO
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
var twoPow1023 = TWOPOW1023
var max = MAXINT1024

var _ = Describe("Int1024", func() {

	Context("when converting from string", func() {
		It("should return the right result for 1024 bit numbers", func() {
			maxFromString := FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137215")
			Ω(maxFromString.Equals(&max)).Should(BeTrue())

			zeroFromString := FromString("0")
			Ω(zeroFromString.Equals(&zero)).Should(BeTrue())

			oneFromString := FromString("1")
			Ω(oneFromString.Equals(&one)).Should(BeTrue())

			Ω(func() { FromString("NOT A STRING") }).Should(Panic())
			Ω(func() { FromString("0x123") }).Should(Panic())
			Ω(func() { FromString("1234I") }).Should(Panic())
		})
	})

	Context("when converting to binary string", func() {
		It("should return the right result for 1024 bit numbers", func() {
			actual := max.ToBinary()
			expected := "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
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
			expected := "179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137215"
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

			// actual := FromBytes([]byte{01})
			// Ω(actual.String()).Should(Equal("1"))
		})
	})
})
