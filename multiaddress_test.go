package identity_test

import (
	"github.com/republicprotocol/go-identity"
	. "github.com/republicprotocol/go-identity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("MultiAddresses with support for Republic Protocol", func() {

	Context("after importing the package", func() {
		It("should expose a protocol called republic", func() {
			Ω(ProtocolWithName("republic").Name).Should(Equal("republic"))
		})

		It("should expose a protocol with the correct constant values", func() {
			Ω(ProtocolWithCode(RepublicCode).Name).Should(Equal("republic"))
			Ω(ProtocolWithCode(RepublicCode).Code).Should(Equal(RepublicCode))
		})
	})

	Context("Creating new multiAddress", func() {
		It("should be able to get new multiAddress from a valid string", func() {
			addr, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			_, err = identity.NewMultiAddressFromString("/republic/" + addr.String())
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should error when trying getting multiAddress from a bad address", func() {
			_, err := identity.NewMultiAddressFromString("bad address")
			Ω(err).Should(HaveOccurred())
			addr, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			_, err = identity.NewMultiAddressFromString("/republic/" + addr.String()+ "bad")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("retrieving values", func() {
		It("should give the right value of specific protocol", func() {
			ip4, tcp, republicAddress := "127.0.0.1", "80","8MGfbzAMS59Gb4cSjpm34soGNYsM2f"
			addresses := fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s",ip4,tcp,republicAddress)
			multiAddress, err := identity.NewMultiAddressFromString(addresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiAddress.ValueForProtocol(RepublicCode)).Should(Equal(republicAddress))
			Ω(multiAddress.ValueForProtocol(TCPCode)).Should(Equal(tcp))
			Ω(multiAddress.ValueForProtocol(IP4Code)).Should(Equal(ip4))
		})
	})

	Context("marshaling to JSON", func() {
		It("should encode and then decode to the same value", func() {
			addr, _, err := identity.NewAddress()
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

		It("should error when trying to decoding wrong-formatted error", func() {
			newMulti := &identity.MultiAddress{}
			badData := []byte("this is not a valid Multiaddress")
			err := newMulti.UnmarshalJSON(badData)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("converting to Address or ID", func() {
		ip4, tcp, republicAddress := "127.0.0.1", "80","8MGfbzAMS59Gb4cSjpm34soGNYsM2f"
		addresses := fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s",ip4,tcp,republicAddress)

		It("should be converted to an Address", func() {
			multiAddress, err := identity.NewMultiAddressFromString(addresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiAddress.Address()).Should(Equal(Address(republicAddress)))
		})

		It("should be converted to an ID", func() {
			multiAddress, err := identity.NewMultiAddressFromString(addresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiAddress.ID()).Should(Equal(multiAddress.Address().ID()))
		})
	})
})
