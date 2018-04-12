package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

//  {"0x00", "0x00", "0x00", "0x00", "0x00", "0x00"},
// 	{"0x00", "0x01", "0x00", "0x01", "0x01", "0x00"},
// 	{"0x01", "0x00", "0x00", "0x01", "0x01", "0x01"},
// 	{"0x01", "0x01", "0x01", "0x01", "0x00", "0x00"},
// 	{"0x07", "0x08", "0x00", "0x0f", "0x0f", "0x07"},
// 	{"0x05", "0x0f", "0x05", "0x0f", "0x0a", "0x00"},
// 	{"0x013ff6", "0x9a4e", "0x1a46", "0x01bffe", "0x01a5b8", "0x0125b0"},
// 	{"-0x013ff6", "0x9a4e", "0x800a", "-0x0125b2", "-0x01a5bc", "-0x01c000"},
// 	{
// 		"0x1000009dc6e3d9822cba04129bcbe3401",
// 		"0xb9bd7d543685789d57cb918e833af352559021483cdb05cc21fd",
// 		"0x1000001186210100001000009048c2001",
// 		"0xb9bd7d543685789d57cb918e8bfeff7fddb2ebe87dfbbdfe35fd",
// 		"0xb9bd7d543685789d57ca918e8ae69d6fcdb2eae87df2b97215fc",
// 		"0x8c40c2d8822caa04120b8321400",
// 	},
// }

var andFn = func(inputs ...Int1024) Int1024 { return inputs[0].AND(&inputs[1]) }
var orFn = func(inputs ...Int1024) Int1024 { return inputs[0].OR(&inputs[1]) }
var xorFn = func(inputs ...Int1024) Int1024 { return inputs[0].XOR(&inputs[1]) }
var notFn = func(inputs ...Int1024) Int1024 { return inputs[0].NOT() }

