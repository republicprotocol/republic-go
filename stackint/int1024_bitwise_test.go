package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Int1024 bitwise operations", func() {
	Context("when left-shifting a number", func() {
		It("should return the right result for 1024 bit numbers", func() {
			oneShiftLOne := one.ShiftLeft(1)
			Ω(oneShiftLOne.Equals(&two)).Should(BeTrue())

			oneShiftL1023 := one.ShiftLeft(1023)
			Ω(oneShiftL1023.Equals(&twoPow1023)).Should(BeTrue())

			zeroShiftLOne := zero.ShiftLeft(1)
			Ω(zeroShiftLOne.Equals(&zero)).Should(BeTrue())
		})

		It("should overflow without wrapping", func() {
			two1023ShiftLOne := twoPow1023.ShiftLeft(1)
			Ω(two1023ShiftLOne.Equals(&zero)).Should(BeTrue())
		})
	})

	Context("when right-shifting a number", func() {
		It("should return the right result for 1024 bit numbers", func() {
			twoShiftROne := two.ShiftRight(1)
			Ω(twoShiftROne.Equals(&one)).Should(BeTrue())

			zeroShiftROne := zero.ShiftRight(1)
			Ω(zeroShiftROne.Equals(&zero)).Should(BeTrue())
		})

		It("should overflow without wrapping", func() {
			oneShiftROne := one.ShiftRight(1)
			Ω(oneShiftROne.Equals(&zero)).Should(BeTrue())
		})
	})

	Context("when taking the bitwise AND", func() {
		It("should return the right result for 1024 bit numbers", func() {
			maxANDMax := max.AND(&max)
			Ω(maxANDMax.Equals(&max)).Should(BeTrue())

			sevenANDEleven := seven.AND(&eleven)
			Ω(sevenANDEleven.Equals(&three)).Should(BeTrue())
		})
	})

	Context("when taking the bitwise XOR", func() {
		It("should return the right result for 1024 bit numbers", func() {
			maxXORMax := max.XOR(&max)
			Ω(maxXORMax.Equals(&zero)).Should(BeTrue())

			sevenXOREleven := seven.XOR(&eleven)
			Ω(sevenXOREleven.Equals(&twelve)).Should(BeTrue())
		})
	})
})
