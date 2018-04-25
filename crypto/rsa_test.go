package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("RSA Test package", func() {

	Context("Testing type conversions", func() {

		It("should be able to encode and decode a publicKey", func() {
			key, err := RandomRsaKey()
			Ω(err).Should(BeNil())
			pubKeyBytes := PublicKeyToBytes(&key.PublicKey)
			pubKey := BytesToPublicKey(pubKeyBytes)
			Ω(pubKey).Should(Equal(key.PublicKey))
		})

		It("should be able to encryot and decrypt a message", func() {
			key, err := RandomRsaKey()
			plaintext := []byte("Loong is bigdragon")
			Ω(err).Should(BeNil())
			cipherText, err := key.Encrypt(plaintext)
			Ω(err).Should(BeNil())
			decryptedPlaintext, err := key.Decrypt(cipherText)
			Ω(err).Should(BeNil())
			Ω(plaintext).Should(Equal(decryptedPlaintext))
		})

	})
})
