package darkocean_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	. "github.com/republicprotocol/republic-go/darkocean"
)

var _ = Describe("Crypter", func() {

	var crypter Crypter
	var message crypto.Hash32

	BeforeEach(func() {
		keystore, err := crypto.RandomKeystore()
		Expect(err).ShouldNot(HaveOccurred())
		crypter = NewCrypter(keystore, testnetEnv.DarknodeRegistry, NumberOfBootstrapDarkNodes, time.Second)
		message = crypto.NewHash32([]byte("REN"))
	})

	Context("when signing", func() {

		It("should produce valid signatures", func() {
			signature, err := crypter.Sign(message)
			Expect(err).ShouldNot(HaveOccurred())
			err = crypto.VerifySignature(message, signature, crypter.Keystore().Address())
			Expect(err).ShouldNot(HaveOccurred())
		})

	})

	Context("when verifying signatures", func() {

		It("should return an error for unregistered addresses", func() {
			keystore, err := crypto.RandomKeystore()
			Expect(err).ShouldNot(HaveOccurred())

			signature, err := keystore.Sign(message)
			Expect(err).ShouldNot(HaveOccurred())
			err = crypter.Verify(message, signature)
			Expect(err).Should(HaveOccurred())
		})

		It("should not return an error for registered addresses", func() {
		})

	})

	Context("when encrypting", func() {

		It("should encrypt messages for registered addresses", func() {
		})

		It("should not encrypt messages for unregistered addresses", func() {
		})

	})

	Context("when decrypting", func() {

		It("should produce the original plain text", func() {
		})

	})

	Context("when caching", func() {

		It("should update registrations after the update period", func() {
		})

		It("should update public keys after the update period", func() {
		})

	})
})
