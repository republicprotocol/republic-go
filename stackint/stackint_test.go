package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var zero = Zero()
var one = Int1024FromUint64(1)
var two = Int1024FromUint64(2)
var three = Int1024FromUint64(3)
var six = Int1024FromUint64(6)
var oneWord = Int1024FromUint64(WORDMAX)

var _ = Describe("Stackint", func() {
	Context("when adding numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			onePlusTwo := one.Add(&two)
			Ω(onePlusTwo.Equals(&three)).Should(BeTrue())

			oneWordPlusOne := oneWord.Add(&one)
			Ω(oneWordPlusOne.Words()[14]).Should(Equal(Word(1)))
		})
	})

	Context("when subtracting numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			threeMinusTwo := three.Sub(&two)
			Ω(threeMinusTwo.Equals(&one)).Should(BeTrue())

			oneWordPlusOne := oneWord.Add(&one)
			alsoOneWord := oneWordPlusOne.Sub(&one)
			Ω(alsoOneWord.Equals(&oneWord)).Should(BeTrue())
		})

		It("should overflow", func() {
			overflow := zero.Sub(&one)
			for i := 0; i < INT1024WORDS; i++ {
				Ω(overflow.Words()[i]).Should(Equal(Word(WORDMAX)))
			}
		})
	})

	Context("when multiplying numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			twoTimesThree := two.Mul(&three)
			Ω(twoTimesThree.Equals(&six)).Should(BeTrue())

			oneWordSquared := oneWord.Mul(&oneWord)
			Ω(oneWordSquared.Words()[INT1024WORDS-1]).Should(Equal(Word(1)))
			Ω(oneWordSquared.Words()[INT1024WORDS-2]).Should(Equal(Word(WORDMAX - 1)))
		})
	})
})
