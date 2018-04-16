package orderbook_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/stackint"
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
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Open(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Match(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Release(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))
		})

		It("should be able to store data and its status", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Open(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Match(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Confirm(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(100))

			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newEntry(order.ID([]byte{uint8(i)}))
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
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Match(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are confirmed directly", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Confirm(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are released directly", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Release(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are settled directly", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				ord := newEntry(order.ID([]byte{uint8(i)}))
				cache.Settle(ord)
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

	})
})

func newEntry(id order.ID) orderbook.Entry {
	ord := order.Order{
		Signature: []byte{},
		ID:        id,
		Type:      order.TypeLimit,
		Parity:    order.ParityBuy,
		Expiry:    time.Now(),
		FstCode:   order.CurrencyCodeBTC,
		SndCode:   order.CurrencyCodeETH,
		Price:     stackint.FromUint(100),
		MaxVolume: stackint.FromUint(100),
		MinVolume: stackint.FromUint(100),
		Nonce:     stackint.FromUint(100),
	}

	var epochHash [32]byte
	return orderbook.Entry{
		Order:     ord,
		Status:    order.Open,
		EpochHash: epochHash,
	}
}
