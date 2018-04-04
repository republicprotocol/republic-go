package smpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ξ-fragment producers", func() {
	Context("when producing ξ-fragments", func() {

		It("should shutdown when the context is canceled", func() {
			Ω(true).Should(BeFalse())
		})

		It("should produce ξ-fragment", func() {
			Ω(true).Should(BeFalse())
		})

	})
})
