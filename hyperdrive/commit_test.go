package hyperdrive_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Commit", func() {

	Context("when processing commits", func() {

		It("should create commit if everything goes fine ", func() {
			ctx, cancel := context.WithCancel(context.Background())
			capacity, threshold := 100, 100
			signer := NewWeakSigner(WeakSignerID)
			verifier := NewWeakVerifier()

			commitChIn := make(chan Commit)
			go func() {
				for i := 0; i < threshold; i++ {
					commit := Commit{
						Prepare: Prepare{
							Proposal: Proposal{
								Block: Block{Height: Height(1)},
							},
						},
						Signatures: Signatures{Signature([65]byte{byte(i)})},
					}
					commitChIn <- commit
				}
			}()

			commitCh, _, _ := ProcessCommits(ctx, commitChIn, &signer, &verifier, capacity, threshold)
			commit, ok := <-commitCh
			Ω(commit).ShouldNot(BeNil())
			Ω(ok).Should(BeTrue())

			close(commitChIn)
			cancel()
		})

		It("should create fault if cannot be verified ", func() {
			ctx, cancel := context.WithCancel(context.Background())
			capacity, threshold := 100, 100
			signer := NewWeakSigner(WeakSignerID)
			verifier := NewErrorVerifier()

			commitChIn := make(chan Commit)
			go func() {
				commit := Commit{
					Prepare: Prepare{
						Proposal: Proposal{
							Block: Block{Height: Height(1)},
						},
					},
					Signatures: Signatures{Signature([65]byte{byte(0)})},
				}
				commitChIn <- commit
			}()

			_, faultCh, _ := ProcessCommits(ctx, commitChIn, &signer, &verifier, capacity, threshold)
			fault, ok := <-faultCh
			Ω(fault).ShouldNot(BeNil())
			Ω(ok).Should(BeTrue())

			close(commitChIn)
			cancel()
		})
	})
})
