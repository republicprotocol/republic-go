package dht_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"

	"time"

	. "github.com/republicprotocol/go-dht"
)

var _ = Describe("Distributed Hash Table", func() {
	var dht *DHT
	var randomAddress identity.Address
	var randomMulti identity.MultiAddress

	Context("updates", func() {
		BeforeEach(func() {
			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht = NewDHT(address)

			randomAddress, _, err = identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())

			randomMulti, err = randomAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

		})

		It("should be able to find address after it is updated", func() {
			err := dht.Update(randomMulti)
			Ω(err).ShouldNot(HaveOccurred())

			multi, err := dht.FindMultiAddress(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*multi).Should(Equal(randomMulti))
		})

		It("should error when the bucket is full", func() {
			var err error
			for i := 0; i < MaxDHTSize; i++ {
				address, _, e := identity.NewAddress()
				Ω(e).ShouldNot(HaveOccurred())
				multi, e := address.MultiAddress()
				Ω(e).ShouldNot(HaveOccurred())
				e = dht.Update(multi)
				if err == nil && e != nil {
					err = e
					break
				}
			}
			Ω(err).Should(HaveOccurred())
		})

		It("should update time stamp for existing addresses", func() {
			err := dht.Update(randomMulti)
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			t := (*bucket)[0].Time

			for i := 0; i < 5; i++ {
				time.Sleep(time.Millisecond)
				err := dht.Update(randomMulti)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(t).ShouldNot(Equal((*bucket)[0].Time))
				t = (*bucket)[0].Time
			}
		})
	})

	Context("Remove ", func() {

		BeforeEach(func() {
			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht = NewDHT(address)

			randomAddress, _, err = identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())

			randomMulti, err = randomAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should remove the node from the dht", func() {
			err := dht.Update(randomMulti)
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(*bucket)).Should(Equal(1))

			dht.Remove(randomMulti)
			Ω(len(*bucket)).Should(Equal(0))
		})
	})
})
