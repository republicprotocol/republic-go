package hyperdrive_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var WeakSignerID = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

var _ = Describe("Prepare", func() {

	Context("when processing prepares", func() {

		It("should produce proposal", func() {
			ctx, cancel := context.WithCancel(context.Background())
			capacity, threshold := 100, 100
			signer := NewWeakSigner(WeakSignerID)
			verifier := NewWeakVerifier()

			prepareChIn := make(chan Prepare)
			go func() {
				for i := 0; i < threshold; i++ {
					prepare := Prepare{
						Proposal: Proposal{
							Block: Block{Height: Height(1)},
						},
						Signatures: Signatures{Signature([65]byte{byte(i)})},
					}
					prepareChIn <- prepare
				}
			}()

			commitCh, _, _ := ProcessPreparation(ctx, prepareChIn, &signer, &verifier, capacity, threshold)
			commit, ok := <-commitCh
			Ω(commit).ShouldNot(BeNil())
			Ω(ok).Should(BeTrue())

			close(prepareChIn)
			cancel()
		})
	})
})
