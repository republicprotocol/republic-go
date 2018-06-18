package crypto_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/crypto"

	"github.com/republicprotocol/republic-go/dispatch"
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
			dispatch.CoForAll(runtime.NumCPU(), func(i int) {
				key, err := RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				data := Keccak256([]byte("REN"))
				signature, err := key.Sign(data)
				Expect(err).ShouldNot(HaveOccurred())

				err = key.Verify(data, signature)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		It("should be able to return an error when verifying random data", func() {
			dispatch.CoForAll(runtime.NumCPU(), func(i int) {
				key, err := RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				random := make([]byte, 32)
				rand.Read(random)

				sigRandom := make([]byte, 65)
				rand.Read(sigRandom)

				err = key.Verify(random, sigRandom)
				Expect(err).Should(HaveOccurred())
			})
		})

		It("should be able to return an error when verifying nil data", func() {
			dispatch.CoForAll(runtime.NumCPU(), func(i int) {
				key, err := RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				random := make([]byte, 32)
				rand.Read(random)

				sigRandom := make([]byte, 65)
				rand.Read(sigRandom)

				err = key.Verify([]byte{}, sigRandom)
				Expect(err).Should(Equal(ErrNilData))

				err = key.Verify(nil, sigRandom)
				Expect(err).Should(Equal(ErrNilData))
			})
		})

		It("should be able to return an error when verifying nil signatures", func() {
			dispatch.CoForAll(runtime.NumCPU(), func(i int) {
				key, err := RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				random := make([]byte, 32)
				rand.Read(random)

				sigRandom := make([]byte, 65)
				rand.Read(sigRandom)

				err = key.Verify(random, []byte{})
				Expect(err).Should(Equal(ErrNilSignature))

				err = key.Verify(random, nil)
				Expect(err).Should(Equal(ErrNilSignature))
			})
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

		It("should panic when marshalling a badly formatted JSON", func() {
			key, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())
			key.PrivateKey.PublicKey = ecdsa.PublicKey{}
			Expect(func() { json.Marshal(key) }).Should(Panic())
		})

		It("should return an error for invalid JSON", func() {
			keyDecoded := EcdsaKey{}
			err := keyDecoded.UnmarshalJSON([]byte{byte(1)})
			Expect(err).Should(HaveOccurred())
			Expect(keyDecoded.Address()).Should(Equal(""))
		})

	})
})
