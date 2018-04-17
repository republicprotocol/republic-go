package compute_test

import (
	"crypto/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"
)

// HeapInt creates a stackint on the heap - temporary convenience method
func heapInt(n uint) *stackint.Int1024 {
	tmp := stackint.FromUint(n)
	return &tmp
}

var _ = Describe("Delta and delta fragments", func() {

	n := int64(8)
	k := int64(6)
	primeVal, _ := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	prime := &primeVal

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
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
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
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
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
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())
			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
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
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

	})

	Context("when orders use different currencies", func() {

		It("should not find a match for the same currencies in reverse", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeBTC, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the first currencies differ", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeREN, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeREN, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the second currencies differ", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeREN, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when both currencies differ", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeREN, order.CurrencyCodeDGD, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

	Context("when prices vary", func() {

		It("should find a match when the buy price is higher", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(12), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should not find a match when the buy price is lower", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(12), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

	Context("when volumes vary", func() {

		It("should find a match when the maximum buy volume is higher than the maximum sell volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(100), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the maximum sell volume is higher than the maximum buy volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(100), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(100), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the minimum buy volume is higher than the minimum sell volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(1000), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should find a match when the minimum sell volume is higher than the maximum buy volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(1000), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(true))
		})

		It("should not find a match when the maximum buy volume is lower than the minimum sell volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(100), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(1000), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})

		It("should not find a match when the maximum sell volume is lower than the minimum buy volume", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(1000), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(100), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			result := computeDelta(lhs, rhs, n, prime)
			Ω(result.IsMatch(prime)).Should(Equal(false))
		})
	})

	Context("when constructing deltas", func() {
		It("should return nil for incompatible orders", func() {
			one, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			two, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeBTC, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			three, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeBTC, heapInt(10), heapInt(1000), heapInt(100), heapInt(1)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			deltaFragment1 := NewDeltaFragment(one[0], two[0], prime)
			deltaFragment2 := NewDeltaFragment(two[0], three[0], prime)

			nilDelta := NewDelta([]*DeltaFragment{deltaFragment1, deltaFragment2}, prime)
			Ω(nilDelta).Should(BeNil())
		})

		It("should return nil for an empty order array", func() {
			nilDelta := NewDelta([]*DeltaFragment{}, prime)
			Ω(nilDelta).Should(BeNil())
		})
	})

	Context("delta fragment methods", func() {
		It("comparing should return the right result", func() {
			one, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			two, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeBTC, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			oneClone, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			deltaFragsOne := computeDeltaFromOrderFragments(one, two, n, prime)
			deltaFragsClone := computeDeltaFromOrderFragments(oneClone, two, n, prime)

			for i := range deltaFragsOne {
				Ω(deltaFragsOne[i].Equals(deltaFragsOne[i])).Should(BeTrue())
				Ω(deltaFragsOne[i].Equals(deltaFragsClone[i])).Should(BeFalse())
			}
		})
	})

	Context("delta fragment IDs", func() {
		It("can be converted to strings", func() {
			lhs := order.NewFragment([]byte("lhs"),
				order.TypeLimit,
				order.ParityBuy,
				shamir.Share{Key: 0, Value: stackint.Zero()},
				shamir.Share{Key: 0, Value: stackint.Zero()},
				shamir.Share{Key: 0, Value: stackint.Zero()},
				shamir.Share{Key: 0, Value: stackint.Zero()},
				shamir.Share{Key: 0, Value: stackint.Zero()},
			)

			rhs := order.NewFragment([]byte("rhs"),
				order.TypeLimit,
				order.ParitySell,
				shamir.Share{Key: 0, Value: stackint.Zero()},
				shamir.Share{Key: 0, Value: stackint.Zero()},
				shamir.Share{Key: 0, Value: stackint.Zero()},
				shamir.Share{Key: 0, Value: stackint.Zero()},
				shamir.Share{Key: 0, Value: stackint.Zero()},
			)

			frag := NewDeltaFragment(lhs, rhs, prime)

			Ω(frag.ID.String()).Should(Equal("B3HMD3WLjgMgPBEvPJJ3yBxuRHsXs7T1DaZyx1QZvnbg"))
		})
	})

	Context("constructing delta fragments", func() {
		It("should return nil for incompatible fragments", func() {
			lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeETH, order.CurrencyCodeBTC, heapInt(10), heapInt(1000), heapInt(100), heapInt(0)).Split(n, k, prime)
			Ω(err).ShouldNot(HaveOccurred())

			nilFrag := NewDeltaFragment(lhs[0], rhs[1], prime)
			Ω(nilFrag).Should(BeNil())
		})
	})

})

func computeRandomDelta(n, k int64, prime *stackint.Int1024) (*Delta, error) {
	max := stackint.FromUint(2 << 32)
	randomPrice, err := stackint.Random(rand.Reader, &max)
	if err != nil {
		return nil, err
	}
	randomNonce, err := stackint.Random(rand.Reader, &max)
	if err != nil {
		return nil, err
	}
	lhs, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, &randomPrice, heapInt(1000), heapInt(100), &randomNonce).Split(n, k, prime)
	if err != nil {
		return nil, err
	}
	rhs, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, &randomPrice, heapInt(1000), heapInt(100), &randomNonce).Split(n, k, prime)
	if err != nil {
		return nil, err
	}
	return computeDelta(lhs, rhs, n, prime), nil
}

func computeDeltaFromOrderFragments(lhs []*order.Fragment, rhs []*order.Fragment, n int64, prime *stackint.Int1024) []*DeltaFragment {
	deltaFragments := make([]*DeltaFragment, n)
	for i := range deltaFragments {
		deltaFragment := NewDeltaFragment(lhs[i], rhs[i], prime)
		if deltaFragment == nil {
			return nil
		}
		deltaFragments[i] = deltaFragment
	}
	return deltaFragments
}

func computeDelta(lhs []*order.Fragment, rhs []*order.Fragment, n int64, prime *stackint.Int1024) *Delta {
	deltaFragments := computeDeltaFromOrderFragments(lhs, rhs, n, prime)
	return NewDelta(deltaFragments, prime)
}
