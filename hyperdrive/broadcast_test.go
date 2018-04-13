package hyper_test

import (
	"context"
	"log"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Broadcast", func() {

	Context("Proposals", func() {

		It("should recieve all the broadcasted proposals", func() {
			ctx, cancel := context.WithCancel(context.Background())
			chanSetIn := EmptyChannelSet(ctx, 100)
			validator, _ := NewTestValidator(NewSharedBlocks(0, 0), 100)
			chanSetOut := ProcessBroadcast(chanSetIn, validator)

			statusMu := new(sync.RWMutex)
			status := map[[32]byte]uint8{}
			var wg sync.WaitGroup

			wg.Add(101)
			go func() {
				defer cancel()
				defer wg.Done()
				defer log.Println("Ended the first function")
				for i := 0; i < 100; i++ {
					proposal := Proposal{
						Height: uint64(i),
					}
					statusMu.Lock()
					status[ProposalHash(proposal)] = 1
					statusMu.Unlock()
					chanSetIn.Proposal <- proposal
				}
			}()

			go func() {
				defer log.Println("Ended the second function")
				for {
					select {
					case proposal, ok := <-chanSetOut.Proposal:
						if !ok {
							return
						}
						statusMu.RLock()
						Ω(status[ProposalHash(proposal)]).Should(Equal(uint8(1)))
						statusMu.RUnlock()

						statusMu.Lock()
						status[ProposalHash(proposal)]++
						statusMu.Unlock()

						wg.Done()
					}
				}
			}()

			wg.Wait()
		})

		It("should only return unique proposals", func() {
			ctx, cancel := context.WithCancel(context.Background())
			chanSetIn := EmptyChannelSet(ctx, 100)
			validator, _ := NewTestValidator(NewSharedBlocks(0, 0), 100)
			chanSetOut := ProcessBroadcast(chanSetIn, validator)

			statusMu := new(sync.RWMutex)
			status := map[[32]byte]uint8{}
			var wg sync.WaitGroup

			wg.Add(100)
			go func() {
				defer cancel()
				defer wg.Done()
				for i := 0; i < 100; i++ {
					for j := 0; j < i; j++ {
						proposal := Proposal{
							Height: uint64(i),
						}
						statusMu.Lock()
						status[ProposalHash(proposal)] = 1
						statusMu.Unlock()
						chanSetIn.Proposal <- proposal
					}
				}
			}()

			go func() {

				for {
					select {
					case proposal, ok := <-chanSetOut.Proposal:
						if !ok {
							return
						}
						statusMu.RLock()
						Ω(status[ProposalHash(proposal)]).Should(Equal(uint8(1)))
						statusMu.RUnlock()

						statusMu.Lock()
						status[ProposalHash(proposal)]++
						statusMu.Unlock()

						wg.Done()
					}
				}
			}()

			wg.Wait()
		})
	})
})
