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
			// dhtAddress, _, err := identity.NewAddress()
			// Ω(err).ShouldNot(HaveOccurred())
			// dht := NewDHT(dhtAddress, maxBucketLength)

			// randomAddress, _, err := identity.NewAddress()
			// Ω(err).ShouldNot(HaveOccurred())
			// randomMultiAddress, err := randomAddress.MultiAddress()
			// Ω(err).ShouldNot(HaveOccurred())
			// err = dht.UpdateMultiAddress(randomMultiAddress)
			// Ω(err).ShouldNot(HaveOccurred())

			// bucket, err := dht.FindBucket(randomAddress)
			// Ω(err).ShouldNot(HaveOccurred())
			// t := bucket.Get(0).Time

			// time.Sleep(time.Millisecond)

			// err = dht.Update(randomMulti)
			// Ω(err).ShouldNot(HaveOccurred())
			// Ω(t).Should(Equal(bucket.Get(0).Time))
			// t = bucket.Get(0).Time
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

			dht.RemoveMultiAddress(randomMultiAddress)
			Ω(bucket.Length()).Should(Equal(0))
		})
	})
})
