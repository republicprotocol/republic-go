package identity_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("", func() {

	Describe("Republic IDs", func() {
		Context("generated from random key pairs", func() {
			key, err := crypto.RandomEcdsaKey()
			id := identity.Address(key.Address()).ID()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should return 20 bytes", func() {
				Ω(len(id)).Should(Equal(identity.IDLength))
			})

		})

		Context("converting to string", func() {
			key, err := crypto.RandomEcdsaKey()
			id := identity.Address(key.Address()).ID()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should be converted to a string", func() {
				Ω(id.String()).Should(Equal(id.Address().String()))
			})
		})

		Context("converting to ID", func() {
			It("should be converted to an Address", func() {
				key, err := crypto.RandomEcdsaKey()
				id := identity.Address(key.Address()).ID()

				Ω(err).ShouldNot(HaveOccurred())
				address := id.Address()
				newID := address.ID()
				Ω(id).Should(Equal(newID))
			})
		})
	})
})
