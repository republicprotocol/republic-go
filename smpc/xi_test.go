package smpc_test

import (
	"context"
	"crypto/rand"
	"sync"

	"github.com/republicprotocol/republic-go/stackint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

var _ = Describe("ξ-fragments", func() {
	Context("when producing ξ-fragments", func() {

		It("should shutdown when the context is canceled", func() {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			// Start the producer in the background
			_, errCh := smpc.ProduceXiFragmentGenerators(ctx, n, k)

			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				// Drain all errors and expect a cancelation
				for err := range errCh {
					Ω(err).Should(HaveOccurred())
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			cancel()
			wg.Wait()
		})

		It("should produce generators", func() {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			// Start the producer in the background
			xiFragmentGeneratorCh, errCh := smpc.ProduceXiFragmentGenerators(ctx, n, k)

			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				// Wait for a XiFragmentGenerator to be produced
				for xiFragmentGenerator := range xiFragmentGeneratorCh {
					Ω(xiFragmentGenerator.N).Should(Equal(n))
					Ω(xiFragmentGenerator.K).Should(Equal(k))
					cancel()
				}
			}()
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				// Drain all errors and only accept cancelations
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

			// Create a set of process channels for receiving XiFragmentGenerators
			xiFragmentGeneratorChs := make([]chan smpc.XiFragmentGenerator, n)
			for i := int64(0); i < n; i++ {
				xiFragmentGeneratorChs[i] = make(chan smpc.XiFragmentGenerator)
			}
			// Send a XiFragmentGenerator to all process channels
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer func() {
					for i := range xiFragmentGeneratorChs {
						close(xiFragmentGeneratorChs[i])
					}
				}()

				for i := range xiFragmentGeneratorChs {
					xiFragmentGeneratorChs[i] <- smpc.XiFragmentGenerator{
						ID: stackint.One(),
						N:  n,
						K:  k,
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

		It("should produce verified multiplicative shares", func() {
			var wg sync.WaitGroup

			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			// Create a set of process channels for receiving XiFragmentAdditiveShares
			xiFragmentAdditiveSharesChs := make([]chan smpc.XiFragmentAdditiveShares, n)
			for i := int64(0); i < n; i++ {
				xiFragmentAdditiveSharesChs[i] = make(chan smpc.XiFragmentAdditiveShares)
			}

			// Generate the additive shares and store the original values for
			// test verification
			as := make([]stackint.Int1024, n)
			bs := make([]stackint.Int1024, n)
			rs := make([]stackint.Int1024, n)
			aShares := make([]shamir.Shares, n)
			bShares := make([]shamir.Shares, n)
			rShares := make([]shamir.Shares, n)
			for i := int64(0); i < n; i++ {
				var err error
				// Generate random A, B and R
				as[i], err = stackint.Random(rand.Reader, &smpc.Prime)
				Ω(err).ShouldNot(HaveOccurred())
				bs[i], err = stackint.Random(rand.Reader, &smpc.Prime)
				Ω(err).ShouldNot(HaveOccurred())
				rs[i], err = stackint.Random(rand.Reader, &smpc.Prime)
				Ω(err).ShouldNot(HaveOccurred())
				// Store shares of A, B, and R
				aShares[i], err = shamir.Split(n, k, &smpc.Prime, &as[i])
				Ω(err).ShouldNot(HaveOccurred())
				bShares[i], err = shamir.Split(n, k, &smpc.Prime, &bs[i])
				Ω(err).ShouldNot(HaveOccurred())
				rShares[i], err = shamir.Split(n, k, &smpc.Prime, &rs[i])
				Ω(err).ShouldNot(HaveOccurred())
			}
			abs := make([]stackint.Int1024, n)
			rSqs := make([]stackint.Int1024, n)
			for j := int64(0); j < n; j++ {
				// Collect all shares of A, B, and R for this key
				aShare := make(shamir.Shares, n)
				bShare := make(shamir.Shares, n)
				rShare := make(shamir.Shares, n)
				for i := int64(0); i < n; i++ {
					aShare[i] = aShares[i][j]
					Ω(aShare[i].Key).Should(Equal(j + 1))
					bShare[i] = bShares[i][j]
					Ω(bShare[i].Key).Should(Equal(j + 1))
					rShare[i] = rShares[i][j]
					Ω(rShare[i].Key).Should(Equal(j + 1))
				}
				// Sum the shares and multiply them for test verification
				a := smpc.SummateShares(aShare, &smpc.Prime)
				b := smpc.SummateShares(bShare, &smpc.Prime)
				r := smpc.SummateShares(rShare, &smpc.Prime)
				abs[j] = a.Value.MulModulo(&b.Value, &smpc.Prime)
				rSqs[j] = r.Value.MulModulo(&r.Value, &smpc.Prime)
			}

			// Send a XiFragmentAdditiveShares to all process channels
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer func() {
					for i := range xiFragmentAdditiveSharesChs {
						close(xiFragmentAdditiveSharesChs[i])
					}
				}()

				for j := range xiFragmentAdditiveSharesChs {
					xiFragmentAdditiveShares := smpc.XiFragmentAdditiveShares{
						ID: stackint.One(),
						N:  n,
						K:  k,
						A:  make(shamir.Shares, n),
						B:  make(shamir.Shares, n),
						R:  make(shamir.Shares, n),
					}
					for i := int64(0); i < n; i++ {
						xiFragmentAdditiveShares.A[i] = aShares[i][j]
						xiFragmentAdditiveShares.B[i] = bShares[i][j]
						xiFragmentAdditiveShares.R[i] = rShares[i][j]
					}
					xiFragmentAdditiveSharesChs[j] <- xiFragmentAdditiveShares
				}
			}()

			// Create N processes that will consume XiFragmentAdditiveShares
			// and produce XiFragmentMultiplicativeShares
			xiFragmentMultiplicativeSharesChs := make([]<-chan smpc.XiFragmentMultiplicativeShares, n)
			for i := int64(0); i < n; i++ {
				xiFragmentMultiplicativeSharesCh, errCh := smpc.ProcessXiFragmentAdditiveShares(ctx, xiFragmentAdditiveSharesChs[i])
				xiFragmentMultiplicativeSharesChs[i] = xiFragmentMultiplicativeSharesCh

				wg.Add(1)
				go func(i int64) {
					defer GinkgoRecover()
					defer wg.Done()

					for err := range errCh {
						Ω(err).Should(Equal(context.Canceled))
					}
				}(i)
			}

			// Read one XiFragmentMultiplicativeShare from each process
			for i := int64(0); i < n; i++ {
				xiFragmentMultiplicativeShares := <-xiFragmentMultiplicativeSharesChs[i]
				ab := shamir.Join(&smpc.Prime, xiFragmentMultiplicativeShares.AB)
				Ω(ab.Cmp(&abs[i])).Should(Equal(0))
				rSq := shamir.Join(&smpc.Prime, xiFragmentMultiplicativeShares.RSquared)
				Ω(rSq.Cmp(&rSqs[i])).Should(Equal(0))
			}

			cancel()
			wg.Wait()
		})

		It("should consume multiplicative shares and produce a verified share", func() {
			var wg sync.WaitGroup

			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			// Create a set of process channels for receiving XiFragmentMultiplicativeShares
			xiFragmentMultiplicativeSharesChs := make([]chan smpc.XiFragmentMultiplicativeShares, n)
			for i := int64(0); i < n; i++ {
				xiFragmentMultiplicativeSharesChs[i] = make(chan smpc.XiFragmentMultiplicativeShares)
			}

			// Create N processes that will consume
			// XiFragmentMultiplicativeShares and produce XiFragments.
			xiFragmentsChs := make([]<-chan smpc.XiFragment, n)
			for i := int64(0); i < n; i++ {
				xiFragmentsCh, errCh := smpc.ProcessXiFragmentMultiplicativeShares(ctx, xiFragmentMultiplicativeSharesChs[i])
				xiFragmentsChs[i] = xiFragmentsCh

				wg.Add(1)
				go func(i int64) {
					defer GinkgoRecover()
					defer wg.Done()

					for err := range errCh {
						Ω(err).Should(Equal(context.Canceled))
					}
				}(i)
			}

			// Create a random A, B, and AB
			a, err := stackint.Random(rand.Reader, &smpc.Prime)
			Ω(err).ShouldNot(HaveOccurred())
			b, err := stackint.Random(rand.Reader, &smpc.Prime)
			Ω(err).ShouldNot(HaveOccurred())
			ab := a.MulModulo(&b, &smpc.Prime)

			// Create a random R, and R²
			r, err := stackint.Random(rand.Reader, &smpc.Prime)
			Ω(err).ShouldNot(HaveOccurred())
			rSq := r.MulModulo(&r, &smpc.Prime)

			// Create shares of A, B, and R
			aShares, err := shamir.Split(n, k, &smpc.Prime, &a)
			Ω(err).ShouldNot(HaveOccurred())
			bShares, err := shamir.Split(n, k, &smpc.Prime, &b)
			Ω(err).ShouldNot(HaveOccurred())
			rShares, err := shamir.Split(n, k, &smpc.Prime, &r)
			Ω(err).ShouldNot(HaveOccurred())

			// Create shares of each share
			abShares := make([]shamir.Shares, n)
			rSqShares := make([]shamir.Shares, n)
			for i := int64(0); i < n; i++ {
				var err error

				ab := aShares[i].Value.MulModulo(&bShares[i].Value, &smpc.Prime)
				abShares[i], err = shamir.Split(n, k, &smpc.Prime, &ab)
				Ω(err).ShouldNot(HaveOccurred())

				rSq := rShares[i].Value.MulModulo(&rShares[i].Value, &smpc.Prime)
				rSqShares[i], err = shamir.Split(n, k, &smpc.Prime, &rSq)
				Ω(err).ShouldNot(HaveOccurred())
			}

			// Send a XiFragmentMultiplicativeShares to each process
			for j := int64(0); j < n; j++ {
				ab := make(shamir.Shares, n)
				for i := int64(0); i < n; i++ {
					ab[i] = abShares[i][j]
					Ω(ab[i].Key).Should(Equal(j + 1))
				}
				rSq := make(shamir.Shares, n)
				for i := int64(0); i < n; i++ {
					rSq[i] = rSqShares[i][j]
					Ω(rSq[i].Key).Should(Equal(j + 1))
				}
				xiFragmentMultiplicativeSharesChs[j] <- smpc.XiFragmentMultiplicativeShares{
					ID:       stackint.One(),
					N:        n,
					K:        k,
					AB:       ab,
					RSquared: rSq,
				}
			}
			for i := range xiFragmentMultiplicativeSharesChs {
				close(xiFragmentMultiplicativeSharesChs[i])
			}

			// Reconstruct the values AB, and  R² from XiFragments
			xiABShares := make(shamir.Shares, n)
			xiRSqShares := make(shamir.Shares, n)
			for i := range xiFragmentsChs {
				xiFragment := <-xiFragmentsChs[i]
				xiABShares[i] = xiFragment.AB
				xiRSqShares[i] = xiFragment.RSquared
			}
			xiAB := shamir.Join(&smpc.Prime, xiABShares)
			xiRq := shamir.Join(&smpc.Prime, xiRSqShares)

			Ω(ab.Cmp(&xiAB)).Should(Equal(0))
			Ω(rSq.Cmp(&xiRq)).Should(Equal(0))

			cancel()
			wg.Wait()
		})

	})
})
