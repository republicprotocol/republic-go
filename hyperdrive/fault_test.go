package hyper_test

import (
	"context"
	"strconv"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Faults", func() {

	blocks := NewSharedBlocks(1, 1)
	validator, _ := NewTestValidator(blocks, 4)

	Context("when processing faults", func() {

		It("should return errors on shutdown", func() {
			ctx, cancel := context.WithCancel(context.Background())
			faultChIn := make(chan Fault)

			_, errCh := ProcessFault(ctx, faultChIn, validator)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()

				for err := range errCh {
					Ω(err).Should(HaveOccurred())
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			cancel()
			wg.Wait()
		})

		It("should return fault after processing a threshold number of prepares", func() {
			ctx, cancel := context.WithCancel(context.Background())
			faultChIn := make(chan Fault, 5)
			faultChOut, errCh := ProcessFault(ctx, faultChIn, validator)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()

				for {
					select {
					case err := <-errCh:
						Ω(err).Should(HaveOccurred())
						Ω(err).Should(Equal(context.Canceled))
						return
					case fault, ok := <-faultChOut:
						if !ok {
							return
						}
						Ω(fault.Rank).Should(Equal(Rank(0)))
						Ω(fault.Height).Should(Equal(uint64(0)))
						if fault.Signature == Signature("Threshold_BLS") {
							cancel()
						}
					}
				}
			}()

			for i := uint8(0); i < validator.Threshold(); i++ {
				go func(i uint8) {
					faultChIn <- Fault{
						Signature: Signature("Signature of " + strconv.Itoa(int(i))),
						Rank:      Rank(0),
						Height:    0,
					}
				}(i)
			}

			wg.Wait()
		})
	})
})
