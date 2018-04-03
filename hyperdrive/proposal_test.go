package hyper_test

import (
	"context"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Processors", func() {

	Context("when processing proposals", func() {

		It("should return errors on shutdown", func() {
			ctx, cancel := context.WithCancel(context.Background())
			proposalChIn := make(chan Proposal)

			_, _, errCh := ProcessProposal(ctx, proposalChIn)

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

		It("should return prepares after processing a valid proposal", func() {
		})

		It("should return faults after processing an invalid proposal", func() {
		})
	})
})
