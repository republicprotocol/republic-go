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

var _ = Describe("Int1024 bitwise operations", func() {
	Context("when left-shifting a number", func() {
		It("should return the right result for 1024 bit numbers", func() {
			oneShiftLOne := one.ShiftLeft(1)
			Ω(oneShiftLOne).Should(Equal(two))

			shiftHalf := one.ShiftLeft(SIZE - 1)
			expected := HalfMax()
			expected = expected.Add(&one)
			Ω(shiftHalf).Should(Equal(expected))

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
			twoShiftROne := two.ShiftRight(1)
			Ω(twoShiftROne).Should(Equal(one))

			zeroShiftROne := zero.ShiftRight(1)
			Ω(zeroShiftROne).Should(Equal(zero))

			elevenShiftRTwo := eleven.ShiftRight(2)
			Ω(elevenShiftRTwo).Should(Equal(two))
		})

		It("should overflow without wrapping", func() {
			oneShiftROne := one.ShiftRight(1)
			Ω(oneShiftROne).Should(Equal(zero))
		})
	})

	Context("when taking the bitwise AND", func() {
		It("should return the right result for 1024 bit numbers", func() {

			RunAllCases(andFn, []TestCase{
				TestCase{inputsStr: []string{"0", "0"}, expectedStr: "0"},
				TestCase{inputsStr: []string{"0", "0"}, expectedStr: "0"},
			})

			maxANDMax := max.AND(&max)
			Ω(maxANDMax).Should(Equal(max))

			sevenANDEleven := seven.AND(&eleven)
			Ω(sevenANDEleven).Should(Equal(three))
		})
	})

	Context("when taking the bitwise XOR", func() {
		It("should return the right result for 1024 bit numbers", func() {

			maxXORMax := max.XOR(&max)
			Ω(maxXORMax).Should(Equal(zero))

			sevenXOREleven := seven.XOR(&eleven)
			Ω(sevenXOREleven).Should(Equal(twelve))
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
