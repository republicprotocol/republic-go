package oracle_test

import (
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/oracle"
	"github.com/republicprotocol/republic-go/testutils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = Describe("MidpointPrice storage", func() {

	Context("when storing and retrieving data", func() {

		It("should be able to get the right data we store", func() {
			storer := NewMidpointPriceStorer()
			emptyPrice, err := storer.MidpointPrice()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(emptyPrice.Equals(MidpointPrice{})).Should(BeTrue())

			price := testutils.RandMidpointPrice()
			err = storer.PutMidpointPrice(price)
			Expect(err).ShouldNot(HaveOccurred())

			storedPrice, err := storer.MidpointPrice()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(price.Equals(storedPrice)).Should(BeTrue())
		})
	})
})
