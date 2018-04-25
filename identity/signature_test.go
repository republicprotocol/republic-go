package identity_test

import (
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
)

type SignableStruct struct {
	value string
}

func (s SignableStruct) Hash() []byte {
	return crypto.Keccak256([]byte(s.value))
}

var _ = Describe("Siging and verifying signatures with KeyPairs", func() {

	Context("basic SignableStruct", func() {

		value := make([]byte, 10)
		rand.Read(value)
		testStruct := SignableStruct{
			value: string(value),
		}

		keyPair, err := identity.NewKeyPair()
		if err != nil {
			panic(err)
		}

		It("should not error", func() {
			_, err = keyPair.Sign(testStruct)
			Ω(err).Should(Not(HaveOccurred()))
		})

		It("signing ID can be retrieved", func() {
			signature, err := keyPair.Sign(testStruct)
			signer, err := identity.RecoverSigner(testStruct, signature)
			Ω(err).Should(Not(HaveOccurred()))
			Ω(signer).Should(Equal(keyPair.ID()))
		})

		It("signature can be verified", func() {
			signature, err := keyPair.Sign(testStruct)
			err = identity.VerifySignature(testStruct, signature, keyPair.ID())
			Ω(err).Should(Not(HaveOccurred()))
		})

		It("signature verification should error for wrong keypair", func() {
			keyPair2, err := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())

			signature, err := keyPair.Sign(testStruct)
			err = identity.VerifySignature(testStruct, signature, keyPair2.ID())
			Ω(err).Should(Equal(identity.ErrInvalidSignature))
		})

		It("signature verification should error for wrong data", func() {
			value := make([]byte, 10)
			rand.Read(value)
			testStruct2 := SignableStruct{
				value: string(value),
			}

			signature, err := keyPair.Sign(testStruct)
			err = identity.VerifySignature(testStruct2, signature, keyPair.ID())
			Ω(err).Should(Equal(identity.ErrInvalidSignature))
		})

	})

})
