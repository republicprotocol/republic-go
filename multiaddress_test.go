package identity_test

import (
	. "github.com/republicprotocol/go-identity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Multiaddress with republic protocol added", func() {

	Context("After import the indentity package", func() {
		It("should have a protocol called republic", func() {
			Ω(ProtocolWithName("republic").Name).Should(Equal("republic"))
		})

		Specify("The republic protocol code should be defined as a constant ", func() {
			Ω(ProtocolWithCode(P_REPUBLIC).Name).Should(Equal("republic"))
			Ω(ProtocolWithCode(P_REPUBLIC).Code).Should(Equal(P_REPUBLIC))
		})
	})
})
