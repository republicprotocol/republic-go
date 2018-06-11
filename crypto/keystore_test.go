package crypto_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("Keystore", func() {

	Context("when generating", func() {
		It("should be able to generate a random Keystore without returning an error", func() {
			_, err := RandomKeystore()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when encrypting and decrypting", func() {

		It("should be able to encrypt into JSON using a passphrase", func() {
			keystore, err := RandomKeystore()
			Expect(err).ShouldNot(HaveOccurred())

			_, err = keystore.EncryptToJSON("password", StandardScryptN, StandardScryptP)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be able to decrypt from JSON using a passphrase", func() {
			keystore, err := RandomKeystore()
			Expect(err).ShouldNot(HaveOccurred())

			data, err := keystore.EncryptToJSON("password", StandardScryptN, StandardScryptP)
			Expect(err).ShouldNot(HaveOccurred())

			keystoreDecrypted := new(Keystore)
			err = keystoreDecrypted.DecryptFromJSON(data, "password")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(keystore.ID).Should(Equal(keystoreDecrypted.ID))
			Expect(keystore.Version).Should(Equal(keystoreDecrypted.Version))
			Expect(keystore.EcdsaKey.Equal(&keystoreDecrypted.EcdsaKey)).Should(BeTrue())
			Expect(keystore.RsaKey.Equal(&keystoreDecrypted.RsaKey)).Should(BeTrue())
		})

		It("should not be able to decrypt from JSON using an incorrect passphrase", func() {
			keystore, err := RandomKeystore()
			Expect(err).ShouldNot(HaveOccurred())

			data, err := keystore.EncryptToJSON("password", StandardScryptN, StandardScryptP)
			Expect(err).ShouldNot(HaveOccurred())

			keystoreDecrypted := new(Keystore)
			err = keystoreDecrypted.DecryptFromJSON(data, "badpassword")
			Expect(err).Should(Equal(ErrPassphraseCannotDecryptKey))
		})

	})

	Context("when marshaling and unmarshaling", func() {

		It("should be able to marshal and unmarshal as JSON", func() {
			keystore, err := RandomKeystore()
			Expect(err).ShouldNot(HaveOccurred())

			data, err := json.Marshal(keystore)
			Expect(err).ShouldNot(HaveOccurred())

			keystoreDecoded := Keystore{}
			err = json.Unmarshal(data, &keystoreDecoded)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(keystore.ID).Should(Equal(keystoreDecoded.ID))
			Expect(keystore.Version).Should(Equal(keystoreDecoded.Version))
			Expect(keystore.EcdsaKey.Equal(&keystoreDecoded.EcdsaKey)).Should(BeTrue())
			Expect(keystore.RsaKey.Equal(&keystoreDecoded.RsaKey)).Should(BeTrue())
		})

	})
})
