package darknode_test

import (
	"bytes"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/darknode"

	"github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("Crypter", func() {

	var crypter Crypter
	var message crypto.Hash32

	BeforeEach(func() {
		keystore, err := crypto.RandomKeystore()
		Expect(err).ShouldNot(HaveOccurred())
		crypter = NewCrypter(keystore, testnetEnv.DarknodeRegistry, NumberOfDarkNodes/2, time.Second)
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
			for i := 0; i < 100; i++ {
				keystore, err := crypto.RandomKeystore()
				Expect(err).ShouldNot(HaveOccurred())

				signature, err := keystore.Sign(message)
				Expect(err).ShouldNot(HaveOccurred())
				err = crypter.Verify(message, signature)
				Expect(err).Should(HaveOccurred())
			}
		})

		It("should not return an error for registered addresses", func() {
			for _, darknode := range testnetEnv.Darknodes {
				signature, err := darknode.Config.Keystore.Sign(message)
				Expect(err).ShouldNot(HaveOccurred())
				err = crypter.Verify(message, signature)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

	})

	Context("when encrypting", func() {

		It("should encrypt messages for registered addresses", func() {
			for _, darknode := range testnetEnv.Darknodes {
				cipherText, err := crypter.Encrypt(darknode.Address().String(), message[:])
				Expect(err).ShouldNot(HaveOccurred())

				plainText, err := darknode.Config.Keystore.Decrypt(cipherText)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Equal(plainText, message[:])).Should(BeTrue())
			}
		})

		It("should not encrypt messages for unregistered addresses", func() {
			for i := 0; i < 100; i++ {
				keystore, err := crypto.RandomKeystore()
				Expect(err).ShouldNot(HaveOccurred())

				_, err = crypter.Encrypt(keystore.Address(), message[:])
				Expect(err).Should(HaveOccurred())
			}
		})

	})

	Context("when decrypting", func() {

		It("should produce the original plain text", func() {
			cipherText, err := crypter.Keystore().Encrypt(message[:])
			Expect(err).ShouldNot(HaveOccurred())
			plainText, err := crypter.Decrypt(cipherText)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bytes.Equal(plainText, message[:])).Should(BeTrue())
		})

	})

	Context("when caching", func() {

		It("should update registrations after the update period", func() {
		})

		It("should update public keys after the update period", func() {
		})

	})
})
