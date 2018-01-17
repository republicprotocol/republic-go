package compute_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Results and result fragments", func() {
	It("should fail", func() {
		Ω(true).Should(Equal(false))
	})
})
