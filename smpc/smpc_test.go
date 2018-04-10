package smpc_test

import (
	"context"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/smpc"
)

var _ = Describe("Smpc Computer", func() {
	Context("when performing secure multiparty computations", func() {

		It("should produce obscure residue fragments", func() {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			n, k := int64(31), int64(16)

			computers := make([]smpc.Computer, n)
			obscureComputeChsIn := make([]smpc.ObscureComputeInput, n)
			for i := int64(0); i < n; i++ {
				computers[i] = smpc.NewComputer([32]byte{byte(i)}, n, k)
				obscureComputeChsIn[i] = smpc.ObscureComputeInput{
					Rng:              make(chan smpc.ObscureRng),
					RngShares:        make(chan smpc.ObscureRngShares),
					RngSharesIndexed: make(chan smpc.ObscureRngSharesIndexed),
					MulShares:        make(chan smpc.ObscureMulShares),
					MulSharesIndexed: make(chan smpc.ObscureMulSharesIndexed),
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
						for j := int64(0); j < n; j++ {
							select {
							case <-ctx.Done():
							case obscureComputeChsIn[j].RngShares <- v:
							}
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
						for j := int64(0); j < n; j++ {
							select {
							case <-ctx.Done():
							case obscureComputeChsIn[j].MulShares <- v:
							}
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
				time.Sleep(time.Second)

				for i := int64(0); i < n; i++ {
					pLocal := computers[i].SharedObscureResidueTable().NumObscureResidues()
					if pLocal > p {
						p = pLocal
					}
				}
				if int64(p) >= n {
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
