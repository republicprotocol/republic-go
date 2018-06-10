package dht_test

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dht"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/testutils"
)

const maxBucketLength = 20

func randomAddress() (*identity.Address, *identity.MultiAddress, error) {
	randomAddress, randomMultiAddress, err := testutils.RandomAddressAndMultiAddress()
	if err != nil {
		return nil, nil, err
	}
	return randomAddress, &randomMultiAddress, nil
}

func randomDHTAndAddress() (*DHT, *identity.Address, *identity.MultiAddress, error) {
	dhtAddress, _, err := testutils.RandomAddressAndMultiAddress()
	if err != nil {
		return nil, nil, nil, err
	}
	randomAddress, randomMultiAddress, err := randomAddress()
	if err != nil {
		return nil, nil, nil, err
	}
	dht := NewDHT(*dhtAddress, maxBucketLength)
	return &dht, randomAddress, randomMultiAddress, nil
}

func isSortedMultiAddresses(multiAddresses identity.MultiAddresses, target identity.Address) (bool, error) {
	var globalErr error
	isSorted := sort.SliceIsSorted(multiAddresses, func(i, j int) bool {
		left := multiAddresses[i].Address()
		right := multiAddresses[j].Address()
		closer, err := identity.Closer(left, right, target)
		if globalErr == nil {
			globalErr = err
		}
		return closer
	})
	return isSorted, globalErr
}

var _ = Describe("Distributed Hash Table", func() {

	Context("when adding and updating multi-addresses", func() {

		It("should be able to find an address after updating an empty bucket", func() {
			bucket := NewBucket(maxBucketLength)

			randomAddress, randomMultiAddress, err := randomAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = bucket.UpdateMultiAddress(*randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			multiAddress, position := bucket.FindMultiAddress(*randomAddress)
			Ω(position >= 0).Should(Equal(true))
			Ω(*multiAddress).Should(Equal(*randomMultiAddress))
			Ω(len(bucket.MultiAddresses)).Should(Equal(1))
		})

		It("should be able to find an address after updating an empty DHT", func() {
			dht, randomAddress, randomMultiAddress, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.UpdateMultiAddress(*randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			multiAddress, err := dht.FindMultiAddress(*randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*multiAddress).Should(Equal(*randomMultiAddress))
		})

		It("should return multi-addresses that have been updated", func() {
			dht, _, _, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < maxBucketLength; i++ {
				multiAddress, err := testutils.RandomMultiAddress()
				Ω(err).ShouldNot(HaveOccurred())
				err = dht.UpdateMultiAddress(multiAddress)
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(dht.MultiAddresses())).Should(Equal(maxBucketLength))
		})

		It("should move an existing address to the end of its bucket", func() {
			dht, randomAddress, randomMultiAddress, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.UpdateMultiAddress(*randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < IDLengthInBits*maxBucketLength; i++ {
				multiAddress, err := testutils.RandomMultiAddress()
				Ω(err).ShouldNot(HaveOccurred())
				dht.UpdateMultiAddress(multiAddress)
			}

			err = dht.UpdateMultiAddress(*randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(*randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bucket.MultiAddresses[bucket.Length()-1]).Should(Equal(*randomMultiAddress))
		})

		It("should not allow duplicates in the DHT", func() {
			dht, _, randomMultiAddress, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < 3; i++ {
				err = dht.UpdateMultiAddress(*randomMultiAddress)
				Ω(err).ShouldNot(HaveOccurred())
			}

			Ω(len(dht.MultiAddresses())).Should(Equal(1))
		})

		It("should return an error when adding the DHT address", func() {
			dht, _, _, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			dhtMultiAddress, err := dht.Address.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.UpdateMultiAddress(dhtMultiAddress)
			Ω(err).Should(HaveOccurred())
			Ω(len(dht.MultiAddresses())).Should(Equal(0))
		})

		It("should return an error when the bucket is full", func() {
			dht, _, _, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < IDLengthInBits*maxBucketLength+1; i++ {
				multiAddress, e := testutils.RandomMultiAddress()
				Ω(e).ShouldNot(HaveOccurred())
				e = dht.UpdateMultiAddress(multiAddress)
				if err == nil && e != nil {
					err = e
					break
				}
			}
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when removing multi-addresses", func() {
		It("should remove a multi-address when it was already added", func() {
			dht, randomAddress, randomMultiAddress, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			err = dht.UpdateMultiAddress(*randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(*randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bucket.Length()).Should(Equal(1))

			dht.RemoveMultiAddress(*randomMultiAddress)
			Ω(bucket.Length()).Should(Equal(0))
		})

		It("should do nothing when the multi-address was not yet added", func() {
			dht, randomAddress, randomMultiAddress, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(*randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bucket.Length()).Should(Equal(0))

			err = dht.RemoveMultiAddress(*randomMultiAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bucket.Length()).Should(Equal(0))
		})

		It("should remove multi-addresses correctly when the DHT is full", func() {
			dht, _, _, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < IDLengthInBits*maxBucketLength+1; i++ {
				multiAddress, e := testutils.RandomMultiAddress()
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
	})

	Context("when finding multi-addresses", func() {

		It("should return multi-address neighbors when there are less than α", func() {
			dht, randomAddress, _, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < 3; i++ {
				multiAddress, err := testutils.RandomMultiAddress()
				Ω(err).ShouldNot(HaveOccurred())
				err = dht.UpdateMultiAddress(multiAddress)
				Ω(err).ShouldNot(HaveOccurred())
			}

			multiAddresses, err := dht.FindMultiAddressNeighbors(*randomAddress, 4)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(multiAddresses)).Should(Equal(3))

			isSorted, err := isSortedMultiAddresses(multiAddresses, *randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(isSorted).Should(Equal(true))
		})

		It("should return multi-address neighbors when there are more than α", func() {
			dht, randomAddress, _, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			for i := 0; i < 100; i++ {
				for j := 0; j < 4; j++ {
					multiAddress, err := testutils.RandomMultiAddress()
					Ω(err).ShouldNot(HaveOccurred())
					err = dht.UpdateMultiAddress(multiAddress)
					Ω(err).ShouldNot(HaveOccurred())
				}

				multiAddresses, err := dht.FindMultiAddressNeighbors(*randomAddress, 3)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(len(multiAddresses)).Should(Equal(3))

				isSorted, err := isSortedMultiAddresses(multiAddresses, *randomAddress)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(isSorted).Should(Equal(true))

				multiAddresses = dht.MultiAddresses()
				for _, multiAddress := range multiAddresses {
					err := dht.RemoveMultiAddress(multiAddress)
					Ω(err).ShouldNot(HaveOccurred())
				}
				Ω(len(dht.MultiAddresses())).Should(Equal(0))
			}
		})

		It("should return an error when finding the DHT address", func() {
			dht, _, _, err := randomDHTAndAddress()
			Ω(err).ShouldNot(HaveOccurred())

			_, err = dht.FindBucket(dht.Address)
			Ω(err).Should(HaveOccurred())

			_, err = dht.FindMultiAddress(dht.Address)
			Ω(err).Should(HaveOccurred())
		})
	})
})
