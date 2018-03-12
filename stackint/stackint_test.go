package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var one = Int1024FromUint64(1)
var two = Int1024FromUint64(2)
var three = Int1024FromUint64(3)

var _ = Describe("Stackint", func() {
	Context("when adding two numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			oneplustwo := one.Add(&two)
			Ω(oneplustwo.Equals(&three)).Should(BeTrue())

			oneword := Int1024FromUint64(MAXUINT64)
			onewordplusone := oneword.Add(&one)
			Ω(onewordplusone.Words()[14]).Should(Equal(uint64(1)))
		})
	})
})
