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

			var orders [NumberOfTestOrders]orderbook.Entry
			for i := 0; i < NumberOfTestOrders; i++ {
				orders[i] = newEntry(order.ID([]byte{uint8(i)}))
			}

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Open(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Match(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Release(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))
		})

		It("should be able to store data and its status", func() {
			var orders [NumberOfTestOrders]orderbook.Entry
			for i := 0; i < NumberOfTestOrders; i++ {
				orders[i] = newEntry(order.ID([]byte{uint8(i)}))
			}

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Open(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Match(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Confirm(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Settle(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))
		})

		It("should be able to cancel open orders", func() {
			var orders [NumberOfTestOrders]orderbook.Entry
			for i := 0; i < NumberOfTestOrders; i++ {
				orders[i] = newEntry(order.ID([]byte{uint8(i)}))
			}

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Open(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Match(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Cancel(orders[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should be able to cancel unconfirmed orders", func() {
			var orders [NumberOfTestOrders]orderbook.Entry
			for i := 0; i < NumberOfTestOrders; i++ {
				orders[i] = newEntry(order.ID([]byte{uint8(i)}))
			}

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Open(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Match(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Cancel(orders[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})
	})

	Context("negative tests", func() {
		var cache orderbook.Cache
		var orders [NumberOfTestOrders]orderbook.Entry

		BeforeEach(func() {
			cache = orderbook.NewCache()
			Ω(len(cache.Blocks())).Should(Equal(0))

			for i := 0; i < NumberOfTestOrders; i++ {
				orders[i] = newEntry(order.ID([]byte{uint8(i)}))
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are already opened", func() {

			// Open first time
			for _, order := range orders {
				err := cache.Open(order)
				Ω(err).ShouldNot(HaveOccurred())
			}

			// Open second time
			for _, order := range orders {
				err := cache.Open(order)
				Ω(err).Should(HaveOccurred())
			}

			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))
		})

		It("should not accepted orders that are matched directly", func() {
			for _, order := range orders {
				err := cache.Match(order)
				Ω(err).Should(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are confirmed directly", func() {
			for _, order := range orders {
				err := cache.Confirm(order)
				Ω(err).Should(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are released directly", func() {
			for _, order := range orders {
				err := cache.Release(order)
				Ω(err).Should(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not accepted orders that are settled directly", func() {
			for _, order := range orders {
				err := cache.Settle(order)
				Ω(err).Should(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("should not cancel orders that haven't been opened", func() {
			for _, order := range orders {
				err := cache.Cancel(order.ID)
				Ω(err).Should(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

	})

	Context("canceling orders", func() {
		var cache orderbook.Cache
		var orders [NumberOfTestOrders]orderbook.Entry

		BeforeEach(func() {
			cache = orderbook.NewCache()
			Ω(len(cache.Blocks())).Should(Equal(0))

			for i := 0; i < NumberOfTestOrders; i++ {
				orders[i] = newEntry(order.ID([]byte{uint8(i)}))
			}

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Open(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))
		})

		It("can cancel open orders", func() {
			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Cancel(orders[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("can cancel unconfirmed orders", func() {

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Match(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Cancel(orders[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))
		})

		It("can't cancel confirmed orders", func() {

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Match(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
				err = cache.Confirm(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Cancel(orders[i].ID)
				Ω(err).Should(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(NumberOfTestOrders))
		})

		It("can't change status of canceled orders", func() {

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Match(orders[i])
				Ω(err).ShouldNot(HaveOccurred())
				err = cache.Cancel(orders[i].ID)
				Ω(err).ShouldNot(HaveOccurred())
			}
			Ω(len(cache.Blocks())).Should(Equal(0))

			for i := 0; i < NumberOfTestOrders; i++ {
				err := cache.Open(orders[i])
				Ω(err).Should(HaveOccurred())

				err = cache.Match(orders[i])
				Ω(err).Should(HaveOccurred())

				err = cache.Confirm(orders[i])
				Ω(err).Should(HaveOccurred())

				err = cache.Release(orders[i])
				Ω(err).Should(HaveOccurred())
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
