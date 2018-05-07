package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("RsaKey", func() {

	Context("when generating", func() {

		It("should be able to generate a random RsaKey without returning an error", func() {
			_, err := RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should equal itself", func() {
			key, err := RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(key.Equal(&key)).Should(BeTrue())
		})

		It("should not equal another randomly generated RsaKey", func() {
			key1, err := RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			key2, err := RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(key1.Equal(&key2)).Should(BeFalse())
		})
	})

	Context("when encrypting and decrypting", func() {

		It("should be able to encrypt a plain text message", func() {
			key, err := RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			plainText := []byte("REN")
			_, err = key.Encrypt(plainText)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be able to decrypt an encrypted cipher text", func() {
			key, err := RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			plainText := []byte("REN")
			cipherText, err := key.Encrypt(plainText)
			Expect(err).ShouldNot(HaveOccurred())
			plainTextDecrypted, err := key.Decrypt(cipherText)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(plainText).Should(Equal(plainTextDecrypted))
		})

	})

	Context("when marshaling and unmarshaling", func() {

		It("should be able to marshal and unmarshal as JSON", func() {
			key, err := RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			data, err := key.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())

			keyDecoded := RsaKey{}
			err = keyDecoded.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(key.Equal(&keyDecoded)).Should(BeTrue())
		})

		It("should be able to marshal and unmarshal public keys as bytes", func() {
			key, err := RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			data, err := BytesFromRsaPublicKey(&key.PublicKey)
			Expect(err).ShouldNot(HaveOccurred())

			publicKeyDecoded, err := RsaPublicKeyFromBytes(data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(key.N).Should(Equal(publicKeyDecoded.N))
			Expect(key.E).Should(Equal(publicKeyDecoded.E))
		})

	})
})
