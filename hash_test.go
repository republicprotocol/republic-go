package x_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hashes", func() {

	Context("when comparing hashes", func() {
		It("should be true", func() {
			Ω(true).Should(Equal(true))
		})
	})

})
