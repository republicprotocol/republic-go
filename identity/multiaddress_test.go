package identity_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Multiaddresses with support for Republic Protocol", func() {

	Context("after importing the package", func() {
		It("should expose a protocol called republic", func() {
			Ω(identity.ProtocolWithName("republic").Name).Should(Equal("republic"))
		})

		It("should expose a protocol with the correct constant values", func() {
			Ω(identity.ProtocolWithCode(identity.RepublicCode).Name).Should(Equal("republic"))
			Ω(identity.ProtocolWithCode(identity.RepublicCode).Code).Should(Equal(identity.RepublicCode))
		})
	})

	Context("Creating new multiaddress", func() {
		It("should be able to get new multiaddress from a valid string", func() {
			addr, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			_, err = identity.NewMultiaddressFromString("/republic/" + addr.String())
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should error when trying getting multiaddress from a bad address", func() {
			_, err := identity.NewMultiaddressFromString("bad address")
			Ω(err).Should(HaveOccurred())
			addr, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			_, err = identity.NewMultiaddressFromString("/republic/" + addr.String() + "bad")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("retrieving values", func() {
		It("should give the right value of specific protocol", func() {
			ip4, tcp, republicAddress := "127.0.0.1", "80", "8MGfbzAMS59Gb4cSjpm34soGNYsM2f"
			addresses := fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ip4, tcp, republicAddress)
			multiaddress, err := identity.NewMultiaddressFromString(addresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiaddress.ValueForProtocol(identity.RepublicCode)).Should(Equal(republicAddress))
			Ω(multiaddress.ValueForProtocol(identity.TCPCode)).Should(Equal(tcp))
			Ω(multiaddress.ValueForProtocol(identity.IP4Code)).Should(Equal(ip4))
		})
	})

	Context("marshaling to JSON", func() {
		It("should encode and then decode to the same value", func() {
			addr, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			multi, err := identity.NewMultiaddressFromString("/republic/" + addr.String())
			Ω(err).ShouldNot(HaveOccurred())
			data, err := multi.MarshalJSON()
			Ω(err).ShouldNot(HaveOccurred())
			newMulti := &identity.Multiaddress{}
			err = newMulti.UnmarshalJSON(data)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multi.String()).Should(Equal(newMulti.String()))
		})

		It("should error when trying to decoding wrong-formatted error", func() {
			newMulti := &identity.Multiaddress{}
			badData := []byte("this is not a valid Multiaddress")
			err := newMulti.UnmarshalJSON(badData)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("converting to Address or ID", func() {
		ip4, tcp, republicAddress := "127.0.0.1", "80", "8MGfbzAMS59Gb4cSjpm34soGNYsM2f"
		addresses := fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ip4, tcp, republicAddress)

		It("should be converted to an Address", func() {
			multiaddress, err := identity.NewMultiaddressFromString(addresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiaddress.Address()).Should(Equal(identity.Address(republicAddress)))
		})

		It("should be converted to an ID", func() {
			multiaddress, err := identity.NewMultiaddressFromString(addresses)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multiaddress.ID()).Should(Equal(multiaddress.Address().ID()))
		})
	})
})
