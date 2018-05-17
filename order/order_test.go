package order_test

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Orders", func() {

	n := int64(17)
	k := int64(12)

	price := createValue(uint32(1), uint32(2))
	minVolume := createValue(uint32(1), uint32(3))
	maxVolume := createValue(uint32(1), uint32(3))

	Context("when testing for equality", func() {

		It("should return true for orders that are equal", func() {
			expiry := time.Now().Add(time.Hour)
			nonce := int64(10)
			lhs := NewOrder(TypeLimit, ParityBuy, expiry, TokensBTCETH, price, maxVolume, minVolume, nonce)
			rhs := NewOrder(TypeLimit, ParityBuy, expiry, TokensBTCETH, price, maxVolume, minVolume, nonce)

			Ω(bytes.Equal(lhs.ID[:], rhs.ID[:])).Should(Equal(true))
			Ω(lhs.Equal(&rhs)).Should(Equal(true))
		})

		It("should return false for orders that are not equal", func() {
			nonce := int64(10)
			lhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)
			nonce = int64(20)
			rhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			Ω(bytes.Equal(lhs.ID[:], rhs.ID[:])).Should(Equal(false))
			Ω(lhs.Equal(&rhs)).Should(Equal(false))
		})
	})

	Context("when splitting orders", func() {

		It("should return the correct number of order fragments", func() {
			nonce := int64(10)
			ord := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			fragments, err := ord.Split(n, k)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(len(fragments)).Should(Equal(int(n)))
		})

		It("should return different order fragments", func() {
			nonce := int64(10)
			ord := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			fragments, err := ord.Split(n, k)
			Ω(err).ShouldNot(HaveOccurred())

			for i := range fragments {
				for j := i + 1; j < len(fragments); j++ {
					Ω(fragments[i].Equal(&fragments[j])).Should(Equal(false))
				}
			}
		})
	})

	Context("when reading and writing orders from files", func() {

		It("should unmarshal and load orders from file", func() {
			nonce := int64(10)
			ord1 := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)
			nonce = int64(20)
			ord2 := NewOrder(TypeLimit, ParitySell, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			err := WriteOrdersToJSONFile("testOrdersFile.json", []*Order{&ord1, &ord2})
			Ω(err).ShouldNot(HaveOccurred())

			orders, err := NewOrdersFromJSONFile("testOrdersFile.json")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(orders)).Should(Equal(int(2)))
		})

		It("should unmarshal and load a single order from file", func() {
			nonce := int64(10)
			ord1 := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			err := writeOrderToJSONFile("testOrdersFile.json", &ord1)
			Ω(err).ShouldNot(HaveOccurred())

			order, err := NewOrderFromJSONFile("testOrdersFile.json")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(order.Nonce).Should(Equal(int64(10)))
		})
	})
})

func createValue(coeff, exp uint32) Value {
	return Value{
		Co:  coeff,
		Exp: exp,
	}
}

// Write a single order into a JSON file.
func writeOrderToJSONFile(fileName string, order *Order) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(&order); err != nil {
		return err
	}
	return nil
}
