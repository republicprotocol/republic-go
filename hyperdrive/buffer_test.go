package hyper_test

import (
	"context"
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

		It("should only return proposals of current static height", func() {
			ctx, cancel := context.WithCancel(context.Background())
			sb := NewSharedBlocks(0, 0)
			chanSetIn := EmptyChannelSet(ctx, 100)
			validator, _ := NewTestValidator(sb, 100)
			chanSetOut := ProcessBuffer(chanSetIn, validator)

			go func() {
				for {
					select {
					case proposal, ok := <-chanSetOut.Proposal:
						if !ok {
							return
						}
						Ω(proposal.Height).Should(Equal(validator.SharedBlocks().ReadHeight()))
					}
				}
			}()

			for i := 0; i < 100; i++ {
				h := rand.Intn(4)
				chanSetIn.Proposal <- Proposal{
					Height: uint64(h),
					Rank:   Rank(1),
					Block:  Block{},
				}
			}

			cancel()

		})

		It("should only return proposals of current dynamic height which changes every second", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			chanSetIn := EmptyChannelSet(ctx, 100)

			sb := NewSharedBlocks(0, 0)
			validator, _ := NewTestValidator(sb, 100)
			chanSetOut := ProcessBuffer(chanSetIn, validator)

			counterMu := new(sync.RWMutex)
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
			}

			wg.Add(1)
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer log.Println("Channel closed")
				for {
					select {
					case proposal, ok := <-chanSetOut.Proposal:
						if !ok {
							return
						}
						counterMu.Lock()
						counter[proposal.Height]++
						Ω(proposal.Height).Should(Equal(validator.SharedBlocks().ReadHeight()))
						counterMu.Unlock()
					}
				}
			}()

			go func() {
				defer wg.Done()
				defer GinkgoRecover()
				for i := 0; i < 5; i++ {
					time.Sleep(2 * time.Second)
					validator.SharedBlocks().IncrementHeight()
				}
			}()

			wg.Wait()
			counterMu.RLock()
			Ω(uint64(100)).Should(Equal(counter[0] + counter[1] + counter[2] + counter[3] + counter[4]))
			counterMu.RUnlock()
		})

	})
})
