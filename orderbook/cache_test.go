package orderbook_test

import (
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
)

const NumberOfTestOrders = 100

var _ = Describe("order book cache", func() {
	Context("order status change event ", func() {
		var cache orderbook.Cache

		BeforeEach(func() {
			cache = orderbook.NewCache()
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should be able to store data and its status", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Open(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Match(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Release(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))
		})

		It("should be able to store data and its status", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Open(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Match(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Confirm(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Settle(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))
		})
	})

	Context("Negative tests", func() {
		var cache orderbook.Cache

		BeforeEach(func() {
			cache = orderbook.NewCache()
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are told matched directly", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Match(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are confirmed directly", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Confirm(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are released directly", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Release(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are settled directly", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newOrder(order.ID([]byte{uint8(i)}))
				cache.Settle(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

	})
})

func newOrder(id order.ID) order.Order {
	return order.Order{
		Signature: []byte{},
		ID:        id,
		Type:      order.TypeLimit,
		Parity:    order.ParityBuy,
		Expiry:    time.Now(),
		FstCode:   order.CurrencyCodeBTC,
		SndCode:   order.CurrencyCodeETH,
		Price:     big.NewInt(100),
		MaxVolume: big.NewInt(100),
		MinVolume: big.NewInt(100),
		Nonce:     big.NewInt(100),
	}
}
