package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("RSA Test package", func() {

	Context("Testing type conversions", func() {

		It("should be able to encode and decode a publicKey", func() {
			keyPair, err := NewRsaKeyPair()
			Ω(err).Should(BeNil())
			pubKeyBytes, err := PublicKeyToBytes(keyPair.PublicKey)
			Ω(err).Should(BeNil())
			pubKey, err := BytesToPublicKey(pubKeyBytes)
			Ω(err).Should(BeNil())
			Ω(*pubKey).Should(Equal(*keyPair.PublicKey))
		})

		It("should be able to encryot and decrypt a message", func() {
			keyPair, err := NewRsaKeyPair()
			encryptionMsg := []byte("Message")
			Ω(err).Should(BeNil())
			cipherText, err := Encrypt(keyPair.PublicKey, encryptionMsg)
			Ω(err).Should(BeNil())
			decryptionMsg, err := Decrypt(keyPair.PrivateKey, cipherText)
			Ω(err).Should(BeNil())
			Ω(encryptionMsg).Should(Equal(decryptionMsg))
		})

	})
})
