package x_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hashes", func() {

	Context("when comparing hashes", func() {
		Ω(false).Should(Equal(true))
	})

})
