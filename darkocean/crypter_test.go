package darkocean_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/crypto"
	. "github.com/republicprotocol/republic-go/darkocean"
)

var _ = Describe("Crypter", func() {

	var dnr dnr.DarknodeRegistry

	BeforeEach(func() {

	})

	Context("when signing", func() {

		It("should produce valid signatures", func() {
			keystore, err := crypto.RandomKeystore()
			Expect(err).ShouldNot(HaveOccurred())
			crypter := NewCrypter(keystore, dnr, 1, time.Second)

			signature, err := crypter.Sign(crypto.NewHash32([]byte("REN")))
			Expect(err).ShouldNot(HaveOccurred())
			err = crypto.VerifySignature(crypto.NewHash32([]byte("REN")), signature, keystore.Address())
			Expect(err).ShouldNot(HaveOccurred())
		})

	})

	Context("when verifying signatures", func() {

		It("should return an error for unregistered addresses", func() {
			Expect(true).To(BeFalse())
		})

		It("should not return an error for registered addresses", func() {
			Expect(true).To(BeFalse())
		})

	})

	Context("when encrypting", func() {

		It("should encrypt messages for registered addresses", func() {
			Expect(true).To(BeFalse())
		})

		It("should not encrypt messages for unregistered addresses", func() {
			Expect(true).To(BeFalse())
		})

	})

	Context("when decrypting", func() {

		It("should produce the original plain text", func() {
			Expect(true).To(BeFalse())
		})

	})

	Context("when caching", func() {

		It("should update registrations after the update period", func() {
			Expect(true).To(BeFalse())
		})

		It("should update public keys after the update period", func() {
			Expect(true).To(BeFalse())
		})

	})
})
