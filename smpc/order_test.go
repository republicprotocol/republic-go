package smpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Order fragment processor", func() {
	Context("when receiving order fragments", func() {

		It("should shutdown when the context is canceled", func() {
			Ω(true).Should(BeFalse())
		})

		It("should produce computations for all order fragment pairs", func() {
			Ω(true).Should(BeFalse())
		})

	})
})
