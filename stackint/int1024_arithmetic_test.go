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

			first := FromString("340282366920938463417257747247494332417")
			second := FromString("340282366920938463454151235394913435648")
			expected := FromString("680564733841876926871408982642407768065")
			actual := first.Add(&second)
			Ω(actual.Equals(&expected)).Should(BeTrue())

			first = FromString("6893488147419103231")
			second = FromString("30000000000000000000")
			expected = FromString("36893488147419103231")
			actual = first.Add(&second)
			Ω(actual.Equals(&expected)).Should(BeTrue())
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

			sqrt := FromString("13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006084095")
			sqrtMulSqrt := sqrt.Mul(&sqrt)
			diff := FromString("26815615859885194199148049996411692254958731641184786755447122887443528060147093953603748596333806855380063716372972101707507765623893139892867298012168190")
			expected := max.Sub(&diff)
			Ω(sqrtMulSqrt.Equals(&expected)).Should(BeTrue())
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

			// sqr := FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474097562152033539671286128252223189553839160721441767298250321715263238814402734379959506792230903356495130620869925267845538430714092411695463462326211969025")
			// sqrt := FromString("13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006084095")
			// maxDivSqrt := sqr.Div(&sqrt)
			// Ω(maxDivSqrt.Equals(&sqrt)).Should(BeTrue())
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

	Context("when taking the multiplicative inverse", func() {
		It("should return the right result for 1024 bit numbers", func() {
			threeInvModSeven := three.ModInverse(&seven)
			Ω(threeInvModSeven.Equals(&five)).Should(BeTrue())

			Ω(func() { two.ModInverse(&four) }).Should(Panic())

			oneInvTwo := one.ModInverse(&two)
			Ω(oneInvTwo.Equals(&one)).Should(BeTrue())

			twoInvEleven := two.ModInverse(&eleven)
			Ω(twoInvEleven.Equals(&six)).Should(BeTrue())

			n := FromUint64(1234567)
			m := FromUint64(458948883992)
			result := n.ModInverse(&m)
			expected := FromUint64(14332777583)
			Ω(result.Equals(&expected)).Should(BeTrue())

			// fmt.Println(one.ModInverse(&one)) // Actual: 1, Expected?
		})
	})

	Context("when raising powers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			oneExpOne := one.Exp(&one)
			Ω(oneExpOne.Equals(&one)).Should(BeTrue())

			twoExpZero := two.Exp(&zero)
			Ω(twoExpZero.Equals(&one)).Should(BeTrue())

			zeroExpZero := zero.Exp(&zero)
			Ω(zeroExpZero.Equals(&one)).Should(BeTrue())

			zeroExpOne := zero.Exp(&one)
			Ω(zeroExpOne.Equals(&zero)).Should(BeTrue())

			threeExpThree := three.Exp(&three)
			expected := FromUint64(27)
			Ω(threeExpThree.Equals(&expected)).Should(BeTrue())

			oneLess := FromUint64(4294967296)
			lessExpTwo := oneLess.Exp(&two)
			expected = oneWord.Add(&ONE)
			Ω(lessExpTwo.Equals(&expected)).Should(BeTrue())

			sqrt := FromString("13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006084095")
			sqrtExpTwo := sqrt.Exp(&TWO)
			diff := FromString("26815615859885194199148049996411692254958731641184786755447122887443528060147093953603748596333806855380063716372972101707507765623893139892867298012168190")
			expected = max.Sub(&diff)
			Ω(sqrtExpTwo.Equals(&expected)).Should(BeTrue())
		})

		It("should overflow", func() {
			//
		})
	})
})
