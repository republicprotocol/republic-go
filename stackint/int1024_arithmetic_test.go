package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var addFn = func(inputs ...Int1024) Int1024 { return inputs[0].Add(&inputs[1]) }
var subFn = func(inputs ...Int1024) Int1024 { return inputs[0].Sub(&inputs[1]) }
var mulFn = func(inputs ...Int1024) Int1024 { return inputs[0].Mul(&inputs[1]) }

var _ = Describe("Int1024 arithmetic", func() {

	// ADDITION
	Context("when adding numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			RunAllCases(addFn, []TestCase{
				TestCase{inputsStr: []string{"1", "2"}, expectedStr: "3"},
				TestCase{inputsInt: []Int1024{one, two}, expectedInt: &three},
				TestCase{inputsInt: []Int1024{oneWord, one}, expectedStr: "18446744073709551616"},
				TestCase{inputsStr: []string{"340282366920938463417257747247494332417", "340282366920938463454151235394913435648"}, expectedStr: "680564733841876926871408982642407768065"},
				TestCase{inputsStr: []string{"6893488147419103231", "30000000000000000000"}, expectedStr: "36893488147419103231"},

				// Overflow:
				TestCase{inputsInt: []Int1024{max, one}, expectedInt: &zero},
			})
		})
	})

	// SUBTRACTION
	Context("when subtracting numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			RunAllCases(subFn, []TestCase{
				TestCase{inputsStr: []string{"3", "2"}, expectedStr: "1"},
				TestCase{inputsInt: []Int1024{three, two}, expectedInt: &one},
				TestCase{inputsStr: []string{"18446744073709551616", "1"}, expectedStr: "18446744073709551615"},
				TestCase{inputsStr: []string{"18446744073709551616", "9223372036854775808"}, expectedStr: "9223372036854775808"},
				TestCase{inputsInt: []Int1024{max, max}, expectedStr: "0"},
				TestCase{inputsStr: []string{"18446744073709551616", "1"}, expectedInt: &oneWord},

				// Overflow:
				TestCase{inputsInt: []Int1024{zero, one}, expectedInt: &max},
			})
		})
	})

	Context("when multiplying numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			RunAllCases(mulFn, []TestCase{
				TestCase{inputsStr: []string{"2", "3"}, expectedStr: "6"},
				TestCase{inputsInt: []Int1024{two, three}, expectedInt: &six},
				TestCase{inputsStr: []string{"18446744073709551615", "18446744073709551615"}, expectedStr: "340282366920938463426481119284349108225"},
				TestCase{inputsStr: []string{"18446744073709551616", "18446744073709551616"}, expectedStr: "340282366920938463463374607431768211456"},
				TestCase{inputsStr: []string{"13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006084095", "13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006084095"}, expectedStr: "179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474097562152033539671286128252223189553839160721441767298250321715263238814402734379959506792230903356495130620869925267845538430714092411695463462326211969025"},
			})
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

type BinaryFn func(inputs ...Int1024) Int1024

func RunCase(fn BinaryFn, test TestCase) {

	if len(test.inputsInt) == 0 {
		test.inputsInt = make([]Int1024, len(test.inputsStr))
		for i, input := range test.inputsStr {
			test.inputsInt[i] = FromString(input)
		}
	}

	if test.expectedStr == "" {
		test.expectedStr = test.expectedInt.String()
	}

	actual := fn(test.inputsInt...)
	actualStr := actual.String()

	Ω(actualStr).Should(Equal(test.expectedStr))
}

func RunAllCases(fn BinaryFn, testcases []TestCase) {
	for _, testcase := range testcases {
		RunCase(fn, testcase)
	}
}

type TestCase struct {
	inputsStr   []string
	inputsInt   []Int1024
	expectedStr string
	expectedInt *Int1024
}
