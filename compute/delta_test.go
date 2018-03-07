package compute_test

import (
	"crypto/rand"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Delta and delta fragments", func() {

	n := int64(8)
	k := int64(6)
	prime, _ := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)

	Context("when representing IDs as strings", func() {

		It("should return different strings for different deltas", func() {
			lhs, err := computeRandomDelta(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhs, err := computeRandomDelta(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(lhs.ID.String()).ShouldNot(Equal(rhs.ID.String()))
		})
	})

	Context("when using a delta builder", func() {

		It("should only return a delta after receiving k delta fragments", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			deltaFragments := make([]*DeltaFragment, n)
			for i := range deltaFragments {
				deltaFragment := NewDeltaFragment(lhs[i], rhs[i], prime)
				deltaFragments[i] = deltaFragment
			}

			builder := NewDeltaBuilder(k, prime)
			for i := int64(0); i < k-1; i++ {
				delta := builder.InsertDeltaFragment(deltaFragments[i])
				Ω(delta).Should(BeNil())
			}
			delta := builder.InsertDeltaFragment(deltaFragments[k-1])
			Ω(delta).ShouldNot(BeNil())
		})

		It("should not return a delta after the first k delta fragments", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			deltaFragments := make([]*DeltaFragment, n)
			for i := range deltaFragments {
				deltaFragment := NewDeltaFragment(lhs[i], rhs[i], prime)
				deltaFragments[i] = deltaFragment
			}

			builder := NewDeltaBuilder(k, prime)
			for i := int64(0); i < k-1; i++ {
				delta := builder.InsertDeltaFragment(deltaFragments[i])
				Ω(delta).Should(BeNil())
			}
			delta := builder.InsertDeltaFragment(deltaFragments[k-1])
			Ω(delta).ShouldNot(BeNil())

			for i := int64(0); i < n; i++ {
				delta := builder.InsertDeltaFragment(deltaFragments[i])
				Ω(delta).Should(BeNil())
			}
		})

		It("should not return a delta using k non-unique fragments", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			deltaFragments := make([]*DeltaFragment, n)
			for i := range deltaFragments {
				deltaFragment := NewDeltaFragment(lhs[i], rhs[i], prime)
				deltaFragments[i] = deltaFragment
			}

			builder := NewDeltaBuilder(k, prime)
			for i := int64(0); i < k-1; i++ {
				delta := builder.InsertDeltaFragment(deltaFragments[i])
				Ω(delta).Should(BeNil())
			}
			for i := int64(0); i < k-1; i++ {
				delta := builder.InsertDeltaFragment(deltaFragments[i])
				Ω(delta).Should(BeNil())
			}
		})
	})

	Context("when orders are an exact match", func() {

		It("should find a match", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

	})

	Context("when orders use different currencies", func() {

		It("should not find a match for the same currencies in reverse", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeBTC, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the first currencies differ", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeREN, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeREN, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the second currencies differ", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeREN, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when both currencies differ", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeREN, order.CurrencyCodeDGD, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

	Context("when prices vary", func() {

		It("should find a match when the buy price is higher", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(12), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should not find a match when the buy price is lower", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(12), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

	Context("when volumes vary", func() {

		It("should find a match when the maximum buy volume is higher than the maximum sell volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the maximum sell volume is higher than the maximum buy volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the minimum buy volume is higher than the minimum sell volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(1000), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the minimum sell volume is higher than the maximum buy volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(1000), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should not find a match when the maximum buy volume is lower than the minimum sell volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(1000), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the maximum sell volume is lower than the minimum buy volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(1000), big.NewInt(1000), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, big.NewInt(10), big.NewInt(100), big.NewInt(100), big.NewInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

})

func computeRandomDelta(n, k int64, prime *big.Int) (*Delta, error) {
	randomPrice, err := rand.Int(rand.Reader, big.NewInt(2<<32))
	if err != nil {
		return nil, err
	}
	randomNonce, err := rand.Int(rand.Reader, big.NewInt(2<<32))
	if err != nil {
		return nil, err
	}
	lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, randomPrice, big.NewInt(1000), big.NewInt(100), randomNonce).Split(n, k, prime)
	if err != nil {
		return nil, err
	}
	rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, randomPrice, big.NewInt(1000), big.NewInt(100), randomNonce).Split(n, k, prime)
	if err != nil {
		return nil, err
	}
	return computeDeltaFromOrderFragments(lhs, rhs, n, prime), nil
}

func computeDeltaFromOrderFragments(lhs []*order.Fragment, rhs []*order.Fragment, n int64, prime *big.Int) *Delta {
	deltaFragments := make([]*DeltaFragment, n)
	for i := range deltaFragments {
		deltaFragment := NewDeltaFragment(lhs[i], rhs[i], prime)
		if deltaFragment == nil {
			return nil
		}
		deltaFragments[i] = deltaFragment
	}
	return NewDelta(deltaFragments, prime)

}
