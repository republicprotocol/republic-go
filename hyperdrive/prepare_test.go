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
			capacity, threshold := 100, 100
			signer := NewWeakSigner(WeakSignerID)
			verifier := NewWeakVerifier()

			prepareChIn := make(chan Prepare)
			go func() {
				defer close(prepareChIn)

				for i := 0; i < threshold; i++ {
					prepare := Prepare{
						Proposal: Proposal{
							Block: Block{},
						},
					}
					prepareChIn <- prepare
				}
			}()

			commitCh, _, _ := ProcessPreparation(context.Background(), prepareChIn, &signer, &verifier, capacity, threshold)
			commit := <- commitCh
			Ω(commit).ShouldNot(BeNil())
		})
	})
})
