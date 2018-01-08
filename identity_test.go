package identity_test

import (
	"github.com/jbenet/go-base58"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"
)

var _ = Describe("Republic identity", func() {

	Describe("Key pairs", func() {
		Context("generating a new key pair", func() {
			keyPair, err := identity.NewKeyPair()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should have non-nil private key and public key", func() {
				Ω(keyPair.PrivateKey).ShouldNot(BeNil())
				Ω(keyPair.PublicKey).ShouldNot(BeNil())
			})
		})

		Context("getting the ID from a key pair", func() {
			keyPair, err := identity.NewKeyPair()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			id := keyPair.ID()

			It("should return 20 bytes", func() {
				Ω(len(id)).Should(Equal(identity.IDLength))
			})
		})

		Context("getting the address from a key pair", func() {
			keyPair, err := identity.NewKeyPair()

			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			address := keyPair.Address()

			It("should have a length of 30", func() {
				Ω(len(address)).Should(Equal(identity.AddressLength))
			})

			decoded := base58.Decode(string(address))

			Specify("the decoded string should not be empty", func() {
				Ω(decoded).ShouldNot(BeEmpty())
			})
			It("should have 0x1B as its first byte", func() {
				Ω(decoded[0]).Should(Equal(uint8(0x1B)))
			})
			It("should have 0x14 as its second byte", func() {
				Ω(decoded[1]).Should(Equal(uint8(identity.IDLength)))
			})

			It("should be a base58 encoding of its public ID after the first two bytes", func() {
				Ω(decoded[2:]).Should(Equal([]byte(keyPair.ID())))
			})
		})

	})

	Describe("Republic ID", func() {
		Context("generating a new ID", func() {
			id, _, err := identity.NewID()
			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should return 20 bytes", func() {
				Ω(len(id)).Should(Equal(identity.IDLength))
			})
		})
	})

	Describe(" Republic Address", func() {
		Context("generating a new address", func() {
			address, _, err := identity.NewAddress()
			It("should not error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should have a length of 30", func() {
				Ω(len(address)).Should(Equal(identity.AddressLength))
			})

			decoded := base58.Decode(string(address))

			Specify("the decoded string should not be empty", func() {
				Ω(decoded).ShouldNot(BeEmpty())
			})
			It("should have 0x1B as its first byte", func() {
				Ω(decoded[0]).Should(Equal(uint8(0x1B)))
			})
			It("should have 0x14 as its second byte", func() {
				Ω(decoded[1]).Should(Equal(uint8(identity.IDLength)))
			})
		})

		Context("getting distance between two addresses", func() {
			address1 := identity.Address("8MK6bwP1ADVPaMQ4Gxfm85KYbEdJ6Y")
			address2 := identity.Address("8MHkhs4aQ7m7mz7rY1HqEcPwHBgikU")
			zeroDistance := []byte{}
			for i := 0; i < 20; i++ {
				zeroDistance = append(zeroDistance, uint8(0))
			}

			Specify("node should have 0 distance to itself", func() {
				distance, err := address1.Distance(address1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(distance).Should(Equal(zeroDistance))
				distance, err = address2.Distance(address2)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(distance).Should(Equal(zeroDistance))
			})

			Specify("distance from a to b should equals to distance from b to a", func() {
				distance1, err := address1.Distance(address2)
				Ω(err).ShouldNot(HaveOccurred())
				distance2, err := address2.Distance(address1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(distance1).Should(Equal(distance2))
			})

			Specify("correctness of the function", func() {
				distance1, err := address1.Distance(address2)
				Ω(err).ShouldNot(HaveOccurred())
				mannuallyCalculatedResult := []byte{160, 232, 172, 153, 9, 57, 197, 82, 23, 48, 72, 85, 64, 91, 251, 207, 200, 78, 138, 192}
				Ω(distance1).Should(Equal(mannuallyCalculatedResult))
			})

		})

		Context("getting same prefix bits of two addresses", func() {
			address1 := identity.Address("8MK6bwP1ADVPaMQ4Gxfm85KYbEdJ6Y")
			address2 := identity.Address("8MHkhs4aQ7m7mz7rY1HqEcPwHBgikU")

			Specify("we should get same result even if we swap the variable", func() {
				same1, err := address1.SamePrefixLength(address2)
				Ω(err).ShouldNot(HaveOccurred())
				same2, err := address2.SamePrefixLength(address1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(same1).Should(Equal(same2))
			})

			Specify("same prefix of itself should be length bits", func() {
				same1, err := address1.SamePrefixLength(address1)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(same1).Should(Equal(identity.IDLength * 8))
				same2, err := address2.SamePrefixLength(address2)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(same2).Should(Equal(identity.IDLength * 8))
			})

			It("should produce the right result", func() {
				same, err := address1.SamePrefixLength(address2)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(same).Should(Equal(0))
			})

		})

		Context("checking if one is closer to the target than the other", func() {
			address1 := identity.Address("8MK6bwP1ADVPaMQ4Gxfm85KYbEdJ6Y")
			target := identity.Address("8MHkhs4aQ7m7mz7rY1HqEcPwHBgikU")

			Specify("no one should be closer to the target than the target itself", func() {
				randomAddress, _, err := identity.NewAddress()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(identity.Closer(address1, randomAddress, address1)).Should(BeTrue())
				Ω(identity.Closer(randomAddress, address1, randomAddress)).Should(BeTrue())
			})

			It("should get different result if we swap the two addresses", func() {
				randomAddress, _, err := identity.NewAddress()
				Ω(err).ShouldNot(HaveOccurred())
				isAddress1Closer, err := identity.Closer(address1, randomAddress, target)
				Ω(err).ShouldNot(HaveOccurred())
				isRandomAddressCloser, err := identity.Closer(randomAddress, address1, target)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(isAddress1Closer).Should(Equal(!isRandomAddressCloser))
			})
		})

		Context("getting multiaddress from the address", func() {
			address := identity.Address("8MK6bwP1ADVPaMQ4Gxfm85KYbEdJ6Y")
			multi, err := address.MultiAddress()
			Specify("error should not occur", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should be able to retrieve the republic address from the multiaddress", func() {
				Ω(multi.String()).Should(Equal("/republic/" + string(address)))
				Ω(multi.ValueForProtocol(identity.RepublicCode)).Should(Equal(string(address)))
			})
		})
	})
})
