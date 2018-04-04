package smpc_test

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stackint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/smpc"
)

var _ = Describe("Order fragment processor", func() {
	Context("when receiving order fragments", func() {

		It("should shutdown when the context is canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())

			orderFragmentCh := make(chan order.Fragment)
			_, errCh := smpc.ProcessOrderFragments(ctx, orderFragmentCh, 0)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for err := range errCh {
					Ω(err).Should(HaveOccurred())
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			cancel()
			wg.Wait()
		})

		It("should produce computations for all order fragment pairs", func() {
			var wg sync.WaitGroup

			buyOrderFragmentCount := rand.Intn(255)
			sellOrderFragmentCount := rand.Intn(255)

			ctx, cancel := context.WithCancel(context.Background())
			orderFragmentCh := make(chan order.Fragment)
			orderFragmentComputationCh, errCh := smpc.ProcessOrderFragments(ctx, orderFragmentCh, buyOrderFragmentCount*sellOrderFragmentCount)

			// Consume OrderFragmentComputations and errors
			orderFragmentComputationCount := 0
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for range orderFragmentComputationCh {
					orderFragmentComputationCount++
					log.Println(orderFragmentComputationCount)
				}
			}()
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for err := range errCh {
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			// Produce order fragments)
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer close(orderFragmentCh)

				zeroShare := shamir.Share{Key: 0, Value: stackint.Zero()}
				for i := 0; i < buyOrderFragmentCount; i++ {
					buy := order.NewFragment(order.ID([]byte{0, byte(i)}), order.TypeLimit, order.ParityBuy, zeroShare, zeroShare, zeroShare, zeroShare, zeroShare)
					orderFragmentCh <- *buy
				}
				for i := 0; i < sellOrderFragmentCount; i++ {
					sell := order.NewFragment(order.ID([]byte{byte(i), 0}), order.TypeLimit, order.ParitySell, zeroShare, zeroShare, zeroShare, zeroShare, zeroShare)
					orderFragmentCh <- *sell
				}
			}()

			wg.Wait()
			time.Sleep(2 * time.Second)
			cancel()

			Ω(orderFragmentComputationCount).Should(Equal(buyOrderFragmentCount * sellOrderFragmentCount))
		})

	})
})
