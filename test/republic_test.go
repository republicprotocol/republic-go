package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Republic Protocol Test", func() {
	Context("large integration test", func() {
		It("should finish the test without any error", func() {
			Ω(true).Should(Equal(true))
			test()
		})
	})
})
