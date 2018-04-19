package hyperdrive_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Proposals", func() {

	Context("when processing proposal", func() {

		It("should produce prepare if everything goes fine ", func() {
			ctx, cancel := context.WithCancel(context.Background())
			capacity := 100
			signer := NewWeakSigner(WeakSignerID)
			verifier := NewWeakVerifier()

			proposalChIn := make(chan Proposal)
			go func() {
				proposal := Proposal{Block: Block{Height: Height(1)}}
				proposalChIn <- proposal
			}()

			prepareCh, _, _ := ProcessProposal(ctx, proposalChIn, &signer, &verifier, capacity)
			prepare, ok := <-prepareCh
			Ω(prepare).ShouldNot(BeNil())
			Ω(ok).Should(BeTrue())

			close(proposalChIn)
			cancel()
		})

		It("should produce a fault if it cannot be verified", func() {
			ctx, cancel := context.WithCancel(context.Background())
			capacity := 100
			signer := NewWeakSigner(WeakSignerID)
			verifier := NewErrorVerifier()

			proposalChIn := make(chan Proposal)
			go func() {
				proposal := Proposal{Block: Block{Height: Height(1)}}
				proposalChIn <- proposal
			}()

			_, faultCh, _ := ProcessProposal(ctx, proposalChIn, &signer, &verifier, capacity)
			fault, ok := <-faultCh
			Ω(fault).ShouldNot(BeNil())
			Ω(ok).Should(BeTrue())

			close(proposalChIn)
			cancel()
		})
	})

})
