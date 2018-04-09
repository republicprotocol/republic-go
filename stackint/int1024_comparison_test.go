package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// . "github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Int1024 bitwise operations", func() {
	Context("when calculating x=y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Ω(one.Equals(&one)).Should(BeTrue())
			Ω(zero.Equals(&zero)).Should(BeTrue())
			Ω(max.Equals(&max)).Should(BeTrue())

			Ω(one.Equals(&two)).Should(BeFalse())
			Ω(two.Equals(&one)).Should(BeFalse())

			Ω(zero.Equals(&max)).Should(BeFalse())
			Ω(max.Equals(&zero)).Should(BeFalse())
		})
	})

	Context("when calculating x=0", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Ω(one.IsZero()).Should(BeFalse())
			Ω(zero.IsZero()).Should(BeTrue())
			Ω(max.IsZero()).Should(BeFalse())
		})
	})

	Context("when calculating x<y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Ω(one.LessThan(&one)).Should(BeFalse())
			Ω(zero.LessThan(&zero)).Should(BeFalse())
			Ω(max.LessThan(&max)).Should(BeFalse())

			Ω(one.LessThan(&two)).Should(BeTrue())
			Ω(two.LessThan(&one)).Should(BeFalse())

			Ω(zero.LessThan(&max)).Should(BeTrue())
			Ω(max.LessThan(&zero)).Should(BeFalse())

			Ω(zero.LessThan(&two64)).Should(BeTrue())
			Ω(two64.LessThan(&zero)).Should(BeFalse())
		})
	})

	Context("when calculating x<=y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Ω(one.LessThanOrEqual(&one)).Should(BeTrue())
			Ω(zero.LessThanOrEqual(&zero)).Should(BeTrue())
			Ω(max.LessThanOrEqual(&max)).Should(BeTrue())

			Ω(one.LessThanOrEqual(&two)).Should(BeTrue())
			Ω(two.LessThanOrEqual(&one)).Should(BeFalse())

			Ω(zero.LessThanOrEqual(&max)).Should(BeTrue())
			Ω(max.LessThanOrEqual(&zero)).Should(BeFalse())

			Ω(zero.LessThanOrEqual(&two64)).Should(BeTrue())
			Ω(two64.LessThanOrEqual(&zero)).Should(BeFalse())
		})
	})

	Context("when calculating x>y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Ω(one.GreaterThan(&one)).Should(BeFalse())
			Ω(zero.GreaterThan(&zero)).Should(BeFalse())
			Ω(max.GreaterThan(&max)).Should(BeFalse())

			Ω(one.GreaterThan(&two)).Should(BeFalse())
			Ω(two.GreaterThan(&one)).Should(BeTrue())

			Ω(zero.GreaterThan(&max)).Should(BeFalse())
			Ω(max.GreaterThan(&zero)).Should(BeTrue())
		})
	})

	Context("when calculating x>=y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Ω(one.GreaterThanOrEqual(&one)).Should(BeTrue())
			Ω(zero.GreaterThanOrEqual(&zero)).Should(BeTrue())
			Ω(max.GreaterThanOrEqual(&max)).Should(BeTrue())

			Ω(one.GreaterThanOrEqual(&two)).Should(BeFalse())
			Ω(two.GreaterThanOrEqual(&one)).Should(BeTrue())

			Ω(zero.GreaterThanOrEqual(&max)).Should(BeFalse())
			Ω(max.GreaterThanOrEqual(&zero)).Should(BeTrue())
		})
	})

	Context("is even", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Ω(zero.IsEven()).Should(BeTrue())
			Ω(one.IsEven()).Should(BeFalse())
			Ω(two.IsEven()).Should(BeTrue())
			Ω(three.IsEven()).Should(BeFalse())
			Ω(max.IsEven()).Should(BeFalse())
		})
	})
})
