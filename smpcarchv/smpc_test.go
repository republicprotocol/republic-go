package smpc_test

import (
	"context"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/smpc"
)

var _ = Describe("Smpc Computer", func() {
	Context("when performing secure multiparty computations", func() {

		PIt("should produce obscure residue fragments", func(done Done) {
			defer close(done)

			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			n, k, numResidues := int64(3), int64(2), 100

			computers := make([]smpc.SmpcComputer, n)
			obscureComputeChsIn := make([]smpc.ObscureComputeInput, n)
			for i := int64(0); i < n; i++ {
				computers[i] = smpc.NewSmpcComputer(identity.ID([]byte{byte(i)}), n, k)
				obscureComputeChsIn[i] = smpc.ObscureComputeInput{
					Rng:              make(chan smpc.ObscureRng, n),
					RngShares:        make(chan smpc.ObscureRngShares, n),
					RngSharesIndexed: make(chan smpc.ObscureRngSharesIndexed, n),
					MulShares:        make(chan smpc.ObscureMulShares, n),
					MulSharesIndexed: make(chan smpc.ObscureMulSharesIndexed, n),
				}
			}

			obscureComputeChsOut := make([]smpc.ObscureComputeOutput, n)
			errChs := make([]<-chan error, n)
			for i := int64(0); i < n; i++ {
				obscureComputeChsOut[i], errChs[i] = computers[i].ComputeObscure(ctx, obscureComputeChsIn[i])
			}

			for i := int64(0); i < n; i++ {

				wg.Add(1)
				go func(i int64) {
					defer GinkgoRecover()
					defer wg.Done()

					for v := range obscureComputeChsOut[i].Rng {
						for j := int64(0); j < n; j++ {
							select {
							case <-ctx.Done():
							case obscureComputeChsIn[j].Rng <- v:
							}
						}
					}
				}(i)

				wg.Add(1)
				go func(i int64) {
					defer GinkgoRecover()
					defer wg.Done()

					for v := range obscureComputeChsOut[i].RngShares {
						// NOTE: In a production deployment a lookup table to
						// associate owners with a sending channel would be
						// used. Here, the test suite acts as this lookup table
						// directly by creating computers with simple IDs.
						owner := computers[i].SharedObscureResidueTable().ObscureResidueOwner(v.ObscureResidueID)
						select {
						case <-ctx.Done():
						case obscureComputeChsIn[owner[i]].RngShares <- v:
						}
					}
				}(i)

				wg.Add(1)
				go func(i int64) {
					defer GinkgoRecover()
					defer wg.Done()

					for v := range obscureComputeChsOut[i].RngSharesIndexed {
						select {
						case <-ctx.Done():
						case obscureComputeChsIn[v.A[0].Key-1].RngSharesIndexed <- v:
						}
					}
				}(i)

				wg.Add(1)
				go func(i int64) {
					defer GinkgoRecover()
					defer wg.Done()

					for v := range obscureComputeChsOut[i].MulShares {
						// NOTE: In a production deployment a lookup table to
						// associate owners with a sending channel would be
						// used. Here, the test suite acts as this lookup table
						// directly by creating computers with simple IDs.
						owner := computers[i].SharedObscureResidueTable().ObscureResidueOwner(v.ObscureResidueID)
						select {
						case <-ctx.Done():
						case obscureComputeChsIn[owner[i]].MulShares <- v:
						}
					}
				}(i)

				wg.Add(1)
				go func(i int64) {
					defer GinkgoRecover()
					defer wg.Done()

					for v := range obscureComputeChsOut[i].MulSharesIndexed {
						for j := int64(0); j < n; j++ {
							select {
							case <-ctx.Done():
							case obscureComputeChsIn[v.AB[0].Key-1].MulSharesIndexed <- v:
							}
						}
					}
				}(i)
			}

			errCh := dispatch.MergeErrors(errChs...)
			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for err := range errCh {
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			p := 0
			for {
				time.Sleep(time.Millisecond)
				for i := int64(0); i < n; i++ {
					pLocal := computers[i].SharedObscureResidueTable().NumObscureResidues()
					if pLocal > p {
						p = pLocal
					}
				}
				if p >= numResidues {
					cancel()
					break
				}
			}

			wg.Wait()

			// Cleanup
			for i := int64(0); i < n; i++ {
				obscureComputeChsIn[i].Close()
			}
		})

	})
})
