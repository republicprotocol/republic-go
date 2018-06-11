package identity_test

import (
	"bytes"
	"fmt"

	"github.com/republicprotocol/republic-go/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("MultiAddresses with support for Republic Protocol", func() {

	Context("when importing the package", func() {
		It("should expose a protocol called republic", func() {
			Ω(identity.ProtocolWithName("republic").Name).Should(Equal("republic"))
		})

		It("should expose a protocol with the correct constant values", func() {
			Ω(identity.ProtocolWithCode(identity.RepublicCode).Name).Should(Equal("republic"))
			Ω(identity.ProtocolWithCode(identity.RepublicCode).Code).Should(Equal(identity.RepublicCode))
		})
	})

	Context("when creating new multi-addresses", func() {
		It("should be able to get new multiAddress from a valid string", func() {
			key, err := crypto.RandomEcdsaKey()
			addr := identity.Address(key.Address())
			Ω(err).ShouldNot(HaveOccurred())
			_, err = identity.NewMultiAddressFromString("/republic/" + addr.String())
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should error when trying getting multiAddress from a bad address", func() {
			_, err := identity.NewMultiAddressFromString("bad address")
			Ω(err).Should(HaveOccurred())
			key, err := crypto.RandomEcdsaKey()
			addr := identity.Address(key.Address())
			Ω(err).ShouldNot(HaveOccurred())
			_, err = identity.NewMultiAddressFromString("/republic/" + addr.String() + "bad")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when retrieving values", func() {
		It("should give the right value of specific protocol", func() {
			ip4, tcp, republicAddress := "127.0.0.1", "80", "8MGfbzAMS59Gb4cSjpm34soGNYsM2f"
			addresses := fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ip4, tcp, republicAddress)
			multiAddress, err := identity.NewMultiAddressFromString(addresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiAddress.ValueForProtocol(identity.RepublicCode)).Should(Equal(republicAddress))
			Ω(multiAddress.ValueForProtocol(identity.TCPCode)).Should(Equal(tcp))
			Ω(multiAddress.ValueForProtocol(identity.IP4Code)).Should(Equal(ip4))
		})
	})

	Context("when marshaling to JSON", func() {
		It("should encode and then decode to the same value", func() {
			key, err := crypto.RandomEcdsaKey()
			addr := identity.Address(key.Address())
			Ω(err).ShouldNot(HaveOccurred())
			multi, err := identity.NewMultiAddressFromString("/republic/" + addr.String())
			Ω(err).ShouldNot(HaveOccurred())
			data, err := multi.MarshalJSON()
			Ω(err).ShouldNot(HaveOccurred())
			newMulti := &identity.MultiAddress{}
			err = newMulti.UnmarshalJSON(data)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multi.String()).Should(Equal(newMulti.String()))
		})

		It("should correctly encode empty Address", func() {
			empty := identity.MultiAddress{}
			data, err := empty.MarshalJSON()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(data)).Should(Equal("\"\""))

			newEmpty := &identity.MultiAddress{}
			err = newEmpty.UnmarshalJSON(data)
			Ω(err).ShouldNot(HaveOccurred())
			fmt.Println(newEmpty.Address())
			Ω(empty.Address()).Should(Equal(newEmpty.Address()))
		})

		It("should error when trying to decoding wrong-formatted error", func() {
			newMulti := &identity.MultiAddress{}
			badData := []byte("this is not a valid MultiAddress")
			err := newMulti.UnmarshalJSON(badData)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when converting to an address or ID", func() {
		ip4, tcp, republicAddress := "127.0.0.1", "80", "8MGfbzAMS59Gb4cSjpm34soGNYsM2f"
		addresses := fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ip4, tcp, republicAddress)

		It("should be converted to an Address", func() {
			multiAddress, err := identity.NewMultiAddressFromString(addresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiAddress.Address()).Should(Equal(identity.Address(republicAddress)))
		})

		It("should be converted to an ID", func() {
			multiAddress, err := identity.NewMultiAddressFromString(addresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiAddress.ID()).Should(Equal(multiAddress.Address().ID()))
		})
	})

	Context("when hashing", func() {
		It("should produce the same hash from the same multi-addresses", func() {
			multiAddress, err := testutils.RandomMultiAddress()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddressOther := multiAddress
			Ω(bytes.Equal(multiAddress.Hash(), multiAddressOther.Hash())).Should(BeTrue())
		})
	})
})
