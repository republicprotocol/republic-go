package order_test

import (
	"time"

	"github.com/republicprotocol/republic-go/stackint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Order fragments", func() {

	n := int64(17)
	k := int64(12)
	primeStr := "179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111"
	prime := stackint.FromString(primeStr)

	price := stackint.FromUint64(10)
	minVolume := stackint.FromUint64(100)
	maxVolume := stackint.FromUint64(1000)

	Context("when representing IDs as strings", func() {

		It("should return the same string for the same order fragments", func() {
			nonce := stackint.FromUint64(0)

			fragments, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range fragments {
				Ω(fragments[i].ID.String()).Should(Equal(fragments[i].ID.String()))
			}
		})

		It("should return different strings for the different order fragments", func() {
			nonce := stackint.FromUint64(0)

			fragments, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range fragments {
				for j := i + 1; j < len(fragments); j++ {
					Ω(fragments[i].ID.String()).ShouldNot(Equal(Equal(fragments[j].ID.String())))
				}
			}
		})
	})

	Context("when testing for equality", func() {

		It("should return true for order fragments IDs that are equal", func() {
			nonce := stackint.FromUint64(0)

			fragments, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range fragments {
				Ω(fragments[i].ID.Equal(fragments[i].ID)).Should(Equal(true))
			}
		})

		It("should return false for order fragments IDs that are not equal", func() {
			nonce := stackint.FromUint64(0)

			fragments, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range fragments {
				for j := i + 1; j < len(fragments); j++ {
					Ω(fragments[i].ID.Equal(fragments[j].ID)).Should(Equal(false))
				}
			}
		})

		It("should return true for orders fragments that are equal", func() {
			nonce := stackint.FromUint64(0)

			fragments, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range fragments {
				Ω(fragments[i].Equal(fragments[i])).Should(Equal(true))
			}
		})

		It("should return false for orders fragments that are not equal", func() {
			nonce := stackint.FromUint64(0)

			fragments, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range fragments {
				for j := i + 1; j < len(fragments); j++ {
					Ω(fragments[i].Equal(fragments[j])).Should(Equal(false))
				}
			}
		})
	})
	Context("when testing for compatibility", func() {

		It("should return true for pairwise order fragments from orders with different parity", func() {
			nonce := stackint.FromUint64(0)

			lhs, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())

			nonce = stackint.FromUint64(1)
			rhs, err := NewOrder(TypeLimit, ParitySell, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := int64(0); i < n; i++ {
				Ω(lhs[i].IsCompatible(rhs[i])).Should(Equal(true))
			}
		})

		It("should return false for pairwise order fragments from orders with equal parity", func() {
			nonce := stackint.FromUint64(0)

			lhs, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())

			nonce = stackint.FromUint64(1)
			rhs, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := int64(0); i < n; i++ {
				Ω(lhs[i].IsCompatible(rhs[i])).Should(Equal(false))
			}
		})

		It("should return false for non-pairwise order fragments from orders with different parity", func() {
			nonce := stackint.FromUint64(0)

			lhs, err := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())

			nonce = stackint.FromUint64(0)
			rhs, err := NewOrder(TypeLimit, ParitySell, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := int64(0); i < n; i++ {
				for j := i + 1; j < n; j++ {
					Ω(lhs[i].IsCompatible(rhs[j])).Should(Equal(false))
				}
			}
		})

		It("should return false for non-pairwise order fragments from orders with equal parity", func() {
			nonce := stackint.FromUint64(0)

			lhs, err := NewOrder(TypeLimit, ParitySell, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())

			nonce = stackint.FromUint64(1)
			rhs, err := NewOrder(TypeLimit, ParitySell, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce).Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := int64(0); i < n; i++ {
				for j := i + 1; j < n; j++ {
					Ω(lhs[i].IsCompatible(rhs[j])).Should(Equal(false))
				}
			}
		})
	})

})
