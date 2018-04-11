package hyper_test

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/republicprotocol/republic-go/hyperdrive"
	"golang.org/x/crypto/sha3"
)

var _ = Describe("Hyperdrive", func() {

	commanderCount := uint8(240)

	Context("Hyperdrive", func() {

		// It("Achieves consensus on a block over 200 commanders with 75% threshold", func() {
		// 	hyper := NewHyperDrive(commanderCount)
		// 	hyper.init()
		// 	proposal := Proposal{
		// 		Signature("Proposal"),
		// 		Block{
		// 			Tuples{},
		// 			Signature("Proposal"),
		// 		},
		// 		Rank(0),
		// 		Height(0),
		// 	}
		// 	hyper.network.propose(proposal)
		// 	ctx, cancel := context.WithCancel(context.Background())
		// 	hyper.run(ctx)
		// 	hyper.network.run(ctx)
		// 	time.Sleep(2 * time.Second)
		// 	cancel()
		// })

		FIt("Achieves consensus 50 blocks over 240 commanders with 2/3 threshold", func() {
			numberOfBlocks := 50
			hyperdrive := NewHyperDrive(commanderCount)
			hyperdrive.init()
			proposals := make([]Proposal, numberOfBlocks)
			for i := 0; i < numberOfBlocks; i++ {
				tuple := Tuple{
					ID: sha3.Sum256([]byte(strconv.Itoa(i))),
				}
				proposals[i] = Proposal{
					Signature("Proposal"),
					Block{
						Tuples:    Tuples{tuple},
						Signature: Signature("Proposal"),
					},
					Rank(1),
					uint64(i),
				}
			}
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer log.Println("Proposed multiple blocks")
				defer GinkgoRecover()
				defer time.Sleep(1 * time.Minute)

				hyperdrive.network.proposeMultiple(proposals)
			}()
			ctx, cancel := context.WithCancel(context.Background())
			go hyperdrive.run(ctx)
			wg.Wait()
			cancel()
		})
	})
})
