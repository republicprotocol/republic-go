package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var zero = Zero
var one = Int1024FromUint64(1)
var two = Int1024FromUint64(2)
var three = Int1024FromUint64(3)
var six = Int1024FromUint64(6)
var seven = Int1024FromUint64(7)
var eleven = Int1024FromUint64(11)
var oneWord = Int1024FromUint64(WORDMAX)
var max = zero.NOT()

var _ = Describe("Int1024", func() {
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

	Context("when dividing numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			sixDivtwo, sixModTwo := six.DivMod(&two)
			Ω(sixDivtwo.Equals(&three)).Should(BeTrue())
			Ω(sixModTwo.Equals(&zero)).Should(BeTrue())

			sevenDivtwo, sevenModTwo := seven.DivMod(&two)
			Ω(sevenDivtwo.Equals(&three)).Should(BeTrue())
			Ω(sevenModTwo.Equals(&one)).Should(BeTrue())

			elevenDivThree, elevenModThree := eleven.DivMod(&three)
			Ω(elevenDivThree.Equals(&three)).Should(BeTrue())
			Ω(elevenModThree.Equals(&two)).Should(BeTrue())

			maxDivMax, maxModMax := max.DivMod(&max)
			Ω(maxDivMax.Equals(&one)).Should(BeTrue())
			Ω(maxModMax.Equals(&zero)).Should(BeTrue())
		})
	})
})
