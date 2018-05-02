package identity_test

import (
	"github.com/jbenet/go-base58"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("", func() {

	Describe("Republic addresses", func() {
		Context("generated from random key pairs", func() {
			key, err := crypto.RandomEcdsaKey()
			address := identity.Address(key.Address())

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should have a length of 20 bytes", func() {
				Ω(len(address)).Should(Equal(identity.AddressLength))
			})

			decoded := base58.Decode(string(address))

			It("should not decode to the empty string", func() {
				Ω(decoded).ShouldNot(BeEmpty())
			})
			It("should have 0x1B as its first byte", func() {
				Ω(decoded[0]).Should(Equal(uint8(0x1B)))
			})
			It("should have 0x14 as its second byte", func() {
				Ω(decoded[1]).Should(Equal(uint8(identity.IDLength)))
			})
		})

		Context("calculating distances", func() {
			address1 := identity.Address("8MK6bwP1ADVPaMQ4Gxfm85KYbEdJ6Y")
			address2 := identity.Address("8MHkhs4aQ7m7mz7rY1HqEcPwHBgikU")
			badAddress := identity.Address("8MHkhs4aQ7m7mz7rY1HqEcPwHBg")
			badAddress2 := identity.Address("8MHkhs4aQ7m7mz7rY1HqEcPwHBgik0")
			zeroDistance := []byte{}
			for i := 0; i < 20; i++ {
				zeroDistance = append(zeroDistance, uint8(0))
			}

			It("should error when calculating distance on wrong formatted address", func() {
				_, err := address1.Distance(badAddress)
				Ω(err).Should(HaveOccurred())

				_, err = badAddress.Distance(address1)
				Ω(err).Should(HaveOccurred())

				_, err = address1.Distance(badAddress2)
				Ω(err).Should(HaveOccurred())

				_, err = badAddress2.Distance(address1)
				Ω(err).Should(HaveOccurred())
			})

			It("should have a distance of 0 from itself", func() {
				distance, err := address1.Distance(address1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(distance).Should(Equal(zeroDistance))
				distance, err = address2.Distance(address2)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(distance).Should(Equal(zeroDistance))
			})

			It("should have symmetrical distances", func() {
				distance1, err := address1.Distance(address2)
				Ω(err).ShouldNot(HaveOccurred())
				distance2, err := address2.Distance(address1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(distance1).Should(Equal(distance2))
			})

			It("should calculate the correct distance compared against a known value", func() {
				distance1, err := address1.Distance(address2)
				Ω(err).ShouldNot(HaveOccurred())
				mannuallyCalculatedResult := []byte{160, 232, 172, 153, 9, 57, 197, 82, 23, 48, 72, 85, 64, 91, 251, 207, 200, 78, 138, 192}
				Ω(distance1).Should(Equal(mannuallyCalculatedResult))
			})
		})

		Context("comparing prefix bits", func() {
			address1 := identity.Address("8MHKbwP1ADVPaMQ4Gxfm85KYbEdJ6Y")
			address2 := identity.Address("8MHkcs4aQ7m7mz7rY1HqEcPwHBgikU")
			badAddress := identity.Address("8MHkhs4aQ7m7mz7rY1HqEcPwHBg")

			It("should error when comparing bits on wrong formatted address", func() {
				_, err := address1.SamePrefixLength(badAddress)
				Ω(err).Should(HaveOccurred())
			})

			It("should have symmetrical prefix lengths", func() {
				same1, err := address1.SamePrefixLength(address2)
				Ω(err).ShouldNot(HaveOccurred())
				same2, err := address2.SamePrefixLength(address1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(same1).Should(Equal(same2))
			})

			It("should have a prefix length of 80 bits against itself", func() {
				same1, err := address1.SamePrefixLength(address1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(same1).Should(Equal(identity.IDLength * 8))
				same2, err := address2.SamePrefixLength(address2)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(same2).Should(Equal(identity.IDLength * 8))
			})

			It("should calculate the correct prefix length compared against a known value", func() {
				same, err := address1.SamePrefixLength(address2)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(same).Should(Equal(2))
			})

		})

		Context("comparing who is closer to the target", func() {
			address1 := identity.Address("8MK6bwP1ADVPaMQ4Gxfm85KYbEdJ6Y")
			target := identity.Address("8MHkhs4aQ7m7mz7rY1HqEcPwHBgikU")
			badAddress := identity.Address("8MHkhs4aQ7m7mz7rY1HqEcPwHBg")

			It("should error when comparing with a bad-formatted address", func() {
				_, err := identity.Closer(address1, badAddress, target)
				Ω(err).Should(HaveOccurred())
				_, err = identity.Closer(badAddress, address1, target)
				Ω(err).Should(HaveOccurred())
			})

			It("should not be possible to be closer to the target than the target", func() {
				key, err := crypto.RandomEcdsaKey()
				randomAddress := identity.Address(key.Address())

				Ω(err).ShouldNot(HaveOccurred())
				Ω(identity.Closer(address1, randomAddress, address1)).Should(BeTrue())
				Ω(identity.Closer(randomAddress, address1, randomAddress)).Should(BeTrue())
			})

			It("should be asymmetrical", func() {
				key, err := crypto.RandomEcdsaKey()
				randomAddress := identity.Address(key.Address())

				Ω(err).ShouldNot(HaveOccurred())
				isAddress1Closer, err := identity.Closer(address1, randomAddress, target)
				Ω(err).ShouldNot(HaveOccurred())
				isRandomAddressCloser, err := identity.Closer(randomAddress, address1, target)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(isAddress1Closer).Should(Equal(!isRandomAddressCloser))
			})

			It("should not return being closer when comparing with itself", func() {
				key, err := crypto.RandomEcdsaKey()
				randomAddress := identity.Address(key.Address())

				Ω(err).ShouldNot(HaveOccurred())
				isCloser, err := identity.Closer(randomAddress, randomAddress, target)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(isCloser).Should(BeFalse())
			})
		})

		Context("getting the multi-address", func() {
			address := identity.Address("8MK6bwP1ADVPaMQ4Gxfm85KYbEdJ6Y")
			multiAddress, err := address.MultiAddress()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should be able to get the address from the multi-address", func() {
				Ω(string(multiAddress.Address())).Should(Equal(string(address)))
			})
		})

		Context("converting to ID", func() {
			It("should be converted to an ID", func() {
				address := identity.Address("8MK6bwP1ADVPaMQ4Gxfm85KYbEdJ6Y")
				ID := address.ID()
				newAddress := ID.Address()
				Ω(address).Should(Equal(newAddress))
			})
		})
	})
})
