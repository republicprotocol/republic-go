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

			keyDecoded := new(RsaKey)
			err = keyDecoded.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(key.D).Should(Equal(keyDecoded.D))
			Expect(key.Primes).Should(HaveLen(len(keyDecoded.Primes)))
			for i := range key.Primes {
				Expect(key.Primes[i]).Should(Equal(keyDecoded.Primes[i]))
			}
			Expect(key.N).Should(Equal(keyDecoded.N))
			Expect(key.E).Should(Equal(keyDecoded.E))
		})

		It("should be able to marshal and unmarshal public keys as bytes", func() {
			key, err := RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			data := BytesFromRsaPublicKey(&key.PublicKey)
			publicKeyDecoded := RsaPublicKeyFromBytes(data)

			Expect(key.N).Should(Equal(publicKeyDecoded.N))
			Expect(key.E).Should(Equal(publicKeyDecoded.E))
		})

	})
})
