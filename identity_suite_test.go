package identity_test

import (
	"testing"

	"github.com/jbenet/go-base58"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
)

func TestAddress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Republic Identity Suite")
}

var _ = Describe("Republic identity", func() {

	Context("when generating a key pair", func() {
		keyPair, err := identity.NewKeyPair()
		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should have non-nil private-key and public key", func() {
			Ω(keyPair.PrivateKey).ShouldNot(BeNil())
			Ω(keyPair.PublicKey).ShouldNot(BeNil())
		})
	})

	Context("get republic ID from the key pair", func() {
		keyPair, err := identity.NewKeyPair()
		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		ID := keyPair.PublicID()
		It("should have 20 bytes", func() {
			Ω(len(ID)).Should(Equal(identity.IDLength))
		})
	})

	Context("get republic address from the key pair", func() {
		keyPair, err := identity.NewKeyPair()
		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		address := keyPair.PublicAddress()
		decoded := base58.Decode(address)
		It("should have 0x1B as its first byte", func() {
			Ω(decoded[0]).Should(Equal(uint8(0x1B)))
		})
		It("should have 0x14 as its second byte", func() {
			Ω(decoded[1]).Should(Equal(uint8(identity.IDLength)))
		})
		It("should be base58 encoding of its public ID aftre the first two bytes", func() {
			Ω(decoded[2:]).Should(Equal(keyPair.PublicID()))
		})
	})

	Context("get republic multiaddress from the key pair", func() {
		keyPair, err := identity.NewKeyPair()
		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		multiaddress, err := keyPair.MultiAddress()
		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should be a string concatenated by '/republic/' and its republic address", func() {
			Ω(multiaddress.String()).Should(Equal("/republic/" + keyPair.PublicAddress()))
		})
	})
})
