package smpc_test

import (
	"context"
	"math/rand"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/stackint"
)

var _ = Describe("Order fragment processor", func() {
	Context("when receiving order fragments", func() {

		It("should shutdown when the context is canceled", func() {
			done := make(chan struct{})
			orderFragmentCh := make(chan order.Fragment)
			sharedOrderTable := smpc.NewSharedOrderTable()
			ch := smpc.OrderFragmentsToOrderTuples(done, orderFragmentCh, &sharedOrderTable, 0)
			close(done)
			for range ch {
			}
		})

		PIt("should produce order tuples for all order pairs", func(done Done) {
			defer close(done)

			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())

			numBuyOrders := rand.Intn(255)
			numSellOrders := rand.Intn(255)
			numOrderTuples := numBuyOrders * numSellOrders

			orderFragmentCh := make(chan order.Fragment)
			sharedOrderTable := smpc.NewSharedOrderTable()
			orderTuplesCh := smpc.OrderFragmentsToOrderTuples(ctx.Done(), orderFragmentCh, &sharedOrderTable, numOrderTuples)

			// Consume OrderTuples and cancel the process once all OrderTuples
			// have been consumed
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for i := 0; i < numOrderTuples; i++ {
					_, ok := <-orderTuplesCh
					if !ok {
						Ω(i).Should(Equal(numOrderTuples - 1))
						return
					}
				}
				cancel()
			}()

			// Produce order fragments for the process
			var orderWg sync.WaitGroup
			orderWg.Add(numBuyOrders + numSellOrders)
			for i := 0; i < numBuyOrders; i++ {
				go func(i int) {
					defer GinkgoRecover()
					defer orderWg.Done()

					zeroShare := shamir.Share{Key: 0, Value: stackint.Zero()}
					buy := order.NewFragment(order.ID([]byte{0, byte(i)}), order.TypeLimit, order.ParityBuy, zeroShare, zeroShare, zeroShare, zeroShare, zeroShare)
					orderFragmentCh <- *buy
				}(i)
			}
			for i := 0; i < numSellOrders; i++ {
				go func(i int) {
					defer GinkgoRecover()
					defer orderWg.Done()

					zeroShare := shamir.Share{Key: 0, Value: stackint.Zero()}
					sell := order.NewFragment(order.ID([]byte{byte(i), 0}), order.TypeLimit, order.ParitySell, zeroShare, zeroShare, zeroShare, zeroShare, zeroShare)
					orderFragmentCh <- *sell
				}(i)
			}
			orderWg.Wait()
			close(orderFragmentCh)

			wg.Wait()
		})
	})
})
