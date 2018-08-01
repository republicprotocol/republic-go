package oracle_test

import (
	"log"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/oracle"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = Describe("Swarm storage", func() {

	Context("when storing and retrieving data", func() {

		It("should be able to get the right data we store", func() {
			storer := NewMidpointPriceStorer()
			iter, err := storer.MidpointPrices()
			Expect(err).ShouldNot(HaveOccurred())
			prices, err := iter.Collect()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(prices)).Should(Equal(0))

			price := randMidpointPrice()
			err = storer.PutMidpointPrice(price)
			Expect(err).ShouldNot(HaveOccurred())

			storedPrice, err := storer.MidpointPrice(price.Tokens)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(price.Equals(storedPrice)).Should(BeTrue())
		})
	})

	Context("when getting data using the iterator", func() {
		var NumberOfMessages = 10

		It("should be able to iterate the collection element by element", func() {
			storer := NewMidpointPriceStorer()
			messages := make([]MidpointPrice, NumberOfMessages)
			for i := 0; i < NumberOfMessages; i++ {
				messages[i] = randMidpointPrice()
				err := storer.PutMidpointPrice(messages[i])
				Expect(err).ShouldNot(HaveOccurred())
				log.Println(i, messages[i])
			}

			iter, err := storer.MidpointPrices()
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < NumberOfMessages; i++ {
				Expect(iter.Next()).Should(BeTrue())
				price, err := iter.Cursor()
				Expect(err).ShouldNot(HaveOccurred())
				log.Println(price)
				log.Println(messages[i])
				Expect(price.Equals(messages[i])).Should(BeTrue())
			}
		})
	})
})

func randMidpointPrice() MidpointPrice {
	return MidpointPrice{
		Signature: []byte{},
		Tokens:    rand.Uint64(),
		Price:     rand.Uint64(),
		Nonce:     rand.Uint64(),
	}
}
