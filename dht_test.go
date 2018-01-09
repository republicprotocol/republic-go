package dht_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"

	. "github.com/republicprotocol/go-dht"
	"time"
)

var _ = Describe("Distributed Hash Table", func() {
	var dht *DHT
	var nodeAddress identity.Address
	var nodeMulti identity.MultiAddress
	var err error

	BeforeEach(func() {
		// Create a new address and its related multiaddress
		address, _, err := identity.NewAddress()
		Ω(err).ShouldNot(HaveOccurred())
		dht = NewDHT(address)

		// Create a new node with random address
		nodeAddress, _, err = identity.NewAddress()
		Ω(err).ShouldNot(HaveOccurred())
		nodeMulti, err = nodeAddress.MultiAddress()
		Ω(err).ShouldNot(HaveOccurred())

	})

	Context("Update nodes", func() {
		It("should find the updated address", func() {
			err = dht.Update(nodeMulti)
			Ω(err).ShouldNot(HaveOccurred())

			finded, err := dht.FindMultiAddress(nodeAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*finded).Should(Equal(nodeMulti))
		})

		It("should error when the bucket is full", func() {
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
	})

	Context("Update same node multi times ", func() {
		Specify("the new time stamp should be different to the old one", func() {
			err = dht.Update(nodeMulti)
			Ω(err).ShouldNot(HaveOccurred())

			// Get the first time stamp
			bucket,err := dht.FindBucket(nodeAddress)
			Ω(err).ShouldNot(HaveOccurred())
			time1 := (*bucket)[0].Time
			for i:=0;i <5;i++{
				// Update the same node after 5 seconds
				time.Sleep(3* time.Second)
				err = dht.Update(nodeMulti)
				Ω(err).ShouldNot(HaveOccurred())
				time2 := (*bucket)[0].Time
				Ω(time1).ShouldNot(Equal(time2))
				time1 = time2
			}
		})
	})

})
