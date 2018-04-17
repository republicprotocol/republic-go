package hyperdrive_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Faults", func() {

	Context("when processing faults", func() {

		It("should produce fault when reach threshold", func() {
			ctx, cancel := context.WithCancel(context.Background())
			capacity := 100
			signer := NewWeakSigner(WeakSignerID)
			verifier := NewWeakVerifier()

			faultChIn := make(chan Fault)
			go func() {
				fault := Fault{
					Rank : 1,
					Height: 1,
				}
				faultChIn <- fault
			}()

			faultCh, _ := ProcessFault(ctx, faultChIn, &signer, &verifier, capacity)
			fault, ok := <-faultCh
			Ω(fault).ShouldNot(BeNil())
			Ω(ok).Should(BeTrue())

			close(faultChIn)
			cancel()
		})
	})
})
