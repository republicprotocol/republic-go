package order_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/order"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/shamir"
)

var _ = Describe("Order fragments", func() {

	orderID := [32]byte{}
	tokens := shamir.Share{}
	price := CoExpShare{
		Co: shamir.Share{
			Index: uint64(5),
			Value: uint64(4),
		},
		Exp: shamir.Share{
			Index: uint64(5),
			Value: uint64(4),
		},
	}
	minVolume := CoExpShare{
		Co: shamir.Share{
			Index: uint64(10),
			Value: uint64(20),
		},
		Exp: shamir.Share{
			Index: uint64(50),
			Value: uint64(40),
		},
	}
	maxVolume := minVolume
	nonce := shamir.Share{}

	Context("when creating new fragments", func() {

		It("should return a new Fragment with order details initialized", func() {
			copy(orderID[:], "orderID")
			fragment, err := NewFragment(orderID, TypeLimit, ParityBuy, SettlementRenEx, time.Now(), tokens, price, maxVolume, minVolume, nonce)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(bytes.Equal(fragment.OrderID[:], orderID[:])).Should(Equal(true))
		})

		It("should return a new Fragment with a keccak256 encrypted 32 byte ID", func() {
			copy(orderID[:], "orderID")
			expiry := time.Now()
			fragment, err := NewFragment(orderID, TypeLimit, ParityBuy, SettlementRenEx, expiry, tokens, price, maxVolume, minVolume, nonce)
			Expect(err).ShouldNot(HaveOccurred())

			expectedFragment := Fragment{
				OrderID:         orderID,
				OrderType:       TypeLimit,
				OrderParity:     ParityBuy,
				OrderSettlement: SettlementRenEx,
				OrderExpiry:     expiry,
				Tokens:          tokens,
				Price:           price,
				Volume:          maxVolume,
				MinimumVolume:   minVolume,
				Nonce:           nonce,
			}
			fragmentBytes, err := expectedFragment.Bytes()
			Expect(err).ShouldNot(HaveOccurred())
			hash := crypto.Keccak256(fragmentBytes)
			expectedFragmentID := [32]byte{}
			copy(expectedFragmentID[:], hash)

			Expect(bytes.Equal(expectedFragmentID[:], fragment.ID[:])).Should(Equal(true))
		})
	})

	Context("when testing for equality", func() {

		It("should return true if order fragments are equal", func() {
			copy(orderID[:], "orderID")
			expiry := time.Now()

			lhs, err := NewFragment(orderID, TypeLimit, ParityBuy, SettlementRenEx, expiry, tokens, price, maxVolume, minVolume, nonce)
			Expect(err).ShouldNot(HaveOccurred())

			rhs, err := NewFragment(orderID, TypeLimit, ParityBuy, SettlementRenEx, expiry, tokens, price, maxVolume, minVolume, nonce)
			Expect(err).ShouldNot(HaveOccurred())

			Ω(lhs.ID.Equal(rhs.ID)).Should(Equal(true))
			Ω(lhs.ID.String()).Should(Equal(rhs.ID.String()))
			Ω(lhs.Equal(&rhs)).Should(Equal(true))

		})

		It("should return false if order fragments are not equal", func() {
			copy(orderID[:], "orderID")
			lhs, err := NewFragment(orderID, TypeLimit, ParityBuy, SettlementRenEx, time.Now(), tokens, price, maxVolume, minVolume, nonce)
			Expect(err).ShouldNot(HaveOccurred())

			copy(orderID[:], "newOrderID")

			rhs, err := NewFragment(orderID, TypeLimit, ParityBuy, SettlementRenEx, time.Now(), tokens, price, maxVolume, minVolume, nonce)
			Expect(err).ShouldNot(HaveOccurred())

			Ω(lhs.ID.Equal(rhs.ID)).Should(Equal(false))
			Ω(lhs.ID.String()).ShouldNot(Equal(rhs.ID.String()))
			Ω(lhs.Equal(&rhs)).Should(Equal(false))
		})
	})

	Context("when encrypting and decrypting fragments", func() {

		It("should return the same fragment after decrypting its encrypted form", func() {
			copy(orderID[:], "orderID")
			fragment, err := NewFragment(orderID, TypeLimit, ParityBuy, SettlementRenEx, time.Now(), tokens, price, maxVolume, minVolume, nonce)
			Expect(err).ShouldNot(HaveOccurred())

			// Generate new RSA key
			rsaKey, err := crypto.RandomRsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			// Encrypting the fragment must not return an error
			encryptedFragment, err := fragment.Encrypt(rsaKey.PublicKey)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(encryptedFragment).ToNot(Equal(fragment))

			// Decrypting an encrypted fragment must return the original fragment
			decryptedFragment, err := encryptedFragment.Decrypt(rsaKey.PrivateKey)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(decryptedFragment).To(Equal(fragment))

			// Decrypting with incorrect private key must return an error
			newRsaKey, err := rsa.GenerateKey(rand.Reader, 512)
			Expect(err).ShouldNot(HaveOccurred())
			decryptedFragment, err = encryptedFragment.Decrypt(newRsaKey)

			Expect(err).Should(HaveOccurred())
			Expect(decryptedFragment).ToNot(Equal(fragment))
		})
	})
})
