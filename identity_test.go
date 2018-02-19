package identity_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
)

var _ = Describe("", func() {

	Describe("Republic IDs", func() {
		Context("generated from random key pairs", func() {
			id, _, err := identity.NewID()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should return 20 bytes", func() {
				Ω(len(id)).Should(Equal(identity.IDLength))
			})

		})

		Context("converting to string", func() {
			id, _, err := identity.NewID()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should be converted to a string", func() {
				stringID := id.String()
				Ω(len(stringID)).Should(Equal(identity.IDLength))
			})
		})

		Context("converting to ID", func() {
			It("should be converted to an Address", func() {
				id, _, err := identity.NewID()
				Ω(err).ShouldNot(HaveOccurred())
				address := id.Address()
				newID := address.ID()
				Ω(id).Should(Equal(newID))
			})
		})
	})
})
