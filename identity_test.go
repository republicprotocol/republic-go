package identity_test

import (
	"github.com/jbenet/go-base58"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
)

var _ = Describe("Republic identity", func() {

	Context("generating a new key pair", func() {
		keyPair, err := identity.NewKeyPair()

		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should have non-nil private key and public key", func() {
			Ω(keyPair.PrivateKey).ShouldNot(BeNil())
			Ω(keyPair.PublicKey).ShouldNot(BeNil())
		})
	})

	Context("getting the ID from a key pair", func() {
		keyPair, _ := identity.NewKeyPair()
		id := keyPair.PublicID()

		It("should return 20 bytes", func() {
			Ω(len(id)).Should(Equal(identity.IDLength))
		})
	})

	Context("getting the address from a key pair", func() {
		keyPair, _ := identity.NewKeyPair()
		address := keyPair.PublicAddress()
		decoded := base58.Decode(string(address))

		It("should have 0x1B as its first byte", func() {
			Ω(decoded[0]).Should(Equal(uint8(0x1B)))
		})
		It("should have 0x14 as its second byte", func() {
			Ω(decoded[1]).Should(Equal(uint8(identity.IDLength)))
		})
		It("should be a base58 encoding of its public ID after the first two bytes", func() {
			Ω(decoded[2:]).Should(Equal([]byte(keyPair.PublicID())))
		})
	})

	Context("getting a multiaddress from a key pair", func() {
		keyPair, _ := identity.NewKeyPair()
		multiaddress, err := keyPair.MultiAddress()

		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should be a string concatenated by '/republic/' and its republic address", func() {
			Ω(multiaddress.String()).Should(Equal("/republic/" + string(keyPair.PublicAddress())))
		})
	})
})