var _ = Describe("Int1024 bitwise operations", func() {
	Context("when left-shifting a number", func() {
		It("should return the right result for 1024 bit numbers", func() {
			oneShiftLOne := one.ShiftLeft(1)
			Ω(oneShiftLOne).Should(Equal(two))

			shiftHalf := one.ShiftLeft(SIZE - 1)
			expected := HalfMax()
			expected = expected.Add(&one)
			Ω(shiftHalf).Should(Equal(expected))

			// Shift into next word by 1
			shifted := oneWord.ShiftLeft(1)
			expected = two64.Add(&oneWord)
			expected = expected.Sub(&one)
			Ω(shifted).Should(Equal(expected))

			// Shift into next word by more than 1
			shifted = two.ShiftLeft(63)
			expected = two64
			Ω(shifted).Should(Equal(expected))

			zeroShiftLOne := zero.ShiftLeft(1)
			Ω(zeroShiftLOne).Should(Equal(zero))
		})

		It("should overflow without wrapping", func() {
			overflow := one.ShiftLeft(SIZE)
			Ω(overflow).Should(Equal(zero))
		})

	})

	Context("when right-shifting a number", func() {
		It("should return the right result for 1024 bit numbers", func() {
			oneShiftROne := one.ShiftRight(1)
			Ω(oneShiftROne).Should(Equal(zero))

			twoShiftROne := two.ShiftRight(1)
			Ω(twoShiftROne).Should(Equal(one))

			zeroShiftROne := zero.ShiftRight(1)
			Ω(zeroShiftROne).Should(Equal(zero))

			// (shift amount) >= (word size)
			shifted := two64.ShiftRight(WORDSIZE)
			expected := one
			Ω(shifted).Should(Equal(expected))

			// 1 < (shift amount) < (word size)
			shifted = two64.ShiftRight(WORDSIZE - 1)
			expected = two
			Ω(shifted).Should(Equal(expected))

			// shift amount == 1
			shifted = two64.ShiftRight(1)
			expected = FromUint(1 << (WORDSIZE - 1))
			Ω(shifted).Should(Equal(expected))

			elevenShiftRTwo := eleven.ShiftRight(2)
			Ω(elevenShiftRTwo).Should(Equal(two))
		})

		It("should overflow without wrapping", func() {
			oneShiftROne := one.ShiftRight(2)
			Ω(oneShiftROne).Should(Equal(zero))

			shifted := one.ShiftRight(WORDSIZE)
			Ω(shifted).Should(Equal(zero))

			shifted = one.ShiftRight(WORDSIZE * 2)
			Ω(shifted).Should(Equal(zero))
		})
	})

	Context("when taking the bitwise AND", func() {
		It("should return the right result for 1024 bit numbers", func() {
			RunAllCases(andFn, []TestCase{
				TestCase{inputsStr: []string{"0", "0"}, expectedStr: "0"},
				TestCase{inputsStr: []string{"1", "0"}, expectedStr: "0"},
				TestCase{inputsStr: []string{"0", "1"}, expectedStr: "0"},
				TestCase{inputsStr: []string{"1", "1"}, expectedStr: "1"},
				TestCase{inputsStr: []string{"10", "11"}, expectedStr: "10"},
				TestCase{inputsStr: []string{"11", "11"}, expectedStr: "11"},
				TestCase{inputsInt: []Int1024{max, max}, expectedInt: &max},
				TestCase{inputsInt: []Int1024{one, max}, expectedInt: &one},
				TestCase{inputsInt: []Int1024{max, one}, expectedInt: &one},
				TestCase{inputsInt: []Int1024{seven, eleven}, expectedInt: &three},
			})
		})
	})

	Context("when taking the bitwise OR", func() {
		It("should return the right result for 1024 bit numbers", func() {
			RunAllCases(orFn, []TestCase{
				TestCase{inputsStr: []string{"0", "0"}, expectedStr: "0"},
				TestCase{inputsStr: []string{"1", "0"}, expectedStr: "1"},
				TestCase{inputsStr: []string{"0", "1"}, expectedStr: "1"},
				TestCase{inputsStr: []string{"1", "1"}, expectedStr: "1"},
				TestCase{inputsStr: []string{"10", "11"}, expectedStr: "11"},
				TestCase{inputsStr: []string{"11", "11"}, expectedStr: "11"},
				TestCase{inputsInt: []Int1024{max, max}, expectedInt: &max},
				TestCase{inputsInt: []Int1024{one, max}, expectedInt: &max},
				TestCase{inputsInt: []Int1024{max, one}, expectedInt: &max},
				TestCase{inputsInt: []Int1024{seven, eleven}, expectedStr: "15"},
			})
		})
	})

	Context("when taking the bitwise XOR", func() {
		It("should return the right result for 1024 bit numbers", func() {
			maxSubOne := max.Sub(&one)
			RunAllCases(xorFn, []TestCase{
				TestCase{inputsStr: []string{"0", "0"}, expectedStr: "0"},
				TestCase{inputsStr: []string{"1", "0"}, expectedStr: "1"},
				TestCase{inputsStr: []string{"0", "1"}, expectedStr: "1"},
				TestCase{inputsStr: []string{"1", "1"}, expectedStr: "0"},
				TestCase{inputsStr: []string{"10", "11"}, expectedStr: "1"},
				TestCase{inputsStr: []string{"11", "11"}, expectedStr: "0"},
				TestCase{inputsInt: []Int1024{max, max}, expectedInt: &zero},
				TestCase{inputsInt: []Int1024{one, max}, expectedInt: &maxSubOne},
				TestCase{inputsInt: []Int1024{max, one}, expectedInt: &maxSubOne},
				TestCase{inputsInt: []Int1024{seven, eleven}, expectedInt: &twelve},
			})
		})
	})

	Context("when taking the bitwise NOT", func() {
		It("should return the right result for 1024 bit numbers", func() {
			maxSubOne := max.Sub(&one)
			maxSubWord := max.Sub(&oneWord)
			maxSubTwo64 := max.Sub(&two64)

			RunAllCases(notFn, []TestCase{
				TestCase{inputsInt: []Int1024{max}, expectedInt: &zero},
				TestCase{inputsInt: []Int1024{zero}, expectedInt: &max},
				TestCase{inputsInt: []Int1024{one}, expectedInt: &maxSubOne},
				TestCase{inputsInt: []Int1024{maxSubOne}, expectedInt: &one},
				TestCase{inputsInt: []Int1024{oneWord}, expectedInt: &maxSubWord},
				TestCase{inputsInt: []Int1024{maxSubWord}, expectedInt: &oneWord},
				TestCase{inputsInt: []Int1024{two64}, expectedInt: &maxSubTwo64},
				TestCase{inputsInt: []Int1024{maxSubTwo64}, expectedInt: &two64},
			})
		})
	})

	Context("when checking bits", func() {
		It("should return the right result for 1024 bit numbers", func() {
			cases := [][]interface{}{
				TC("0b1", 0, true),
				TC("0b1", -1, false),
				TC("0b0", 0, false),
				TC("0b10", 1, true),
				TC("0b10", 0, false),
				TC("0b10000000", 7, true),
				TC("0b10000000", 8, false),
				TC("0b10000000", 6, false),
			}

			for _, tc := range cases {
				tmp, err := FromString(tc[0].(string))
				Ω(err).Should(BeNil())
				Ω(tmp.IsBitSet(tc[1].(int))).Should(Equal(tc[2].(bool)))
			}

			for i := 0; i < SIZE; i++ {
				Ω(max.IsBitSet(i))
				Ω(zero.IsBitSet(i)).Should(BeFalse())
			}
		})
	})
})
