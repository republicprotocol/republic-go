package hyper_test

import (
	"context"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
	"github.com/republicprotocol/republic-go/identity"
)

type TestSigner struct {
	identity.ID
	identity.KeyPair
}

func (t *TestSigner) Sign(b []byte) (Signature, error) {
	return t.ID.String
}

func NewTestSigner() (TestSigner, error) {
	id, kp, err := identity.NewID()
	if err != nil {
		return TestSigner{}, err
	}
	return TestSigner{
		id,
		kp,
	}
}

var _ = Describe("Processors", func() {

	signer, err := 

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
