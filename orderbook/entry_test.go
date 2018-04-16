package orderbook_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("entry", func() {
	Context("creation", func() {
		It("shouldn't error", func() {

			ord := order.Order{
				Signature: []byte{},
				ID:        order.ID([]byte{0}),
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
			entry := orderbook.NewEntry(ord, order.Open, epochHash)
			Ω(entry).ShouldNot(BeNil())
		})
	})
})
