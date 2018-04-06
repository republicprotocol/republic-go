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
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())

			orderFragmentCh := make(chan order.Fragment)
			_, errCh := smpc.ProcessOrderFragments(ctx, orderFragmentCh, 0)

			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for err := range errCh {
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			cancel()
			wg.Wait()
		})

		It("should produce order tuples for all order pairs", func() {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())

			numBuyOrders := rand.Intn(255)
			numSellOrders := rand.Intn(255)
			numOrderTuples := numBuyOrders * numSellOrders

			orderFragmentCh := make(chan order.Fragment)
			orderTuplesCh, errCh := smpc.ProcessOrderFragments(ctx, orderFragmentCh, numOrderTuples)

			// Consume OrderTuples and errors
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
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for err := range errCh {
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			// Produce order fragments
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer close(orderFragmentCh)

				zeroShare := shamir.Share{Key: 0, Value: stackint.Zero()}
				for i := 0; i < numBuyOrders; i++ {
					buy := order.NewFragment(order.ID([]byte{0, byte(i)}), order.TypeLimit, order.ParityBuy, zeroShare, zeroShare, zeroShare, zeroShare, zeroShare)
					orderFragmentCh <- *buy

				}
				for i := 0; i < numSellOrders; i++ {
					sell := order.NewFragment(order.ID([]byte{byte(i), 0}), order.TypeLimit, order.ParitySell, zeroShare, zeroShare, zeroShare, zeroShare, zeroShare)
					orderFragmentCh <- *sell
				}
			}()

			wg.Wait()
		})
	})
})
