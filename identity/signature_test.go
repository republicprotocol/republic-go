package identity_test

import (
	"crypto/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
)

type SignableStruct struct {
	value string
}

func (s SignableStruct) SerializeForSigning() []byte {
	return []byte(s.value)
}

var _ = Describe("Siging and verifying signatures with KeyPairs", func() {

	Context("basic SignableStruct", func() {

		value := make([]byte, 10)
		rand.Read(value)

		testStruct := SignableStruct{
			value: string(value),
		}

		It("should not error", func() {
			keyPair, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())

			_, err = keyPair.Sign(testStruct)
			Ω(err).Should(Not(HaveOccurred()))
		})

		It("signing ID can be retrieved", func() {
			keyPair, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())

			signature, err := keyPair.Sign(testStruct)
			signer, err := identity.RecoverSigner(testStruct, signature)
			Ω(err).Should(Not(HaveOccurred()))
			Ω(signer).Should(Equal(keyPair.ID()))
		})

		It("signature can be verified", func() {
			keyPair, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())

			signature, err := keyPair.Sign(testStruct)
			err = identity.VerifySignature(testStruct, signature, keyPair.ID())
			Ω(err).Should(Not(HaveOccurred()))
		})

		It("signature verification should error for wrong keypair", func() {
			keyPair1, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())

			keyPair2, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())

			signature, err := keyPair1.Sign(testStruct)
			err = identity.VerifySignature(testStruct, signature, keyPair2.ID())
			Ω(err).Should(Equal(identity.ErrInvalidSignature))
		})

	})

})
