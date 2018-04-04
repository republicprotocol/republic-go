package smpc_test

import (
	"context"
	"sync"

	"github.com/republicprotocol/republic-go/stackint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

var _ = Describe("ξ-fragment producers", func() {
	Context("when producing ξ-fragments", func() {

		It("should shutdown when the context is canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			_, errCh := smpc.ProduceXiFragmentGenerators(ctx, n, k)

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

		It("should produce generators", func() {

			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			xiFragmentGeneratorCh, errCh := smpc.ProduceXiFragmentGenerators(ctx, n, k)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				quit := false
				for !quit {
					select {
					case xiFragmentGenerator, ok := <-xiFragmentGeneratorCh:
						if !ok {
							quit = true
							break
						}
						Ω(xiFragmentGenerator.N).Should(Equal(n))
						Ω(xiFragmentGenerator.K).Should(Equal(k))
						cancel()
					case err, ok := <-errCh:
						if !ok {
							quit = true
							break
						}
						Ω(err).Should(Equal(context.Canceled))
					}
				}
				for err := range errCh {
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			wg.Wait()
		})

		It("should produce verified additive shares", func() {
			var wg sync.WaitGroup

			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			// Produce XiFragmentGenerators
			xiFragmentGeneratorCh, errCh := smpc.ProduceXiFragmentGenerators(ctx, n, k)
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for err := range errCh {
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			// Create a set of process channels for receiving XiFragmentGenerators
			xiFragmentGeneratorChs := make([]chan smpc.XiFragmentGenerator, n)
			for i := int64(0); i < n; i++ {
				xiFragmentGeneratorChs[i] = make(chan smpc.XiFragmentGenerator)
			}
			// Split all XiFragmentGenerators to all process channels
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer func() {
					for i := range xiFragmentGeneratorChs {
						close(xiFragmentGeneratorChs[i])
					}
				}()

				for xiFragmentGenerator := range xiFragmentGeneratorCh {
					for i := range xiFragmentGeneratorChs {
						xiFragmentGeneratorChs[i] <- xiFragmentGenerator
					}
				}
			}()

			// Create N processes that will consume XiFragmentGenerators and
			// produce XiAdditiveFragmentShares
			xiFragmentAdditiveSharesChs := make([]<-chan smpc.XiFragmentAdditiveShares, n)
			for i := int64(0); i < n; i++ {
				xiFragmentAdditiveSharesCh, errCh := smpc.ProcessXiFragmentGenerators(ctx, xiFragmentGeneratorChs[i])
				xiFragmentAdditiveSharesChs[i] = xiFragmentAdditiveSharesCh

				wg.Add(1)
				go func(i int64) {
					defer GinkgoRecover()
					defer wg.Done()

					for err := range errCh {
						Ω(err).Should(Equal(context.Canceled))
					}
				}(i)
			}

			// Read one XiFragmentAdditiveShare from each process
			xiFragmentAdditiveSharesCollection := make([]smpc.XiFragmentAdditiveShares, n)
			for i := int64(0); i < n; i++ {
				xiFragmentAdditiveSharesCollection[i] = <-xiFragmentAdditiveSharesChs[i]
			}

			// Reconstruct the expected random number
			randomNumber := stackint.Zero()
			for i := int64(0); i < n; i++ {
				r := shamir.Join(&smpc.Prime, xiFragmentAdditiveSharesCollection[i].A)
				randomNumber = randomNumber.AddModulo(&r, &smpc.Prime)
			}

			// Reconstruct the random number from the additive shares
			randomNumberShares := make(shamir.Shares, k)
			for j := int64(0); j < k; j++ {
				rs := make(shamir.Shares, n)
				for i := int64(0); i < n; i++ {
					rs[i] = xiFragmentAdditiveSharesCollection[i].A[j]
				}
				randomNumberShares[j] = smpc.SummateShares(rs, &smpc.Prime)
				Ω(randomNumberShares[j].Key).Should(Equal(j + 1))
			}
			randomNumberCmp := shamir.Join(&smpc.Prime, randomNumberShares)

			Ω(randomNumber.Cmp(&randomNumberCmp)).Should(Equal(0))

			cancel()
			wg.Wait()
		})

	})
})
