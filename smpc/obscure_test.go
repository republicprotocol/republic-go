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

var _ = Describe("Obscure residue fragments", func() {
	Context("when producing obscure residue fragments", func() {

		It("should produce obscure random number generators", func() {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			// Start the producer in the background
			sharedOrderResidueTable := smpc.NewSharedObscureResidueTable([32]byte{})
			obscureRngCh, errCh := smpc.ProduceObscureRngs(ctx, n, k, &sharedOrderResidueTable)

			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				// Wait for a ObscureRng to be produced
				for obscureRng := range obscureRngCh {
					Ω(obscureRng.N).Should(Equal(n))
					Ω(obscureRng.K).Should(Equal(k))
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

		It("should produce verifiable obscure random numbers", func() {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			// Create a set of process channels for receiving ObscureRngs
			obscureRngChs := make([]chan smpc.ObscureRng, n)
			for i := int64(0); i < n; i++ {
				obscureRngChs[i] = make(chan smpc.ObscureRng)
			}
			// Send a ObscureRng to all process channels
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer func() {
					for i := range obscureRngChs {
						close(obscureRngChs[i])
					}
				}()

				for i := range obscureRngChs {
					obscureRngChs[i] <- smpc.ObscureRng{
						ObscureResidueID: [32]byte{1},
						Owner:            [32]byte{1},
						Signature:        [32]byte{},
						N:                n,
						K:                k,
					}
				}
			}()

			// Create N processes that will consume ObscureRngs and
			// produce ObscureRngSharesIndexed
			sharedOrderResidueTable := smpc.NewSharedObscureResidueTable([32]byte{})
			obscureRngSharesChs := make([]<-chan smpc.ObscureRngShares, n)
			for i := int64(0); i < n; i++ {
				obscureRngShares, errCh := smpc.ProcessObscureRngs(ctx, obscureRngChs[i], &sharedOrderResidueTable)
				obscureRngSharesChs[i] = obscureRngShares

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
			obscureRngShares := make([]smpc.ObscureRngShares, n)
			for i := int64(0); i < n; i++ {
				obscureRngShares[i] = <-obscureRngSharesChs[i]
			}

			// Reconstruct the expected random number
			randomNumber := stackint.Zero()
			for i := int64(0); i < n; i++ {
				a := shamir.Join(&smpc.Prime, obscureRngShares[i].A)
				randomNumber = randomNumber.AddModulo(&a, &smpc.Prime)
			}

			// Reconstruct the random number from the additive shares
			randomNumberShares := make(shamir.Shares, k)
			for j := int64(0); j < k; j++ {
				rs := make(shamir.Shares, n)
				for i := int64(0); i < n; i++ {
					rs[i] = obscureRngShares[i].A[j]
				}
				randomNumberShares[j] = smpc.SummateShares(rs, &smpc.Prime)
				Ω(randomNumberShares[j].Key).Should(Equal(j + 1))
			}
			randomNumberCmp := shamir.Join(&smpc.Prime, randomNumberShares)

			Ω(randomNumber.Cmp(&randomNumberCmp)).Should(Equal(0))

			cancel()
			wg.Wait()
		})

		It("should produce verifiable obscure multiplication shares", func() {
			var wg sync.WaitGroup

			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			// Create a set of process channels for receiving ObscureRngSharesIndexed
			obscureRngSharesChs := make([]chan smpc.ObscureRngSharesIndexed, n)
			for i := int64(0); i < n; i++ {
				obscureRngSharesChs[i] = make(chan smpc.ObscureRngSharesIndexed)
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

			// Send a ObscureRngSharesIndexed to all process channels
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer func() {
					for i := range obscureRngSharesChs {
						close(obscureRngSharesChs[i])
					}
				}()

				for j := range obscureRngSharesChs {
					obscureRngSharesIndexed := smpc.ObscureRngSharesIndexed{
						ObscureResidueID: [32]byte{1},
						N:                n,
						K:                k,
						A:                make(shamir.Shares, n),
						B:                make(shamir.Shares, n),
						R:                make(shamir.Shares, n),
					}
					for i := int64(0); i < n; i++ {
						obscureRngSharesIndexed.A[i] = aShares[i][j]
						obscureRngSharesIndexed.B[i] = bShares[i][j]
						obscureRngSharesIndexed.R[i] = rShares[i][j]
					}
					obscureRngSharesChs[j] <- obscureRngSharesIndexed
				}
			}()

			// Create N processes that will consume ObscureRngSharesIndexed
			// and produce ObscureMulShares
			obscureMulSharesChs := make([]<-chan smpc.ObscureMulShares, n)
			for i := int64(0); i < n; i++ {
				obscureMulSharesCh, errCh := smpc.ProcessObscureRngSharesIndexed(ctx, obscureRngSharesChs[i])
				obscureMulSharesChs[i] = obscureMulSharesCh

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
				obscureMulShares := <-obscureMulSharesChs[i]
				ab := shamir.Join(&smpc.Prime, obscureMulShares.AB)
				Ω(ab.Cmp(&abs[i])).Should(Equal(0))
				rSq := shamir.Join(&smpc.Prime, obscureMulShares.RSq)
				Ω(rSq.Cmp(&rSqs[i])).Should(Equal(0))
			}

			cancel()
			wg.Wait()
		})

		It("should produce verifiable obscure residue fragments", func() {
			var wg sync.WaitGroup

			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			// Create a set of process channels for receiving ObscureMulSharesIndexed
			obscureMulSharesIndexedChs := make([]chan smpc.ObscureMulSharesIndexed, n)
			for i := int64(0); i < n; i++ {
				obscureMulSharesIndexedChs[i] = make(chan smpc.ObscureMulSharesIndexed)
			}

			// Create N processes that will consume
			// ObscureMulSharesIndexed and produce XiFragments.
			obscureResidueFragmentChs := make([]<-chan smpc.ObscureResidueFragment, n)
			for i := int64(0); i < n; i++ {
				obscureResidueFragmentCh, errCh := smpc.ProcessObscureMulSharesIndexed(ctx, obscureMulSharesIndexedChs[i])
				obscureResidueFragmentChs[i] = obscureResidueFragmentCh

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

			// Send a ObscureMulSharesIndexed to each process
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

				obscureMulSharesIndexedChs[j] <- smpc.ObscureMulSharesIndexed{
					ObscureResidueID: [32]byte{1},
					N:                n,
					K:                k,
					AB:               ab,
					RSq:              rSq,
				}
			}
			for i := range obscureMulSharesIndexedChs {
				close(obscureMulSharesIndexedChs[i])
			}

			// Reconstruct the values AB, and  R² from ObscureResidueFragments
			obscureABShares := make(shamir.Shares, n)
			obscureRSqShares := make(shamir.Shares, n)
			for i := range obscureResidueFragmentChs {
				obscureResidueFragment := <-obscureResidueFragmentChs[i]
				obscureABShares[i] = obscureResidueFragment.AB
				obscureRSqShares[i] = obscureResidueFragment.RSq
			}
			obscureAB := shamir.Join(&smpc.Prime, obscureABShares)
			obscureRq := shamir.Join(&smpc.Prime, obscureRSqShares)

			Ω(ab.Cmp(&obscureAB)).Should(Equal(0))
			Ω(rSq.Cmp(&obscureRq)).Should(Equal(0))

			cancel()
			wg.Wait()
		})

	})
})
