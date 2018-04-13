package hyper_test

import (
	"context"
	"log"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Hyperdrive", func() {

	threshold := uint8(160)
	commanderCount := uint8(240)

	Context("Hyperdrive", func() {

		FIt("Achieves consensus on a block over 240 commanders with 75% threshold", func() {
			ctx, cancel := context.WithCancel(context.Background())

			// Network
			ingress := make([]ChannelSet, commanderCount)
			egress := make([]ChannelSet, commanderCount)

			for i := uint8(0); i < commanderCount; i++ {
				ingress[i] = EmptyChannelSet(ctx, commanderCount)
				egress[i] = EmptyChannelSet(ctx, commanderCount)
			}

			for i := uint8(0); i < commanderCount; i++ {
				go egress[i].Split(ingress)
			}

			log.Println("Network initialized.... ")

			// Hyperdrive

			// Initialize replicas
			replicas := make([]Replica, commanderCount)
			for i := uint8(0); i < commanderCount; i++ {
				blocks := NewSharedBlocks(0, 0)
				validator, _ := NewTestValidator(blocks, threshold)
				replicas[i] = NewReplica(ctx, validator, ingress[i])
			}

			// Run replicas
			for i := uint8(0); i < commanderCount; i++ {
				go egress[i].Copy(replicas[i].Run())
			}

			log.Println("Starting the hyperdrive.... ")

			// Broadcast proposal to all the nodes
			proposal := Proposal{
				Signature("Proposal"),
				Block{
					Tuples{},
					Signature("Proposal"),
				},
				Rank(0),
				0,
			}

			for i := 0; i < len(replicas); i++ {
				ingress[i].Proposal <- proposal
			}

			log.Println("Broadcasted the proposals")

			// Wait for the blocks from all the nodes
			var wg sync.WaitGroup
			for i := uint8(0); i < commanderCount; i++ {
				wg.Add(1)
				go func(i uint8) {
					defer wg.Done()
					_ = <-egress[i].Block
					log.Println("Block recieved on", i)
				}(i)
			}

			log.Println("Waiting for the blocks")
			wg.Wait()
			log.Println("Success!!!!!")
			cancel()

		})

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
