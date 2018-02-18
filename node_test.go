package node_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dark nodes", func() {
	It("should fail", func() {
		Ω(true).Should(Equal(false))
	})
})
