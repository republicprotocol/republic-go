package order_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Orders", func() {

	n := int64(17)
	k := int64(12)

	prime := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

	price := stackint.FromUint64(10)
	minVolume := stackint.FromUint64(100)
	maxVolume := stackint.FromUint64(1000)

	Context("when representing IDs as strings", func() {

		It("should return the same string for the same orders", func() {
			nonce := stackint.Zero()
			lhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			rhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			Ω(lhs.ID.String()).Should(Equal(rhs.ID.String()))
		})

		It("should return different strings for different orders", func() {
			nonce := stackint.Zero()
			lhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			nonce = stackint.One()
			rhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			Ω(lhs.ID.String()).ShouldNot(Equal(rhs.ID.String()))
		})
	})

	Context("when testing for equality", func() {

		It("should return true for order IDs that are equal", func() {
			nonce := stackint.Zero()
			lhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			rhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			Ω(lhs.ID.Equal(rhs.ID)).Should(Equal(true))
		})

		It("should return false for order IDs that are not equal", func() {
			nonce := stackint.Zero()
			lhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			nonce = stackint.One()
			rhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			Ω(lhs.ID.Equal(rhs.ID)).Should(Equal(false))
		})

		It("should return true for orders that are equal", func() {
			expiry := time.Now().Add(time.Hour)
			nonce := stackint.Zero()
			lhs := NewOrder(TypeLimit, ParityBuy, expiry, CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			rhs := NewOrder(TypeLimit, ParityBuy, expiry, CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			println(lhs)
			println(rhs)
			Ω(lhs.Equal(rhs)).Should(Equal(true))
		})

		It("should return true for orders that are not equal", func() {
			nonce := stackint.Zero()
			lhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			nonce = stackint.One()
			rhs := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			Ω(lhs.Equal(rhs)).Should(Equal(false))
		})
	})

	Context("when splitting orders", func() {

		It("should return the correct number of order fragments", func() {
			nonce := stackint.Zero()
			order := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			fragments, err := order.Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(fragments)).Should(Equal(int(n)))
		})

		It("should return different order fragments", func() {
			nonce := stackint.Zero()
			order := NewOrder(TypeLimit, ParityBuy, time.Now().Add(time.Hour), CurrencyCodeBTC, CurrencyCodeETH, &price, &maxVolume, &minVolume, &nonce)
			fragments, err := order.Split(n, k, &prime)
			Ω(err).ShouldNot(HaveOccurred())
			for i := range fragments {
				for j := i + 1; j < len(fragments); j++ {
					Ω(fragments[i].Equal(fragments[j])).Should(Equal(false))
				}
			}
		})
	})
})
