package identity_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
)

var _ = Describe("", func() {

	Describe("Key pair", func() {
		Context("generation", func() {
			keyPair, err := identity.NewKeyPair()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should have non-nil private key and public key", func() {
				Ω(keyPair.PrivateKey).ShouldNot(BeNil())
				Ω(keyPair.PublicKey).ShouldNot(BeNil())
			})
		})

		Context("generation from private key", func() {
			var key identity.KeyPair
			privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should generate key pair from the private key", func() {
				key, err = identity.NewKeyPairFromPrivateKey(privateKey)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should have the same private and public key as the 'father' private key", func() {
				Ω(*privateKey).Should(Equal(*(key.PrivateKey)))
				Ω(privateKey.Public()).Should(Equal(key.PublicKey))
			})
		})

		Context("IDs", func() {
			keyPair, err := identity.NewKeyPair()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should return 20 bytes", func() {
				id := keyPair.ID()
				Ω(len(id)).Should(Equal(identity.IDLength))
			})
		})

		Context("addresses", func() {
			keyPair, err := identity.NewKeyPair()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			address := keyPair.Address()

			It("should have a length of 20 bytes", func() {
				Ω(len(address)).Should(Equal(identity.AddressLength))
			})

			decoded := base58.Decode(string(address))

			It("should not decode to the empty string", func() {
				Ω(decoded).ShouldNot(BeEmpty())
			})
			It("should have 0x1B as their first byte", func() {
				Ω(decoded[0]).Should(Equal(uint8(0x1B)))
			})
			It("should have 0x14 as their second byte", func() {
				Ω(decoded[1]).Should(Equal(uint8(identity.IDLength)))
			})
			It("should be a base58 encoding of the ID after the first two bytes", func() {
				Ω(decoded[2:]).Should(Equal([]byte(keyPair.ID())))
			})
		})

	})
})
