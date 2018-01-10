package x_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {

	Describe("Well formed configuration files", func() {
		It("should be finished in the future", func() {
			Ω(false).Should(Equal(true))
		})
	})

	Describe("Malformed configuration files", func() {
		It("should be finished in the future", func() {
			Ω(false).Should(Equal(true))
		})
	})

})
