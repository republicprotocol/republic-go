package identity_test

import (
	"github.com/republicprotocol/go-identity"
	. "github.com/republicprotocol/go-identity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MultiAddresses with support for Republic Protocol", func() {

	Context("after importing the package", func() {

		It("should expose a protocol called republic", func() {
			Ω(ProtocolWithName("republic").Name).Should(Equal("republic"))
		})

		It("should expose a protocol with the correct constant values", func() {
			Ω(ProtocolWithCode(RepublicCode).Name).Should(Equal("republic"))
			Ω(ProtocolWithCode(RepublicCode).Code).Should(Equal(RepublicCode))
		})
	})

	Context("marshaling to JSON", func() {
		It("should encode and then decode to the same value", func() {
			addr, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			multi, err := identity.NewMultiAddressFromString("/republic/" + addr.String())
			Ω(err).ShouldNot(HaveOccurred())
			data, err := multi.MarshalJSON()
			Ω(err).ShouldNot(HaveOccurred())
			newMulti := &identity.MultiAddress{}
			err = newMulti.UnmarshalJSON(data)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multi.String()).Should(Equal(newMulti.String()))
		})
	})
})
