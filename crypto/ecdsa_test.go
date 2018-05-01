package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("EcdsaKey", func() {

	Context("when generating", func() {
		It("should be able to generate a random EcdsaKey without returning an error", func() {
			_, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when signing and verifyinh", func() {

		It("should be able to sign a hash", func() {
			key, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			hash32 := NewHash32([]byte("REN"))
			_, err = key.Sign(hash32)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be able to verify a signature", func() {
			key, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			hash32 := NewHash32([]byte("REN"))
			sigHash32, err := key.Sign(hash32)
			Expect(err).ShouldNot(HaveOccurred())

			err = VerifySignature(hash32, sigHash32, key.Address())
			Expect(err).ShouldNot(HaveOccurred())
		})

	})

	Context("when marshaling and unmarshaling", func() {

		It("should be able to marshal and unmarshal as JSON", func() {
			key, err := RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			data, err := key.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())

			keyDecoded := new(EcdsaKey)
			err = keyDecoded.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(key.D).Should(Equal(keyDecoded.D))
			Expect(key.X).Should(Equal(keyDecoded.X))
			Expect(key.Y).Should(Equal(keyDecoded.Y))
			Expect(key.Curve.Params().P).Should(Equal(keyDecoded.Curve.Params().P))
			Expect(key.Curve.Params().N).Should(Equal(keyDecoded.Curve.Params().N))
			Expect(key.Curve.Params().B).Should(Equal(keyDecoded.Curve.Params().B))
			Expect(key.Curve.Params().Gx).Should(Equal(keyDecoded.Curve.Params().Gx))
			Expect(key.Curve.Params().Gy).Should(Equal(keyDecoded.Curve.Params().Gy))
			Expect(key.Curve.Params().BitSize).Should(Equal(keyDecoded.Curve.Params().BitSize))
			Expect("s256").Should(Equal(keyDecoded.Curve.Params().Name))
		})

	})
})
