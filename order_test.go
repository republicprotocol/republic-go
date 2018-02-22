package compute_test

import (
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	do "github.com/republicprotocol/go-do"
	compute "github.com/republicprotocol/go-order-compute"
)

var _ = Describe("Orders and order fragments", func() {

	n := int64(17)
	k := int64(12)
	prime, _ := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

	Context("when serializing IDs to strings", func() {

		It("should return the string representation of the ID", func() {
			order := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0))
			Ω(order.ID.String()).Should(Equal(string(order.ID)))
		})
	})

	Context("when testing for equality", func() {

		It("should return true for equal orders", func() {
			order := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0))
			Ω(order.Equals(order)).Should(Equal(true))

			orderFragments, err := order.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			for _, orderFragment := range orderFragments {
				Ω(orderFragment.ID.Equals(orderFragment.ID)).Should(Equal(true))
			}
		})

		It("should return false for orders that are not equal", func() {
			lhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0))
			rhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(1))
			Ω(lhs.Equals(rhs)).Should(Equal(false))

			lhsFragments, err := lhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhsFragments, err := rhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range lhsFragments {
				Ω(lhsFragments[i].ID.Equals(rhsFragments[i].ID)).Should(Equal(false))
			}

		})
	})

	Context("when adding order fragments", func() {

		It("should not return an error for compatible order fragments", func() {
			lhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0))
			rhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(1))
			Ω(lhs.Equals(rhs)).Should(Equal(false))

			lhsFragments, err := lhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhsFragments, err := rhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range lhsFragments {
				_, err := lhsFragments[i].Add(rhsFragments[i], prime)
				Ω(err).ShouldNot(HaveOccurred())
			}
		})

		It("should return an error for incompatible order fragments", func() {
			lhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0))
			rhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(1))
			Ω(lhs.Equals(rhs)).Should(Equal(false))

			lhsFragments, err := lhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhsFragments, err := rhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range lhsFragments {
				_, err := lhsFragments[i].Add(rhsFragments[i], prime)
				Ω(err).Should(HaveOccurred())
			}
		})
	})

	Context("when subtracting order fragments", func() {

		It("should always subtract the sell from the buy", func() {
			lhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0))
			rhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(1))
			Ω(lhs.Equals(rhs)).Should(Equal(false))

			lhsFragments, err := lhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhsFragments, err := rhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range lhsFragments {
				fst, err := lhsFragments[i].Sub(rhsFragments[i], prime)
				Ω(err).ShouldNot(HaveOccurred())
				snd, err := rhsFragments[i].Sub(lhsFragments[i], prime)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(fst.Equals(snd)).Should(Equal(true))
			}
		})

		It("should not return an error for compatible order fragments", func() {
			lhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0))
			rhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(1))
			Ω(lhs.Equals(rhs)).Should(Equal(false))

			lhsFragments, err := lhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhsFragments, err := rhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range lhsFragments {
				_, err := lhsFragments[i].Sub(rhsFragments[i], prime)
				Ω(err).ShouldNot(HaveOccurred())
			}
		})

		It("should return an error for incompatible order fragments", func() {
			lhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0))
			rhs := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(1))
			Ω(lhs.Equals(rhs)).Should(Equal(false))

			lhsFragments, err := lhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhsFragments, err := rhs.Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range lhsFragments {
				_, err := lhsFragments[i].Sub(rhsFragments[i], prime)
				Ω(err).Should(HaveOccurred())
			}
		})
	})

	Context("when using a hidden order book", func() {

		It("should wait for a full block of order comparisons", func() {

			orderBook := compute.NewHiddenOrderBook(4)
			blockChan := do.Process(func() do.Option {
				return do.Ok(orderBook.WaitForShard())
			})

			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			orderBook.AddOrderFragment(lhs[0], prime)
			for i := 0; i < 4; i++ {
				rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, time.Now().Add(time.Hour), compute.CurrencyCodeBTC, compute.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
				Ω(err).ShouldNot(HaveOccurred())
				orderBook.AddOrderFragment(rhs[0], prime)
			}

			shard := (<-blockChan).Ok.(compute.Shard)
			Ω(len(shard.Deltas)).Should(Equal(4))
		})
	})
})
