package leveldb_test

import (
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/leveldb"

	"github.com/republicprotocol/republic-go/oracle"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = Describe("MidpointPrice storage", func() {

	Context("when storing and retrieving data", func() {

		It("should be able to get the right data we store", func() {
			storer := NewMidpointPriceStorer()
			emptyPrice, err := storer.MidpointPrices()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(emptyPrice.Equals(oracle.MidpointPrice{})).Should(BeTrue())

			prices := testutils.RandMidpointPrice()
			err = storer.PutMidpointPrice(prices)
			Expect(err).ShouldNot(HaveOccurred())

			storedPrice, err := storer.MidpointPrices()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(prices.Equals(storedPrice)).Should(BeTrue())

			for token, price := range prices.Prices {
				storedPrice, err := storer.MidpointPrice(order.Tokens(token))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(storedPrice).Should(Equal(price))
			}

			// Error when mid-point price details for invalid token is requested.
			_, err = storer.MidpointPrice(order.Tokens(len(prices.Prices) + 1))
			Expect(err).Should(HaveOccurred())
		})
	})
})
