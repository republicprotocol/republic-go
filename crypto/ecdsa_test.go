package crypto_test

import (
	"crypto/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("Ecdsa keys", func() {

	Context("when generating", func() {

		It("should be able to generate a random EcdsaKey without returning an error", func() {
			_, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should equal itself", func() {
			key, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(key.Equal(&key)).Should(BeTrue())
		})

		It("should not equal another randomly generated EcdsaKey", func() {
			key1, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			key2, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(key1.Equal(&key2)).Should(BeFalse())
		})
	})

	Context("when signing and verifying", func() {

		It("should be able to sign a hash", func() {
			key, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			_, err = key.Sign(Keccak256([]byte("REN")))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be able to verify a signature", func() {
			for i := 0; i < 1000; i++ {
				key, err := RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				data := Keccak256([]byte("REN"))
				signature, err := key.Sign(data)
				Expect(err).ShouldNot(HaveOccurred())

				err = key.Verify(data, signature)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		It("should be able to return an error when verifying random data", func() {
			for i := 0; i < 1000; i++ {
				key, err := RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				random := make([]byte, 32)
				rand.Read(random)

				sigRandom := make([]byte, 65)
				rand.Read(sigRandom)

				err = key.Verify(random, sigRandom)
				Expect(err).Should(HaveOccurred())
			}
		})

	})

	Context("when marshaling and unmarshaling", func() {

		It("should be able to marshal and unmarshal as JSON", func() {
			key, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			data, err := key.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())

			keyDecoded := EcdsaKey{}
			err = keyDecoded.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(key.Equal(&keyDecoded)).Should(BeTrue())
			Expect("s256").Should(Equal(keyDecoded.Curve.Params().Name)) // We explicitly name the curve here because the ethSecp256k1.S256() curve implementation does not include a name
		})

	})
})
