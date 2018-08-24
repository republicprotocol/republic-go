package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// . "github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Int1024 bitwise operations", func() {
	Context("when calculating x=y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Expect(one.Equals(&one)).Should(BeTrue())
			Expect(zero.Equals(&zero)).Should(BeTrue())
			Expect(max.Equals(&max)).Should(BeTrue())

			Expect(one.Equals(&two)).Should(BeFalse())
			Expect(two.Equals(&one)).Should(BeFalse())

			Expect(zero.Equals(&max)).Should(BeFalse())
			Expect(max.Equals(&zero)).Should(BeFalse())
		})
	})

	Context("when calculating x=0", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Expect(one.IsZero()).Should(BeFalse())
			Expect(zero.IsZero()).Should(BeTrue())
			Expect(max.IsZero()).Should(BeFalse())
		})
	})

	Context("when calculating x<y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Expect(one.LessThan(&one)).Should(BeFalse())
			Expect(zero.LessThan(&zero)).Should(BeFalse())
			Expect(max.LessThan(&max)).Should(BeFalse())

			Expect(one.LessThan(&two)).Should(BeTrue())
			Expect(two.LessThan(&one)).Should(BeFalse())

			Expect(zero.LessThan(&max)).Should(BeTrue())
			Expect(max.LessThan(&zero)).Should(BeFalse())

			Expect(zero.LessThan(&two64)).Should(BeTrue())
			Expect(two64.LessThan(&zero)).Should(BeFalse())
		})
	})

	Context("when calculating x<=y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Expect(one.LessThanOrEqual(&one)).Should(BeTrue())
			Expect(zero.LessThanOrEqual(&zero)).Should(BeTrue())
			Expect(max.LessThanOrEqual(&max)).Should(BeTrue())

			Expect(one.LessThanOrEqual(&two)).Should(BeTrue())
			Expect(two.LessThanOrEqual(&one)).Should(BeFalse())

			Expect(zero.LessThanOrEqual(&max)).Should(BeTrue())
			Expect(max.LessThanOrEqual(&zero)).Should(BeFalse())

			Expect(zero.LessThanOrEqual(&two64)).Should(BeTrue())
			Expect(two64.LessThanOrEqual(&zero)).Should(BeFalse())
		})
	})

	Context("when calculating x>y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Expect(one.GreaterThan(&one)).Should(BeFalse())
			Expect(zero.GreaterThan(&zero)).Should(BeFalse())
			Expect(max.GreaterThan(&max)).Should(BeFalse())

			Expect(one.GreaterThan(&two)).Should(BeFalse())
			Expect(two.GreaterThan(&one)).Should(BeTrue())

			Expect(zero.GreaterThan(&max)).Should(BeFalse())
			Expect(max.GreaterThan(&zero)).Should(BeTrue())
		})
	})

	Context("when calculating x>=y", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Expect(one.GreaterThanOrEqual(&one)).Should(BeTrue())
			Expect(zero.GreaterThanOrEqual(&zero)).Should(BeTrue())
			Expect(max.GreaterThanOrEqual(&max)).Should(BeTrue())

			Expect(one.GreaterThanOrEqual(&two)).Should(BeFalse())
			Expect(two.GreaterThanOrEqual(&one)).Should(BeTrue())

			Expect(zero.GreaterThanOrEqual(&max)).Should(BeFalse())
			Expect(max.GreaterThanOrEqual(&zero)).Should(BeTrue())
		})
	})

	Context("is even", func() {
		It("should return the right result for 1024 bit numbers", func() {
			Expect(zero.IsEven()).Should(BeTrue())
			Expect(one.IsEven()).Should(BeFalse())
			Expect(two.IsEven()).Should(BeTrue())
			Expect(three.IsEven()).Should(BeFalse())
			Expect(max.IsEven()).Should(BeFalse())
		})
	})
})
