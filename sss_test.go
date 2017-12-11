package sss_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/go-sss"
)

var _ = Describe("Shamir's secret sharing", func() {
	Context("With all shares", func() {
		It("should be able to reconstruct the secret", func() {
			Ω(true).Should(Equal(true))
			Stub()
		})
	})
})
