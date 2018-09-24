package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("AES Cipher", func() {

	Context("when generating", func() {

		It("should be able to generate a random EcdsaKey without returning an error", func() {
			_, err := RandomAESCipher()
			Expect(err).ShouldNot(HaveOccurred())
		})

	})

	Context("when encrypting and decrypting", func() {

		It("should be able to encrypt and decrypt text with random AesCiper", func() {
			secrets := []string{
				"republicprotocol",
				"18514",
				"!@#$%^&*()",
			}

			aes, err := RandomAESCipher()
			Expect(err).ShouldNot(HaveOccurred())

			for _, secret := range secrets {
				cipherText, err := aes.Encrypt([]byte(secret))
				Expect(err).ShouldNot(HaveOccurred())

				plainText, err := aes.Decrypt(cipherText)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(string(plainText)).Should(Equal(secret))
			}
		})

		It("should be able to encrypt and decrypt text with provided AesCiper", func() {
			// Create a secret of the AesCipher
			aesCiperSecret := [16]byte{}
			for i := 0; i < 16; i++ {
				aesCiperSecret[i] = byte(i)
			}
			aes := NewAESCipher(aesCiperSecret[:])

			secrets := []string{
				"republicprotocol",
				"18514",
				"!@#$%^&*()",
			}
			for _, secret := range secrets {
				cipherText, err := aes.Encrypt([]byte(secret))
				Expect(err).ShouldNot(HaveOccurred())

				plainText, err := aes.Decrypt(cipherText)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(string(plainText)).Should(Equal(secret))
			}
		})
	})

	Context("when giving a wrong secret for the aescipher", func() {

		It("should error when trying to encrypt with a cipher having incorrect number of bytes", func() {
			plainText := "republicprotocol"

			// Create a secret of the AesCipher
			for i := 0; i < 64; i++ {
				aesCiperSecret := []byte{}
				for j := 0; j < i; j++ {
					aesCiperSecret = append(aesCiperSecret, byte(j))
				}
				aes := NewAESCipher(aesCiperSecret[:])
				cipherBytes, err := aes.Encrypt([]byte(plainText))

				switch i {
				case 16, 24, 32:
					Expect(err).ShouldNot(HaveOccurred())
					plainBytes, err := aes.Decrypt(cipherBytes)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(string(plainBytes)).Should(Equal(plainText))

				default:
					Expect(err).Should(HaveOccurred())
					_, err := aes.Decrypt(cipherBytes)
					Expect(err).Should(HaveOccurred())
				}
			}
		})
	})
})
