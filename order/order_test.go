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

	price := NewCoExp(uint64(1), uint64(2))
	minVolume := NewCoExp(uint64(1), uint64(3))
	maxVolume := NewCoExp(uint64(1), uint64(3))

	Context("when testing for equality", func() {

		It("should return true for orders that are equal", func() {
			expiry := time.Now().Add(time.Hour)
			nonce := uint64(10)
			lhs := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, expiry, TokensBTCETH, price, maxVolume, minVolume, nonce)
			rhs := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, expiry, TokensBTCETH, price, maxVolume, minVolume, nonce)

			Expect(bytes.Equal(lhs.ID[:], rhs.ID[:])).Should(Equal(true))
			Expect(lhs.Equal(&rhs)).Should(Equal(true))
		})

		It("should return false for orders that are not equal", func() {
			nonce := uint64(10)
			lhs := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)
			nonce = uint64(20)
			rhs := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			Expect(bytes.Equal(lhs.ID[:], rhs.ID[:])).Should(Equal(false))
			Expect(lhs.Equal(&rhs)).Should(Equal(false))
		})
	})

	Context("when splitting orders", func() {

		It("should return the correct number of order fragments", func() {
			nonce := uint64(10)
			ord := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			fragments, err := ord.Split(n, k)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(len(fragments)).Should(Equal(int(n)))
		})

		It("should return different order fragments", func() {
			nonce := uint64(10)
			ord := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			fragments, err := ord.Split(n, k)
			Expect(err).ShouldNot(HaveOccurred())

			for i := range fragments {
				for j := i + 1; j < len(fragments); j++ {
					Expect(fragments[i].Equal(&fragments[j])).Should(Equal(false))
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
			Expect(err).ShouldNot(HaveOccurred())

			orders, err := NewOrdersFromJSONFile("orders.out")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(orders)).Should(Equal(int(2)))
		})

		It("should unmarshal and load a single order from file", func() {
			nonce := uint64(10)
			ord1 := NewOrder(TypeLimit, ParityBuy, SettlementRenEx, time.Now().Add(time.Hour), TokensBTCETH, price, maxVolume, minVolume, nonce)

			err := writeOrderToJSONFile("orders.out", &ord1)
			Expect(err).ShouldNot(HaveOccurred())

			order, err := NewOrderFromJSONFile("orders.out")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(order.Nonce).Should(Equal(uint64(10)))
		})
	})

	Context("when handling token and tokens", func() {
		It("should return token name as a string", func() {
			Expect(TokenBTC.String()).Should(Equal("BTC"))
			Expect(TokenETH.String()).Should(Equal("ETH"))
			Expect(TokenDGX.String()).Should(Equal("DGX"))
			Expect(TokenREN.String()).Should(Equal("REN"))
			Expect(Token(100).String()).Should(Equal("unexpected token"))
		})

		It("should return token pair as a string", func() {
			Expect(TokensBTCETH.String()).Should(Equal("BTC-ETH"))
			Expect(TokensBTCDGX.String()).Should(Equal("BTC-DGX"))
			Expect(TokensBTCREN.String()).Should(Equal("BTC-REN"))
			Expect(TokensETHDGX.String()).Should(Equal("ETH-DGX"))
			Expect(TokensETHREN.String()).Should(Equal("ETH-REN"))
			Expect(TokensDGXREN.String()).Should(Equal("DGX-REN"))
			Expect(Tokens(100).String()).Should(Equal("unexpected tokens"))
		})

		It("should be able to extract the first and second token from a token pair", func() {
			Expect(TokensBTCETH.PriorityToken()).Should(Equal(TokenETH))
			Expect(TokensBTCETH.NonPriorityToken()).Should(Equal(TokenBTC))

			Expect(TokensBTCDGX.PriorityToken()).Should(Equal(TokenDGX))
			Expect(TokensBTCDGX.NonPriorityToken()).Should(Equal(TokenBTC))

			Expect(TokensBTCREN.PriorityToken()).Should(Equal(TokenREN))
			Expect(TokensBTCREN.NonPriorityToken()).Should(Equal(TokenBTC))

			Expect(TokensETHDGX.PriorityToken()).Should(Equal(TokenDGX))
			Expect(TokensETHDGX.NonPriorityToken()).Should(Equal(TokenETH))

			Expect(TokensETHREN.PriorityToken()).Should(Equal(TokenREN))
			Expect(TokensETHREN.NonPriorityToken()).Should(Equal(TokenETH))

			Expect(TokensDGXREN.PriorityToken()).Should(Equal(TokenREN))
			Expect(TokensDGXREN.NonPriorityToken()).Should(Equal(TokenDGX))
		})
	})

	Context("when handling parity", func() {
		It("should return parity as a string", func() {
			Expect(ParityBuy.String()).Should(Equal("buy"))
			Expect(ParitySell.String()).Should(Equal("sell"))
			Expect(Parity(100).String()).Should(Equal("unexpected parity"))
		})
	})

	Context("when handling status", func() {
		It("should return status as a string", func() {
			Expect(Open.String()).Should(Equal("open"))
			Expect(Confirmed.String()).Should(Equal("confirmed"))
			Expect(Canceled.String()).Should(Equal("canceled"))
			Expect(Nil.String()).Should(Equal("nil"))
			Expect(Status(100).String()).Should(Equal("unexpected order status"))
		})
	})

	Context("when handling settlement", func() {
		It("should return settlement as a string", func() {
			Expect(SettlementRenEx.String()).Should(Equal("RenEx"))
			Expect(SettlementRenExAtomic.String()).Should(Equal("RenEx Atomic"))
			Expect(SettlementNil.String()).Should(Equal("unexpected order settlement"))
			Expect(Settlement(100).String()).Should(Equal("unexpected order settlement"))
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
