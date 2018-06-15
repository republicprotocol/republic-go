package order_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Orders", func() {

	n := int64(17)
	k := int64(12)

	price := NewCoExp(uint64(1), uint64(2))
	minVolume := NewCoExp(uint64(1), uint64(3))
	maxVolume := NewCoExp(uint64(1), uint64(3))

	Context("when testing for equality", func() {

		It("should return true for orders that are equal", func() {
			expiry := time.Now().Add(time.Hour)
			nonce := uint64(10)
			lhs := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, expiry, TokensBTCETH, price, maxVolume, minVolume, nonce)
			rhs := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, expiry, TokensBTCETH, price, maxVolume, minVolume, nonce)

			Ω(bytes.Equal(lhs.ID[:], rhs.ID[:])).Should(Equal(true))
			Ω(lhs.Equal(&rhs)).Should(Equal(true))
		})

		It("should return false for orders that are not equal", func() {
			nonce := uint64(10)
			lhs := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)
			nonce = uint64(20)
			rhs := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			Ω(bytes.Equal(lhs.ID[:], rhs.ID[:])).Should(Equal(false))
			Ω(lhs.Equal(&rhs)).Should(Equal(false))
		})
	})

	Context("when splitting orders", func() {

		It("should return the correct number of order fragments", func() {
			nonce := uint64(10)
			ord := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			fragments, err := ord.Split(n, k)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(len(fragments)).Should(Equal(int(n)))
		})

		It("should return different order fragments", func() {
			nonce := uint64(10)
			ord := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

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
			nonce := uint64(10)
			ord1 := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)
			nonce = uint64(20)
			ord2 := NewOrder(TypeLimit, ParitySell, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			err := WriteOrdersToJSONFile("orders.out", []*Order{&ord1, &ord2})
			Ω(err).ShouldNot(HaveOccurred())

			orders, err := NewOrdersFromJSONFile("orders.out")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(orders)).Should(Equal(int(2)))
		})

		It("should unmarshal and load a single order from file", func() {
			nonce := uint64(10)
			ord1 := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			err := writeOrderToJSONFile("orders.out", &ord1)
			Ω(err).ShouldNot(HaveOccurred())

			order, err := NewOrderFromJSONFile("orders.out")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(order.Nonce).Should(Equal(uint64(10)))
		})
	})

	Context("token  and tokens ", func() {
		It("should able to to show token name when format printing", func() {
			Ω(fmt.Sprintf("%s", TokenBTC)).Should(Equal("BTC"))
			Ω(fmt.Sprintf("%s", TokenETH)).Should(Equal("ETH"))
			Ω(fmt.Sprintf("%s", TokenDGX)).Should(Equal("DGX"))
			Ω(fmt.Sprintf("%s", TokenREN)).Should(Equal("REN"))
		})

		It("should be able to extract the first and second token from a token pair", func() {
			Ω(TokensBTCETH.PriorityToken()).Should(Equal(TokenETH))
			Ω(TokensBTCETH.NonPriorityToken()).Should(Equal(TokenBTC))

			Ω(TokensBTCDGX.PriorityToken()).Should(Equal(TokenDGX))
			Ω(TokensBTCDGX.NonPriorityToken()).Should(Equal(TokenBTC))

			Ω(TokensBTCREN.PriorityToken()).Should(Equal(TokenREN))
			Ω(TokensBTCREN.NonPriorityToken()).Should(Equal(TokenBTC))

			Ω(TokensETHDGX.PriorityToken()).Should(Equal(TokenDGX))
			Ω(TokensETHDGX.NonPriorityToken()).Should(Equal(TokenETH))

			Ω(TokensETHREN.PriorityToken()).Should(Equal(TokenREN))
			Ω(TokensETHREN.NonPriorityToken()).Should(Equal(TokenETH))

			Ω(TokensDGXREN.PriorityToken()).Should(Equal(TokenREN))
			Ω(TokensDGXREN.NonPriorityToken()).Should(Equal(TokenDGX))
		})
	})
})

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
