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
			Expect(identity.ProtocolWithName("republic").Name).Should(Equal("republic"))
		})

		It("should expose a protocol with the correct constant values", func() {
			Expect(identity.ProtocolWithCode(identity.RepublicCode).Name).Should(Equal("republic"))
			Expect(identity.ProtocolWithCode(identity.RepublicCode).Code).Should(Equal(identity.RepublicCode))
		})
	})

	Context("when creating new multi-addresses", func() {
		It("should be able to get new multiAddress from a valid string", func() {
			key, err := crypto.RandomEcdsaKey()
			addr := identity.Address(key.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = identity.NewMultiAddressFromString("/republic/" + addr.String())
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should error when trying getting multiAddress from a bad address", func() {
			_, err := identity.NewMultiAddressFromString("bad address")
			Expect(err).Should(HaveOccurred())
			key, err := crypto.RandomEcdsaKey()
			addr := identity.Address(key.Address())
			Expect(err).ShouldNot(HaveOccurred())
			_, err = identity.NewMultiAddressFromString("/republic/" + addr.String() + "bad")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when retrieving values", func() {
		It("should give the right value of specific protocol", func() {
			ip4, tcp, republicAddress := "127.0.0.1", "80", "8MGfbzAMS59Gb4cSjpm34soGNYsM2f"
			addresses := fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ip4, tcp, republicAddress)
			multiAddress, err := identity.NewMultiAddressFromString(addresses)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(multiAddress.ValueForProtocol(identity.RepublicCode)).Should(Equal(republicAddress))
			Expect(multiAddress.ValueForProtocol(identity.TCPCode)).Should(Equal(tcp))
			Expect(multiAddress.ValueForProtocol(identity.IP4Code)).Should(Equal(ip4))
		})
	})

	Context("when marshaling to JSON", func() {
		It("should encode and then decode to the same value", func() {
			key, err := crypto.RandomEcdsaKey()
			addr := identity.Address(key.Address())
			Expect(err).ShouldNot(HaveOccurred())
			multi, err := identity.NewMultiAddressFromString("/republic/" + addr.String())
			Expect(err).ShouldNot(HaveOccurred())
			data, err := multi.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())
			newMulti := &identity.MultiAddress{}
			err = newMulti.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(multi.String()).Should(Equal(newMulti.String()))
		})

		It("should correctly encode empty Address", func() {
			empty := identity.MultiAddress{}
			data, err := empty.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal("\"\""))

			newEmpty := &identity.MultiAddress{}
			err = newEmpty.UnmarshalJSON(data)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(newEmpty.Address())
			Expect(empty.Address()).Should(Equal(newEmpty.Address()))
		})

		It("should error when trying to decoding wrong-formatted error", func() {
			newMulti := &identity.MultiAddress{}
			badData := []byte("this is not a valid MultiAddress")
			err := newMulti.UnmarshalJSON(badData)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when converting to an address or ID", func() {
		ip4, tcp, republicAddress := "127.0.0.1", "80", "8MGfbzAMS59Gb4cSjpm34soGNYsM2f"
		addresses := fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ip4, tcp, republicAddress)

		It("should be converted to an Address", func() {
			multiAddress, err := identity.NewMultiAddressFromString(addresses)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(multiAddress.Address()).Should(Equal(identity.Address(republicAddress)))
		})

		It("should be converted to an ID", func() {
			multiAddress, err := identity.NewMultiAddressFromString(addresses)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(multiAddress.ID()).Should(Equal(multiAddress.Address().ID()))
		})
	})

	Context("when hashing", func() {
		It("should produce the same hash from the same multi-addresses", func() {
			multiAddress, err := testutils.RandomMultiAddress()
			Expect(err).ShouldNot(HaveOccurred())
			multiAddressOther := multiAddress
			Expect(bytes.Equal(multiAddress.Hash(), multiAddressOther.Hash())).Should(BeTrue())
		})
	})
})
