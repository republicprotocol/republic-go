package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Stackint", func() {
	Context("when adding two numbers", func() {
		It("should return the right result for 1024 bit numbers", func() {
			x := NewInt1024(1)
			y := NewInt1024(2)
			actual := x.Add(&y)
			expected := NewInt1024(3)
			Ω(actual.Equals(&expected)).Should(BeTrue())
		})
	})
})
