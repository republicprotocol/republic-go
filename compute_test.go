package compute_test

import (
	"math/big"

	"github.com/republicprotocol/go-order-compute"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Computations", func() {

	// Parameters for secret sharing.
	n := int64(35)
	k := int64(24)
	prime, _ := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

	Context("when orders are an exact match", func() {

		It("should find a match", func() {
			// Split a buy order into fragments.
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Split a sell order into fragments.
			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Combine fragments into a result.
			result, err := computeResultFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})
	})

	Context("when prices vary", func() {

		It("should find a match when the buy price is higher", func() {
			// Split a buy order into fragments.
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 12, 1000, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Split a sell order into fragments.
			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Combine fragments into a result.
			result, err := computeResultFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should not find a match when the buy price is lower", func() {
			// Split a buy order into fragments.
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Split a sell order into fragments.
			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 12, 1000, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Combine fragments into a result.
			result, err := computeResultFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

	Context("when volumes vary", func() {

		It("should find a match when the maximum buy volume is higher than the maximum sell volume", func() {
			// Split a buy order into fragments.
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Split a sell order into fragments.
			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 100, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Combine fragments into a result.
			result, err := computeResultFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the maximum sell volume is higher than the maximum buy volume", func() {
			// Split a buy order into fragments.
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 100, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Split a sell order into fragments.
			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 100, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Combine fragments into a result.
			result, err := computeResultFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the minimum buy volume is higher than the minimum sell volume", func() {
			// Split a buy order into fragments.
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 1000, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Split a sell order into fragments.
			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Combine fragments into a result.
			result, err := computeResultFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the minimum sell volume is higher than the maximum buy volume", func() {
			// Split a buy order into fragments.
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Split a sell order into fragments.
			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 1000, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Combine fragments into a result.
			result, err := computeResultFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should not find a match when the maximum buy volume is lower than the minimum sell volume", func() {
			// Split a buy order into fragments.
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 100, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Split a sell order into fragments.
			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 1000, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Combine fragments into a result.
			result, err := computeResultFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the maximum sell volume is lower than the minimum buy volume", func() {
			// Split a buy order into fragments.
			lhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParityBuy, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 1000, 1000, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Split a sell order into fragments.
			rhs, err := compute.NewOrder(compute.OrderTypeLimit, compute.OrderParitySell, compute.CurrencyCodeBTC, compute.CurrencyCodeETH, 10, 100, 100, 0).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			// Combine fragments into a result.
			result, err := computeResultFromOrderFragments(lhs, rhs, n, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

})

func computeResultFromOrderFragments(lhs []*compute.OrderFragment, rhs []*compute.OrderFragment, n int64, prime *big.Int) (*compute.Result, error) {
	// Generate pairwise computations for each fragment class.
	computations := make([]*compute.Computation, n)
	for i := range computations {
		computation, err := compute.NewComputation(lhs[i], rhs[i])
		if err != nil {
			return nil, err
		}
		computations[i] = computation
	}
	// Compute all result fragments.
	resultFragments := make([]*compute.ResultFragment, n)
	for i := range resultFragments {
		resultFragment, err := computations[i].Sub(prime)
		if err != nil {
			return nil, err
		}
		resultFragments[i] = resultFragment
	}
	// Combine them into a final result.
	return compute.NewResult(resultFragments, prime)

}
