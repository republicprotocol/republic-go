package hyper_test

import (
	"log"
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Buffer", func() {

	Context("Proposals", func() {

		// FIt("should only return proposals of current static height", func() {
		// 	sb := NewSharedBlocks(0, 0)
		// 	chanSetIn := EmptyChannelSet(100)
		// 	validator, _ := NewTestValidator(sb, 100)
		// 	chanSetOut := ProcessBuffer(chanSetIn, validator)

		// 	go func() {
		// 		for {
		// 			select {
		// 			case proposal, ok := <-chanSetOut.Proposal:
		// 				if !ok {
		// 					return
		// 				}
		// 				Ω(proposal.Height).Should(Equal(validator.SharedBlocks().ReadHeight()))
		// 			}
		// 		}
		// 	}()

		// 	for i := 0; i < 100; i++ {
		// 		h := rand.Intn(4)
		// 		chanSetIn.Proposal <- Proposal{
		// 			Height: uint64(h),
		// 			Rank:   Rank(1),
		// 			Block:  Block{},
		// 		}
		// 	}

		// 	chanSetIn.Close()

		// })

		FIt("should only return proposals of current dynamic height which changes every second", func() {
			chanSetIn := EmptyChannelSet(100)
			defer chanSetIn.Close()

			sb := NewSharedBlocks(0, 0)
			validator, _ := NewTestValidator(sb, 100)
			chanSetOut := ProcessBuffer(chanSetIn, validator)
			counter := map[uint64]uint64{}
			var wg sync.WaitGroup

			randcounter := map[int]int{}
			for i := 0; i < 100; i++ {
				h := rand.Intn(5)
				randcounter[h]++
				chanSetIn.Proposal <- Proposal{
					Height: uint64(h),
					Rank:   Rank(1),
					Block:  Block{},
				}
				log.Println("Random Counter", randcounter[0], randcounter[1], randcounter[2], randcounter[3], randcounter[4], randcounter[0]+randcounter[1]+randcounter[2]+randcounter[3]+randcounter[4])
			}

			go func() {
				defer GinkgoRecover()
				defer log.Println("Channel closed")
				for {
					select {
					case proposal, ok := <-chanSetOut.Proposal:
						if !ok {
							return
						}
						counter[proposal.Height]++
						Ω(proposal.Height).Should(Equal(validator.SharedBlocks().ReadHeight()))
						log.Println("Counter", counter[0], counter[1], counter[2], counter[3], counter[4], counter[0]+counter[1]+counter[2]+counter[3]+counter[4])
					}
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				defer GinkgoRecover()
				defer time.Sleep(1 * time.Second)
				for i := 0; i < 5; i++ {
					time.Sleep(2 * time.Second)
					validator.SharedBlocks().IncrementHeight()
				}
			}()

			wg.Wait()

			Ω(uint64(100)).Should(Equal(counter[0] + counter[1] + counter[2] + counter[3] + counter[4]))

		})

	})
})
