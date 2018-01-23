package dht_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"

	. "github.com/republicprotocol/go-dht"
)

var _ = Describe("Distributed Hash Table", func() {

	maxBucketLength := 20

	Context("when updating multi-addresses", func() {

		It("should be able to find an address after updating an empty bucket", func() {
			bucket := NewBucket(maxBucketLength)

			randomAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			randomMultiAddress, err := randomAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = bucket.UpdateMultiAddress(randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			multiAddress, position := bucket.FindMultiAddress(randomAddress)
			Ω(position >= 0).Should(Equal(true))
			Ω(*multiAddress).Should(Equal(randomMultiAddress))
			Ω(len(bucket.MultiAddresses)).Should(Equal(1))
		})

		It("should be able to find an address after updating an empty DHT", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			randomAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			randomMultiAddress, err := randomAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.UpdateMultiAddress(randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			multiAddress, err := dht.FindMultiAddress(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*multiAddress).Should(Equal(randomMultiAddress))
		})

		It("should return multi-addresses that have been updated", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			for i := 0; i < maxBucketLength; i++ {
				address, _, err := identity.NewAddress()
				Ω(err).ShouldNot(HaveOccurred())
				multiAddress, err := address.MultiAddress()
				Ω(err).ShouldNot(HaveOccurred())
				err = dht.UpdateMultiAddress(multiAddress)
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(dht.MultiAddresses())).Should(Equal(maxBucketLength))
		})

		It("should move an existing address to the end", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddress, err := address.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())
			err = dht.UpdateMultiAddress(multiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < IDLengthInBits*maxBucketLength; i++ {
				address, _, e := identity.NewAddress()
				Ω(e).ShouldNot(HaveOccurred())
				multiAddress, e := address.MultiAddress()
				Ω(e).ShouldNot(HaveOccurred())
				e = dht.UpdateMultiAddress(multiAddress)
				if err == nil && e != nil {
					err = e
					break
				}
			}
			Ω(err).Should(HaveOccurred())

			err = dht.UpdateMultiAddress(multiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(address)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bucket.MultiAddresses[bucket.Length()-1]).Should(Equal(multiAddress))
		})

		It("should return an error when adding the DHT address", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			dhtMultiAddress, err := dhtAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.UpdateMultiAddress(dhtMultiAddress)
			Ω(err).Should(Equal(ErrDHTAddress))
		})

		It("should return an error when the bucket is full", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			for i := 0; i < IDLengthInBits*maxBucketLength+1; i++ {
				address, _, e := identity.NewAddress()
				Ω(e).ShouldNot(HaveOccurred())
				multiAddress, e := address.MultiAddress()
				Ω(e).ShouldNot(HaveOccurred())
				e = dht.UpdateMultiAddress(multiAddress)
				if err == nil && e != nil {
					err = e
					break
				}
			}
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error updating a bucket with a bad multi-address", func() {
			bucket := NewBucket(maxBucketLength)

			badMultiAddress, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80")
			Ω(err).ShouldNot(HaveOccurred())

			err = bucket.UpdateMultiAddress(badMultiAddress)
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error updating a DHT with a bad multi-address", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			badMultiAddress, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80")
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.UpdateMultiAddress(badMultiAddress)
			Ω(err).Should(HaveOccurred())
		})

		It("should be able to find an address after updating an empty DHT", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			randomAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			randomMultiAddress, err := randomAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.UpdateMultiAddress(randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			multiAddress, err := dht.FindMultiAddress(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*multiAddress).Should(Equal(randomMultiAddress))
		})
	})

	Context("when removing multi-addresses", func() {
		It("should remove a multi-address when it was already added", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			randomAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			randomMultiAddress, err := randomAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.UpdateMultiAddress(randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bucket.Length()).Should(Equal(1))

			dht.RemoveMultiAddress(randomMultiAddress)
			Ω(bucket.Length()).Should(Equal(0))
		})

		It("should do nothing when the multi-address was not yet added", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			randomAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			randomMultiAddress, err := randomAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bucket.Length()).Should(Equal(0))

			err = dht.RemoveMultiAddress(randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bucket.Length()).Should(Equal(0))
		})

		It("should remove multi-addresses correctly when the DHT is full", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			for i := 0; i < IDLengthInBits*maxBucketLength+1; i++ {
				address, _, e := identity.NewAddress()
				Ω(e).ShouldNot(HaveOccurred())
				multiAddress, e := address.MultiAddress()
				Ω(e).ShouldNot(HaveOccurred())
				e = dht.UpdateMultiAddress(multiAddress)
				if err == nil && e != nil {
					err = e
					break
				}
			}
			Ω(err).Should(HaveOccurred())

			multiAddresses := dht.MultiAddresses()
			for _, multiAddress := range multiAddresses {
				err := dht.RemoveMultiAddress(multiAddress)
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(dht.MultiAddresses())).Should(Equal(0))
		})

		It("should return an error removing a bad multi-address from a DHT", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			badMultiAddress, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/80")
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.RemoveMultiAddress(badMultiAddress)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when finding multi-addresses", func() {

		It("should return multi-address neighbors", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			multiAddress, err := address.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())
			err = dht.UpdateMultiAddress(multiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < IDLengthInBits*maxBucketLength; i++ {
				address, _, e := identity.NewAddress()
				Ω(e).ShouldNot(HaveOccurred())
				multiAddress, e := address.MultiAddress()
				Ω(e).ShouldNot(HaveOccurred())
				e = dht.UpdateMultiAddress(multiAddress)
				if err == nil && e != nil {
					err = e
					break
				}
			}
			Ω(err).Should(HaveOccurred())

			multiAddresses, err := dht.FindMultiAddressNeighbors(address, 3)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(multiAddresses)).Should(Equal(3))
		})

		It("should return an error when finding the DHT address", func() {
			dhtAddress, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(dhtAddress, maxBucketLength)

			_, err = dht.FindBucket(dhtAddress)
			Ω(err).Should(HaveOccurred())

			_, err = dht.FindMultiAddress(dhtAddress)
			Ω(err).Should(HaveOccurred())

			_, err = dht.FindMultiAddressNeighbors(dhtAddress, 3)
			Ω(err).Should(HaveOccurred())
		})
	})
})
