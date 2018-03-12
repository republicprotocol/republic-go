package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Int1024 arithmetic", func() {
	Context("when adding numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			onePlusTwo := one.Add(&two)
			Ω(onePlusTwo.Equals(&three)).Should(BeTrue())

			oneWordPlusOne := oneWord.Add(&one)
			Ω(oneWordPlusOne.Words()[14]).Should(Equal(Word(1)))
		})

		It("should overflow", func() {
			overflow := max.Add(&one)
			Ω(overflow.Equals(&zero)).Should(BeTrue())
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
			Ω(overflow.Equals(&max)).Should(BeTrue())
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
			sixDivtwo := six.Div(&two)
			Ω(sixDivtwo.Equals(&three)).Should(BeTrue())

			sevenDivtwo := seven.Div(&two)
			Ω(sevenDivtwo.Equals(&three)).Should(BeTrue())

			elevenDivThree := eleven.Div(&three)
			Ω(elevenDivThree.Equals(&three)).Should(BeTrue())

			maxDivMax := max.Div(&max)
			Ω(maxDivMax.Equals(&one)).Should(BeTrue())
		})

		It("panic when dividing by zero", func() {
			Ω(func() { one.Div(&zero) }).Should(Panic())
		})
	})

	Context("when taking the modulus", func() {
		It("should return the right result for 1024 bit numbers", func() {
			sixModTwo := six.Mod(&two)
			Ω(sixModTwo.Equals(&zero)).Should(BeTrue())

			sevenModTwo := seven.Mod(&two)
			Ω(sevenModTwo.Equals(&one)).Should(BeTrue())

			elevenModThree := eleven.Mod(&three)
			Ω(elevenModThree.Equals(&two)).Should(BeTrue())

			maxModMax := max.Mod(&max)
			Ω(maxModMax.Equals(&zero)).Should(BeTrue())
		})

		It("panic when taking the modulus by zero", func() {
			Ω(func() { one.Div(&zero) }).Should(Panic())
		})
	})
})
