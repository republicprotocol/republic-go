package dht_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"

	. "github.com/republicprotocol/go-dht"
)

var _ = Describe("Distributed Hash Table", func() {

	Context("updates", func() {

		It("should find the updated address", func() {
			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(address)

			address, _, err = identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			multi, err := address.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())
			err = dht.Update(multi)
			Ω(err).ShouldNot(HaveOccurred())

			finded, err := dht.FindMultiAddress(address)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(finded).Should(Equal(multi))
		})

		It("should error when the bucket is full", func() {
			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht := NewDHT(address)
			for i := 0; i < MaxDHTSize; i++ {
				address, _, e := identity.NewAddress()
				Ω(e).ShouldNot(HaveOccurred())
				multi, e := address.MultiAddress()
				Ω(e).ShouldNot(HaveOccurred())
				e = dht.Update(multi)
				if err == nil {
					err = e
				}
			}
			Ω(err).Should(HaveOccurred())
		})

	})

})
