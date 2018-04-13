package hyper_test

import (
	. "github.com/onsi/ginkgo"
	//	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Hyperdrive", func() {

	// commanderCount := uint8(240)

	Context("Hyperdrive", func() {

		// It("Achieves consensus on a block over 240 commanders with 75% threshold", func() {
		// 	ctx, cancel := context.WithCancel(context.Background())

		// 	hyper := NewHyperDrive(ctx, commanderCount)
		// 	hyper.init()
		// 	var wg sync.WaitGroup
		// 	proposal := Proposal{
		// 		Signature("Proposal"),
		// 		Block{
		// 			Tuples{},
		// 			Signature("Proposal"),
		// 		},
		// 		Rank(0),
		// 		0,
		// 	}

		// 	wg.Add(1)
		// 	go func() {
		// 		defer wg.Done()
		// 		hyper.run()
		// 	}()

		// 	wg.Add(1)
		// 	go func() {
		// 		defer wg.Done()
		// 		defer log.Println("Finished sending proposals")
		// 		hyper.network.propose(proposal)
		// 	}()

		// 	go func() {
		// 		defer cancel()
		// 		defer log.Println("Success!!")
		// 		for i := uint8(0); i < commanderCount; i++ {
		// 			wg.Add(1)
		// 			go func(i uint8) {
		// 				defer wg.Done()
		// 				_ = <-hyper.network.Egress[i].Block
		// 			}(i)
		// 		}
		// 		wg.Wait()
		// 	}()
		// 	time.Sleep(1 * time.Minute)
		// 	cancel()
		// 	log.Println("Waiting here")
		// 	wg.Wait()
		// })

		// 	FIt("Achieves consensus 50 blocks over 240 commanders with 2/3 threshold", func() {
		// 		numberOfBlocks := 50
		// 		hyperdrive := NewHyperDrive(commanderCount)
		// 		hyperdrive.init()
		// 		proposals := make([]Proposal, numberOfBlocks)
		// 		for i := 0; i < numberOfBlocks; i++ {
		// 			tuple := Tuple{
		// 				ID: sha3.Sum256([]byte(strconv.Itoa(i))),
		// 			}
		// 			proposals[i] = Proposal{
		// 				Signature("Proposal"),
		// 				Block{
		// 					Tuples:    Tuples{tuple},
		// 					Signature: Signature("Proposal"),
		// 				},
		// 				Rank(1),
		// 				uint64(i),
		// 			}
		// 		}
		// 		var wg sync.WaitGroup
		// 		wg.Add(1)
		// 		go func() {
		// 			defer wg.Done()
		// 			defer log.Println("Proposed multiple blocks")
		// 			defer GinkgoRecover()
		// 			defer time.Sleep(1 * time.Minute)

		// 			hyperdrive.network.proposeMultiple(proposals)
		// 		}()
		// 		ctx, cancel := context.WithCancel(context.Background())
		// 		go hyperdrive.run(ctx)
		// 		wg.Wait()
		// 		cancel()
		// 	})
	})
})
